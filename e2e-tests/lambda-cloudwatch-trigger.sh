#!/bin/bash

source common-scripts/stack-scripts.sh
source common-scripts/resource-scripts.sh
source common-scripts/logs-scripts.sh
source common-scripts/test-configs.cfg

test_logs_with_filter_pattern() {
  delete_stack "$CLOUDWATCH_FILTER_PATTERN_CASE"
  delete_stack "$CLOUDWATCH_SECRET_MANAGER_CASE"
  delete_stack "$CLOUDWATCH_INVALID_LOG_GROUP_CASE"
}

#test_logs_for_secret_manager() {
#cat <<EOF > cloudwatch-parameter.json
#'[{"LogGroupName":"$LOG_GROUP_NAME_SECRET_MANAGER","FilterPattern":"$LOG_GROUP_FILTER_PATTERN"}]'
#EOF
#
#cat <<EOF > common_attribute.json
#'[{"AttributeName":"$CUSTOM_ATTRIBUTE_KEY","AttributeValue":"$CUSTOM_ATTRIBUTE_VALUE"}]'
#EOF
#COMMON_ATTRIBUTES=$(<common_attribute.json)
#
#  LOG_GROUP_NAMES=$(<cloudwatch-parameter.json)
#  echo "Testing for log group configuration JSON: $(<cloudwatch-parameter.json)"
#
#  log_message=$(create_log_message "$LOG_MESSAGE_CLOUDWATCH" "$LOG_GROUP_FILTER_PATTERN")
#
#  deploy_cloudwatch_trigger_stack "$CLOUDWATCH_SECRET_MANAGER_CASE" "true" "$LOG_GROUP_NAMES" "$COMMON_ATTRIBUTES"
#  validate_stack_deployment_status "$CLOUDWATCH_SECRET_MANAGER_CASE"
#  validate_lambda_subscription_created "$CLOUDWATCH_SECRET_MANAGER_CASE" "$LOG_GROUP_NAME_SECRET_MANAGER" "$LOG_GROUP_FILTER_PATTERN"
#  create_cloudwatch_log_event "$LOG_GROUP_NAME_SECRET_MANAGER" "$LOG_STREAM_NAME" "$log_message"
#  validate_logs_in_new_relic "$NEW_RELIC_USER_KEY" "$NEW_RELIC_ACCOUNT_ID" "$ATTRIBUTE_KEY_CLOUDWATCH" "$LOG_STREAM_NAME" "$log_message"
#
#}

#test_logs_for_invalid_log_group() {
#cat <<EOF > cloudwatch-parameter.json
#'[{"LogGroupName":"$LOG_GROUP_NAME_INVALID","FilterPattern":"$LOG_GROUP_FILTER_PATTERN"}]'
#EOF
#
#  LOG_GROUP_NAMES=$(<cloudwatch-parameter.json)
#  echo "Testing for log group configuration JSON: $(<cloudwatch-parameter.json)"
#
#  deploy_cloudwatch_trigger_stack "$CLOUDWATCH_INVALID_LOG_GROUP_CASE" "false" "$LOG_GROUP_NAMES" "''"
#  validate_stack_deployment_status "$CLOUDWATCH_INVALID_LOG_GROUP_CASE"
#
#  # validate that lambda subscription is not created
#  validate_lambda_subscription_not_created "$CLOUDWATCH_INVALID_LOG_GROUP_CASE" "$LOG_GROUP_NAME_INVALID" "$LOG_GROUP_FILTER_PATTERN"
#
#}

case $1 in
  test_logs_with_filter_pattern)
    test_logs_with_filter_pattern
    ;;
#  test_logs_for_secret_manager)
#    test_logs_for_secret_manager
#    ;;
#  test_logs_for_invalid_log_group)
#    test_logs_for_invalid_log_group
#    ;;
  *)
  echo "Invalid test case specified."
  exit 1
  ;;
esac