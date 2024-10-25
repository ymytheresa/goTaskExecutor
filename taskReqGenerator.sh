# #!/bin/bash

# url="http://localhost:8080/task"

# for i in {1..100}
# do
#   request_id="$i"

#   # Add timestamp
#   timestamp=$(date +"%Y-%m-%d %H:%M:%S")
#   echo "Sending request: {\"request_id\": \"$request_id\"} at $timestamp"

#   curl -X POST "$url" \
#     -H "Content-Type: application/json" \
#     -d "{\"request_id\": \"$request_id\"}"

# done

# echo "100 tasks submitted."
#!/bin/bash

# URL to send POST requests to
url="http://localhost:8080/task"

# Total number of requests to send
total_requests=100

# Maximum number of concurrent requests
max_concurrent=10

# Function to send a single POST request
send_request() {
  local request_id=$1
  local url=$2

  # Add timestamp
  timestamp=$(date +"%Y-%m-%d %H:%M:%S")
  echo "Sending request: {\"request_id\": \"$request_id\"} at $timestamp"

  # Send POST request
  curl -X POST "$url" \
    -H "Content-Type: application/json" \
    -d "{\"request_id\": \"$request_id\"}" &
}

# Initialize a counter for concurrent jobs
current_jobs=0

# Loop to send requests
for i in $(seq 1 "$total_requests"); do
  send_request "$i" "$url"
  current_jobs=$((current_jobs + 1))

  # If the maximum number of concurrent jobs is reached, wait for any to finish
  if [ "$current_jobs" -ge "$max_concurrent" ]; then
    wait -n
    current_jobs=$((current_jobs - 1))
  fi
done

# Wait for all remaining background jobs to finish
wait

echo "$total_requests tasks submitted."