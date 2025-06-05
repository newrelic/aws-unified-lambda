#!/bin/bash

source common-scripts/stack-scripts.sh
source common-scripts/resource-scripts.sh
source common-scripts/logs-scripts.sh
source common-scripts/test-configs.cfg

test_logs_for_prefix() {
cat <<EOF > s3-parameter.json
'[{"bucket":"$S3_BUCKET_NAME","prefix":"$S3_BUCKET_PREFIX"}]'
EOF

cat <<EOF > common_attribute.json
'[{"AttributeName":"$CUSTOM_ATTRIBUTE_KEY","AttributeValue":"$CUSTOM_ATTRIBUTE_VALUE"}]'
EOF
COMMON_ATTRIBUTES=$(<common_attribute.json)

  S3_BUCKET_NAMES=$(<s3-parameter.json)
  echo "Testing for s3 bucket configuration JSON: $(<s3-parameter.json)"

  log_message=$(create_log_message "$LOG_MESSAGE_S3" "")

  
  delete_stack "$S3_PREFIX_CASE"
}

test_logs_for_secret_manager() {
cat <<EOF > s3-parameter.json
'[{"bucket":"$S3_BUCKET_NAME_SECRET_MANAGER","prefix":"$S3_BUCKET_PREFIX"}]'
EOF

cat <<EOF > common_attribute.json
'[{"AttributeName":"$CUSTOM_ATTRIBUTE_KEY","AttributeValue":"$CUSTOM_ATTRIBUTE_VALUE"}]'
EOF
COMMON_ATTRIBUTES=$(<common_attribute.json)

  S3_BUCKET_NAMES=$(<s3-parameter.json)
  echo "Testing for s3 bucket configuration JSON: $(<s3-parameter.json)"

  log_message=$(create_log_message "$LOG_MESSAGE_S3" "")

  delete_stack "$S3_SECRET_MANAGER_CASE"
}

test_logs_for_invalid_bucket_name() {
cat <<EOF > s3-parameter.json
'[{"bucket":"$S3_BUCKET_NAME_INVALID","prefix":"$S3_BUCKET_PREFIX"}]'
EOF

  S3_BUCKET_NAMES=$(<s3-parameter.json)
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
