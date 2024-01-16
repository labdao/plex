#!/bin/bash

# Build the Docker image
docker build -t colabdesign-noninteractive-autoprompt .

OUTPUT_DIR="test-runs/outputs_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$PWD/$OUTPUT_DIR"
echo "Output directory is $OUTPUT_DIR"

PLEX_JOB_INPUTS=$(cat mock_plex_user_input_autoprompt.json)
docker run --gpus=all \
-e PLEX_JOB_INPUTS="$PLEX_JOB_INPUTS" \
-v "$PWD/testdata/inputs":/inputs/ \
-v "$PWD/$OUTPUT_DIR":/outputs colabdesign-noninteractive-autoprompt
