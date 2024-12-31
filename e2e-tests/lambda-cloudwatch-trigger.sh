#!/bin/bash

source common-scripts.sh
source config-file.cfg

# test case constants
CLOUDWATCH_TRIGGER_CASE=cloudwatch-trigger-stack-1

validate_stack_resources() {
  stack_name=$1
  log_group_name=$2
  log_group_filter=$3

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

  subscriptions=$(aws logs describe-subscription-filters --log-group-name "$log_group_name" --query 'subscriptionFilters[*].[destinationArn, filterPattern]' --output text)

  if echo "$subscriptions" | grep -q "$lambda_function_arn" && echo "$subscriptions" | grep -q "$log_group_filter"; then
    echo "Lambda function $lambda_function_arn is subscribed to log group: $log_group_name with filter: $log_group_filter"
  else
    exit_with_error "Lambda function $lambda_function_arn is not subscribed to log group: $log_group_name"
  fi

}

cat <<EOF > parameter.json
'[{"LogGroupName":"$LOG_GROUP_NAME","FilterPattern":"$LOG_GROUP_FILTER_PATTERN"}]'
EOF
LOG_GROUP_NAMES=$(<parameter.json)

deploy_stack "$LAMBDA_TEMPLATE_BUILD_DIR/$LAMBDA_TEMPLATE" "$CLOUDWATCH_TRIGGER_CASE" "$NEW_RELIC_LICENSE_KEY" "$NEW_RELIC_REGION" "$NEW_RELIC_ACCOUNT_ID" "false" "''" "$LOG_GROUP_NAMES" "''"
validate_stack_deployment_status "$CLOUDWATCH_TRIGGER_CASE"
validate_stack_resources "$CLOUDWATCH_TRIGGER_CASE" "$LOG_GROUP_NAME" "$LOG_GROUP_FILTER_PATTERN"
delete_stack "$CLOUDWATCH_TRIGGER_CASE"