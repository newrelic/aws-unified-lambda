#!/bin/bash

source common-scripts/stack-scripts.sh
source common-scripts/resource-scripts.sh
source common-scripts/logs-scripts.sh
source common-scripts/config-file.cfg

test_for_lambda_firehose_stack() {
cat <<EOF > cloudwatch-parameter.json
'[{"LogGroupName":"$LOG_GROUP_NAME_METRICS_CASE","FilterPattern":"$LOG_GROUP_FILTER_PATTERN"}]'
EOF

cat <<EOF > s3-parameter.json
'[{"bucket":"$S3_BUCKET_NAME_METRICS_CASE","prefix":"$S3_BUCKET_PREFIX"}]'
EOF

cat <<EOF > common_attribute.json
'[{"AttributeName":"$COMMON_ATTRIBUTE_KEY","AttributeValue":"$COMMON_ATTRIBUTE_VALUE"}]'
EOF

  COMMON_ATTRIBUTES=$(<common_attribute.json)
  LOG_GROUP_NAMES=$(<cloudwatch-parameter.json)
  S3_BUCKET_NAMES=$(<s3-parameter.json)

  deploy_lambda_firehose_stack "$S3_BUCKET_NAMES" "$LOG_GROUP_NAMES" "$COMMON_ATTRIBUTES" "false" "false"
  validate_stack_deployment_status "$LAMBDA_FIREHOSE_CASE"
  delete_stack "$LAMBDA_FIREHOSE_CASE"
}

test_for_lambda_metrics_polling_stack() {
cat <<EOF > cloudwatch-parameter.json
'[{"LogGroupName":"$LOG_GROUP_NAME_METRICS_CASE","FilterPattern":"$LOG_GROUP_FILTER_PATTERN"}]'
EOF

cat <<EOF > s3-parameter.json
'[{"bucket":"$S3_BUCKET_NAME_METRICS_CASE","prefix":"$S3_BUCKET_PREFIX"}]'
EOF

cat <<EOF > common_attribute.json
'[{"AttributeName":"$COMMON_ATTRIBUTE_KEY","AttributeValue":"$COMMON_ATTRIBUTE_VALUE"}]'
EOF

  COMMON_ATTRIBUTES=$(<common_attribute.json)
  LOG_GROUP_NAMES=$(<cloudwatch-parameter.json)
  S3_BUCKET_NAMES=$(<s3-parameter.json)

  deploy_lambda_metric_polling_stack "$LOG_GROUP_NAMES" "$COMMON_ATTRIBUTES" "false" "$S3_BUCKET_NAMES"
  validate_stack_deployment_status "$LAMBDA_METRIC_POLLING_CASE"
  delete_stack "$LAMBDA_METRIC_POLLING_CASE"
}

test_for_lambda_metrics_streaming_stack() {
cat <<EOF > cloudwatch-parameter.json
'[{"LogGroupName":"$LOG_GROUP_NAME_METRICS_CASE","FilterPattern":"$LOG_GROUP_FILTER_PATTERN"}]'
EOF

cat <<EOF > s3-parameter.json
'[{"bucket":"$S3_BUCKET_NAME_METRICS_CASE","prefix":"$S3_BUCKET_PREFIX"}]'
EOF

cat <<EOF > common_attribute.json
'[{"AttributeName":"$COMMON_ATTRIBUTE_KEY","AttributeValue":"$COMMON_ATTRIBUTE_VALUE"}]'
EOF

  COMMON_ATTRIBUTES=$(<common_attribute.json)
  LOG_GROUP_NAMES=$(<cloudwatch-parameter.json)
  S3_BUCKET_NAMES=$(<s3-parameter.json)

  deploy_lambda_metric_streaming_stack "$LOG_GROUP_NAMES" "$COMMON_ATTRIBUTES" "false" "$S3_BUCKET_NAMES"
  validate_stack_deployment_status "$LAMBDA_METRIC_STREAMING_CASE"
  delete_stack "$LAMBDA_METRIC_STREAMING_CASE"
}

test_for_firehose_metric_polling_stack() {
cat <<EOF > cloudwatch-parameter.json
'[{"LogGroupName":"$LOG_GROUP_NAME_METRICS_CASE","FilterPattern":"$LOG_GROUP_FILTER_PATTERN"}]'
EOF

cat <<EOF > common_attribute.json
'[{"AttributeName":"$COMMON_ATTRIBUTE_KEY","AttributeValue":"$COMMON_ATTRIBUTE_VALUE"}]'
EOF

  COMMON_ATTRIBUTES=$(<common_attribute.json)
  LOG_GROUP_NAMES=$(<cloudwatch-parameter.json)

  deploy_firehose_metric_polling_stack "$LOG_GROUP_NAMES" "$COMMON_ATTRIBUTES" "false"
  validate_stack_deployment_status "$FIREHOSE_METRIC_POLLING_CASE"
  delete_stack "$FIREHOSE_METRIC_POLLING_CASE"
}

