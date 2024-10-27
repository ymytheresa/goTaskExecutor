#!/bin/bash

url="http://localhost:8080/task"

total_requests=10

output_file="responses_$(date +"%Y%m%d_%H%M%S").txt"

send_request() {
  local request_id=$1
  local url=$2

  # Capture the timestamp when the request is sent
  send_timestamp=$(date +"%Y-%m-%d %H:%M:%S")

  request_payload="{\"request_id\": \"$request_id\"}"

  echo "[$send_timestamp] Sending request: $request_payload" | tee -a "$output_file"

  # Send the request and capture the response along with the HTTP status code
  response=$(curl -s -w "\n%{http_code}" -X POST "$url" \
    -H "Content-Type: application/json" \
    -d "$request_payload")

  # Capture the timestamp when the response is received
  receive_timestamp=$(date +"%Y-%m-%d %H:%M:%S")

  http_code=$(echo "$response" | tail -n1)
  response_body=$(echo "$response" | sed '$d')

  echo "[$receive_timestamp] | HTTP Status: $http_code | Request ID: $request_id | Response : $response_body" | tee -a "$output_file"
}

for i in $(seq 1 "$total_requests"); do
  send_request "$i" "$url" &
done

wait

echo -e "\n$(date +"%Y-%m-%d %H:%M:%S") $total_requests tasks submitted.\n" | tee -a "$output_file"

echo "Analyzing results..." | tee -a "$output_file"

completed_tasks=$(grep -o "completed successfully" "$output_file" | wc -l)
failed_tasks=$(grep -o "Task .* failed" "$output_file" | wc -l)
already_completed=$(grep -o "task already completed" "$output_file" | wc -l)
queue_full=$(grep -o "Task .* is not accepted due to queue full" "$output_file" | wc -l)
echo "Summary:" | tee -a "$output_file"
echo "Total queue full tasks: $queue_full" | tee -a "$output_file"
echo "Total completed tasks: $completed_tasks" | tee -a "$output_file"
echo "Total failed tasks: $failed_tasks" | tee -a "$output_file"
echo "Total duplicated tasks: $already_completed" | tee -a "$output_file"
