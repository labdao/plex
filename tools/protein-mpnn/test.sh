#!/bin/bash

# Build the Docker image
docker build -t pf/protein-mpnn .

PLEX_JOB_INPUTS=$(cat mock_plex_user_inputs_protein-mpnn.json)
# Run the Docker container
docker run -it --gpus=all \
-e PLEX_JOB_INPUTS="$PLEX_JOB_INPUTS" \
-v /home/convexity-research/philipp/lab-exchange/tools/protein-mpnn/:/inputs \
-v /home/convexity-research/philipp/lab-exchange/tools/protein-mpnn/outputs/:/app/outputs \
pf/protein-mpnn:latest

# /app
