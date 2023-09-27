#!/bin/bash
set -e

if [ "${IPFS_DEBUG}" == "true" ]; then
  set -x
fi

if [ "${PRIVATE_IPFS}" == "true" ]; then
  echo "Running in private mode, removing bootstrap"
  ipfs bootstrap rm --all
fi
