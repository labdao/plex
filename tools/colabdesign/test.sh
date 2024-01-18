#!/bin/bash

# Build the Docker image
docker build -t colabdesign-noninteractive .

# Create a unique output directory
OUTPUT_DIR="outputs_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$PWD/$OUTPUT_DIR"
echo "Output directory is $OUTPUT_DIR"

# Run the Docker container
docker run --gpus=all \
-e PLEX_JOB_INPUTS='{"binder_length":80,"hotspot":"B58,B80,B139","number_of_binders":2,"target_chain":"B","target_start_residue":17,"target_end_residue":209,"target_protein":"/inputs/target_protein/3di3.pdb","contigs_override":""}' \
-v $PWD/testdata/inputs:/inputs/ \
-v "$PWD/$OUTPUT_DIR":/outputs colabdesign-noninteractive

# Check if output default_best.pdb exists in $OUTPUT_DIR
if [ -f "$PWD/$OUTPUT_DIR/default_best.pdb" ]; then
    echo "Output file default_best.pdb found in $PWD/$OUTPUT_DIR."
else
    echo "Output file default_best.pdb not found in $PWD/$OUTPUT_DIR."
fi

# Check if output default_scores.csv exists in $OUTPUT_DIR
if [ -f "$PWD/$OUTPUT_DIR/default_scores.csv" ]; then
    echo "Output file default_scores.csv found in $PWD/$OUTPUT_DIR."
else
    echo "Output file default_scores.csv not found in $PWD/$OUTPUT_DIR."
fi
