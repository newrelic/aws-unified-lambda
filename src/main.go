// Package: main
// This package is the entry point for the Lambda function. It initializes the New Relic client for logging and processes the incoming event.
package main

import (
	"context"
	"sync"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/newrelic/aws-unified-lambda-logging/cloudwatch"
	"github.com/newrelic/aws-unified-lambda-logging/common"
	"github.com/newrelic/aws-unified-lambda-logging/logger"
	"github.com/newrelic/aws-unified-lambda-logging/s3"
	"github.com/newrelic/aws-unified-lambda-logging/unmarshal"
	"github.com/newrelic/aws-unified-lambda-logging/util"
)

var log = logger.NewLogrusLogger(logger.WithDebugLevel())

// handlerWithArgs is the main Lambda handler function.
// It processes the incoming event and sends the logs to New Relic for logging.
// It supports CloudWatch and S3 events.
// It tracks the consumer go routines using a WaitGroup.
func handlerWithArgs(ctx context.Context, event unmarshal.Event, nrClient util.NewRelicClientAPI) error {
	channel := make(chan common.DetailedLogsBatch)
	var wg sync.WaitGroup
	wg.Add(1)

	go util.ConsumeLogBatches(ctx, channel, &wg, nrClient)

	awsConfiguration, err := util.GetAWSConfiguration(ctx)

	if err != nil {
		log.Fatal("Error getting AWS configuration")
	}
	switch event.EventType {
	case unmarshal.CLOUDWATCH:
		log.Debugf("processing cloudwatch event: %v", event.CloudwatchLogsData)
		err = cloudwatch.GetLogs(event.CloudwatchLogsData, awsConfiguration, channel)
	case unmarshal.S3:
		log.Debugf("processing s3 event: %v", event.S3Event)
		var s3Client s3.ObjectClient
		s3Client, err = s3.NewS3Client(ctx)
		if err != nil {
			log.Fatalf("error creating s3 client: %v", err)
		}
		err = s3.GetLogsFromS3Event(ctx, event.S3Event, awsConfiguration, channel, s3Client, s3.DefaultReaderFactory)
	case unmarshal.SNS:
		log.Debugf("processing sns event: %v", event.SNSEvent)
		var s3Client s3.ObjectClient
		s3Client, err = s3.NewS3Client(ctx)
		if err != nil {
			log.Fatalf("error creating s3 client: %v", err)
		}
		err = s3.GetLogsFromSNSEvent(ctx, event.SNSEvent, awsConfiguration, channel, s3Client, s3.DefaultReaderFactory)
	default:
		log.Error("unable to process unknown event type. Supported event types are cloudwatch and s3")
		return nil
	}

	if err != nil {
		log.Fatalf("error processing event: %v", err)
	}

	close(channel)

	wg.Wait()
	return nil
}

// main is the entry point of the program.
// It initializes a new New Relic client and starts a Lambda handler.
func main() {
	nrClient, err := util.NewNRClient()
	if err != nil {
		log.Fatalf("error initializing newrelic client: %v", err)
	} else {
		handler := func(ctx context.Context, event unmarshal.Event) error {
			return handlerWithArgs(ctx, event, nrClient)
		}
		lambda.Start(handler)
	}
}
