package util

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/lambdacontext"
)

// mockLambdaContext creates a mock Lambda context with the given ARN.
// It returns a context.Context object that can be used for testing purposes.
func mockLambdaContext(arn string) context.Context {
	lc := &lambdacontext.LambdaContext{
		InvokedFunctionArn: arn,
	}
	return lambdacontext.NewContext(context.Background(), lc)
}

// TestGetAWSConfiguration is a unit test function that tests the GetAWSConfiguration function.
// It verifies the behavior of GetAWSConfiguration by providing different test cases.
// Each test case includes a name, an ARN, the expected AWSConfiguration, and an expectation of whether an error is expected.
func TestGetAWSConfiguration(t *testing.T) {
	tests := []struct {
		name        string           // Name of the test case
		arn         string           // ARN of the Lambda function
		want        AWSConfiguration // Expected AWSConfiguration
		expectError bool             // Whether an error is expected
	}{
		{
			name:        "Valid ARN",
			arn:         "arn:aws:lambda:us-west-1:123456789012:function:my-function",
			want:        AWSConfiguration{Realm: "aws", AccountID: "123456789012", Region: "us-west-1"},
			expectError: false,
		},
		{
			name:        "Invalid ARN - Too few parts",
			arn:         "arn:aws:lambda:us-west-1:123456789012",
			want:        AWSConfiguration{},
			expectError: true,
		},
		{
			name:        "Invalid ARN - Too many parts",
			arn:         "arn:aws:lambda:us-west-1:123456789012:function:my-function:version:extra",
			want:        AWSConfiguration{},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := mockLambdaContext(tt.arn)
			got, err := GetAWSConfiguration(ctx)
			if (err != nil) != tt.expectError {
				t.Errorf("GetAWSConfiguration() error = %v, expectError %v", err, tt.expectError)
				return
			}
			if !tt.expectError && got != tt.want {
				t.Errorf("GetAWSConfiguration() got = %v, want %v", got, tt.want)
			}
		})
	}
}
