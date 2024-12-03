// Package unmarshal deals provides functions to unmarshal events to various event type such as Cloudwatch, S3
package unmarshal

import (
	"encoding/json"
	"github.com/newrelic/aws-unified-lambda-logging/logger"

	"github.com/aws/aws-lambda-go/events"
)

// Defines the event types
const (
	CLOUDWATCH = "cloudwatch" // CLOUDWATCH represents the event type for CloudWatch logs.
	S3         = "s3"         // S3 represents the event type for S3 events.
)

var log = logger.NewLogrusLogger(logger.WithDebugLevel())

// Event represents the unified event structure.
type Event struct {
	EventType          string                    // EventType represents the type of the event.
	CloudwatchLogsData events.CloudwatchLogsData // CloudwatchLogsData represents the CloudWatch logs data.
	S3Event            events.S3Event            // S3Event represents the S3 event data.
}

// UnmarshalJSON unmarshals the JSON data into the Event struct.
func (event *Event) UnmarshalJSON(data []byte) error {
	log.Debugf("event : %v", string(data[:]))
	var err error

	// Try to unmarshal the event as CloudwatchLogsEvent
	var cloudWatchEvent events.CloudwatchLogsEvent
	json.Unmarshal(data, &cloudWatchEvent)
	var cloudwatchLogsData events.CloudwatchLogsData
	cloudwatchLogsData, err = cloudWatchEvent.AWSLogs.Parse()
	if err == nil {
		event.EventType = CLOUDWATCH
		event.CloudwatchLogsData = cloudwatchLogsData

		return err
	}

	// Try to unmarshal the event as S3Event
	var s3Event events.S3Event
	err = json.Unmarshal(data, &s3Event)
	if err == nil && len(s3Event.Records) != 0 && s3Event.Records[0].EventName != "" {
		event.EventType = S3
		event.S3Event = s3Event

		return err
	}

	return nil
}
