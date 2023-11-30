#!/bin/bash

# Build the Docker image
docker build -t csv-aggregate .

# Create a unique output directory
OUTPUT_DIR="outputs_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$PWD/$OUTPUT_DIR"
echo "Output directory is $OUTPUT_DIR"

docker run \
-e PLEX_JOB_INPUTS='{"input_csvs": ["/inputs/input_csvs/0/result_1.csv", "/inputs/input_csvs/1/result_2.csv", "/inputs/input_csvs/2/result_3.csv"]}' \
-v $PWD/testdata/inputs:/inputs/ \
-v "$PWD/$OUTPUT_DIR":/outputs csv-aggregate

# Check if output default_best.pdb exists in $OUTPUT_DIR
if [ -f "$PWD/$OUTPUT_DIR/aggregated_results.csv" ]; then
    echo "Output file aggregated_results.csv found in $PWD/$OUTPUT_DIR."
else
    echo "Output file aggregated_results.csv not found in $PWD/$OUTPUT_DIR."
fi
