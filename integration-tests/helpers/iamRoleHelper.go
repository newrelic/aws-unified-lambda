package helpers

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/iam"
	"log"
)

// CreateIAMRole function to create IAM role
func CreateIAMRole(iamClient *iam.IAM, lambdaPolicyName string) (string, error) {
	roleName := lambdaPolicyName
	rolePolicy := `{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": { "Service": "lambda.amazonaws.com" },
				"Action": "sts:AssumeRole"
			},
			{
				"Effect": "Allow",
				"Action": [
					"s3:GetObject",
					"s3:ListBucket",
					"s3:GetBucketLocation",
					"s3:GetObjectVersion",
					"s3:GetLifecycleConfiguration"
				],
				"Resource": [
					"arn:aws:s3:::*",
                	"arn:aws:s3:::*/*"
				]
			},
			{
				"Effect": "Allow",
				"Action": [
					"secretsmanager:GetSecretValue",
					"secretsmanager:DescribeSecret"
				],
				"Resource": "arn:aws:secretsmanager:*:*:secret:*"
			},
			{
				"Effect": "Allow",
				"Action": [
					"sqs:SendMessage",
					"sqs:ChangeMessageVisibility",
					"sqs:ChangeMessageVisibilityBatch",
					"sqs:DeleteMessage",
					"sqs:DeleteMessageBatch",
					"sqs:GetQueueAttributes",
					"sqs:ReceiveMessage"
				],
				"Resource": "arn:aws:sqs:*:*:*"
			},
			{
				"Effect": "Allow",
				"Action": [
					"logs:CreateLogGroup",
                	"logs:CreateLogStream",
                	"logs:PutLogEvents"
				],
				"Resource": "*"
			}
		]
	}`

	role, err := iamClient.CreateRole(&iam.CreateRoleInput{
		RoleName:                 aws.String(roleName),
		AssumeRolePolicyDocument: aws.String(rolePolicy),
	})
	if err != nil {
		return "", fmt.Errorf("failed to create IAM role: %v", err)
	}
	log.Printf("created IAM role for lambda")
	return *role.Role.Arn, nil
}
