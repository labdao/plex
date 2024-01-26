#!/bin/bash

# Build the Docker image
docker build -t pf/af2-update .

PLEX_JOB_INPUTS=$(cat mock_plex_user_inputs_af2.json)

echo "Running the Docker container..."
docker run -it --rm --gpus all \
-e PLEX_JOB_INPUTS="$PLEX_JOB_INPUTS" \
-v /home/convexity-research/philipp/lab-exchange/tools/af2:/inputs \
-v /home/convexity-research/philipp/lab-exchange/tools/af2:/app \
pf/af2-update:latest
echo "Docker container has finished running."