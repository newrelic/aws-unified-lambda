package s3

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/dsnet/compress/bzip2"
	"github.com/newrelic/aws-unified-lambda-logging/common"
	"github.com/newrelic/aws-unified-lambda-logging/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// s3_test provides unit tests for the s3 package.
// This file contains a mock implementation of the S3 API for testing purposes.

// MockAPI is a struct that represents a mock implementation of the S3 API.
type MockAPI struct {
	mock.Mock
}

// GetObject provides a mock response for the GetObject function of the S3 API.
// It returns the mock GetObjectOutput and an error if any.
func (m *MockAPI) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*s3.GetObjectOutput), args.Error(1)
}

// MockReaderFactory is a mock implementation of the ReaderFactory function type
type MockReaderFactory struct {
	mock.Mock
}

// Create creates a new io.Reader from the given input and filename.
// It returns the created io.Reader and any error encountered during the creation process.
func (m *MockReaderFactory) Create(input io.ReadCloser, filename string) (io.Reader, error) {
	args := m.Called(input, filename)
	return args.Get(0).(io.Reader), args.Error(1)
}

// generateLogsOnCount generates a string containing log entries.
// It takes an integer logCount as input and returns a string with log entries.
func generateLogsOnCount(logCount int) string {
	var logs string

	for i := 1; i <= logCount; i++ {
		logs = logs + fmt.Sprintf("Log entry %d\n", i)
	}

	return logs
}

// generateLogOnSize generates a string of a given size.
func generateLogOnSize(size int) string {
	if size <= 0 {
		return ""
	}
	return strings.Repeat("a", size)
}

// CloudTrailRecord represents the structure of a CloudTrail event.
type CloudTrailTestRecord struct {
	EventVersion string `json:"eventVersion"`
	EventName    string `json:"eventName"`
	EventTime    string `json:"eventTime"`
	AWSRegion    string `json:"awsRegion"`
}

// CloudTrailRecords represents the structure of the JSON containing the "Records" array.
type CloudTrailRecords struct {
	Records []CloudTrailTestRecord `json:"Records"`
}

// generateCloudTrailLogs creates a CloudTrail log with a specified number of records.
func generateCloudTrailTestLogs(size int) string {
	records := make([]CloudTrailTestRecord, size)
	for i := 0; i < size; i++ {
		records[i] = CloudTrailTestRecord{
			EventVersion: "1.0",
			EventName:    "ExampleEventName",
			EventTime:    time.Now().UTC().Format(time.RFC3339),
			AWSRegion:    "us-east-1",
		}
	}

	cloudTrailLogs := CloudTrailRecords{
		Records: records,
	}

	jsonData, _ := json.Marshal(cloudTrailLogs)
	return string(jsonData)
}

// TestGetLogsFromS3Event is a unit test function that tests the GetLogsFromS3Event function.
// It tests different scenarios of S3 event processing and verifies the expected results.
func TestGetLogsFromS3Event(t *testing.T) {
	tests := []struct {
		name          string                   // Name of the test case
		setupS3Mock   func(*MockAPI)           // Function to set up the S3 mock
		setupRFMock   func(*MockReaderFactory) // Function to set up the ReaderFactory mock
		expectedError error                    // Expected error from the function
		batchSize     int                      // Expected number of batches
		URLDecodedKey string                   // URLDecodedKey of the S3 object
	}{
		{
			name: "Successful S3 event processing",
			setupS3Mock: func(m *MockAPI) {
				m.On("GetObject", mock.Anything, mock.Anything).Return(&s3.GetObjectOutput{
					Body: io.NopCloser(bytes.NewReader([]byte("log content"))),
				}, nil)
			},
			setupRFMock: func(m *MockReaderFactory) {
				m.On("Create", mock.Anything, "test-key.gz").Return(strings.NewReader("log content"), nil)
			},
			expectedError: nil,
			batchSize:     1,
			URLDecodedKey: "test-key.gz",
		},
		{
			name: "Error fetching S3 object",
			setupS3Mock: func(m *MockAPI) {
				m.On("GetObject", mock.Anything, mock.Anything).Return(&s3.GetObjectOutput{}, errors.New("s3 error"))
			},
			setupRFMock:   func(m *MockReaderFactory) {},
			expectedError: errors.New("s3 error"),
			URLDecodedKey: "test-key.gz",
		},
		{
			name: "Successful S3 event processing. Used to test the maximum number of messages in a batch.",
			setupS3Mock: func(m *MockAPI) {
				m.On("GetObject", mock.Anything, mock.Anything).Return(&s3.GetObjectOutput{
					Body: io.NopCloser(bytes.NewReader([]byte(generateLogsOnCount(common.MaxPayloadMessages + 10)))),
				}, nil)
			},
			setupRFMock: func(m *MockReaderFactory) {
				m.On("Create", mock.Anything, "test-key.gz").Return(strings.NewReader(generateLogsOnCount(common.MaxPayloadMessages+10)), nil)
			},
			expectedError: nil,
			batchSize:     2,
			URLDecodedKey: "test-key.gz",
		},
		{
			name: "Successful S3 event processing. Used to test the maximum payload size in a batch.",
			setupS3Mock: func(m *MockAPI) {
				m.On("GetObject", mock.Anything, mock.Anything).Return(&s3.GetObjectOutput{
					Body: io.NopCloser(bytes.NewReader([]byte(generateLogOnSize(1024*1024*1 + 10)))),
				}, nil)
			},
			setupRFMock: func(m *MockReaderFactory) {
				m.On("Create", mock.Anything, "test-key.gz").Return(strings.NewReader(generateLogOnSize(1024*1024*1+10)), nil)
			},
			expectedError: nil,
			batchSize:     2,
			URLDecodedKey: "test-key.gz",
		},
		{
			name:          "CloudTrail Digest Ignore Scenario.",
			setupS3Mock:   func(m *MockAPI) {},
			setupRFMock:   func(m *MockReaderFactory) {},
			expectedError: nil,
			batchSize:     0,
			URLDecodedKey: "test-key_CloudTrail-Digest_2021-09-01T00-00-00Z.json.gz",
		},
		{
			name: "Reading CloudTrail logs from S3 event",
			setupS3Mock: func(m *MockAPI) {
				m.On("GetObject", mock.Anything, mock.Anything).Return(&s3.GetObjectOutput{
					Body: io.NopCloser(bytes.NewReader([]byte(generateCloudTrailTestLogs(4)))),
				}, nil)
			},
			setupRFMock: func(m *MockReaderFactory) {
				m.On("Create", mock.Anything, "test-key_CloudTrail_2021-09-01T00-00-00Z.json.gz").Return(strings.NewReader(generateCloudTrailTestLogs(4)), nil)
			},
			expectedError: nil,
			batchSize:     1,
			URLDecodedKey: "test-key_CloudTrail_2021-09-01T00-00-00Z.json.gz",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			awsConfiguration := util.AWSConfiguration{
				AccountID: "123456789012",
				Realm:     "aws",
				Region:    "us-west-2",
			}

			channel := make(chan common.DetailedLogsBatch, 2)
			s3Event := events.S3Event{
				Records: []events.S3EventRecord{
					{
						S3: events.S3Entity{
							Bucket: events.S3Bucket{
								Name: "test-bucket",
							},
							Object: events.S3Object{
								URLDecodedKey: tc.URLDecodedKey,
							},
						},
					},
				},
			}

			mockS3Client := new(MockAPI)
			tc.setupS3Mock(mockS3Client)

			mockReaderFactory := new(MockReaderFactory)
			tc.setupRFMock(mockReaderFactory)

			// Call the GetLogsFromS3Event function
			err := GetLogsFromS3Event(ctx, s3Event, awsConfiguration, channel, mockS3Client, mockReaderFactory.Create)
			close(channel)

			// Check for expected errors
			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.EqualError(t, err, tc.expectedError.Error())
			} else {
				assert.NoError(t, err)
				batchCount := 0
				for batch := range channel {
					assert.NotEmpty(t, batch)
					batchCount++
				}
				assert.Equal(t, tc.batchSize, batchCount)
			}

			// Assert that all expectations were met
			mockS3Client.AssertExpectations(t)
			mockReaderFactory.AssertExpectations(t)
		})
	}
}

