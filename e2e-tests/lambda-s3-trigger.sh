#!/bin/bash

source common-scripts/stack-scripts.sh
source common-scripts/resource-scripts.sh
source common-scripts/logs-scripts.sh
source common-scripts/test-configs.cfg

test_logs_for_prefix() {
  delete_stack "$S3_SECRET_MANAGER_CASE"
  delete_stack "$S3_PREFIX_CASE"
  delete_stack "$S3_INVALID_BUCKET_CASE"
}

case $1 in
  test_logs_for_prefix)
    test_logs_for_prefix
    ;;
  *)
  echo "Invalid test case specified."
  exit 1
  ;;
esac