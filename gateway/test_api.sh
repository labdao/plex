#!/bin/bash

# Function to perform a GET request with and without the API key
test_endpoint() {
    local endpoint=$1
    local api_key=$2
    local status_without_key status_with_key response_without_key response_with_key

    # Request without the API key
    status_without_key=$(curl -o /dev/null -s -w "%{http_code}\n" http://localhost:8080/${endpoint})
    response_without_key=$(format_response "$status_without_key")

    # Request with the API key
    status_with_key=$(curl -o /dev/null -s -w "%{http_code}\n" -H "Authorization: Bearer ${api_key}" http://localhost:8080/${endpoint})
    response_with_key=$(format_response "$status_with_key")

    # Print the endpoint and the responses without showing the API key
    printf "%-20s %-20s %-20s\n" "GET /$endpoint" "$response_without_key" "$response_with_key"
}

# Function to format the response
format_response() {
    case $1 in
        200) echo "Success" ;;
        401) echo "Unauthorized" ;;
        *) echo "$1" ;; # Default case to just return the status code
    esac
}

# Health check
echo "Performing health check..."
curl -X GET http://localhost:8080/healthcheck
echo -e "\n"

status_code=$(curl -o /dev/null -s -w "%{http_code}\n" http://localhost:8080/healthcheck)
if [ "$status_code" != "200" ]; then
    echo "Health check failed with status code $status_code. Exiting."
    exit 1
fi

# Prompt the user for their API key but do not display it
read -sp "Please enter your API key: " api_key
echo # Move to a new line

# Array of GET endpoints that do not require additional parameters
declare -a endpoints=("tools" "datafiles" "flows" "queue-summary" "tags")

# Display the table headers
printf "%-20s %-20s %-20s\n" "Endpoint" "Without API Key" "With API Key"

# Test each endpoint
for endpoint in "${endpoints[@]}"; do
    test_endpoint "$endpoint" "$api_key"
done