// TestGetLogsFromSNSEvent is a unit test function that tests the GetLogsFromSNSEvent function.
// It tests different scenarios of SNS event processing and verifies the expected results.
// func TestGetLogsFromSNSEvent(t *testing.T) {
//todo
// }

// compressBzip2 compresses the given data using the bzip2 algorithm and returns the compressed data.
// It uses the bzip2.WriterConfig with the BestCompression level for optimal compression.
func compressBzip2(data []byte) []byte {
	var buf bytes.Buffer
	bw, _ := bzip2.NewWriter(&buf, &bzip2.WriterConfig{Level: bzip2.BestCompression})
	_, _ = bw.Write(data)
	_ = bw.Close()
	return buf.Bytes()
}

// compressGzip compresses the given data using gzip compression algorithm.
// It takes a byte slice as input and returns the compressed data as a byte slice.
func compressGzip(data []byte) []byte {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	_, _ = gw.Write(data)
	_ = gw.Close()
	return buf.Bytes()
}

// TestDefaultReaderFactory is a test function that tests the behavior of the DefaultReaderFactory function.
// It tests the creation of readers for different types of files, including uncompressed files, GZIP compressed files,
// BZIP2 compressed files, and invalid GZIP files. It verifies that the readers are created correctly and that the
// content read from the readers matches the expected content.
func TestDefaultReaderFactory(t *testing.T) {
	// Test cases for different file types and scenarios
	tests := []struct {
		name        string              // Name of the test case
		filename    string              // Filename of the test file
		content     []byte              // Content of the test file
		compress    func([]byte) []byte // Compression function for the test file content
		expectError bool                // Flag indicating whether an error is expected
	}{
		{
			name:     "Uncompressed file",
			filename: "testfile.txt",
			content:  []byte("uncompressed content"),
		},
		{
			name:     "GZIP compressed file",
			filename: "testfile.gz",
			content:  []byte("gzip compressed content"),
			compress: compressGzip,
		},
		{
			name:     "BZIP2 compressed file",
			filename: "testfile.bz2",
			content:  []byte("bzip2 compressed content"),
			compress: compressBzip2,
		},
		{
			name:        "Invalid GZIP file",
			filename:    "invalid.gz",
			content:     []byte("invalid gzip content"),
			compress:    nil,
			expectError: true,
		},
	}

	// Iterate over the test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var input io.ReadCloser
			// Create the input reader based on the compression function
			if tc.compress != nil {
				compressed := tc.compress(tc.content)
				input = io.NopCloser(bytes.NewReader(compressed))
			} else {
				input = io.NopCloser(bytes.NewReader(tc.content))
			}

			// Call the DefaultReaderFactory function
			reader, err := DefaultReaderFactory(input, tc.filename)

			// Check if an error is expected
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				// Read the content from the reader
				output, readErr := io.ReadAll(reader)
				assert.NoError(t, readErr)
				// Verify that the content matches the expected content
				assert.Equal(t, tc.content, output)
			}
		})
	}
}
