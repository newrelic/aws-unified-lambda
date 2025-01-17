#!/bin/bash

source common-scripts.sh
source config-file.cfg

# test case constants
S3_TRIGGER_CASE=e2e-s3-trigger-stack

deploy_s3_trigger_stack() {
  template_file=$1
  stack_name=$2
  license_key=$3
  new_relic_region=$4
  new_relic_account_id=$5
  secret_license_key=$6
  s3_bucket_names=$7
  common_attributes=$8

  echo "Deploying s3 trigger stack with name: $stack_name"

  sam deploy \
    --template-file "$template_file" \
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

validate_lambda_s3_trigger_created() {
  # this function fetches bucket configurations and
  # validates if lambda event notification is configured
  stack_name=$1
  bucket_name=$2
  bucket_prefix=$3

  echo "Validating s3 lambda trigger event for stack name: $stack_name, bucket name: $bucket_name, and prefix: $bucket_prefix"

  lambda_function_arn=$(get_lambda_function_arn "$stack_name")

  notification_configuration=$(aws s3api get-bucket-notification-configuration --bucket "$bucket_name")

  lambda_configurations=$(echo "$notification_configuration" |
      jq --arg lambda "$lambda_function_arn" --arg prefix "$bucket_prefix" '
      .LambdaFunctionConfigurations[] |
      select(.LambdaFunctionArn == $lambda and (.Filter.Key.FilterRules[]? | select(.Name == "Prefix" and .Value == $prefix)?))')

  if [ -n "$lambda_configurations" ]; then
      echo "S3 triggers with prefix '$bucket_prefix' found for Lambda function $lambda_function_arn on bucket $bucket_name:"
      echo "$lambda_configurations" | jq '.'
  else
      exit_with_error "No S3 triggers with prefix '$bucket_prefix' found for Lambda function $lambda_function_arn on bucket $bucket_name."
  fi
}

test_logs_s3() {
  echo "came"
  cat <<EOF > s3-parameter.json
  '[{"bucket":"$S3_BUCKET_NAME","prefix":"$S3_BUCKET_PREFIX"}]'
EOF
  S3_BUCKET_NAMES=$(<s3-parameter.json)
  echo "Testing for s3 bucket configuration JSON: $(<s3-parameter.json)"

  log_message=$(create_log_message "$LOG_MESSAGE_CLOUDWATCH" "$LOG_GROUP_FILTER_PATTERN")

  deploy_s3_trigger_stack "$LAMBDA_TEMPLATE_BUILD_DIR/$LAMBDA_TEMPLATE" "$S3_TRIGGER_CASE" "$NEW_RELIC_LICENSE_KEY" "$NEW_RELIC_REGION" "$NEW_RELIC_ACCOUNT_ID" "false" "$S3_BUCKET_NAMES" "''"
  validate_stack_deployment_status "$S3_TRIGGER_CASE"
  validate_lambda_s3_trigger_created "$S3_TRIGGER_CASE" "$S3_BUCKET_NAME" "$S3_BUCKET_PREFIX"
  upload_file_to_s3_bucket "$S3_BUCKET_NAME" "$S3_BUCKET_OBJECT_NAME" "$S3_BUCKET_PREFIX" "$log_message"
  validate_logs_in_new_relic "$NEW_RELIC_USER_KEY" "$NEW_RELIC_ACCOUNT_ID" "$ATTRIBUTE_KEY_S3" "$S3_BUCKET_PREFIX" "$LOG_MESSAGE_S3"
  delete_stack "$S3_TRIGGER_CASE"
}
