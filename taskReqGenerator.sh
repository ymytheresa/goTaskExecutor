#!/bin/bash

url="http://localhost:8080/task"

for i in {1..100}
do
  request_id="$i"

  echo "Sending request: {\"request_id\": \"$request_id\"}"

  curl -X POST "$url" \
    -H "Content-Type: application/json" \
    -d "{\"request_id\": \"$request_id\"}"

  sleep 0.1
done

echo "100 tasks submitted."