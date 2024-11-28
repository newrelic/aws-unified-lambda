package helpers

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	integrationtests "github.com/newrelic/aws-unified-lambda-logging/common"
	"log"
)

// CreateLambdaFunction function to create lambda
func CreateLambdaFunction(lambdaClient *lambda.Lambda, iamRoleArn string, zipFileData []byte, newRelicAPIKey string, lambdaName string, secretName *string, dlqArn string, includeAPIKeyInEnv bool) error {
	envVariables := map[string]*string{
		"NEW_RELIC_LICENSE_KEY_SECRET_NAME": aws.String(""),
		"NEW_RELIC_REGION":                  aws.String(integrationtests.NewRelicRegion),
		"DEBUG_ENABLED":                     aws.String("false"),
	}

	if includeAPIKeyInEnv {
		envVariables["LICENSE_KEY"] = aws.String(newRelicAPIKey)
	} else {
		envVariables["NEW_RELIC_LICENSE_KEY_SECRET_NAME"] = aws.String(*secretName)
	}

	_, err := lambdaClient.CreateFunction(&lambda.CreateFunctionInput{
		FunctionName: aws.String(lambdaName),
		Handler:      aws.String(integrationtests.LambdaHandler),
		Role:         aws.String(iamRoleArn),
		Runtime:      aws.String(integrationtests.LambdaRuntime),
		Code: &lambda.FunctionCode{
			ZipFile: zipFileData,
		},
		Timeout:    aws.Int64(240),
		MemorySize: aws.Int64(5120),
		Architectures: []*string{
			aws.String("x86_64"),
		},
		Environment: &lambda.Environment{
			Variables: envVariables,
		},
		DeadLetterConfig: &lambda.DeadLetterConfig{
			TargetArn: &dlqArn,
		},
	})

	if err != nil {
		return fmt.Errorf("unable to create lambda function: %v", err)
	}
	log.Printf("Lambda function %s created.", lambdaName)
	return nil
}
