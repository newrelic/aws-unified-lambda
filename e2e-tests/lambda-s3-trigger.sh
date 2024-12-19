#!/bin/bash

source ./common-stack-scripts.sh

exit_with_error() {
  echo "Error: $1"
  exit 1
}

deploy_stack() {
  template_file=$1
  stack_name=$2
  license_key=$3
  new_relic_region=$4
  new_relic_account_id=$5
  secret_license_key=$6
  s3_bucket_names=$7
  log_group_config=$8
  common_attributes=$9

  sam deploy \
    --template-file "$template_file" \
    --stack-name "$stack_name" \
    --parameter-overrides \
      LicenseKey="$license_key" \
      NewRelicRegion="$new_relic_region" \
      NewRelicAccountId="$new_relic_account_id" \
      StoreNRLicenseKeyInSecretManager="$secret_license_key" \
      S3BucketNames="$s3_bucket_names" \
      LogGroupConfig="$log_group_config" \
      CommonAttributes="$common_attributes" \
    --capabilities CAPABILITY_IAM
}

validate_stack_resources() {
  stack_name=$1
  bucket_name=$2
  bucket_prefix=$3

  resources=$(aws cloudformation describe-stack-resources --stack-name "$stack_name")

  lambda_function_arn=$(echo "$resources" | \
      jq -r '.StackResources[] | select(.ResourceType == "AWS::Lambda::Function") | .PhysicalResourceId' | \
      grep "$LAMBDA_NAME" | \
      xargs -I {} aws lambda get-function --function-name {} --query 'Configuration.FunctionArn' --output text)

  notification_configuration=$(aws s3api get-bucket-notification-configuration --bucket "$bucket_name")

  lambda_configurations=$(echo "$notification_configuration" |
      jq --arg lambda "$lambda_function_arn" --arg prefix "$bucket_prefix" '
      .LambdaFunctionConfigurations[] |
      select(.LambdaFunctionArn == $lambda and (.Filter.Key.FilterRules[]? | select(.Name == "Prefix" and .Value == $prefix)?))')

  if [ -n "$lambda_configurations" ]; then
      echo "S3 triggers with prefix '$bucket_prefix' found for Lambda function $lambda_function_arn on bucket $bucket_name:"
      echo "$lambda_configurations" | jq '.'
  else
      exit_with_error "No S3 triggers with prefix '$bucket_prefix' found for Lambda function $lambda_function_arn on bucket $bucket_name."
  fi
}

source config-file.cfg

cat <<EOF > parameter.json
'[{"bucket":"$S3_BUCKET_TRIGGER","prefix":"$S3_BUCKET_PREFIX"}]'
EOF
S3_BUCKET_NAMES=$(<parameter.json)

deploy_stack "$LAMBDA_TEMPLATE_BUILD_DIR/$LAMBDA_TEMPLATE" "e2eTests" "$NEW_RELIC_LICENSE_KEY" "$NEW_RELIC_REGION" "$NEW_RELIC_ACCOUNT_ID" "false" "$S3_BUCKET_NAMES" "''" "''"
validate_stack_deployment_status "e2eTests"
validate_stack_resources "e2eTests" "$S3_BUCKET_TRIGGER" "$S3_BUCKET_PREFIX"
delete_stack "e2eTests"