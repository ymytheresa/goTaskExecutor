#!/bin/bash

url="http://localhost:8080/task"

for i in {1..100}
do
  request_id="$i"

  # Add timestamp
  timestamp=$(date +"%Y-%m-%d %H:%M:%S")
  echo "Sending request: {\"request_id\": \"$request_id\"} at $timestamp"

  curl -X POST "$url" \
    -H "Content-Type: application/json" \
    -d "{\"request_id\": \"$request_id\"}"

  sleep 0.1
done

echo "100 tasks submitted."
