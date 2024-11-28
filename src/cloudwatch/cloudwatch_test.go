package cloudwatch

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/newrelic/aws-unified-lambda-logging/common"
	"github.com/newrelic/aws-unified-lambda-logging/util"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
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
		logEvents       []events.CloudwatchLogsLogEvent // Log events to process
		expectedBatches int                             // Expected number of batches
	}{
		{
			name: "Success with single batch",
			logEvents: []events.CloudwatchLogsLogEvent{
				{Message: "test message 1", Timestamp: time.Now().Unix()},
				{Message: "test message 2", Timestamp: time.Now().Unix()},
			},
			expectedBatches: 1,
		},
		{
			name: "Success with multiple batches",
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
			logEvents:       []events.CloudwatchLogsLogEvent{},
			expectedBatches: 0,
		},
		{
			name: "log event with a single message to check if batching works",
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
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			cloudwatchLogsData := events.CloudwatchLogsData{
				LogGroup:  "test-log-group",
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
