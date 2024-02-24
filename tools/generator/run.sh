#!/bin/bash
USER_NAME="philipp"
CONTAINER_NAME="protein-binder-designer"
PLEX_JOB_INPUTS=$(cat user_input.json)
HOST_CACHE_DIR=/home/convexity-research/$USER_NAME/lab-exchange/tools/generator/cache

docker build --rm -t $CONTAINER_NAME .

# Run the Docker container
mkdir -p "$HOST_CACHE_DIR"
docker run -it --gpus=all \
-e PLEX_JOB_INPUTS="$PLEX_JOB_INPUTS" \
-e TRANSFORMERS_CACHE=/transformers_cache \
-v "$HOST_CACHE_DIR":/transformers_cache \
-v "/home/convexity-research/$USER_NAME/lab-exchange/tools/generator/":/inputs \
-v "/home/convexity-research/$USER_NAME/lab-exchange/tools/generator/outputs/":/app/outputs \
$CONTAINER_NAME:latest

# -v "$HOST_CACHE_DIR":/app/cache \

# # Run the Docker container
# mkdir -p "$HOST_CACHE_DIR"
# docker run -it --gpus=all \
# -e PLEX_JOB_INPUTS="$PLEX_JOB_INPUTS" \
# -e TRANSFORMERS_CACHE=/transformers_cache \
# -v "$HOST_CACHE_DIR":/transformers_cache \
# -v "/home/convexity-research/$USER_NAME/lab-exchange/tools/generator/":/inputs \
# -v "/home/convexity-research/$USER_NAME/lab-exchange/tools/generator/outputs/":/app/outputs \
# $CONTAINER_NAME:latest

# #!/bin/bash
# CONTAINER_NAME="protein-binder-designer"
# PLEX_JOB_INPUTS=$(cat user_input.json)
# HOST_CACHE_DIR=$(pwd)/cache

# docker build --rm -t $CONTAINER_NAME .

# mkdir -p "$HOST_CACHE_DIR"
# docker run -it --gpus=all \
# -e PLEX_JOB_INPUTS="$PLEX_JOB_INPUTS" \
# -e HF_HOME=/transformers_cache \
# -v "$HOST_CACHE_DIR":/transformers_cache \
# -v "$(pwd)/outputs":/app/outputs \
# $CONTAINER_NAME:latest