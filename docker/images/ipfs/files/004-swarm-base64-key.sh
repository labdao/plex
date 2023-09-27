#!/bin/bash

set -e
user=ipfs
repo="$IPFS_PATH"

if [ "${IPFS_DEBUG}" == "true" ]; then
  set -x
fi


# Set up the swarm key, if provided
SWARM_KEY_FILE="$repo/swarm.key"
SWARM_KEY_PERM=0400

# Create a swarm key from a given environment variable
if [ -n "$IPFS_SWARM_KEY_BASE64" ] && [ "${PRIVATE_IPFS}" == "true" ]; then
  echo "Copying swarm key from variable IPFS_SWARM_KEY_BASE64..."
  echo "$IPFS_SWARM_KEY_BASE64" | base64 -d >"$SWARM_KEY_FILE" || exit 1
  chmod $SWARM_KEY_PERM "$SWARM_KEY_FILE"
fi
