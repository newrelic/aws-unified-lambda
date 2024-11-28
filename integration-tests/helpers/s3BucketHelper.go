package helpers

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"log"
	"os"
	"path/filepath"
)

// CreateBucket function to create bucket
func CreateBucket(s3Client *s3.S3, s3BucketName string) error {
	_, err := s3Client.CreateBucket(&s3.CreateBucketInput{
		Bucket: aws.String(s3BucketName),
	})
	if err != nil {
		return fmt.Errorf("unable to create bucket: %v", err)
	}
	log.Printf("Bucket %s created.", s3BucketName)
	return nil
}

// CreateS3EventSource function to create event source
func CreateS3EventSource(s3Client *s3.S3, s3BucketName, lambdaName string) error {
	_, err := s3Client.PutBucketNotificationConfiguration(&s3.PutBucketNotificationConfigurationInput{
		Bucket: aws.String(s3BucketName),
		NotificationConfiguration: &s3.NotificationConfiguration{
			LambdaFunctionConfigurations: []*s3.LambdaFunctionConfiguration{
				{
					LambdaFunctionArn: aws.String(fmt.Sprintf("arn:aws:lambda:us-east-1:000000000000:function:%s", lambdaName)),
					Events:            []*string{aws.String("s3:ObjectCreated:*")},
				},
			},
		},
	})

	if err != nil {
		return fmt.Errorf("unable to configure bucket notification: %v", err)
	}
	log.Printf("S3 event source created for bucket %s to invoke Lambda function %s.", s3BucketName, lambdaName)
	return nil
}

// SimulateS3Event function to trigger event
func SimulateS3Event(filePath string, sess *session.Session, s3BucketName string) error {
	uploader := s3manager.NewUploader(sess)

	file, err := os.Open(filePath)
	defer file.Close()

	bucket := s3BucketName
	key := filepath.Base(filePath)

	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   file,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return nil
}
