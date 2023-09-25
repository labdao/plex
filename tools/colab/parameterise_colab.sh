#!/bin/bash

# Input notebook file
NOTEBOOK_FILE="$1"

# Output notebook file
OUTPUT_NOTEBOOK_FILE="$2"

# Tags and markers
PARAMETER_TAG="parameters"
PARAMETER_MARKER="#%%param"

# Check if output file name is provided, else default to "_modified"
if [ -z "$OUTPUT_NOTEBOOK_FILE" ]; then
  OUTPUT_NOTEBOOK_FILE="${NOTEBOOK_FILE%.ipynb}_modified.ipynb"
fi

# Extracting cells with our specific marker
TAG_CELL_INDICES=$(jq -r --arg MARKER "$PARAMETER_MARKER" '
  [.cells | to_entries[] | select(.value.source | join("") | contains($MARKER)) | .key]
' "$NOTEBOOK_FILE")

# Check for number of tagged cells
NUM_TAGGED_CELLS=$(echo "$TAG_CELL_INDICES" | jq '. | length')
if (( NUM_TAGGED_CELLS > 1 )); then
  echo "Error: More than one cell contains the parameter marker. Exiting."
  exit 1
fi

# Only proceed if a valid cell was found
if (( NUM_TAGGED_CELLS == 1 )); then
  # Modify the notebook
  jq --arg TAG "$PARAMETER_TAG" --argjson INDEX "$(echo "$TAG_CELL_INDICES" | jq '.[0]')" '
    if .cells[$INDEX].metadata.tags then
      .cells[$INDEX].metadata.tags += [$TAG]
    else
      .cells[$INDEX].metadata = {tags: [$TAG]}
    end
  ' "$NOTEBOOK_FILE" > "$OUTPUT_NOTEBOOK_FILE"
  echo "Modifications done. New notebook is: $OUTPUT_NOTEBOOK_FILE"
else
  echo "No cells with the specified parameter marker found."
fi
