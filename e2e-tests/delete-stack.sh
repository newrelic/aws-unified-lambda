#!/bin/bash

source common-scripts/stack-scripts.sh
source common-scripts/resource-scripts.sh
source common-scripts/logs-scripts.sh
source common-scripts/test-configs.cfg

delete_stack "$LAMBDA_FIREHOSE_CASE"
delete_stack "$LAMBDA_METRIC_POLLING_CASE"
delete_stack "$LAMBDA_METRIC_STREAMING_CASE"
delete_stack "$FIREHOSE_METRIC_POLLING_CASE"
delete_stack "$FIREHOSE_METRIC_STREAMING_CASE"
delete_stack "$LAMBDA_FIREHOSE_METRIC_POLLING_CASE"
delete_stack "$LAMBDA_FIREHOSE_METRIC_STREAMING_CASE"
delete_stack "$S3_SECRET_MANAGER_CASE"
delete_stack "$S3_PREFIX_CASE"
delete_stack "$S3_INVALID_BUCKET_CASE"
delete_stack "$CLOUDWATCH_FILTER_PATTERN_CASE"
delete_stack "$CLOUDWATCH_SECRET_MANAGER_CASE"
delete_stack "$CLOUDWATCH_INVALID_LOG_GROUP_CASE"