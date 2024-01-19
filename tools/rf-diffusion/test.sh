#!/bin/bash

# Build the Docker image
docker build -t pf/rf-diffusion .

# Run the Docker container
docker run -it --gpus all \
  -v /home/convexity-research/philipp/lab-exchange/tools/rf-diffusion/:/inputs \
  -v /home/convexity-research/philipp/lab-exchange/tools/rf-diffusion/outputs/:/app/outputs \
