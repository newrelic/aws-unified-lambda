package util

import (
	"encoding/json"
	"reflect"
	"regexp"
	"testing"

	"github.com/newrelic/aws-unified-lambda-logging/common"
	"github.com/stretchr/testify/assert"
)

// TestAddCustomMetaData tests adding custom meta data utility function.
func TestAddCustomMetaData(t *testing.T) {
	tests := []struct {
		name       string                 // Test case name
		jsonString string                 // JSON string to parse
		attributes map[string]interface{} // Attributes to update
		wantErr    bool                   // Expected error
		expected   map[string]interface{} // Expected attributes
	}{
		{
			name:       "Invalid JSON",
			jsonString: "[{\"AttributeName\": \"name\", \"AttributeValue\": \"John\"}, {\"AttributeName\": \"surName\", \"AttributeValue\": \"Doe\"}",
			attributes: make(map[string]interface{}),
			expected:   map[string]interface{}{},
			wantErr:    false,
		},
		{
			name:       "Empty JSON",
			jsonString: "",
			attributes: make(map[string]interface{}),
			expected:   map[string]interface{}{},
			wantErr:    false,
		},
		{
			name:       "Valid JSON",
			jsonString: "[{\"AttributeName\": \"name\", \"AttributeValue\": \"John\"}, {\"AttributeName\": \"surName\", \"AttributeValue\": \"Doe\"}]",
			attributes: make(map[string]interface{}),
			expected: map[string]interface{}{
				"name":    "John",
				"surName": "Doe",
			},
			wantErr: false,
		},
	}

	// Iterate over the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the AddCustomMetaData function
			err := AddCustomMetaData(tt.jsonString, tt.attributes)

			// Check if the error matches the expected error
			if (err != nil) != tt.wantErr {
				t.Errorf("AddCustomMetaData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check if the attributes are correctly updated
			if err == nil && tt.jsonString != "" {
				if !reflect.DeepEqual(tt.attributes, tt.expected) {
					t.Errorf("AddCustomMetaData() attributes = %v, want %v", tt.attributes, tt.expected)
				}
			}
		})
	}
}

// generateTestString generates a test string of the specified size.
func generateTestString(size int) string {
	return string(make([]byte, size))
}

