#!/bin/bash

source config-file.cfg

deploy_cloudwatch_trigger_stack() {
  template_file=$1
  stack_name=$2
  license_key=$3
  new_relic_region=$4
  new_relic_account_id=$5
  secret_license_key=$6
  log_group_config=$7
  common_attributes=$8

  echo "Deploying cloudwatch trigger stack with name: $stack_name"

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

validate_stack_deployment_status() {
  stack_name=$1

  echo "Validating stack deployment status for stack name: $stack_name"

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