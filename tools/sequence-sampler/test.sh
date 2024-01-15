#!/bin/bash

# Build the Docker image
docker build -t seqoptim .

# Run the Docker container
docker run -it --gpus all \
  -v /home/convexity-research/philipp/lab-exchange/tools/sequence-sampler/:/app/seq-sampler \
  -v /home/convexity-research/philipp/lab-exchange/tools/sequence-sampler/:/inputs \
  -v /home/convexity-research/philipp/lab-exchange/tools/sequence-sampler/outputs/:/app/outputs \
  seqoptim:latest