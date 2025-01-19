#!/bin/bash

source common-scripts/stack-scripts.sh
source common-scripts/resource-scripts.sh
source common-scripts/logs-scripts.sh
source common-scripts/config-file.cfg

# test case constants
CLOUDWATCH_TRIGGER_CASE=e2e-cloudwatch-trigger-stack

deploy_cloudwatch_trigger_stack() {
  template_file=$1
  stack_name=$2
  license_key=$3
  new_relic_region=$4
  new_relic_account_id=$5
  secret_license_key=$6
  log_group_config=$7
  common_attributes=$8

  echo "Deploying cloudwatch trigger stack with name: $stack_name"

  sam deploy \
    --template-file "$template_file" \
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

test_logs_cloudwatch() {
cat <<EOF > cloudwatch-parameter.json
'[{"LogGroupName":"$LOG_GROUP_NAME","FilterPattern":"$LOG_GROUP_FILTER_PATTERN"}]'
EOF
  LOG_GROUP_NAMES=$(<cloudwatch-parameter.json)
  echo "Testing for log group configuration JSON: $(<cloudwatch-parameter.json)"

  log_message=$(create_log_message "$LOG_MESSAGE_CLOUDWATCH" "$LOG_GROUP_FILTER_PATTERN")

  deploy_cloudwatch_trigger_stack "$LAMBDA_TEMPLATE_BUILD_DIR/$LAMBDA_TEMPLATE" "$CLOUDWATCH_TRIGGER_CASE" "$NEW_RELIC_LICENSE_KEY" "$NEW_RELIC_REGION" "$NEW_RELIC_ACCOUNT_ID" "false"  "$LOG_GROUP_NAMES" "''"
  validate_stack_deployment_status "$CLOUDWATCH_TRIGGER_CASE"
  validate_lambda_subscription_created "$CLOUDWATCH_TRIGGER_CASE" "$LOG_GROUP_NAME" "$LOG_GROUP_FILTER_PATTERN"
  create_cloudwatch_log_event "$LOG_GROUP_NAME" "$LOG_STREAM_NAME" "$log_message"
  validate_logs_in_new_relic "$NEW_RELIC_USER_KEY" "$NEW_RELIC_ACCOUNT_ID" "$ATTRIBUTE_KEY_CLOUDWATCH" "$LOG_STREAM_NAME" "$log_message"
  delete_stack "$CLOUDWATCH_TRIGGER_CASE"
}

case $1 in
  test_logs_cloudwatch)
    test_logs_cloudwatch
    ;;
  *)
  echo "Invalid test case specified."
  exit 1
  ;;
esac