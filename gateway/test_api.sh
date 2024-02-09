#!/bin/bash

# Health check
echo "Performing health check..."
curl -X GET http://localhost:8080/healthcheck
echo -e "\n"

status_code=$(curl -o /dev/null -s -w "%{http_code}\n" http://localhost:8080/healthcheck)
if [ "$status_code" != "200" ]; then
    echo "Health check failed with status code $status_code. Exiting."
    exit 1
fi

# Prompt the user for their API key
read -p "Please enter your API key: " api_key

# First request without the API key
echo "Making request without API key..."
curl -X GET http://localhost:8080/datafiles
echo -e "\n\n"

# Second request with the API key
echo "Making request with API key..."
curl -X GET -H "Authorization: Bearer $api_key" http://localhost:8080/datafiles
echo -e "\n"

