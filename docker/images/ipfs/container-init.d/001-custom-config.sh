#!/bin/sh
set -e

if [ "${IPFS_DEBUG}" == "true" ]; then
  set -x
fi

ipfs config --json API.HTTPHeaders.Access-Control-Allow-Methods '["PUT", "POST", "GET"]'
ipfs config Pinning.Recursive true
ipfs config --json Swarm.RelayClient.Enabled false
