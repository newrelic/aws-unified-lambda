# AWS unified logging

AWS unified logging application sends logs data from CloudWatch and S3 to New Relic.


## Features


- S3 file processing: Handles the gzip and bzip2 compression formats. Other than these file formats are treated as uncompressed.
- CloudWatch logs processing
- DLQ support to handle events after that fail after two retries.


## Limitations

- Supports uncompressed files up to 400 MB.
- Supports gzip and bzip2 compressed files up to 200 MB.
- Does not parse log lines in S3 log files, such as extracting timestamps from log lines.
- Log lines exceeding 8 MB will cause event processing to fail.


## Requirements

- AWS CLI must be installed and configured with Administrator permission
- Docker must be installed. Refer [Docker documentation.](https://www.docker.com/community-edition)
- Golang must be installed. Refer [Golang documentation](https://golang.org)
- Install the AWS SAM CLI. Refer [SAM CLI Documentation](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/serverless-sam-cli-install.html)


## Deployment

To test this Lambda function, deploy the CloudFormation template (`template.yml`) using the AWS SAM CLI. Ensure AWS is authenticated with your chosen account.


### CloudFormation parameters

| Parameters      | Description         | 
|-----------------|---------------------|
| `LicenseKey`    | The key for forwarding logs to New Relic.  |
| `NewRelicRegion`  | Set to 'US' or 'EU' to specify the New Relic endpoint for log forwarding.    | 
| `NewRelicAccountId`    | New Relic account ID where logs will be sent. | 
| `StoreNRLicenseKeyInSecretManager` | Set to `true` to store the license key in AWS Secrets Manager or `false` to keep it in environment variables.   |
| `S3BucketNames`  | A JSON array of your S3 bucket names and prefixes for Lambda triggers. For example, `[{"bucket1":"prefix1"}, {"bucket2":"prefix2"}]`   |
| `LogGroupConfig`      | A JSON array of CloudWatch LogGroup names and filters for Lambda triggers. For example, `[{"LogGroupName":"group1"}, {"LogGroupName":"group2", "FilterPattern":"ERROR"}, {"LogGroupName":"group3", "FilterPattern":"INFO"}]`   |
| `CommonAttributes`     | JSON object of common attributes to add to all logs. For example, `[{"AttributeName": "name1", "AttributeValue": "value1"}, {"AttributeName": "name2", "AttributeValue": "value2"}]`  |


### Lambda environment variables

| Parameter    | Description   |
|--------------|---------------|
| `LICENSE_KEY`    | Your New Relic license key when `StoreNRLicenseKeyInSecretManager` is set to `false`.  |
| `NEW_RELIC_LICENSE_KEY_SECRET_NAME`  | The name of the AWS secret when the `StoreNRLicenseKeyInSecretManager` is set to `true`. |
| `NEW_RELIC_REGION`  | The New Relic region to which data will be sent (set to the specified value for `NRRegion`). |
| `DEBUG_ENABLED`   | Enables debug logging for the Lambda function (modifiable in the AWS console). By default this field is set to `false`. |
| `CUSTOM_META_DATA` | Custom metadata set to the specified value for `CommonAttributes`.  |

**Note:**
- An S3 bucket will be created to store the packaged Lambda function.
- A secret will be created in AWS Secrets Manager to store the New Relic license key if `LICENSE_KEY_FETCH_FROM_SECRET_MANAGER` is set to `true`.
- Creating an AWS secret may incur additional costs as reads during every cold start of this Lambda function.
- IAM roles and policies will be created as needed.


#### Commands for deployment:


- To create a build:

    ```shell
    make build 
    ````

- To deploy:
  
    ```shell
    make deploy 
    ````

- To delete a cloudformation stack :

    ```shell
    make delete STACK_NAME=<stack-name>
    ```

**Note**: All the above commands run with the expectation that the AWS default configuration is available on the machine.


## Templates

### For logs and metrics:

1. [logging-firehose-metric-polling.yaml](/logging-firehose-metric-polling.yaml): To fetch logs using Firehose and metrics using polling.
2. [logging-firehose-metric-stream.yaml](/logging-firehose-metric-stream.yaml): To fetch logs using Firehose and metrics using streaming.
3. [logging-lambda-firehose-metric-polling.yaml](/logging-lambda-firehose-metric-polling.yaml): To fetch logs using Lambda and Firehose, and metrics using polling.
4. [logging-lambda-firehose-metric-stream.yaml](/logging-lambda-firehose-metric-stream.yaml): To fetch logs using Lambda and Firehose, and metrics using streaming.
5. [logging-lambda-metric-polling.yaml](/logging-lambda-metric-polling.yaml): To fetch logs using Lambda and metrics using polling.
6. [logging-lambda-metric-stream.yaml](/logging-lambda-metric-stream.yaml):To fetch logs using Lambda and metrics using streaming.

### For logs:

1. [lambda-template.yaml](/lambda-template.yaml): To fetch logs using Lambda.
2. [logging-lambda-firehose-template.yaml](/logging-lambda-firehose-template.yaml): To fetch logs using Firehose and Lambda.


## Building and packaging

To build and package, follow these steps for each template:

1. Authenticate your AWS account.
2. Create an S3 bucket with name. For example, `test123`.
3. To create a build, run: `sam build -u --template-file fileName.yaml`
    **Note:** By default, build will be available at `.aws-sam/build` with the generated `template.yaml`
4. To package the build, run: `sam package --s3-bucket test123 --template-file .aws-sam/build/template.yaml  --output-template-file fileName.yaml --region us-east-2`
5. Copy the main template file to the S3 bucket using : `aws s3 cp .aws-sam/build/fileName.yaml s3://test123/fileName.yaml`

