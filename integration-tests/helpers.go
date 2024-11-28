package integrationtests

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/endpoints"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/aws/aws-sdk-go/service/sqs"
	integrationtests "github.com/newrelic/aws-unified-lambda-logging/common"
	"github.com/newrelic/aws-unified-lambda-logging/helpers"
	"github.com/newrelic/newrelic-client-go/v2/pkg/testhelpers"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// CreateAWSSession function to create AWS session
func CreateAWSSession() *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Credentials:      credentials.NewStaticCredentials("test", "test", ""),
		Endpoint:         aws.String(integrationtests.LocalEndpoint),
		Region:           aws.String(endpoints.UsEast1RegionID),
		S3ForcePathStyle: aws.Bool(true),
	}))
}

// BuildAndDeployResources function to build and deploy resources (IAM, Lambda, S3 Bucket and S3 trigger event)
func BuildAndDeployResources(newRelicAPIKey string, secretName *string, awsSession *session.Session) (string, string, error) {
	iamClient := iam.New(awsSession)
	lamdaClient := lambda.New(awsSession)
	s3Client := s3.New(awsSession)
	smClient := secretsmanager.New(awsSession)
	sqsClient := sqs.New(awsSession)

	s3BucketName := "s3-bucket-" + testhelpers.GenerateRandomName(5)
	lambdaName := "lambda-" + testhelpers.GenerateRandomName(5)
	lambdaPolicyName := "lambda-execution-policy-" + testhelpers.GenerateRandomName(5)

	// Create IAM Role
	iamRoleArn, err := helpers.CreateIAMRole(iamClient, lambdaPolicyName)
	if err != nil {
		fmt.Println("failed to create IAM role:", err)
		os.Exit(1)
	}

	dlqName := "my-dlq"
	createQueueOutput, err := sqsClient.CreateQueue(&sqs.CreateQueueInput{
		QueueName: aws.String(dlqName),
	})

	if secretName != nil {
		println("created secret")
		if err = helpers.CreateSecret(smClient, *secretName, newRelicAPIKey); err != nil {
			fmt.Println("failed to create IAM role:", err)
			os.Exit(1)
		}
	}

	// Build lambda binary
	err = buildExecutableAndZipFile()
	if err != nil {
		fmt.Println("failed to build executable:", err)
		os.Exit(1)
	}
	defer os.Remove(filepath.Join("../src", integrationtests.ExecutableFileName))
	defer os.Remove(integrationtests.ZipFileName)

	zipFileData, err := ioutil.ReadFile(integrationtests.ZipFileName)

	// Create Lambda Function
	dlqArn := createQueueOutput.QueueUrl
	if secretName != nil {
		if err = helpers.CreateLambdaFunction(lamdaClient, iamRoleArn, zipFileData, newRelicAPIKey, lambdaName, secretName, *dlqArn, false); err != nil {
			fmt.Println("failed to deploy Lambda function:", err)
			os.Exit(1)
		}
	} else {
		if err = helpers.CreateLambdaFunction(lamdaClient, iamRoleArn, zipFileData, newRelicAPIKey, lambdaName, secretName, *dlqArn, true); err != nil {
			fmt.Println("failed to deploy Lambda function:", err)
			os.Exit(1)
		}
	}

	// Create S3 Bucket
	if err = helpers.CreateBucket(s3Client, s3BucketName); err != nil {
		fmt.Println("Failed to create S3 bucket: ", err)
		os.Exit(1)
	}

	// Create S3 Event Source
	for i := 0; i < integrationtests.NoOfRetriesForEventCreation; i++ {
		if err = helpers.CreateS3EventSource(s3Client, s3BucketName, lambdaName); err != nil {
			errorString := err.Error()
			if strings.Contains(errorString, integrationtests.RetryError) {
				fmt.Println("Failed to create S3 event source. Retrying for error: ", err)
				time.Sleep(integrationtests.WaitTimeForResourceCreation)
				continue
			} else {
				fmt.Println("Failed to create S3 event source: ", err)
				os.Exit(1)
			}
		}
		break
	}

	return s3BucketName, lambdaName, nil
}

func buildExecutableAndZipFile() error {
	sourceDir := "../src"
	srcPath, _ := filepath.Abs(sourceDir)
	outExecutable := filepath.Join(srcPath, integrationtests.ExecutableFileName)

	cmd := exec.Command("go", "build", "-o", outExecutable, filepath.Join(srcPath, "main.go"))
	cmd.Env = append(os.Environ(),
		"GOOS=linux",
		"GOARCH=amd64",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = srcPath
	err := cmd.Run()
	if err != nil {
		return err
	}

	// Define the zip file path
	integrationTestsDir, _ := filepath.Abs(".")
	zipFilePath := filepath.Join(integrationTestsDir, integrationtests.ZipFileName)

	cmd = exec.Command("zip", "-r", zipFilePath, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = srcPath
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}
