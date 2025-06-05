#!/bin/bash

source common-scripts/stack-scripts.sh
source common-scripts/resource-scripts.sh
source common-scripts/logs-scripts.sh
source common-scripts/test-configs.cfg

test_logs_with_filter_pattern() {
cat <<EOF > cloudwatch-parameter.json
'[{"LogGroupName":"$LOG_GROUP_NAME","FilterPattern":"$LOG_GROUP_FILTER_PATTERN"}]'
EOF

cat <<EOF > common_attribute.json
'[{"AttributeName":"$CUSTOM_ATTRIBUTE_KEY","AttributeValue":"$CUSTOM_ATTRIBUTE_VALUE"}]'
EOF
COMMON_ATTRIBUTES=$(<common_attribute.json)

  LOG_GROUP_NAMES=$(<cloudwatch-parameter.json)
  echo "Testing for log group configuration JSON: $(<cloudwatch-parameter.json)"

  delete_stack "$CLOUDWATCH_FILTER_PATTERN_CASE"
}

test_logs_for_secret_manager() {
cat <<EOF > cloudwatch-parameter.json
'[{"LogGroupName":"$LOG_GROUP_NAME_SECRET_MANAGER","FilterPattern":"$LOG_GROUP_FILTER_PATTERN"}]'
EOF

cat <<EOF > common_attribute.json
'[{"AttributeName":"$CUSTOM_ATTRIBUTE_KEY","AttributeValue":"$CUSTOM_ATTRIBUTE_VALUE"}]'
EOF
COMMON_ATTRIBUTES=$(<common_attribute.json)

  LOG_GROUP_NAMES=$(<cloudwatch-parameter.json)
  delete_stack "$CLOUDWATCH_SECRET_MANAGER_CASE"
}

test_logs_for_invalid_log_group() {
cat <<EOF > cloudwatch-parameter.json
'[{"LogGroupName":"$LOG_GROUP_NAME_INVALID","FilterPattern":"$LOG_GROUP_FILTER_PATTERN"}]'
EOF

  LOG_GROUP_NAMES=$(<cloudwatch-parameter.json)
  echo "Testing for log group configuration JSON: $(<cloudwatch-parameter.json)"

  
  delete_stack "$CLOUDWATCH_INVALID_LOG_GROUP_CASE"
}

case $1 in
  test_logs_with_filter_pattern)
    test_logs_with_filter_pattern
    ;;
  test_logs_for_secret_manager)
    test_logs_for_secret_manager
    ;;
  test_logs_for_invalid_log_group)
    test_logs_for_invalid_log_group
    ;;
  *)
  echo "Invalid test case specified."
  exit 1
  ;;
esac