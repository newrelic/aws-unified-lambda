package cloudwatch

import (
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/newrelic/aws-unified-lambda-logging/common"
	"github.com/newrelic/aws-unified-lambda-logging/util"
	"github.com/stretchr/testify/assert"
)

// mockAWSConfiguration returns a mock AWSConfiguration object for testing purposes.
func mockAWSConfiguration() util.AWSConfiguration {
	return util.AWSConfiguration{
		AccountID: "123456789012",
		Realm:     "aws",
		Region:    "us-west-2",
	}
}

// TestGetLogs is a unit test function that tests the GetLogs function.
// It verifies the behavior of GetLogs by running multiple test cases.
// Each test case consists of a set of log events and an expected number of batches.
func TestGetLogs(t *testing.T) {
	tests := []struct {
		name            string                          // Name of the test case
		logGroup        string                          // Log group name
		logEvents       []events.CloudwatchLogsLogEvent // Log events to process
		expectedBatches int                             // Expected number of batches
		expectedReqID   string                          // Expected request ID
	}{
		{
			name:     "Success with single batch",
			logGroup: "test-log-group",
			logEvents: []events.CloudwatchLogsLogEvent{
				{Message: "test message 1", Timestamp: time.Now().Unix()},
				{Message: "test message 2", Timestamp: time.Now().Unix()},
			},
			expectedBatches: 1,
		},
		{
			name:     "Success with multiple batches",
			logGroup: "test-log-group",
			logEvents: func() []events.CloudwatchLogsLogEvent {
				var logEvents []events.CloudwatchLogsLogEvent
				for i := 0; i < common.MaxPayloadMessages+1; i++ {
					logEvents = append(logEvents, events.CloudwatchLogsLogEvent{
						Message:   "test message",
						Timestamp: time.Now().Unix(),
					})
				}
				return logEvents
			}(),
			expectedBatches: 2,
		},
		{
			name:            "Empty log data",
			logGroup:        "test-log-group",
			logEvents:       []events.CloudwatchLogsLogEvent{},
			expectedBatches: 0,
		},
		{
			name:     "log event with a single message to check if batching works",
			logGroup: "test-log-group",
			logEvents: func() []events.CloudwatchLogsLogEvent {
				var logEvents []events.CloudwatchLogsLogEvent
				logEvents = append(logEvents, events.CloudwatchLogsLogEvent{
					Message:   strings.Repeat("a", 1024*1024*1+10),
					Timestamp: time.Now().Unix(),
				})
				return logEvents
			}(),
			expectedBatches: 2,
		},
		{
			name:     "Lambda log group with request ID",
			logGroup: "/aws/lambda/test-log-group",
			logEvents: []events.CloudwatchLogsLogEvent{
				{Message: "RequestId: d653fb2c-0234-46ff-ae6b-9a418b888420 Start of request", Timestamp: time.Now().Unix()},
				{Message: "Processing request", Timestamp: time.Now().Unix()},
				{Message: "RequestId: d653fb2c-0234-46ff-ae6b-9a418b888420 End of request", Timestamp: time.Now().Unix()},
			},
			expectedBatches: 1,
			expectedReqID:   "d653fb2c-0234-46ff-ae6b-9a418b888420",
		},
		{
			name:     "Non-Lambda log group",
			logGroup: "test-log-group",
			logEvents: []events.CloudwatchLogsLogEvent{
				{Message: "Some log message", Timestamp: time.Now().Unix()},
			},
			expectedBatches: 1,
			expectedReqID:   "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cloudwatchLogsData := events.CloudwatchLogsData{
				LogGroup:  tc.logGroup,
				LogStream: "test-log-stream",
				LogEvents: tc.logEvents,
			}
			awsConfig := mockAWSConfiguration()
			// create a channel to produce messages
			channel := make(chan common.DetailedLogsBatch, 2) // Buffer size of 2 to prevent blocking

			err := GetLogs(cloudwatchLogsData, awsConfig, channel)
			assert.NoError(t, err)

			close(channel)
			var batches []common.DetailedLogsBatch
			for batch := range channel {
				batches = append(batches, batch)
			}

			assert.Equal(t, tc.expectedBatches, len(batches), "Expected number of batches does not match")

			// Check the last request ID if expected
			if tc.expectedReqID != "" {
				for _, batch := range batches {
					for _, log := range batch {
						assert.Equal(t, tc.expectedReqID, log.Entries[0].Attributes["requestId"], "RequestId does not match")
					}
				}
			}

			// iterate through batches consumed from the channel
			for _, batch := range batches {
				for _, log := range batch {
					assert.Equal(t, cloudwatchLogsData.LogGroup, log.CommonData.Attributes["logGroup"], "LogGroup does not match")
					assert.Equal(t, cloudwatchLogsData.LogStream, log.CommonData.Attributes["logStream"], "LogStream does not match")
					assert.NotEmpty(t, log.Entries, "Expected non-empty log entries")
				}
			}
		})
	}
}
