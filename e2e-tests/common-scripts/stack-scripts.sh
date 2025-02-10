#!/bin/bash

source config-file.cfg

deploy_cloudwatch_trigger_stack() {
  stack_name=$1
  secret_license_key=$2
  log_group_config=$3
  common_attributes=$4

  echo "Deploying cloudwatch trigger stack with name: $stack_name"

  sam deploy \
    --template-file "$TEMPLATE_BUILD_DIR/$LAMBDA_TEMPLATE" \
    --stack-name "$stack_name" \
    --parameter-overrides \
      LicenseKey="$NEW_RELIC_LICENSE_KEY" \
      NewRelicRegion="$NEW_RELIC_REGION" \
      NewRelicAccountId="$NEW_RELIC_ACCOUNT_ID" \
      StoreNRLicenseKeyInSecretManager="$secret_license_key" \
      S3BucketNames="''" \
      LogGroupConfig="$log_group_config" \
      CommonAttributes="$common_attributes" \
    --capabilities CAPABILITY_IAM
}

deploy_s3_trigger_stack() {
  stack_name=$1
  secret_license_key=$2
  s3_bucket_names=$3
  common_attributes=$4

  echo "Deploying s3 trigger stack with name: $stack_name"

  sam deploy \
    --template-file "$TEMPLATE_BUILD_DIR/$LAMBDA_TEMPLATE" \
    --stack-name "$stack_name" \
    --parameter-overrides \
      LicenseKey="$NEW_RELIC_LICENSE_KEY" \
      NewRelicRegion="$NEW_RELIC_REGION" \
      NewRelicAccountId="$NEW_RELIC_ACCOUNT_ID" \
      StoreNRLicenseKeyInSecretManager="$secret_license_key" \
      S3BucketNames="$s3_bucket_names" \
      LogGroupConfig="''" \
      CommonAttributes="$common_attributes" \
    --capabilities CAPABILITY_IAM
}

deploy_lambda_firehose_stack() {
  s3_bucket_names=$1
  log_group_config=$2
  common_attributes=$3
  enable_cloudwatch_logging_for_firehose=$4
  store_license_key_in_secret_manager=$5

  echo "Deploying lambda-firehose stack with name: $LAMBDA_FIREHOSE_TEMPLATE_CASE"

  aws cloudformation deploy \
    --template-file "$TEMPLATE_BUILD_DIR/$LAMBDA_FIREHOSE_TEMPLATE" \
    --stack-name "$LAMBDA_FIREHOSE_CASE" \
    --parameter-overrides \
      LicenseKey="$NEW_RELIC_LICENSE_KEY" \
      NewRelicRegion="$NEW_RELIC_REGION" \
      NewRelicAccountId="$NEW_RELIC_ACCOUNT_ID" \
      S3BucketNames="$s3_bucket_names" \
      LogGroupConfig="$log_group_config" \
      CommonAttributes="$common_attributes" \
      LoggingFirehoseStreamName="$FIREHOSE_STREAM_NAME_LAMBDA_FIREHOSE_CASE" \
      LoggingS3BackupBucketName="$BACKUP_BUCKET_NAME" \
      EnableCloudWatchLoggingForFirehose="$enable_cloudwatch_logging_for_firehose" \
      StoreNRLicenseKeyInSecretManager="$store_license_key_in_secret_manager" \
    --capabilities CAPABILITY_IAM CAPABILITY_AUTO_EXPAND CAPABILITY_NAMED_IAM
}

deploy_lambda_metric_polling_stack() {
  log_group_config=$1
  common_attributes=$2
  store_license_key_in_secret_manager=$3
  s3_bucket_names=$4

  echo "Deploying lambda metric polling stack with name: $LAMBDA_METRIC_POLLING_CASE"

  aws cloudformation deploy \
    --template-file "$TEMPLATE_BUILD_DIR/$LAMBDA_METRIC_POLLING_TEMPLATE" \
    --stack-name "$LAMBDA_METRIC_POLLING_CASE" \
    --parameter-overrides \
      IAMRoleName="$IAM_ROLE_NAME" \
      NewRelicAccountId="$NEW_RELIC_ACCOUNT_ID" \
      IntegrationName="$LAMBDA_METRIC_POLLING_CASE" \
      NewRelicAPIKey="$NEW_RELIC_USER_KEY" \
      PollingIntegrationSlugs="$POLLING_INTEGRATION_SLUGS" \
      NewRelicLicenseKey="$NEW_RELIC_LICENSE_KEY" \
      NewRelicRegion="$NEW_RELIC_REGION" \
      S3BucketNames="$s3_bucket_names" \
      LogGroupConfig="$log_group_config" \
      CommonAttributes="$common_attributes" \
      StoreNRLicenseKeyInSecretManager="$store_license_key_in_secret_manager" \
    --capabilities CAPABILITY_IAM CAPABILITY_AUTO_EXPAND CAPABILITY_NAMED_IAM
}

