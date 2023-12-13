#!/bin/bash

# Build the Docker image
docker build -t protein-to-dna-icodon .

# Create a unique output directory
OUTPUT_DIR="test-runs/outputs_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$PWD/$OUTPUT_DIR"
echo "Output directory is $OUTPUT_DIR"

# input pdb file
docker run \
-e PLEX_JOB_INPUTS='{"input_file": "/inputs/input_file/VTNC_HUMAN_2-41.pdb", "specie": "human", "iterations": "20", "make_more_optimal": "T"}' \
-v $PWD/testdata/inputs:/inputs/ \
-v "$PWD/$OUTPUT_DIR":/outputs protein-to-dna-icodon

# Create a unique output directory
OUTPUT_DIR="test-runs/outputs_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$PWD/$OUTPUT_DIR"
echo "Output directory is $OUTPUT_DIR"

# input a txt file
docker run \
-e PLEX_JOB_INPUTS='{"input_file": "/inputs/input_file/amino_acid_sequence.txt", "specie": "human", "iterations": "10", "make_more_optimal": "T"}' \
-v $PWD/testdata/inputs:/inputs/ \
-v "$PWD/$OUTPUT_DIR":/outputs protein-to-dna-icodon