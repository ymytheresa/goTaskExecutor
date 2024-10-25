#!/bin/bash

url="http://localhost:8080/task"

total_requests=100

output_file="responses_$(date +"%Y%m%d_%H%M%S").txt"

send_request() {
  local request_id=$1
  local url=$2

  timestamp=$(date +"%Y-%m-%d %H:%M:%S")

  request_payload="{\"request_id\": \"$request_id\"}"

  echo "[$timestamp] Sending request: $request_payload" | tee -a "$output_file"

  response=$(curl -s -w "\n%{http_code}" -X POST "$url" \
    -H "Content-Type: application/json" \
    -d "$request_payload")

    http_code=$(echo "$response" | tail -n1)
    response_body=$(echo "$response" | sed '$d')

    echo "[$timestamp] | HTTP Status: $http_code | Request ID: $request_id | Response : $response_body" | tee -a "$output_file"
    
}

for i in $(seq 1 "$total_requests"); do
  send_request "$i" "$url" &
done

wait

echo "$(date +"%Y-%m-%d %H:%M:%S") $total_requests tasks submitted." | tee -a "$output_file"
