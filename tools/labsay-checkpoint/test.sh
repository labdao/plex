#!/bin/bash

docker build -t labsay .

OUTPUT_DIR="test-runs/outputs_$(date +%Y%m%d_%H%M%S)"
mkdir -p "$PWD/$OUTPUT_DIR"
echo "Output directory is $OUTPUT_DIR"

docker run \
-e PLEX_JOB_INPUTS='{"file_example": "/inputs/file_example/message.txt", "string_example": "hello world", "number_example": 196883}' \
-e JOB_UUID='1234' \
--env-file ~/aws.env \
-v $PWD/testdata/inputs:/inputs/ \
-v "$PWD/$OUTPUT_DIR":/outputs labsay

if [ -f "$PWD/$OUTPUT_DIR/result.txt" ]; then
    echo "Output file result.txt found in $PWD/$OUTPUT_DIR."
else
    echo "Output file result.txt not found in $PWD/$OUTPUT_DIR."
fi
