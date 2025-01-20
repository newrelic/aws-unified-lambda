#!/bin/bash

source common-scripts/stack-scripts.sh
source common-scripts/resource-scripts.sh
source common-scripts/logs-scripts.sh
source common-scripts/config-file.cfg

test_logs_cloudwatch() {
cat <<EOF > cloudwatch-parameter.json
'[{"LogGroupName":"$LOG_GROUP_NAME","FilterPattern":"$LOG_GROUP_FILTER_PATTERN"}]'
EOF

cat <<EOF > common_attribute.json
'[{"AttributeName":"$COMMON_ATTRIBUTE_KEY","AttributeValue":"$COMMON_ATTRIBUTE_VALUE"}]'
EOF
COMMON_ATTRIBUTES=$(<common_attribute.json)

  LOG_GROUP_NAMES=$(<cloudwatch-parameter.json)
  echo "Testing for log group configuration JSON: $(<cloudwatch-parameter.json)"

  log_message=$(create_log_message "$LOG_MESSAGE_CLOUDWATCH" "$LOG_GROUP_FILTER_PATTERN")

  deploy_cloudwatch_trigger_stack "$LAMBDA_TEMPLATE_BUILD_DIR/$LAMBDA_TEMPLATE" "$CLOUDWATCH_TRIGGER_CASE" "$NEW_RELIC_LICENSE_KEY" "$NEW_RELIC_REGION" "$NEW_RELIC_ACCOUNT_ID" "false"  "$LOG_GROUP_NAMES" "$COMMON_ATTRIBUTES"
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