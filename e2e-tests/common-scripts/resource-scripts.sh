#!/bin/bash

source test-configs.cfg
source stack-scripts.sh

get_lambda_function_arn() {
  stack_name=$1

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

  echo "$lambda_function_arn"
}

create_cloudwatch_log_event() {
  log_group_name=$1
  log_stream_name=$2
  log_message=$3

  echo "Creating log event in CloudWatch Log Group"

  # Check if the log stream exists else create one
  log_stream_exists=$(aws logs describe-log-streams --log-group-name "$log_group_name" --log-stream-name-prefix "$log_stream_name" --query "logStreams[?logStreamName=='$log_stream_name'] | length(@)" --output text)

  if [ -n "$log_stream_exists" ] && [ "$log_stream_exists" -eq 0 ]; then
    echo "Log stream does not exist. Creating log stream: $log_stream_name"
    aws logs create-log-stream --log-group-name "$log_group_name" --log-stream-name "$log_stream_name"
  fi

  aws logs put-log-events \
    --log-group-name "$log_group_name" \
    --log-stream-name "$log_stream_name" \
    --log-events "timestamp=$(date +%s000),message=\"$log_message\""

  echo "Log event with message: $log_message created successfully."
}

upload_file_to_s3_bucket() {
  bucket_name=$1
  log_file=$2
  prefix=$3
  log_message=$4

  echo "$log_message" >> "$log_file"
  if [[ $? -ne 0 ]]; then
      echo "Failed to write log to file."
      return 1
  fi

  aws s3 cp "$log_file" "s3://$bucket_name/$prefix"
  if [[ $? -ne 0 ]]; then
      echo "Failed to upload log file to S3."
      return 1
  fi

  echo "Log successfully uploaded as s3://$bucket_name/$prefix"
}

validate_lambda_subscription_created() {
  stack_name=$1
  log_group_name=$2
  log_group_filter=$3

  echo "Validating cloudwatch lambda subscription for stack name: $stack_name, log group name: $log_group_name, and log group filter: $log_group_filter"

  lambda_function_arn=$(get_lambda_function_arn "$stack_name")

  subscriptions=$(aws logs describe-subscription-filters --log-group-name "$log_group_name" --query 'subscriptionFilters[*].[destinationArn, filterPattern]' --output text)

  if echo "$subscriptions" | grep -q "$lambda_function_arn" && echo "$subscriptions" | grep -q "$log_group_filter"; then
    echo "Lambda function $lambda_function_arn is subscribed to log group: $log_group_name with filter: $log_group_filter"
  else
    exit_with_error "Lambda function $lambda_function_arn is not subscribed to log group: $log_group_name"
  fi

}

validate_lambda_subscription_not_created() {
  stack_name=$1
  log_group_name=$2
  log_group_filter=$3

  echo "Validating cloudwatch lambda subscription for stack name: $stack_name, log group name: $log_group_name, and log group filter: $log_group_filter"

  lambda_function_arn=$(get_lambda_function_arn "$stack_name")

  event_source_mappings=$(aws lambda list-event-source-mappings --function-name lambda_function_arn --query "EventSourceMappings" --output text)

  # Check if the output is empty
  if [ -z "$event_source_mappings" ]; then
      echo "No event source mappings found for the Lambda function $lambda_function_arn. Validation successful"
  else
      exit_with_error "Event source mappings exist for the Lambda function '$lambda_function_arn'. Validation failed"
  fi
}

validate_lambda_s3_trigger_created() {
  stack_name=$1
  bucket_name=$2
  bucket_prefix=$3

  echo "Validating s3 lambda trigger event for stack name: $stack_name, bucket name: $bucket_name, and prefix: $bucket_prefix"

  lambda_function_arn=$(get_lambda_function_arn "$stack_name")

  notification_configuration=$(aws s3api get-bucket-notification-configuration --bucket "$bucket_name")

  lambda_configurations=$(echo "$notification_configuration" |
      jq --arg lambda "$lambda_function_arn" --arg prefix "$bucket_prefix" '
      .LambdaFunctionConfigurations[] |
      select(.LambdaFunctionArn == $lambda and (.Filter.Key.FilterRules[]? | select(.Name == "Prefix" and .Value == $prefix)?))')

  if [ -n "$lambda_configurations" ]; then
      echo "S3 triggers with prefix '$bucket_prefix' found for Lambda function $lambda_function_arn on bucket $bucket_name:"
      echo "$lambda_configurations" | jq '.'
  else
      exit_with_error "No S3 triggers with prefix '$bucket_prefix' found for Lambda function $lambda_function_arn on bucket: $bucket_name."
  fi
}

validate_lambda_s3_trigger_not_created() {
  stack_name=$1
  bucket_name=$2
  bucket_prefix=$3

  echo "Validating s3 lambda trigger event for stack name: $stack_name, bucket name: $bucket_name, and prefix: $bucket_prefix"

  lambda_function_arn=$(get_lambda_function_arn "$stack_name")

  event_source_mappings=$(aws lambda list-event-source-mappings --function-name lambda_function_arn --query "EventSourceMappings" --output text)

  # Check if the output is empty
  if [ -z "$event_source_mappings" ]; then
      echo "No event source mappings found for the Lambda function $lambda_function_arn. Validation successful"
  else
      exit_with_error "Event source mappings exist for the Lambda function '$lambda_function_arn'. Validation failed"
  fi
}