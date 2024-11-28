package unmarshal

import (
	"encoding/json"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
)

// TestUnmarshalJSONS3Event is a unit test function that tests the unmarshaling of a JSON S3 event.
// It verifies that the unmarshaled event matches the expected event.
func TestUnmarshalJSONS3Event(t *testing.T) {
	input := []byte(`{
		"Records": [
			{
				"eventName": "ObjectCreated:Put"
			}
		]
	}`)
	expected := Event{
		EventType: S3,
		S3Event: events.S3Event{
			Records: []events.S3EventRecord{
				{
					EventName: "ObjectCreated:Put",
				},
			},
		},
	}

	var event Event
	err := json.Unmarshal(input, &event)

	assert.NoError(t, err)
	assert.Equal(t, expected.EventType, event.EventType)
	assert.Equal(t, expected.S3Event, event.S3Event)
}

// TestUnmarshalJSONCloudWatchLogsData is a unit test function that tests the unmarshaling of a JSON CloudWatch Logs Data event.
// It verifies that the unmarshaled event does not match the expected event when the input is incorrect.
func TestUnmarshalJSONCloudWatchLogsData(t *testing.T) {
	input := []byte(`{
		"awslogs": {
			"message": "test message"
		}
	}`)

	expected := Event{
		EventType: CLOUDWATCH,
		CloudwatchLogsData: events.CloudwatchLogsData{
			LogEvents: []events.CloudwatchLogsLogEvent{
				{
					Message: "test message",
				},
			},
		},
	}

	var event Event
	json.Unmarshal(input, &event)

	assert.NotEqual(t, expected.EventType, event.EventType)
	assert.NotEqual(t, expected.CloudwatchLogsData, event.CloudwatchLogsData)
}
