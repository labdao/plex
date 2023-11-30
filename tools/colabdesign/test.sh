#!/bin/bash

# Build the Docker image
docker build -t colabdesign-noninteractive .

# Create a unique output directory
OUTPUT_DIR="outputs_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$PWD/$OUTPUT_DIR"
echo "Output directory is $OUTPUT_DIR"

# Run the Docker container
# docker run --gpus=all -v $PWD/testdata/pdc_upar_1_target.pdb:/inputs/target_protein/pdc_upar_1_target.pdb \
# -v $PWD/testdata/VTNCparams1.yaml:/app/conf/params/VTNCparams1.yaml \
# -v "$PWD/$OUTPUT_DIR":/outputs colabdesign-noninteractive python -u main.py inputs=container outputs=container params=VTNCparams1

docker run --gpus=all \
-e PLEX_JOB_INPUTS='{"binder_length":50,"hotspot":"","number_of_binders":2,"target_chain":"B","target_start_residue":50,"target_end_residue":200,"target_protein":"/inputs/target_protein/pdc_upar_1_target.pdb","contigs_override":"A1-283:11/2/5/11/11"}' \
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