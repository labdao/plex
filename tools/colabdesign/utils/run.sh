#!/bin/bash

docker build -t generate_record_png .

INPUT_DIR="."
OUTPUT_DIR="test-runs/outputs_$(date +%Y%m%d_%H%M%S)"

mkdir -p "$PWD/$OUTPUT_DIR"
echo "Output directory is $OUTPUT_DIR"

INPUT_FILE="Condition_52_Design_0.pdb"

mkdir -p "$OUTPUT_DIR"

start_time=$(date +%s)
docker run --gpus all \
-e INPUT_FILE="/app/Condition_52_Design_0.pdb" \
-e OUTPUT_FILE="/app/output/output_image.png" \
-v .:/app/input \
-v "$PWD/$OUTPUT_DIR":/app/output generate_record_png

end_time=$(date +%s)

duration=$((end_time - start_time))
echo "Time taken: ${duration} seconds"