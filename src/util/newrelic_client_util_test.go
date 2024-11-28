// Package util provides utility functions for the AWS Unified Logging library.

package util

import (
	"context"
	"github.com/newrelic/aws-unified-lambda-logging/common"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockNRClient is a mock type for the Logs interface.
type MockNRClient struct {
	mock.Mock
}

// CreateLogEntry is a mock method that satisfies the Logs interface.
func (m *MockNRClient) CreateLogEntry(batch interface{}) error {
	args := m.Called(batch)
	return args.Error(0)
}

// newNRClientTestCase represents a test case for the NewNRClient function.
type newNRClientTestCase struct {
	name             string // Name of the test case
	envDebug         string // Environment variable for debug
	envRegion        string // Environment variable for region
	envLicenseKey    string // Environment variable for license key
	expectedLogLevel string // Expected log level
	expectError      bool   // Whether an error is expected
}

// TestNewNRClient tests the NewNRClient function with different scenarios.
func TestNewNRClient(t *testing.T) {
	testCases := []newNRClientTestCase{
		{
			name:             "Debug enabled",
			envDebug:         "true",
			envRegion:        "us",
			envLicenseKey:    "valid_license_key",
			expectedLogLevel: "debug",
			expectError:      false,
		},
		{
			name:             "Debug disabled",
			envRegion:        "us",
			envLicenseKey:    "valid_license_key",
			expectedLogLevel: "info",
			expectError:      false,
		},
		{
			name:          "Invalid region",
			envRegion:     "invalid",
			envLicenseKey: "valid_license_key",
			expectError:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup environment variables
			if tc.envDebug != "" {
				os.Setenv(common.DebugEnabled, tc.envDebug)
				defer os.Unsetenv(common.DebugEnabled)
			}
			if tc.envRegion != "" {
				os.Setenv(common.NewRelicRegion, tc.envRegion)
				defer os.Unsetenv(common.NewRelicRegion)
			}
			if tc.envLicenseKey != "" {
				os.Setenv(common.EnvLicenseKey, tc.envLicenseKey)
				defer os.Unsetenv(common.EnvLicenseKey)
			}

			nrClient, err := NewNRClient()

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, nrClient)
			}
		})
	}
}

// TestConsumeLogBatches tests the ConsumeLogBatches function.
func TestConsumeLogBatches(t *testing.T) {
	mockNRClient := new(MockNRClient)
	mockNRClient.On("CreateLogEntry", mock.Anything).Return(nil)

	channel := make(chan common.DetailedLogsBatch, 1)
	wg := new(sync.WaitGroup)

	logBatch := []common.DetailedLog{{
		CommonData: common.Common{
			Attributes: common.LogAttributes{
				"awsAccountId": "123456789012",
			},
		},
		Entries: []common.Log{{
			Timestamp: "1234567890",
			Log:       "test log",
		}},
	}}

	channel <- logBatch

	ctx := context.TODO()
	wg.Add(1)
	go ConsumeLogBatches(ctx, channel, wg, mockNRClient)
	close(channel)
	wg.Wait()
	mockNRClient.AssertNumberOfCalls(t, "CreateLogEntry", 1)
}
