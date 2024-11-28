// Package s3 provides functions for retrieving logs from Amazon S3 buckets, process them into batches of Detailed Json logs and send them to a channel.
package s3

import (
	"bufio"
	"compress/bzip2"
	"compress/gzip"
	"context"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/newrelic/aws-unified-lambda-logging/common"
	"github.com/newrelic/aws-unified-lambda-logging/logger"
	"github.com/newrelic/aws-unified-lambda-logging/util"
)

// ObjectClient is an interface that defines the methods for interacting with the S3 service.
type ObjectClient interface {
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

// ReaderFactory defines a function type that creates a new io.Reader based on the input reader and file extension.
// Checkout DefaultReaderFactory for an example implementation.
type ReaderFactory func(io.ReadCloser, string) (io.Reader, error)

var log = logger.NewLogrusLogger(logger.WithDebugLevel())

// GetLogsFromS3Event batches logs from S3 into DetailedJson format and sends them to the specified channel.
// It returns an error if there is a problem retrieving or sending the logs.
func GetLogsFromS3Event(ctx context.Context, s3Event events.S3Event, awsConfiguration util.AWSConfiguration, channel chan common.DetailedLogsBatch, s3Client ObjectClient, readerFactory ReaderFactory) error {
	for _, record := range s3Event.Records {

		// The Following are the common attributes for all log messages.
		// New Relic uses these common attributes to generate Unique Entity ID.
		attributes := common.LogAttributes{
			"aws.accountId":            awsConfiguration.AccountID,
			"logBucketName":            record.S3.Bucket.Name,
			"logObjectURI":             record.S3.Object.URLDecodedKey,
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

		if err := buildMeltLogsFromS3Bucket(ctx, record.S3.Bucket.Name, record.S3.Object.URLDecodedKey, channel, attributes, s3Client, readerFactory); err != nil {
			return err
		}
	}

	return nil
}

// fetchS3Reader fetches an S3 object from the specified bucket and returns an io.ReadCloser for reading its contents.
// It returns the io.ReadCloser and any error encountered during the operation.
func fetchS3Reader(ctx context.Context, bucketName string, objectName string, s3Client ObjectClient) (io.ReadCloser, error) {
	resp, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		log.Errorf("failed to get S3 object reader: %v", err)
		return nil, err
	}

	return resp.Body, nil
}

// buildMeltLogsFromS3Bucket reads the contents of an S3 object line by line,
// splits large messages, and produces log data batches to a channel.
func buildMeltLogsFromS3Bucket(ctx context.Context, bucketName string, objectName string, channel chan common.DetailedLogsBatch, attributes common.LogAttributes, s3Client ObjectClient, readerFactory ReaderFactory) error {
	s3Reader, err := fetchS3Reader(ctx, bucketName, objectName, s3Client)
	if err != nil {
		return err
	}
	defer s3Reader.Close()

	reader, err := readerFactory(s3Reader, objectName)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(reader)
	scanner.Buffer(make([]byte, 0, common.MaxBufferSize), common.MaxBufferSize)
	batchSize := 0

	messageCount := 0

	var currentBatch common.LogData

	log.Debug("Reading file line by line")

	for scanner.Scan() {
		line := scanner.Text()
		messages := util.SplitLargeMessages(line)

		for _, message := range messages {
			entry := common.Log{
				Log: message,
			}

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

	log.Debug("Finished reading file line by line")

	if len(currentBatch) > 0 {
		util.ProduceMessageToChannel(channel, currentBatch, attributes)
	}

	if err := scanner.Err(); err != nil {
		log.Errorf("failed to read line by line for object %s in bucket %s: %v", objectName, bucketName, err)
		return nil
	}

	return nil
}

// DefaultReaderFactory returns an io.Reader that can be used to read the contents of the input file.
func DefaultReaderFactory(input io.ReadCloser, filename string) (io.Reader, error) {
	var reader io.Reader
	// If the file has a ".gz" extension, it creates a gzip reader.
	// If the file has a ".bz2" extension, it creates a bzip2 reader.
	// For all other file extensions, it returns the input reader as is.
	switch {
	case strings.HasSuffix(filename, ".gz"):
		gzipReader, err := gzip.NewReader(input)
		if err != nil {
			log.Errorf("failed to create gzip reader: %v", err)
			return nil, err
		}
		reader = gzipReader
	case strings.HasSuffix(filename, ".bz2"):
		reader = bzip2.NewReader(input)
	default:
		reader = input
	}

	return reader, nil
}

// NewS3Client creates a new S3 client using the provided context and returns the client.
func NewS3Client(ctx context.Context) (ObjectClient, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		log.Errorf("unable to load SDK config: %v", err)
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)
	return s3Client, nil
}
