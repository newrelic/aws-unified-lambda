#!/bin/bash

source common-scripts/stack-scripts.sh
source common-scripts/resource-scripts.sh
source common-scripts/logs-scripts.sh
source common-scripts/config-file.cfg

test_logs_s3() {
cat <<EOF > s3-parameter.json
'[{"bucket":"$S3_BUCKET_NAME","prefix":"$S3_BUCKET_PREFIX"}]'
EOF

cat <<EOF > common_attribute.json
'[{"AttributeName":"$COMMON_ATTRIBUTE_KEY","AttributeValue":"$COMMON_ATTRIBUTE_VALUE"}]'
EOF
COMMON_ATTRIBUTES=$(<common_attribute.json)

  S3_BUCKET_NAMES=$(<s3-parameter.json)
  echo "Testing for s3 bucket configuration JSON: $(<s3-parameter.json)"

  log_message=$(create_log_message "$LOG_MESSAGE_S3" "")

  deploy_s3_trigger_stack "$LAMBDA_TEMPLATE_BUILD_DIR/$LAMBDA_TEMPLATE" "$S3_TRIGGER_CASE" "$NEW_RELIC_LICENSE_KEY" "$NEW_RELIC_REGION" "$NEW_RELIC_ACCOUNT_ID" "false" "$S3_BUCKET_NAMES" "$COMMON_ATTRIBUTES"
  validate_stack_deployment_status "$S3_TRIGGER_CASE"
  validate_lambda_s3_trigger_created "$S3_TRIGGER_CASE" "$S3_BUCKET_NAME" "$S3_BUCKET_PREFIX"
  upload_file_to_s3_bucket "$S3_BUCKET_NAME" "$S3_BUCKET_OBJECT_NAME" "$S3_BUCKET_PREFIX" "$log_message"
  validate_logs_in_new_relic "$NEW_RELIC_USER_KEY" "$NEW_RELIC_ACCOUNT_ID" "$ATTRIBUTE_KEY_S3" "$S3_BUCKET_PREFIX" "$log_message"
  delete_stack "$S3_TRIGGER_CASE"
}

case $1 in
  test_logs_s3)
    test_logs_s3
    ;;
  *)
  echo "Invalid test case specified."
  exit 1
  ;;
esac
