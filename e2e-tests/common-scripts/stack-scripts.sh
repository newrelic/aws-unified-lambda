#!/bin/bash

source config-file.cfg

deploy_cloudwatch_trigger_stack() {
  stack_name=$1
  license_key=$2
  new_relic_region=$3
  new_relic_account_id=$4
  secret_license_key=$5
  log_group_config=$6
  common_attributes=$7

  echo "Deploying cloudwatch trigger stack with name: $stack_name"

  sam deploy \
    --template-file "$TEMPLATES_BUILD_DIR/$LAMBDA_TEMPLATE" \
    --stack-name "$stack_name" \
    --parameter-overrides \
      LicenseKey="$license_key" \
      NewRelicRegion="$new_relic_region" \
      NewRelicAccountId="$new_relic_account_id" \
      StoreNRLicenseKeyInSecretManager="$secret_license_key" \
      S3BucketNames="''" \
      LogGroupConfig="$log_group_config" \
      CommonAttributes="$common_attributes" \
    --capabilities CAPABILITY_IAM
}

deploy_s3_trigger_stack() {
  stack_name=$1
  license_key=$2
  new_relic_region=$3
  new_relic_account_id=$4
  secret_license_key=$5
  s3_bucket_names=$6
  common_attributes=$7

  echo "Deploying s3 trigger stack with name: $stack_name"

  sam deploy \
    --template-file "$TEMPLATES_BUILD_DIR/$LAMBDA_TEMPLATE" \
    --stack-name "$stack_name" \
    --parameter-overrides \
      LicenseKey="$license_key" \
      NewRelicRegion="$new_relic_region" \
      NewRelicAccountId="$new_relic_account_id" \
      StoreNRLicenseKeyInSecretManager="$secret_license_key" \
      S3BucketNames="$s3_bucket_names" \
      LogGroupConfig="''" \
      CommonAttributes="$common_attributes" \
    --capabilities CAPABILITY_IAM
}

deploy_firehose_metric_polling_stack() {
  polling_integration_slugs=$1
  s3_bucket_names=$2
  log_group_config=$3
  common_attributes=$4
  store_license_key_in_secret_manager=$5

  echo "Deploying New Relic stack with name: $stack_name"

  aws cloudformation deploy \
    --template-file "$TEMPLATES_BUILD_DIR/$FIREHOSE_METRIC_POLLING_TEMPLATE" \
    --stack-name "$FIREHOSE_METRIC_POLLING_CASE" \
    --parameter-overrides \
      IAMRoleName="$IAM_ROLE_NAME" \
      NewRelicAccountId="$NEW_RELIC_ACCOUNT_ID" \
      NewRelicRegion="$NEW_RELIC_REGION" \
      IntegrationName="$FIREHOSE_METRIC_POLLING_CASE" \
      NewRelicAPIKey="$NEW_RELIC_USER_KEY" \
      PollingIntegrationSlugs="$polling_integration_slugs" \
      NewRelicLicenseKey="$NEW_RELIC_LICENSE_KEY" \
      S3BucketNames="$s3_bucket_names" \
      LogGroupConfig="$log_group_config" \
      CommonAttributes="$common_attributes" \
      StoreNRLicenseKeyInSecretManager="$store_license_key_in_secret_manager" \
    --capabilities CAPABILITY_IAM
}

