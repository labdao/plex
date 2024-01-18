#!/bin/bash

# Build the Docker image
docker build -t labsay .

# Create a unique output directory
OUTPUT_DIR="test-runs/outputs_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$PWD/$OUTPUT_DIR"
echo "Output directory is $OUTPUT_DIR"

docker run \
-e PLEX_JOB_INPUTS='{"file_example": "/inputs/file_example/message.txt", "string_example": "hello world", "number_example": 196883}' \
-v $PWD/testdata/inputs:/inputs/ \
-v "$PWD/$OUTPUT_DIR":/outputs labsay

# Check if output default_best.pdb exists in $OUTPUT_DIR
if [ -f "$PWD/$OUTPUT_DIR/result.txt" ]; then
    echo "Output file result.txt found in $PWD/$OUTPUT_DIR."
else
    echo "Output file result.txt not found in $PWD/$OUTPUT_DIR."
fi
