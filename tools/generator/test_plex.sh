#!/bin/bash

# Build the Docker image
docker build --rm -t pf/generator .

PLEX_JOB_INPUTS=$(cat mock_plex_user_input_generator.json)

# Define a directory on the host for the transformers cache
HOST_CACHE_DIR=/home/convexity-research/philipp/lab-exchange/tools/generator/cache

# Ensure the cache directory exists on the host
mkdir -p "$HOST_CACHE_DIR"

# Run the Docker container
docker run -it --gpus=all \
-e PLEX_JOB_INPUTS="$PLEX_JOB_INPUTS" \
-e TRANSFORMERS_CACHE=/transformers_cache \
-v "$HOST_CACHE_DIR":/transformers_cache \
-v /home/convexity-research/philipp/lab-exchange/tools/generator/:/inputs \
-v /home/convexity-research/philipp/lab-exchange/tools/generator/outputs/:/app/outputs \
pf/generator:latest
