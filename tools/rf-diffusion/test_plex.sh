#!/bin/bash

# Build the Docker image
docker build -t pf/sequence-sampler-updated .

PLEX_JOB_INPUTS=$(cat mock_plex_user_input_sequence_sampler.json)
# Run the Docker container
docker run -it --gpus=all \
-e PLEX_JOB_INPUTS="$PLEX_JOB_INPUTS" \
-v /home/convexity-research/philipp/lab-exchange/tools/sequence-sampler/:/inputs \
-v /home/convexity-research/philipp/lab-exchange/tools/sequence-sampler/outputs/:/app/outputs \
pf/sequence-sampler-updated:latest
