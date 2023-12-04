#!/bin/bash

# Variables
IMAGE="ghcr.io/sokrypton/colabfold:1.5.3-cuda11.8.0"
CACHE_DIR="$PWD/cache"
INPUT_DIR="$PWD"
INPUT_FILE="input_sequence.fasta"
OUTPUT_DIR="$PWD/output"

# Create directories
mkdir -p "$CACHE_DIR"
mkdir -p "$OUTPUT_DIR"
echo "Cache directory is $CACHE_DIR"
echo "Output directory is $OUTPUT_DIR"

# Pull the Docker Image
echo "Pulling Docker image..."
docker pull $IMAGE

# Download AlphaFold2 Weights
echo "Downloading AlphaFold2 weights..."
docker run --user $(id -u) -ti --rm \
  -v "$CACHE_DIR":/cache:rw \
  $IMAGE \
  python -m colabfold.download

# Run Prediction
echo "Running prediction job..."
docker run --user $(id -u) -ti --rm --runtime=nvidia --gpus 1 \
  -v "$CACHE_DIR":/cache:rw \
  -v "$INPUT_DIR":/work:rw \
  $IMAGE \
  colabfold_batch /work/$INPUT_FILE /work/output

echo "Prediction job complete. Results are in $OUTPUT_DIR"
