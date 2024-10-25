#!/bin/bash

# URL of the API endpoint
url="http://localhost:8080/task"  # Corrected URL

# Loop to send 100 tasks
for i in {1..100}
do
  # Generate a random request ID or use the loop index
  request_id="task-$i"

  # Output the request being sent
  echo "Sending request: {\"request_id\": \"$request_id\"}"

  # Print the current integer value of i
  echo "Current integer: $i"

  # Print the URL being used
  echo "Using URL: $url"

  # Send the POST request with a JSON payload
  curl -X POST "$url" \
    -H "Content-Type: application/json" \
    -d "{\"request_id\": \"$request_id\"}"

  # Optional: Add a delay if needed between requests
  sleep 0.1
done

echo "100 tasks submitted."