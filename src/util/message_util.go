// Package util provides generic utility functions.
package util

import (
	"encoding/json"
	"regexp"
	"time"

	"github.com/newrelic/aws-unified-lambda-logging/common"
)

// SplitLargeMessages splits a large message into smaller messages if its length exceeds the maximum message size.
func SplitLargeMessages(message string) []string {
	var result []string
	if len(message) > common.MaxMessageSize {
		// recursive call to split the messages.
		result = append(result, SplitLargeMessages(message[:len(message)/2])...)
		result = append(result, SplitLargeMessages(message[len(message)/2:])...)
	} else {
		result = append(result, message)
	}
	return result
}

// CustomAttribute represents a custom attribute with a name and value.
type CustomAttribute struct {
	AttributeName  string `json:"AttributeName"`  // Name of the custom attribute
	AttributeValue string `json:"AttributeValue"` // Value of the custom attribute
}

// AddCustomMetaData adds custom metadata attributes to a provided map.
func AddCustomMetaData(jsonString string, attributes map[string]interface{}) error {
	if jsonString == "" {
		return nil
	}
	var customAttributes []CustomAttribute
	err := json.Unmarshal([]byte(jsonString), &customAttributes)

	if err != nil {
		log.Errorf("failed to unmarshal custom metadata: %v", err)
		return nil
	}

	for _, customAttribute := range customAttributes {
		// Adding the check to avoid overwriting the existing attribute - specifically introduced to not override entity synthesis parameters
		if _, exists := attributes[customAttribute.AttributeName]; !exists {
			// Add attribute to the map if the key is not present
			attributes[customAttribute.AttributeName] = customAttribute.AttributeValue
		}
	}

	return nil
}

// ProduceMessageToChannel sends a log batch to a channel for further processing.
func ProduceMessageToChannel(channel chan common.DetailedLogsBatch, currentBatch common.LogData, attributes common.LogAttributes) {
	channel <- []common.DetailedLog{{
		CommonData: common.Common{
			Attributes: attributes,
		},
		Entries: currentBatch,
	}}
}

// CloudTrailRecords represents a list of CloudTrail records.
type CloudTrailRecords struct {
	Records []map[string]interface{} `json:"Records"`
}

// ParseCloudTrailEvents parses a CloudTrail message and returns a list of log records as strings.
func ParseCloudTrailEvents(message string) ([]string, error) {
	var cloudTrailRecords CloudTrailRecords
	err := json.Unmarshal([]byte(message), &cloudTrailRecords)

	if err != nil {
		return nil, err
	}

	// Serialize each record into a JSON string.
	var records []string
	for _, record := range cloudTrailRecords.Records {
		if record["eventTime"] != nil {
			parsedTime, err := time.Parse(time.RFC3339, record["eventTime"].(string))
			if err == nil {
				record["timestamp"] = parsedTime.UnixMilli()
			}
		}

		recordJSON, err := json.Marshal(record)
		if err != nil {
			log.Errorf("Error marshaling record to JSON: %v while parsing %v", err, record)
			continue
		}
		records = append(records, string(recordJSON))
	}
	return records, nil
}

// AddRequestID extracts the requestId from the message and updates the attributes map.
// It returns the last requestId found to keep track of the requestId across log messages.
func AddRequestID(message string, logAttribute common.LogAttributes, lastRequestID string, regularExpression *regexp.Regexp) string {

	// Check if the message is in JSON format
	var jsonMessage map[string]interface{}
	if err := json.Unmarshal([]byte(message), &jsonMessage); err == nil {
		// Message is in JSON format
		if requestID, ok := ExtractRequestIDFromJSON(jsonMessage); ok {
			lastRequestID = requestID
			logAttribute["requestId"] = lastRequestID
		} else if lastRequestID != "" {
			logAttribute["requestId"] = lastRequestID
		}
	} else {
		// Message is in text format
		matches := regularExpression.FindStringSubmatch(message)
		// if message has valid regular expression then matches will have 2 elements, first element is the whole string and second element is the requestId
		if len(matches) == 2 {
			lastRequestID = matches[1]
			logAttribute["requestId"] = lastRequestID
		} else if lastRequestID != "" {
			logAttribute["requestId"] = lastRequestID
		}

	}
	return lastRequestID
}

// ExtractRequestIDFromJSON extracts the requestId from the JSON message.
// It returns the requestId and a boolean indicating if the requestId was found.
// The function handle both the cases where the requestId is at the top level or inside the record.requestId.
func ExtractRequestIDFromJSON(jsonMessage map[string]interface{}) (string, bool) {
	if requestID, ok := jsonMessage["requestId"].(string); ok {
		return requestID, true
	}
	if record, ok := jsonMessage["record"].(map[string]interface{}); ok {
		if requestID, ok := record["requestId"].(string); ok {
			return requestID, true
		}
	}
	return "", false
}
