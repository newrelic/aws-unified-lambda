package util

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

// AWSConfiguration represents the AWS configuration information.
type AWSConfiguration struct {
	Realm     string // Realm represents the realm of the AWS configuration.
	AccountID string // AccountID represents the AWS account ID.
	Region    string // Region represents the AWS region.
}

// GetAWSConfiguration retrieves the AWS configuration from the lambda context.
// It returns the AWSConfiguration struct and an error if the lambda context is not available or the function ARN is invalid.
func GetAWSConfiguration(ctx context.Context) (AWSConfiguration, error) {
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		return AWSConfiguration{}, errors.New("failed to get lambda context")
	}
	arn := lc.InvokedFunctionArn
	parts := strings.Split(arn, ":")

	if len(parts) < 7 || len(parts) > 8 {
		return AWSConfiguration{}, errors.New("invalid function ARN")
	}

	return AWSConfiguration{
		Realm:     parts[1],
		AccountID: parts[4],
		Region:    parts[3],
	}, nil
}
