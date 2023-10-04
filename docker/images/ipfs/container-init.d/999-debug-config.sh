#!/bin/bash
set -e

if [ "${IPFS_DEBUG}" == "true" ]; then
  set -x
  echo "Dumping env"
  env
  cat /data/ipfs/config
  ls -ltra /data/ipfs/repo.lock || true
fi
