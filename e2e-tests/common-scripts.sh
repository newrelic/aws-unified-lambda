#!/bin/bash

source config-file.cfg

validate_stack_deployment_status() {
  stack_name=$1

  stack_status=$(aws cloudformation describe-stacks --stack-name "$stack_name" --query "Stacks[0].StackStatus" --output text)
  if [[ "$stack_status" == "ROLLBACK_COMPLETE" || "$stack_status" == "ROLLBACK_FAILED" || "$stack_status" == "CREATE_FAILED"  || "$stack_status" == "UPDATE_FAILED" ]]; then
    echo "Stack $stack_name failed to be created and rolled back."
    failure_reason=$(aws cloudformation describe-stack-events --stack-name "$stack_name" --query "StackEvents[?ResourceStatus==\`$stack_status\`].ResourceStatusReason" --output text)
    exit_with_error "Stack $stack_name failed to be created. Failure reason: $failure_reason"
  else
    echo "Stack $stack_name was created successfully."
  fi
}

delete_stack() {
  stack_name=$1

  aws cloudformation delete-stack --stack-name "$stack_name"

  echo "Initiated deletion of stack: $stack_name"

  stack_status=$(aws cloudformation describe-stacks --stack-name "$stack_name" --query 'Stacks[0].StackStatus' --output text)

  # delete stack with exponential back off retires
  max_sleep_time=300  # Cap sleep time to 5 minutes
  sleep_time=30

  while [[ $stack_status == "DELETE_IN_PROGRESS" ]]; do
    echo "Stack $stack_name is still being deleted..."

    sleep $sleep_time
    if (( sleep_time < max_sleep_time )); then
      sleep_time=$(( sleep_time * 2 ))
    fi

    stack_status=$(aws cloudformation describe-stacks --stack-name "$stack_name" --query 'Stacks[0].StackStatus' --output text 2>/dev/null || true)
  done

  if [ -z "$stack_status" ]; then
    echo "Stack $stack_name has been successfully deleted."
  elif [ "$stack_status" == "DELETE_FAILED" ]; then
    echo "Failed to delete stack $stack_name."
  else
    echo "Unexpected stack status: $stack_status."
  fi
}

exit_with_error() {
  echo "Error: $1"
  exit 1
}

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

  timestamp=$(($(date +%s) * 1000 + $(date +%N) / 1000000))

  aws logs put-log-events \
    --log-group-name "$log_group_name" \
    --log-stream-name "$log_stream_name" \
    --log-events timestamp=$timestamp,message="$log_message"

  echo "Log event with message: $log_message created successfully."
}

validate_logs_in_new_relic() {
  user_key=$1
  account_id=$2
  stream_name=$3
  log_message=$4

  sleep_time=$SLEEP_TIME
  i=1

  while i < "$MAX_RETRIES"; do
    echo "Fetching logs from new relic for stream name: $stream_name"

    response=$(fetch_new_relic_logs_api "$stack_name")

    if echo "$response" | grep -q "$log_message"; then
      echo "Log event successfully found in New Relic."
      return 0
    else
      echo "Log event not found in New Relic. Retrying in $sleep_time seconds..."
      sleep "$sleep_time"
      if (( sleep_time < MAX_SLEEP_TIME )); then
        sleep_time=$(( sleep_time * 2 ))
      fi
    fi

    i=$((i + 1))
  done
  exit_with_error "Log event with stream name: $stream_name not found in New Relic."
}

fetch_new_relic_logs_api() {
  user_key=$1
  account_id=$2
  stream_name=$3

  nrql_query="SELECT * FROM Log WHERE $ATTRIBUTE_KEY LIKE '%$stream_name%' SINCE $TIME_RANGE ago"
  query='{"query":"query($id: Int!, $nrql: Nrql!) { actor { account(id: $id) { nrql(query: $nrql) { results } } } }","variables":{"id":'$account_id',"nrql":"'$nrql_query'"}}'

  response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -H "API-Key: $user_key" \
    -d "$query" \
    https://api.newrelic.com/graphql)

  echo "$response"
}