// Package util provides generic utility functions.
package util

import (
	"encoding/json"
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
