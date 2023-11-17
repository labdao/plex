#!/bin/bash

# Build the Docker image
docker build -t colabdesign-noninteractive .

# Create a unique output directory
OUTPUT_DIR="outputs_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$PWD/$OUTPUT_DIR"
echo "Output directory is $OUTPUT_DIR"

# Run the Docker container
docker run --gpus=all -v $PWD/testdata/pdc_upar_1_target.pdb:/inputs/target_protein/pdc_upar_1_target.pdb \
-v $PWD/testdata/VTNCparams1.yaml:/app/conf/params/VTNCparams1.yaml \
-v "$PWD/$OUTPUT_DIR":/outputs colabdesign-noninteractive python -u main.py inputs=container outputs=container params=VTNCparams1

# Check if output default1.pdb exists in $OUTPUT_DIR
if [ -f "$PWD/$OUTPUT_DIR/default_best.pdb" ]; then
    echo "Output file default_best.pdb found in $PWD/$OUTPUT_DIR."
else
    echo "Output file default_best.pdb not found in $PWD/$OUTPUT_DIR."
fi
