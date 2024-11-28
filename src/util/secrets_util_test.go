// Package util provides utility functions for working with secrets in AWS Lambda functions.
package util

import (
	"context"
	"errors"
	"github.com/newrelic/aws-unified-lambda-logging/common"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSecretManagerClient is a mock implementation of the SecretManagerAPI
type MockSecretManagerClient struct {
	mock.Mock
}

// GetSecretValue provides a mock implementation to retrieve the value of a secret from the secret manager.
// It takes a context, input parameters, and optional functional options.
func (m *MockSecretManagerClient) GetSecretValue(ctx context.Context, params *secretsmanager.GetSecretValueInput, optFns ...func(*secretsmanager.Options)) (*secretsmanager.GetSecretValueOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*secretsmanager.GetSecretValueOutput), args.Error(1)
}

// TestGetSecretFromSecretManager is a unit test function that tests the GetSecretFromSecretManager function.
// The function uses a table-driven approach to run multiple test cases.
func TestGetSecretFromSecretManager(t *testing.T) {
	// Test cases for different scenarios
	tests := []struct {
		name           string                         // Name of the test case
		secretName     string                         // Name of the secret
		setupMock      func(*MockSecretManagerClient) // Setup function to mock the SecretManagerClient
		expectedResult map[string]string              // Expected result of the function
		expectedError  error                          // Expected error from the function
	}{
		// Test case: Empty secret name
		{
			name:          "Empty secret name",
			secretName:    "",
			setupMock:     func(m *MockSecretManagerClient) {},
			expectedError: errors.New("secret name is empty"),
		},
		// Test case: Secrets manager error
		{
			name:       "Secrets manager error",
			secretName: "test-secret",
			setupMock: func(m *MockSecretManagerClient) {
				m.On("GetSecretValue", mock.Anything, mock.Anything).Return(&secretsmanager.GetSecretValueOutput{}, errors.New("secrets manager error"))
			},
			expectedError: errors.New("secrets manager error"),
		},
		// Test case: Successful secret retrieval
		{
			name:       "Successful secret retrieval",
			secretName: "test-secret",
			setupMock: func(m *MockSecretManagerClient) {
				secretString := `{"username":"admin","password":"secret"}` // JSON representation of the secret
				m.On("GetSecretValue", mock.Anything, mock.Anything).Return(&secretsmanager.GetSecretValueOutput{
					SecretString: &secretString,
				}, nil)
			},
			expectedResult: map[string]string{"username": "admin", "password": "secret"},
		},
		// Test case: JSON unmarshal error
		{
			name:       "JSON unmarshal error",
			secretName: "test-secret",
			setupMock: func(m *MockSecretManagerClient) {
				secretString := `{"username":admin,"password":"secret"}` // Malformed JSON
				m.On("GetSecretValue", mock.Anything, mock.Anything).Return(&secretsmanager.GetSecretValueOutput{
					SecretString: &secretString,
				}, nil)
			},
			expectedError: errors.New("invalid character 'a' looking for beginning of value"),
		},
		// Test case: Secret in binary format
		{
			name:       "Secret in binary format",
			secretName: "test-secret",
			setupMock: func(m *MockSecretManagerClient) {
				m.On("GetSecretValue", mock.Anything, mock.Anything).Return(&secretsmanager.GetSecretValueOutput{
					SecretString: nil,
				}, nil)
			},
			expectedError: errors.New("secret is in binary format (likely encrypted)"),
		},
	}

	// Iterate over the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSecretManagerClient := new(MockSecretManagerClient)
			tt.setupMock(mockSecretManagerClient)

			// Call the function under test
			result, err := GetSecretFromSecretManager(context.Background(), mockSecretManagerClient, tt.secretName)

			// Check the expected result and error
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			mockSecretManagerClient.AssertExpectations(t)
		})
	}
}

// TestGetLicenseKey is a unit test function that tests the GetLicenseKey function.
// It tests various scenarios such as empty secret name, environment variable value, and expected result.
// The function uses a table-driven approach to run multiple test cases.
// Each test case includes a name, secret name, setup function for mocking the SecretManagerClient,
// expected result, and expected error.
// The function asserts the expected error and result using the testify/assert package.
func TestGetLicenseKey(t *testing.T) {
	tests := []struct {
		name           string                         // Name of the test case
		secretName     string                         // Name of the secret
		setupMock      func(*MockSecretManagerClient) // Setup function to mock the SecretManagerClient
		expectedResult string                         // Expected result of the function
		expectedError  error                          // Expected error from the function
		envSecretValue string                         // Value of the environment variable
	}{
		{
			name:           "Empty secret name",
			secretName:     "",
			setupMock:      func(m *MockSecretManagerClient) {},
			envSecretValue: "secret_value",
			expectedResult: "secret_value",
		},
		{
			name:          "Empty secret name",
			secretName:    "",
			setupMock:     func(m *MockSecretManagerClient) {},
			expectedError: errors.New("secret name is empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envSecretValue != "" {
				os.Setenv(common.EnvLicenseKey, tt.envSecretValue)
				defer os.Unsetenv(common.EnvLicenseKey)
			}
			if tt.secretName != "" {
				os.Setenv(common.NewRelicLicenseKeySecretName, tt.secretName)
				defer os.Unsetenv(common.NewRelicLicenseKeySecretName)
			}
			result, err := GetLicenseKey()
			if tt.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}
