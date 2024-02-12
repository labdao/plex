#!/bin/bash
USER_NAME="philipp"
CONTAINER_NAME="pf/generator-deployment-test"
PLEX_JOB_INPUTS=$(cat mock_plex_user_input_generator.json)
HOST_CACHE_DIR=/home/convexity-research/$USER_NAME/lab-exchange/tools/generator/cache

# Build the Docker image
docker build --rm -t $CONTAINER_NAME .

# Run the Docker container
mkdir -p "$HOST_CACHE_DIR"
docker run -it --gpus=all \
-e PLEX_JOB_INPUTS="$PLEX_JOB_INPUTS" \
-e TRANSFORMERS_CACHE=/transformers_cache \
-v "$HOST_CACHE_DIR":/transformers_cache \
-v "/home/convexity-research/$USER_NAME/lab-exchange/tools/generator/":/inputs \
-v "/home/convexity-research/$USER_NAME/lab-exchange/tools/generator/outputs/":/app/outputs \
$CONTAINER_NAME:latest