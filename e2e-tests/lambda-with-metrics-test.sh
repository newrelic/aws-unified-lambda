#!/bin/bash

source common-scripts/stack-scripts.sh
source common-scripts/resource-scripts.sh
source common-scripts/logs-scripts.sh
source common-scripts/config-file.cfg

test_for_lambda_firehose_stack() {
cat <<EOF > cloudwatch-parameter.json
'[{"LogGroupName":"$LOG_GROUP_NAME","FilterPattern":"$LOG_GROUP_FILTER_PATTERN"}]'
EOF

cat <<EOF > s3-parameter.json
'[{"bucket":"$S3_BUCKET_NAME","prefix":"$S3_BUCKET_PREFIX"}]'
EOF

  LOG_GROUP_NAMES=$(<cloudwatch-parameter.json)
  S3_BUCKET_NAMES=$(<s3-parameter.json)

  deploy_lambda_firehose_stack "$S3_BUCKET_NAMES" "$LOG_GROUP_NAMES" "$COMMON_ATTRIBUTES" "false" "false"
  validate_stack_deployment_status "$LAMBDA_FIREHOSE_CASE"
  delete_stack "$LAMBDA_FIREHOSE_CASE"
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