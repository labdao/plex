#!/bin/bash

# Build the Docker image
docker build -t colabdesign .

# Create a unique output directory that is gitignored
OUTPUT_DIR="test-runs/outputs_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$PWD/$OUTPUT_DIR"
echo "Output directory is $OUTPUT_DIR"

FLOW_UUID="test-flow_uuid_$(date +%y%m%d)"
JOB_UUID="test-job_uuid_$(date +%y%m%d_%H%M%S)"
# Run the Docker container
# docker run --gpus=all -v $PWD/testdata/pdc_upar_1_target.pdb:/inputs/target_protein/pdc_upar_1_target.pdb \
# -v $PWD/testdata/VTNCparams1.yaml:/app/conf/params/VTNCparams1.yaml \
# -v "$PWD/$OUTPUT_DIR":/outputs colabdesign python -u main.py inputs=container outputs=container params=VTNCparams1

docker run --gpus=all \
-e PLEX_JOB_INPUTS='{"binder_length":10,"hotspot":"","number_of_binders":1,"target_chain":"B","target_start_residue":50,"target_end_residue":100,"target_protein":"/inputs/target_protein/pdc_upar_1_target.pdb","contigs_override":"A1-283:11/2/5/11/11"}' \
--env-file ~/aws.env \
-e FLOW_UUID="$FLOW_UUID" \
-e JOB_UUID="$JOB_UUID" \
-e CHECKPOINT_COMPATIBLE="False" \
-v $PWD/testdata/inputs:/inputs/ \
-v "$PWD/$OUTPUT_DIR":/outputs colabdesign

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