test_for_firehose_metric_streaming_stack() {
cat <<EOF > cloudwatch-parameter.json
'[{"LogGroupName":"$LOG_GROUP_NAME_METRICS_CASE","FilterPattern":"$LOG_GROUP_FILTER_PATTERN"}]'
EOF

cat <<EOF > common_attribute.json
'[{"AttributeName":"$COMMON_ATTRIBUTE_KEY","AttributeValue":"$COMMON_ATTRIBUTE_VALUE"}]'
EOF

  COMMON_ATTRIBUTES=$(<common_attribute.json)
  LOG_GROUP_NAMES=$(<cloudwatch-parameter.json)

  deploy_firehose_metric_streaming_stack "$LOG_GROUP_NAMES" "$COMMON_ATTRIBUTES" "false"
  validate_stack_deployment_status "$FIREHOSE_METRIC_STREAMING_CASE"
  delete_stack "$FIREHOSE_METRIC_STREAMING_CASE"
}


test_for_lambda_firehose_metric_polling_stack() {
cat <<EOF > cloudwatch-parameter.json
'[{"LogGroupName":"$LOG_GROUP_NAME_METRICS_CASE","FilterPattern":"$LOG_GROUP_FILTER_PATTERN"}]'
EOF

cat <<EOF > s3-parameter.json
'[{"bucket":"$S3_BUCKET_NAME_METRICS_CASE","prefix":"$S3_BUCKET_PREFIX"}]'
EOF

cat <<EOF > common_attribute.json
'[{"AttributeName":"$COMMON_ATTRIBUTE_KEY","AttributeValue":"$COMMON_ATTRIBUTE_VALUE"}]'
EOF

  COMMON_ATTRIBUTES=$(<common_attribute.json)
  LOG_GROUP_NAMES=$(<cloudwatch-parameter.json)
  S3_BUCKET_NAMES=$(<s3-parameter.json)

  deploy_lambda_firehose_metric_polling_stack "$LOG_GROUP_NAMES" "$COMMON_ATTRIBUTES" "false" "$S3_BUCKET_NAMES"
  validate_stack_deployment_status "$LAMBDA_FIREHOSE_METRIC_POLLING_CASE"
  delete_stack "$LAMBDA_FIREHOSE_METRIC_POLLING_CASE"
}

test_for_lambda_firehose_metric_streaming_stack() {
cat <<EOF > cloudwatch-parameter.json
'[{"LogGroupName":"$LOG_GROUP_NAME_METRICS_CASE","FilterPattern":"$LOG_GROUP_FILTER_PATTERN"}]'
EOF

cat <<EOF > s3-parameter.json
'[{"bucket":"$S3_BUCKET_NAME_METRICS_CASE","prefix":"$S3_BUCKET_PREFIX"}]'
EOF

cat <<EOF > common_attribute.json
'[{"AttributeName":"$COMMON_ATTRIBUTE_KEY","AttributeValue":"$COMMON_ATTRIBUTE_VALUE"}]'
EOF

  COMMON_ATTRIBUTES=$(<common_attribute.json)
  LOG_GROUP_NAMES=$(<cloudwatch-parameter.json)
  S3_BUCKET_NAMES=$(<s3-parameter.json)

  deploy_lambda_firehose_metric_streaming_stack "$LOG_GROUP_NAMES" "$COMMON_ATTRIBUTES" "false" "$S3_BUCKET_NAMES"
  validate_stack_deployment_status "$LAMBDA_FIREHOSE_METRIC_STREAMING_CASE"
  delete_stack "$LAMBDA_FIREHOSE_METRIC_STREAMING_CASE"
}

case $1 in
  test_for_lambda_firehose_stack)
    test_for_lambda_firehose_stack
    ;;
  test_for_lambda_metrics_polling_stack)
    test_for_lambda_metrics_polling_stack
    ;;
  test_for_lambda_metrics_streaming_stack)
    test_for_lambda_metrics_streaming_stack
    ;;
  test_for_firehose_metric_polling_stack)
    test_for_firehose_metric_polling_stack
    ;;
  test_for_firehose_metric_streaming_stack)
    test_for_firehose_metric_streaming_stack
    ;;
  test_for_lambda_firehose_metric_polling_stack)
    test_for_lambda_firehose_metric_polling_stack
    ;;
  test_for_lambda_firehose_metric_streaming_stack)
    test_for_lambda_firehose_metric_streaming_stack
    ;;
  *)
  echo "Invalid test case specified."
  exit 1
  ;;
esac