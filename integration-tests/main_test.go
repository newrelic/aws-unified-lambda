package integrationtests

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joho/godotenv"
	"github.com/newrelic/aws-unified-lambda-logging/common"
	"github.com/newrelic/aws-unified-lambda-logging/helpers"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

var (
	newRelicAPIKey                string
	newRelicAccountID             string
	secretName                    string
	awsSession                    *session.Session
	userKey                       string
	s3BucketName                  string
	s3BucketNameSecretManagerCase string
	err                           error
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, expecting environment variables to be set")
	}

	newRelicAPIKey = os.Getenv("NEW_RELIC_INGEST_KEY")
	newRelicAccountID = os.Getenv("NEW_RELIC_ACCOUNT_ID")
	secretName = os.Getenv("NEW_RELIC_LICENSE_KEY_SECRET_NAME")
	userKey = os.Getenv("NEW_RELIC_USER_KEY")
}

func TestMain(m *testing.M) {
	// setup
	awsSession = CreateAWSSession()
	s3BucketName, _, err = BuildAndDeployResources(newRelicAPIKey, nil, awsSession)

	// Run tests
	m.Run()

	// Assert
	validateResults(&testing.T{})
}

func TestLambdaGzipCompressedFiles(t *testing.T) {
	if err = helpers.SimulateS3Event("test-files/integ-test-logs.gz", awsSession, s3BucketName); err != nil {
		t.Fatalf("Failed to simulate S3 event: %v", err)
	}
}

func TestLambdaBzip2CompressedFiles(t *testing.T) {
	if err = helpers.SimulateS3Event("test-files/integ-test-logs.bz2", awsSession, s3BucketName); err != nil {
		t.Fatalf("Failed to simulate S3 event: %v", err)
	}
}

func TestLambdaUncompressedFiles(t *testing.T) {
	if err = helpers.SimulateS3Event("test-files/integ-test-logs.txt", awsSession, s3BucketName); err != nil {
		t.Fatalf("Failed to simulate S3 event: %v", err)
	}
}

func TestLambdaMultipleLogFiles(t *testing.T) {
	if err = helpers.SimulateS3Event("test-files/integ-test-multiple-logs.txt", awsSession, s3BucketName); err != nil {
		t.Fatalf("Failed to simulate S3 event: %v", err)
	}
}

func TestLambdaLicenseKeyIsFetchedFromAWSSecretsManager(t *testing.T) {
	s3BucketNameSecretManagerCase, _, err = BuildAndDeployResources(newRelicAPIKey, &secretName, awsSession)
	if err != nil {
		fmt.Println("Resources creation failed:", err)
	}

	if err = helpers.SimulateS3Event("test-files/integ-test-logs-secrets-manager.txt", awsSession, s3BucketNameSecretManagerCase); err != nil {
		t.Fatalf("Failed to simulate S3 event: %v", err)
	}
}

func validateResults(t *testing.T) {
	// gzip file
	result, err := fetchLogsWithRetry("integ-test-logs.gz")

	assert.NoError(t, err)
	assert.Equal(t, "This is gzip file", result[0].Message)

	// bzip file
	result, err = fetchLogsWithRetry("integ-test-logs.bz2")

	assert.NoError(t, err)
	assert.Equal(t, "This is bz2 compressed file", result[0].Message)

	// txt file
	result, err = fetchLogsWithRetry("integ-test-logs.txt")

	assert.NoError(t, err)
	assert.Equal(t, "This is txt file", result[0].Message)

	// multiple-line file
	result, err = fetchLogsWithRetry("integ-test-multiple-logs.txt")

	assert.NoError(t, err)
	assert.Equal(t, "This is txt file with multiple lines - line 1", result[0].Message)
	assert.Equal(t, "This is txt file with multiple lines - line 2", result[1].Message)

	// fetch secret from secret manager
	result, err = fetchLogsWithRetry("integ-test-logs-secrets-manager.txt")

	assert.NoError(t, err)
	assert.Equal(t, "This is txt file for secrets manager. License key is fetched from secrets manager", result[0].Message)
}

func fetchLogsWithRetry(fileName string) ([]common.LogEvent, error) {
	for i := 0; i < common.NoOfRetriesForLogsFetch; i++ {
		result, err := helpers.FetchLogsFromNewRelic(userKey, newRelicAccountID, fileName)
		if result == nil || len(result) == 0 {
			fmt.Println("Retrying validation...")
			time.Sleep(common.WaitTimeForLogsFetch)
			continue
		} else {
			return result, err
		}
	}
	return nil, nil
}
