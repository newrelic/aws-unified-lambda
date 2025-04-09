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

case $1 in
  test_logs_with_filter_pattern)
    test_logs_with_filter_pattern
    ;;
  *)
  echo "Invalid test case specified."
  exit 1
  ;;
esac