// TestSplitLargeMessages tests the SplitLargeMessages function.
// It verifies the behavior of splitting large messages into smaller chunks.
func TestSplitLargeMessages(t *testing.T) {
	// Generate test messages
	smallMessage := generateTestString(common.MaxMessageSize - 1) // Just under the limit
	exactMessage := generateTestString(common.MaxMessageSize)     // Exactly at the limit
	largeMessage := generateTestString(common.MaxMessageSize + 1) // Just over the limit

	tests := []struct {
		name    string   // Test case name
		message string   // Message to split
		want    []string // Expected split messages
	}{
		{
			name:    "Empty string",
			message: "",
			want:    []string{""},
		},
		{
			name:    "Small message",
			message: smallMessage,
			want:    []string{smallMessage},
		},
		{
			name:    "Exact size message",
			message: exactMessage,
			want:    []string{exactMessage},
		},
		{
			name:    "Large message",
			message: largeMessage,
			want:    []string{largeMessage[:common.MaxMessageSize/2], largeMessage[common.MaxMessageSize/2:]},
		},
	}

	// Iterate over the test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the SplitLargeMessages function
			if got := SplitLargeMessages(tt.message); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SplitLargeMessages() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestProduceMessageToChannel tests the ProduceMessageToChannel function
func TestProduceMessageToChannel(t *testing.T) {
	// Create a channel for DetailedLogsBatch
	channel := make(chan common.DetailedLogsBatch, 1) // Buffered channel

	// Create a sample log data and attributes
	currentBatch := []common.Log{{
		Timestamp: "1234567890",
		Log:       "test log",
	}}

	attributes := common.LogAttributes{
		"awsAccountId": "123456789012",
	}

	expectedDetailedLog := common.DetailedLogsBatch{{
		CommonData: common.Common{
			Attributes: attributes,
		},
		Entries: currentBatch,
	}}
	ProduceMessageToChannel(channel, currentBatch, attributes)
	receivedDetailedLog := <-channel

	assert.Equal(t, expectedDetailedLog, receivedDetailedLog)

	// Close the channel
	close(channel)
}

// / TestParseCTEvents tests the ParseCloudTrailEvents function with different CloudTrail messages.
func TestParseCTEvents(t *testing.T) {
	// Define test cases
	tests := []struct {
		name    string
		message string
		want    []map[string]interface{} // Use a slice of maps for the expected result
		wantErr bool
	}{
		{
			name: "Valid CloudTrail message",
			message: `{
				"Records": [
					{"eventVersion": "1.05", "eventName": "ConsoleLogin", "eventTime" : "2024-12-03T08:38:47Z"},
					{"eventVersion": "1.05", "eventName": "StartInstances"}
				]
			}`,
			want: []map[string]interface{}{
				{"eventVersion": "1.05", "eventName": "ConsoleLogin", "eventTime": "2024-12-03T08:38:47Z", "timestamp": float64(1733215127000)},
				{"eventVersion": "1.05", "eventName": "StartInstances"},
			},
			wantErr: false,
		},
		{
			name:    "Invalid CloudTrail message",
			message: `{"Records": "not an array"}`,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Empty CloudTrail message",
			message: `{"Records": []}`,
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotJSONStrings, err := ParseCloudTrailEvents(tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseCloudTrailEvents error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Unmarshal the JSON strings back into maps for comparison
			var got []map[string]interface{}
			for _, jsonString := range gotJSONStrings {
				var record map[string]interface{}
				if err := json.Unmarshal([]byte(jsonString), &record); err != nil {
					t.Errorf("Error unmarshaling result JSON: %v", err)
					return
				}
				got = append(got, record)
			}

			// Compare the resulting maps
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseCloudTrailEvents got = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestAddRequestId tests the AddRequestId function.
func TestAddRequestId(t *testing.T) {
	tests := []struct {
		name           string
		messages       []string
		logAttribute   common.LogAttributes
		lastRequestID  string
		expectedReqID  string
		expectedAttrib common.LogAttributes
	}{
		{
			name: "Case where RequestId is already present",
			messages: []string{
				"second message",
			},
			logAttribute:  common.LogAttributes{},
			lastRequestID: "d653fb2c-0234-46ff-ae6b-9a418b888420",
			expectedReqID: "d653fb2c-0234-46ff-ae6b-9a418b888420",
			expectedAttrib: common.LogAttributes{
				"requestId": "d653fb2c-0234-46ff-ae6b-9a418b888420",
			},
		},
		{
			name: "Valid RequestId in text message",
			messages: []string{
				"RequestId: d653fb2c-0234-46ff-ae6b-9a418b888420 hello world",
			},
			logAttribute:  common.LogAttributes{},
			lastRequestID: "",
			expectedReqID: "d653fb2c-0234-46ff-ae6b-9a418b888420",
			expectedAttrib: common.LogAttributes{
				"requestId": "d653fb2c-0234-46ff-ae6b-9a418b888420",
			},
		},
		{
			name: "No RequestId in text message",
			messages: []string{
				"hello world",
			},
			logAttribute:  common.LogAttributes{},
			lastRequestID: "d653fb2c-0234-46ff-ae6b-9a418b888420",
			expectedReqID: "d653fb2c-0234-46ff-ae6b-9a418b888420",
			expectedAttrib: common.LogAttributes{
				"requestId": "d653fb2c-0234-46ff-ae6b-9a418b888420",
			},
		},
		{
			name: "Valid RequestId in JSON message (top level)",
			messages: []string{
				`{
                    "timestamp": "2024-12-22T18:20:56Z",
                    "level": "INFO",
                    "message": "hello-world message",
                    "logger": "root",
                    "requestId": "22df9e13-6523-4567-b7a0-08ab5b41c06f"
                }`,
			},
			logAttribute:  common.LogAttributes{},
			lastRequestID: "",
			expectedReqID: "22df9e13-6523-4567-b7a0-08ab5b41c06f",
			expectedAttrib: common.LogAttributes{
				"requestId": "22df9e13-6523-4567-b7a0-08ab5b41c06f",
			},
		},
		{
			name: "Valid RequestId in JSON message (nested)",
			messages: []string{
				`{
                    "time": "2024-12-22T18:20:56.039Z",
                    "type": "platform.start",
                    "record": {
                        "requestId": "22df9e13-6523-4567-b7a0-08ab5b41c06f",
                        "version": "$LATEST"
                    }
                }`,
			},
			logAttribute:  common.LogAttributes{},
			lastRequestID: "",
			expectedReqID: "22df9e13-6523-4567-b7a0-08ab5b41c06f",
			expectedAttrib: common.LogAttributes{
				"requestId": "22df9e13-6523-4567-b7a0-08ab5b41c06f",
			},
		},
		{
			name: "Mixed JSON and text messages",
			messages: []string{
				`{
                    "time": "2025-01-15T02:10:15.287Z",
                    "type": "platform.start",
                    "record": {
                        "requestId": "4e3e651f-f7a7-46c4-acec-a07ca9876ba6",
                        "version": "$LATEST"
                    }
                }`,
				"2025/01/15 02:10:15 hello world",
				"2025/01/15 02:10:15 hello from New relic",
				`{
                    "time": "2025-01-15T02:10:15.296Z",
                    "type": "platform.report",
                    "record": {
                        "requestId": "4e3e651f-f7a7-46c4-acec-a07ca9876ba6",
                        "metrics": {
                            "durationMs": 1.42,
                            "billedDurationMs": 59,
                            "memorySizeMB": 128,
                            "maxMemoryUsedMB": 18,
                            "initDurationMs": 57.078
                        },
                        "status": "success"
                    }
                }`,
			},
			logAttribute:  common.LogAttributes{},
			lastRequestID: "",
			expectedReqID: "4e3e651f-f7a7-46c4-acec-a07ca9876ba6",
			expectedAttrib: common.LogAttributes{
				"requestId": "4e3e651f-f7a7-46c4-acec-a07ca9876ba6",
			},
		},
		{
			name: "Testing scenario with requestId already present for JSON message",
			messages: []string{
				`{
                    "time": "2025-01-15T02:10:15.287Z",
                    "type": "platform.start",
                    "record": {
                        "version": "$LATEST"
                    }
                }`,
			},
			logAttribute:  common.LogAttributes{},
			lastRequestID: "4e3e651f-f7a7-46c4-acec-a07ca9876ba6",
			expectedReqID: "4e3e651f-f7a7-46c4-acec-a07ca9876ba6",
			expectedAttrib: common.LogAttributes{
				"requestId": "4e3e651f-f7a7-46c4-acec-a07ca9876ba6",
			},
		},
	}

	// Regular expression to match the pattern "RequestId: <UUID>"
	re := regexp.MustCompile(`RequestId:\s([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12})`)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lastRequestID := tt.lastRequestID
			for _, message := range tt.messages {
				lastRequestID = AddRequestID(message, tt.logAttribute, lastRequestID, re)
			}
			assert.Equal(t, tt.expectedReqID, lastRequestID)
			assert.Equal(t, tt.expectedAttrib, tt.logAttribute)
		})
	}
}
