#!/bin/bash

docker build -t labsay .

OUTPUT_DIR="test-runs/outputs_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$PWD/$OUTPUT_DIR"
echo "Output directory is $OUTPUT_DIR"

docker run \
-e PLEX_JOB_INPUTS='{"file_example": "/inputs/file_example/message.txt", "string_example": "hello world", "number_example": 196883, "pdb_checkpoint_0": "/inputs/pdb_checkpoints/example.pdb", "pdb_checkpoint_1": "/inputs/pdb_checkpoints/design_1.pdb", "pdb_checkpoint_2": "/inputs/pdb_checkpoints/BioCD202b18_aa_7fd4f_unrelaxed_rank_003_alphafold2_multimer_v3_model_2_seed_000.pdb"}' \
-e JOB_UUID='1234' \
--env-file ~/aws.env \
-v $PWD/testdata/inputs:/inputs/ \
-v "$PWD/$OUTPUT_DIR":/outputs labsay

if [ -f "$PWD/$OUTPUT_DIR/result.txt" ]; then
    echo "Output file result.txt found in $PWD/$OUTPUT_DIR."
else
    echo "Output file result.txt not found in $PWD/$OUTPUT_DIR."
fi
