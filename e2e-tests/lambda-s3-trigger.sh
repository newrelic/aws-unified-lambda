#!/bin/bash

source common-scripts/stack-scripts.sh
source common-scripts/resource-scripts.sh
source common-scripts/logs-scripts.sh
source common-scripts/config-file.cfg

# test case constants
S3_TRIGGER_CASE=e2e-s3-trigger-stack

deploy_s3_trigger_stack() {
  template_file=$1
  stack_name=$2
  license_key=$3
  new_relic_region=$4
  new_relic_account_id=$5
  secret_license_key=$6
  s3_bucket_names=$7
  common_attributes=$8

  echo "Deploying s3 trigger stack with name: $stack_name"

  sam deploy \
    --template-file "$template_file" \
    --stack-name "$stack_name" \
    --parameter-overrides \
      LicenseKey="$license_key" \
      NewRelicRegion="$new_relic_region" \
      NewRelicAccountId="$new_relic_account_id" \
      StoreNRLicenseKeyInSecretManager="$secret_license_key" \
      S3BucketNames="$s3_bucket_names" \
      LogGroupConfig="''" \
      CommonAttributes="$common_attributes" \
    --capabilities CAPABILITY_IAM
}

test_logs_s3() {
cat <<EOF > s3-parameter.json
'[{"bucket":"$S3_BUCKET_NAME","prefix":"$S3_BUCKET_PREFIX"}]'
EOF
  S3_BUCKET_NAMES=$(<s3-parameter.json)
  echo "Testing for s3 bucket configuration JSON: $(<s3-parameter.json)"

  log_message=$(create_log_message "$LOG_MESSAGE_S3" "N/A")

  deploy_s3_trigger_stack "$LAMBDA_TEMPLATE_BUILD_DIR/$LAMBDA_TEMPLATE" "$S3_TRIGGER_CASE" "$NEW_RELIC_LICENSE_KEY" "$NEW_RELIC_REGION" "$NEW_RELIC_ACCOUNT_ID" "false" "$S3_BUCKET_NAMES" "''"
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
