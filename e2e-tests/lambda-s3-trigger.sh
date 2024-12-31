#!/bin/bash

source common-scripts.sh
source config-file.cfg

# test case constants
S3_TRIGGER_CASE=s3-trigger-stack-1

validate_stack_resources() {
  stack_name=$1
  bucket_name=$2
  bucket_prefix=$3

  lambda_physical_id=$(aws cloudformation describe-stack-resources \
                  --stack-name "$stack_name" \
                  --logical-resource-id "$LAMBDA_LOGICAL_RESOURCE_ID" \
                  --query "StackResources[0].PhysicalResourceId" \
                  --output text
  )
  lambda_function_arn=$(aws lambda get-function --function-name "$lambda_physical_id" \
                  --query "Configuration.FunctionArn" \
                  --output text
  )

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

cat <<EOF > parameter.json
'[{"bucket":"$S3_BUCKET_NAME","prefix":"$S3_BUCKET_PREFIX"}]'
EOF
S3_BUCKET_NAMES=$(<parameter.json)

deploy_stack "$LAMBDA_TEMPLATE_BUILD_DIR/$LAMBDA_TEMPLATE" "$S3_TRIGGER_CASE" "$NEW_RELIC_LICENSE_KEY" "$NEW_RELIC_REGION" "$NEW_RELIC_ACCOUNT_ID" "false" "$S3_BUCKET_NAMES" "''" "''"
validate_stack_deployment_status "$S3_TRIGGER_CASE"
validate_stack_resources "$S3_TRIGGER_CASE" "$S3_BUCKET_NAME" "$S3_BUCKET_PREFIX"
delete_stack "$S3_TRIGGER_CASE"