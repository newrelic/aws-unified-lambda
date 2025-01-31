#!/bin/bash

source common-scripts/stack-scripts.sh
source common-scripts/resource-scripts.sh
source common-scripts/logs-scripts.sh
source common-scripts/config-file.cfg

test_logs_for_prefix() {
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

  deploy_s3_trigger_stack "$TEMPLATES_BUILD_DIR/$LAMBDA_TEMPLATE" "$S3_PREFIX_CASE" "$NEW_RELIC_LICENSE_KEY" "$NEW_RELIC_REGION" "$NEW_RELIC_ACCOUNT_ID" "false" "$S3_BUCKET_NAMES" "$COMMON_ATTRIBUTES"
  validate_stack_deployment_status "$S3_PREFIX_CASE"
  validate_lambda_s3_trigger_created "$S3_PREFIX_CASE" "$S3_BUCKET_NAME" "$S3_BUCKET_PREFIX"

  upload_file_to_s3_bucket "$S3_BUCKET_NAME" "$S3_BUCKET_OBJECT_NAME" "$S3_BUCKET_PREFIX" "$log_message"
  validate_logs_in_new_relic "$NEW_RELIC_USER_KEY" "$NEW_RELIC_ACCOUNT_ID" "$ATTRIBUTE_KEY_S3" "$S3_BUCKET_PREFIX" "$log_message"

  upload_file_to_s3_bucket "$S3_BUCKET_NAME" "$S3_BUCKET_OBJECT_NAME_FOR_INVALID_CASE" "$S3_BUCKET_PREFIX_INVALID" "$log_message"
  validate_logs_not_present "$NEW_RELIC_USER_KEY" "$NEW_RELIC_ACCOUNT_ID" "$ATTRIBUTE_KEY_S3" "$S3_BUCKET_PREFIX_INVALID" "$log_message"

  delete_stack "$S3_PREFIX_CASE"
}

test_logs_for_secret_manager() {
cat <<EOF > s3-parameter.json
'[{"bucket":"$S3_BUCKET_NAME_SECRET_MANAGER","prefix":"$S3_BUCKET_PREFIX"}]'
EOF

cat <<EOF > common_attribute.json
'[{"AttributeName":"$COMMON_ATTRIBUTE_KEY","AttributeValue":"$COMMON_ATTRIBUTE_VALUE"}]'
EOF
COMMON_ATTRIBUTES=$(<common_attribute.json)

  S3_BUCKET_NAMES=$(<s3-parameter.json)
  echo "Testing for s3 bucket configuration JSON: $(<s3-parameter.json)"

  log_message=$(create_log_message "$LOG_MESSAGE_S3" "")

  deploy_s3_trigger_stack "$TEMPLATES_BUILD_DIR/$LAMBDA_TEMPLATE" "$S3_SECRET_MANAGER_CASE" "$NEW_RELIC_LICENSE_KEY" "$NEW_RELIC_REGION" "$NEW_RELIC_ACCOUNT_ID" "true" "$S3_BUCKET_NAMES" "$COMMON_ATTRIBUTES"
  validate_stack_deployment_status "$S3_SECRET_MANAGER_CASE"
  validate_lambda_s3_trigger_created "$S3_SECRET_MANAGER_CASE" "$S3_BUCKET_NAME_SECRET_MANAGER" "$S3_BUCKET_PREFIX"
  upload_file_to_s3_bucket "$S3_BUCKET_NAME_SECRET_MANAGER" "$S3_BUCKET_OBJECT_NAME" "$S3_BUCKET_PREFIX" "$log_message"
  validate_logs_in_new_relic "$NEW_RELIC_USER_KEY" "$NEW_RELIC_ACCOUNT_ID" "$ATTRIBUTE_KEY_S3" "$S3_BUCKET_PREFIX" "$log_message"
  delete_stack "$S3_SECRET_MANAGER_CASE"
}

test_logs_for_invalid_bucket_name() {
cat <<EOF > s3-parameter.json
'[{"bucket":"$S3_BUCKET_NAME_INVALID","prefix":"$S3_BUCKET_PREFIX"}]'
EOF

  S3_BUCKET_NAMES=$(<s3-parameter.json)
  echo "Testing for s3 bucket configuration JSON: $(<s3-parameter.json)"

  deploy_s3_trigger_stack "$TEMPLATES_BUILD_DIR/$LAMBDA_TEMPLATE" "$S3_INVALID_BUCKET_CASE" "$NEW_RELIC_LICENSE_KEY" "$NEW_RELIC_REGION" "$NEW_RELIC_ACCOUNT_ID" "false" "$S3_BUCKET_NAMES" "''"
  validate_stack_deployment_status "$S3_INVALID_BUCKET_CASE"

  # validate that lambda trigger is not created
  validate_lambda_s3_trigger_not_created "$S3_INVALID_BUCKET_CASE" "$S3_BUCKET_NAME_INVALID" "$S3_BUCKET_PREFIX"
  delete_stack "$S3_INVALID_BUCKET_CASE"
}

case $1 in
  test_logs_for_prefix)
    test_logs_for_prefix
    ;;
  test_logs_for_secret_manager)
    test_logs_for_secret_manager
    ;;
  test_logs_for_invalid_bucket_name)
    test_logs_for_invalid_bucket_name
    ;;
  *)
  echo "Invalid test case specified."
  exit 1
  ;;
esac