deploy_firehose_metric_streaming_stack() {
  stack_name=$1
  iam_role_name=$2
  new_relic_account_id=$3
  integration_name=$4
  new_relic_api_key=$5
  polling_integration_slugs=$6
  metric_collection_mode=$7
  firehose_stream_name=$8
  cloudwatch_metric_stream_name=$9
  s3_backup_bucket_name=${10}
  create_config_service=${11}
  s3_config_bucket_name=${12}
  new_relic_license_key=${13}
  new_relic_region=${14}
  log_group_config=${15}
  logging_firehose_stream_name=${16}
  logging_s3_backup_bucket_name=${17}
  enable_cloudwatch_logging_for_firehose=${18}
  common_attributes=${19}
  store_license_key_in_secret_manager=${20}

  echo "Deploying New Relic integration stack with name: $stack_name"

  aws cloudformation deploy \
    --template-file "$TEMPLATES_BUILD_DIR/$FIREHOSE_METRIC_STREAMING_TEMPLATE" \
    --stack-name "$stack_name" \
    --parameter-overrides \
      IAMRoleName="$iam_role_name" \
      NewRelicAccountId="$new_relic_account_id" \
      IntegrationName="$integration_name" \
      NewRelicAPIKey="$new_relic_api_key" \
      PollingIntegrationSlugs="$polling_integration_slugs" \
      MetricCollectionMode="$metric_collection_mode" \
      FirehoseStreamName="$firehose_stream_name" \
      CloudWatchMetricStreamName="$cloudwatch_metric_stream_name" \
      S3BackupBucketName="$s3_backup_bucket_name" \
      CreateConfigService="$create_config_service" \
      S3ConfigBucketName="$s3_config_bucket_name" \
      NewRelicLicenseKey="$new_relic_license_key" \
      NewRelicRegion="$new_relic_region" \
      LogGroupConfig="$log_group_config" \
      LoggingFirehoseStreamName="$logging_firehose_stream_name" \
      LoggingS3BackupBucketName="$logging_s3_backup_bucket_name" \
      EnableCloudWatchLoggingForFirehose="$enable_cloudwatch_logging_for_firehose" \
      CommonAttributes="$common_attributes" \
      StoreNRLicenseKeyInSecretManager="$store_license_key_in_secret_manager" \
    --capabilities CAPABILITY_IAM
}

deploy_lambda_firehose_metric_polling_stack() {
  stack_name=$1
  iam_role_name=$2
  new_relic_account_id=$3
  new_relic_region=$4
  integration_name=$5
  new_relic_api_key=$6
  polling_integration_slugs=$7
  new_relic_license_key=$8
  s3_bucket_names=$9
  log_group_config=${10}
  common_attributes=${11}
  logging_firehose_stream_name=${12}
  logging_s3_backup_bucket_name=${13}
  enable_cloudwatch_logging_for_firehose=${14}
  store_license_key_in_secret_manager=${15}

  echo "Deploying New Relic CloudFormation stack with name: $stack_name"

  aws cloudformation deploy \
    --template-file "$TEMPLATES_BUILD_DIR/$LAMBDA_FIREHOSE_METRIC_POLLING_TEMPLATE" \
    --stack-name "$stack_name" \
    --parameter-overrides \
      IAMRoleName="$iam_role_name" \
      NewRelicAccountId="$new_relic_account_id" \
      NewRelicRegion="$new_relic_region" \
      IntegrationName="$integration_name" \
      NewRelicAPIKey="$new_relic_api_key" \
      PollingIntegrationSlugs="$polling_integration_slugs" \
      NewRelicLicenseKey="$new_relic_license_key" \
      S3BucketNames="$s3_bucket_names" \
      LogGroupConfig="$log_group_config" \
      CommonAttributes="$common_attributes" \
      LoggingFirehoseStreamName="$logging_firehose_stream_name" \
      LoggingS3BackupBucketName="$logging_s3_backup_bucket_name" \
      EnableCloudWatchLoggingForFirehose="$enable_cloudwatch_logging_for_firehose" \
      StoreNRLicenseKeyInSecretManager="$store_license_key_in_secret_manager" \
    --capabilities CAPABILITY_IAM
}

