#!/bin/bash
# Create a unique output directory
OUTPUT_DIR="outputs"
mkdir -p "$PWD/$OUTPUT_DIR"
echo "Output directory is $OUTPUT_DIR"

PLEX_JOB_INPUTS=$(cat mock_plex_user_input_autoprompt.json)
docker run -it --gpus=all \
-e PLEX_JOB_INPUTS="$PLEX_JOB_INPUTS" \
-v $PWD:/inputs/ \
-v "$PWD/$OUTPUT_DIR":/outputs quay.io/labdao/colabdesign:sha-631ebc2 /bin/bash -c "cp /inputs/main.py . && cp /inputs/prompt_generator.py . && cp -r /inputs/conf . \
&& ls /inputs && echo 'Attempting main.py...' && python3 -u main.py"