// Package cloudwatch provides functions for processing Cloudwatch logs, batch them into batches of Detailed Json logs and send them to a channel.
package cloudwatch

import (
	"os"
	"regexp"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
	"github.com/newrelic/aws-unified-lambda-logging/common"
	"github.com/newrelic/aws-unified-lambda-logging/logger"
	"github.com/newrelic/aws-unified-lambda-logging/util"
)

// log is a logger instance used for logging messages.
var log = logger.NewLogrusLogger(logger.WithDebugLevel())

// GetLogs batches logs from CloudWatch into DetailedJson format and sends them to the specified channel.
// It returns an error if there is a problem retrieving or sending the logs.
func GetLogs(cloudwatchLogsData events.CloudwatchLogsData, awsConfiguration util.AWSConfiguration, channel chan common.DetailedLogsBatch) error {

	// Following are the common attributes for all log messages.
	// All the attributes are compulsory for New Relic to generate Unique Entity ID.
	attributes := common.LogAttributes{
		"logGroup":                 cloudwatchLogsData.LogGroup,
		"logStream":                cloudwatchLogsData.LogStream,
		"aws.accountId":            awsConfiguration.AccountID,
		"aws.realm":                awsConfiguration.Realm,
		"aws.region":               awsConfiguration.Region,
		"instrumentation.provider": common.InstrumentationProvider,
		"instrumentation.name":     common.InstrumentationName,
		"instrumentation.version":  common.InstrumentationVersion,
	}

	if err := util.AddCustomMetaData(os.Getenv(common.CustomMetaData), attributes); err != nil {
		log.Errorf("failed to add custom metadata %v", err)
		return err
	}

	if err := batchLogEntries(cloudwatchLogsData, channel, attributes); err != nil {
		return err
	}
	return nil
}

// batchLogEntries processes a batch of CloudWatch log entries and splits them into smaller batches based on payload size and message count
// and produces log data batches to a channel.
// The function returns an error if any.
func batchLogEntries(cloudwatchLogsData events.CloudwatchLogsData, channel chan common.DetailedLogsBatch, attributes common.LogAttributes) error {
	batchSize := 0

	messageCount := 0

	var currentBatch common.LogData

	// Regular expression to match the pattern "RequestId: <UUID> <message>"
	regularExpression := regexp.MustCompile(common.RequestIDRegex)

	// Variable to keep track of the last requestId found
	var lastRequestID string

	for _, record := range cloudwatchLogsData.LogEvents {
		messages := util.SplitLargeMessages(record.Message)

		for _, message := range messages {
			entry := common.Log{
				Timestamp: strconv.FormatInt(record.Timestamp, 10),
				Log:       message,
			}

			lastRequestID = util.AddRequestID(cloudwatchLogsData.LogGroup, message, attributes, lastRequestID, regularExpression)

			if batchSize+len(message) > common.MaxPayloadSize || messageCount >= common.MaxPayloadMessages {
				util.ProduceMessageToChannel(channel, currentBatch, attributes)
				currentBatch = nil
				batchSize = 0
				messageCount = 0
			}

			currentBatch = append(currentBatch, entry)
			batchSize = batchSize + len(message)
			messageCount = messageCount + 1
		}
	}

	if len(currentBatch) > 0 {
		util.ProduceMessageToChannel(channel, currentBatch, attributes)
	}

	log.Debug("Finished processing all cloudwatch logs")

	return nil
}