deploy_lambda_metric_streaming_stack() {
  log_group_config=$1
  common_attributes=$2
  store_license_key_in_secret_manager=$3
  s3_bucket_names=$4

  echo "Deploying lambda metric streaming stack with name: $LAMBDA_METRIC_STREAMING_CASE"

  aws cloudformation deploy \
    --template-file "$TEMPLATE_BUILD_DIR/$LAMBDA_METRIC_STREAMING_TEMPLATE" \
    --stack-name "$LAMBDA_METRIC_STREAMING_CASE" \
    --parameter-overrides \
      IAMRoleName="$IAM_ROLE_NAME" \
      NewRelicAccountId="$NEW_RELIC_ACCOUNT_ID" \
      IntegrationName="$LAMBDA_METRIC_STREAMING_CASE" \
      NewRelicAPIKey="$NEW_RELIC_USER_KEY" \
      PollingIntegrationSlugs="$POLLING_INTEGRATION_SLUGS" \
      NewRelicLicenseKey="$NEW_RELIC_LICENSE_KEY" \
      NewRelicRegion="$NEW_RELIC_REGION" \
      MetricCollectionMode="$METRIC_COLLECTION_MODE" \
      FirehoseStreamName="$METRIC_FIREHOSE_STREAM_NAME_LAMBDA_METRICS_STREAMING_CASE" \
      CloudWatchMetricStreamName="$CLOUDWATCH_METRIC_STREAM_NAME_LAMBDA_METRICS_STREAMING_CASE" \
      S3BackupBucketName="$BACKUP_BUCKET_NAME" \
      CreateConfigService="false" \
      S3ConfigBucketName="$S3_CONFIG_BUCKET_NAME" \
      S3BucketNames="$s3_bucket_names" \
      LogGroupConfig="$log_group_config" \
      CommonAttributes="$common_attributes" \
      StoreNRLicenseKeyInSecretManager="$store_license_key_in_secret_manager" \
    --capabilities CAPABILITY_IAM CAPABILITY_AUTO_EXPAND CAPABILITY_NAMED_IAM
}

validate_stack_deployment_status() {
  stack_name=$1

  echo "Validating stack deployment status for stack name: $stack_name"

  stack_status=$(aws cloudformation describe-stacks --stack-name "$stack_name" --query "Stacks[0].StackStatus" --output text)
  if [[ "$stack_status" == "CREATE_COMPLETE" ]]; then
    echo "Stack $stack_name was created successfully."
  else
    echo "Stack $stack_name failed to be created and rolled back."
    failure_reason=$(aws cloudformation describe-stack-events --stack-name "$stack_name" --query "StackEvents[?ResourceStatus==\`$stack_status\`].ResourceStatusReason" --output text)
    exit_with_error "Stack $stack_name failed to be created. Failure reason: $failure_reason"
  fi
}

delete_stack() {
  stack_name=$1

  aws cloudformation delete-stack --stack-name "$stack_name"

  echo "Initiated deletion of stack: $stack_name"

  stack_status=$(aws cloudformation describe-stacks --stack-name "$stack_name" --query 'Stacks[0].StackStatus' --output text)

  # delete stack with exponential back off retires
  max_sleep_time=300  # Cap sleep time to 5 minutes
  sleep_time=30

  while [[ $stack_status == "DELETE_IN_PROGRESS" ]]; do
    echo "Stack $stack_name is still being deleted..."

    sleep $sleep_time
    if (( sleep_time < max_sleep_time )); then
      sleep_time=$(( sleep_time * 2 ))
    fi

    stack_status=$(aws cloudformation describe-stacks --stack-name "$stack_name" --query 'Stacks[0].StackStatus' --output text 2>/dev/null || true)
  done

  if [ -z "$stack_status" ]; then
    echo "Stack $stack_name has been successfully deleted."
  else
    exit_with_error "Failed to delete stack $stack_name."
  fi
}

exit_with_error() {
  echo "Error: $1"
  exit 1
}