deploy_lambda_firehose_metric_streaming_stack() {
  stack_name=$1
  iam_role_name=$2
  new_relic_account_id=$3
  new_relic_region=$4
  integration_name=$5
  new_relic_api_key=$6
  polling_integration_slugs=$7
  metric_collection_mode=$8
  new_relic_license_key=$9
  firehose_stream_name=${10}
  cloudwatch_metric_stream_name=${11}
  s3_backup_bucket_name=${12}
  create_config_service=${13}
  s3_config_bucket_name=${14}
  s3_bucket_names=${15}
  log_group_config=${16}
  common_attributes=${17}
  logging_firehose_stream_name=${18}
  logging_s3_backup_bucket_name=${19}
  enable_cloudwatch_logging_for_firehose=${20}
  store_license_key_in_secret_manager=${21}

  echo "Deploying New Relic AWS stack with name: $stack_name"

  aws cloudformation deploy \
    --template-file "$TEMPLATES_BUILD_DIR/$LAMBDA_FIREHOSE_METRIC_STREAMING_TEMPLATE" \
    --stack-name "$stack_name" \
    --parameter-overrides \
      IAMRoleName="$iam_role_name" \
      NewRelicAccountId="$new_relic_account_id" \
      NewRelicRegion="$new_relic_region" \
      IntegrationName="$integration_name" \
      NewRelicAPIKey="$new_relic_api_key" \
      PollingIntegrationSlugs="$polling_integration_slugs" \
      MetricCollectionMode="$metric_collection_mode" \
      NewRelicLicenseKey="$new_relic_license_key" \
      FirehoseStreamName="$firehose_stream_name" \
      CloudWatchMetricStreamName="$cloudwatch_metric_stream_name" \
      S3BackupBucketName="$s3_backup_bucket_name" \
      CreateConfigService="$create_config_service" \
      S3ConfigBucketName="$s3_config_bucket_name" \
      S3BucketNames="$s3_bucket_names" \
      LogGroupConfig="$log_group_config" \
      CommonAttributes="$common_attributes" \
      LoggingFirehoseStreamName="$logging_firehose_stream_name" \
      LoggingS3BackupBucketName="$logging_s3_backup_bucket_name" \
      EnableCloudWatchLoggingForFirehose="$enable_cloudwatch_logging_for_firehose" \
      StoreNRLicenseKeyInSecretManager="$store_license_key_in_secret_manager" \
    --capabilities CAPABILITY_IAM
}

deploy_lambda_firehose_stack() {
  s3_bucket_names=$1
  log_group_config=$2
  common_attributes=$3
  enable_cloudwatch_logging_for_firehose=$4
  store_license_key_in_secret_manager=$5

  echo "Deploying New Relic Logging stack with name: $stack_name"

  aws cloudformation deploy \
    --template-file "$TEMPLATES_BUILD_DIR/$LAMBDA_FIREHOSE_TEMPLATE" \
    --stack-name "$LAMBDA_FIREHOSE_CASE" \
    --parameter-overrides \
      LicenseKey="$NEW_RELIC_LICENSE_KEY" \
      NewRelicRegion="$NEW_RELIC_REGION" \
      NewRelicAccountId="$NEW_RELIC_ACCOUNT_ID" \
      S3BucketNames="$s3_bucket_names" \
      LogGroupConfig="$log_group_config" \
      CommonAttributes="$common_attributes" \
      LoggingFirehoseStreamName="$FIREHOSE_STREAM_NAME" \
      LoggingS3BackupBucketName="$BACKUP_BUCKET_NAME" \
      EnableCloudWatchLoggingForFirehose="$enable_cloudwatch_logging_for_firehose" \
      StoreNRLicenseKeyInSecretManager="$store_license_key_in_secret_manager" \
    --capabilities CAPABILITY_IAM
}

validate_stack_deployment_status() {
  stack_name=$1

  echo "Validating stack deployment status for stack name: $stack_name"

  stack_status=$(aws cloudformation describe-stacks --stack-name "$stack_name" --query "Stacks[0].StackStatus" --output text)
  if [[ "$stack_status" == "ROLLBACK_COMPLETE" || "$stack_status" == "ROLLBACK_FAILED" || "$stack_status" == "CREATE_FAILED"  || "$stack_status" == "UPDATE_FAILED" ]]; then
    echo "Stack $stack_name failed to be created and rolled back."
    failure_reason=$(aws cloudformation describe-stack-events --stack-name "$stack_name" --query "StackEvents[?ResourceStatus==\`$stack_status\`].ResourceStatusReason" --output text)
    exit_with_error "Stack $stack_name failed to be created. Failure reason: $failure_reason"
  else
    echo "Stack $stack_name was created successfully."
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