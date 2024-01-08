#!/bin/bash

# Build the Docker image
docker build -t aggregater .

# Create a unique output directory
OUTPUT_DIR="test-runs/outputs_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$PWD/$OUTPUT_DIR"
echo "Output directory is $OUTPUT_DIR"

docker run \
-e PLEX_JOB_INPUTS='{"csv_result_files": ["/inputs/csv_result_files/0/example_default_scores.csv", "/inputs/csv_result_files/1/example_default_scores.csv", "/inputs/csv_result_files/2/example_default_scores.csv"]}' \
-v $PWD/testdata/inputs:/inputs/ \
-v "$PWD/$OUTPUT_DIR":/outputs aggregater

if [ -f "$PWD/$OUTPUT_DIR/distribution.png" ]; then
    echo "Output file distribution.png found in $PWD/$OUTPUT_DIR."
else
    echo "Output file distribution.png not found in $PWD/$OUTPUT_DIR."
fi

if [ -f "$PWD/$OUTPUT_DIR/aggregated.csv" ]; then
    echo "Output file aggregated.csv found in $PWD/$OUTPUT_DIR."
else
    echo "Output file aggregated.csv not found in $PWD/$OUTPUT_DIR."
fi
