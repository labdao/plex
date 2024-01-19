#!/bin/bash

# Build the Docker image
docker build -t pf/rf-diffusion-updated .

PLEX_JOB_INPUTS=$(cat mock_plex_user_inputs_rf-diffusion.json)
# Run the Docker container
docker run -it --gpus=all \
-e PLEX_JOB_INPUTS="$PLEX_JOB_INPUTS" \
-v /home/convexity-research/philipp/lab-exchange/tools/rf-diffusion/:/inputs \
-v /home/convexity-research/philipp/lab-exchange/tools/rf-diffusion/outputs/:/app/outputs \
pf/rf-diffusion-updated:latest

# /app