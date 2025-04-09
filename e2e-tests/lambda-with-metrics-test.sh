#!/bin/bash

source common-scripts/stack-scripts.sh
source common-scripts/resource-scripts.sh
source common-scripts/logs-scripts.sh
source common-scripts/test-configs.cfg

test_for_lambda_firehose_stack() {
  delete_stack "$LAMBDA_FIREHOSE_CASE"
  delete_stack "$LAMBDA_METRIC_POLLING_CASE"
  delete_stack "$LAMBDA_METRIC_STREAMING_CASE"
  delete_stack "$FIREHOSE_METRIC_POLLING_CASE"
  delete_stack "$FIREHOSE_METRIC_STREAMING_CASE"
  delete_stack "$LAMBDA_FIREHOSE_METRIC_POLLING_CASE"
  delete_stack "$LAMBDA_FIREHOSE_METRIC_STREAMING_CASE"
}

case $1 in
  test_for_lambda_firehose_stack)
    test_for_lambda_firehose_stack
    ;;
  *)
  echo "Invalid test case specified."
  exit 1
  ;;
esac