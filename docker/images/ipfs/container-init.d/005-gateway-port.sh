#!/bin/bash
set -e

if [ "${IPFS_DEBUG}" == "true" ]; then
  set -x
fi

if [ -n "${IPFS_GATEWAY_PORT}" ]; then
  ipfs config Addresses.Gateway /ip4/0.0.0.0/tcp/"${IPFS_GATEWAY_PORT}"
fi
