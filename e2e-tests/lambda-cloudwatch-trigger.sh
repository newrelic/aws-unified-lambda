#!/bin/bash

source common-scripts.sh
source config-file.cfg

# test case constants
CLOUDWATCH_TRIGGER_CASE=e2e-cloudwatch-trigger-stack

deploy_cloudwatch_trigger_stack() {
  template_file=$1
  stack_name=$2
  license_key=$3
  new_relic_region=$4
  new_relic_account_id=$5
  secret_license_key=$6
  log_group_config=$7
  common_attributes=$8

  sam deploy \
    --template-file "$template_file" \
    --stack-name "$stack_name" \
    --parameter-overrides \
      LicenseKey="$license_key" \
      NewRelicRegion="$new_relic_region" \
      NewRelicAccountId="$new_relic_account_id" \
      StoreNRLicenseKeyInSecretManager="$secret_license_key" \
      S3BucketNames="''" \
      LogGroupConfig="$log_group_config" \
      CommonAttributes="$common_attributes" \
    --capabilities CAPABILITY_IAM
}

validate_lambda_subscription_created() {
  stack_name=$1
  log_group_name=$2
  log_group_filter=$3

  lambda_function_arn=$(./common-scripts.sh get_lambda_function_arn "$stack_name")

  subscriptions=$(aws logs describe-subscription-filters --log-group-name "$log_group_name" --query 'subscriptionFilters[*].[destinationArn, filterPattern]' --output text)

  if echo "$subscriptions" | grep -q "$lambda_function_arn" && echo "$subscriptions" | grep -q "$log_group_filter"; then
    echo "Lambda function $lambda_function_arn is subscribed to log group: $log_group_name with filter: $log_group_filter"
  else
    exit_with_error "Lambda function $lambda_function_arn is not subscribed to log group: $log_group_name"
  fi

}

cat <<EOF > cloudwatch-parameter.json
'[{"LogGroupName":"$LOG_GROUP_NAME","FilterPattern":"$LOG_GROUP_FILTER_PATTERN"}]'
EOF
LOG_GROUP_NAMES=$(<cloudwatch-parameter.json)

deploy_cloudwatch_trigger_stack "$LAMBDA_TEMPLATE_BUILD_DIR/$LAMBDA_TEMPLATE" "$CLOUDWATCH_TRIGGER_CASE" "$NEW_RELIC_LICENSE_KEY" "$NEW_RELIC_REGION" "$NEW_RELIC_ACCOUNT_ID" "false"  "$LOG_GROUP_NAMES" "''"
validate_stack_deployment_status "$CLOUDWATCH_TRIGGER_CASE"
validate_lambda_subscription_created "$CLOUDWATCH_TRIGGER_CASE" "$LOG_GROUP_NAME" "$LOG_GROUP_FILTER_PATTERN"
#delete_stack "$CLOUDWATCH_TRIGGER_CASE"