#!/bin/bash
PLEX_JOB_INPUTS=$(cat user_input.json)

docker build -t labsay .

FLOW_UUID="test-flow_uuid_$(date +%y%m%d)"
JOB_UUID="test-job_uuid_$(date +%y%m%d_%H%M%S)"

OUTPUT_DIR="test-runs/outputs_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$PWD/$OUTPUT_DIR"
echo "Output directory is $OUTPUT_DIR"

docker run \
-e PLEX_JOB_INPUTS="$PLEX_JOB_INPUTS" \
-e FLOW_UUID="$FLOW_UUID" \
-e JOB_UUID="$JOB_UUID" \
-e CHECKPOINT_COMPATIBLE="False" \
--env-file ~/aws.env \
-v $PWD/testdata/inputs:/inputs/ \
-v "$PWD/$OUTPUT_DIR":/outputs labsay
