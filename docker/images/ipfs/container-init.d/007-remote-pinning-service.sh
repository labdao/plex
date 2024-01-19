#!/bin/bash
# Script to add remote pinning service
# required args:
# IPFS_ADD_REMOTE_PINNING_SERVICE=true/false
# IPFS_REMOTE_PINNING_SERVICE_NAME=userfiendlyname
# IPFS_REMOTE_PINNING_SERVICE_ENDPOINT=endpointurl
# IPFS_REMOTE_PINNING_SERVICE_ACCESS_TOKEN=accesstoken

set -e

if [ "${IPFS_DEBUG}" == "true" ]; then
  set -x
fi

if [ "${IPFS_ADD_REMOTE_PINNING_SERVICE}" == "true" ]; then
  ipfs pin remote service add ${IPFS_REMOTE_PINNING_SERVICE_NAME} ${IPFS_REMOTE_PINNING_SERVICE_ENDPOINT} ${IPFS_REMOTE_PINNING_SERVICE_ACCESS_TOKEN}
fi
