#!/bin/bash

source config-file.cfg
source stack-scripts.sh

validate_logs_in_new_relic() {
  user_key=$1
  account_id=$2
  attribute_key=$3
  attribute_value=$4
  log_message=$5

  sleep_time=$SLEEP_TIME
  attempt=1

  while [[ $attempt -lt $MAX_RETRIES ]]; do
    echo "Fetching logs from new relic for $attribute_key: $attribute_value"
    sleep "$sleep_time"
    response=$(fetch_new_relic_logs_api "$user_key" "$account_id" "$attribute_key" "$attribute_value")

    if echo "$response" | grep -q "$log_message"; then
      echo "Log event successfully found in New Relic."
      return 0
    fi

    if (( sleep_time < MAX_SLEEP_TIME )); then
      sleep_time=$(( sleep_time * 2 ))
    fi
    echo "Log event not found in New Relic. Retrying in $sleep_time seconds..."
    attempt=$((attempt + 1))
  done

  exit_with_error "Log event with $attribute_key: $attribute_value not found in New Relic. Error Received: $response"
}

fetch_new_relic_logs_api() {
  user_key=$1
  account_id=$2
  attribute_key=$3
  attribute_value=$4

  nrql_query="SELECT * FROM Log WHERE $attribute_key LIKE '%$attribute_value%' SINCE $TIME_RANGE ago"
  query='{"query":"query($id: Int!, $nrql: Nrql!) { actor { account(id: $id) { nrql(query: $nrql) { results } } } }","variables":{"id":'$account_id',"nrql":"'$nrql_query'"}}'

  response=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -H "API-Key: $user_key" \
    -d "$query" \
    https://api.newrelic.com/graphql)

  echo "$response"
}

create_log_message() {
  log_message=$1
  filter_pattern=$2

  UUID=$(uuidgen)
  echo "RequestId: $UUID, message: $log_message, filter: $filter_pattern"
}