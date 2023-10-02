#!/bin/sh

set -e

if [ "${IPFS_DEBUG}" = "true" ]; then
  set -x
fi

echo "[]" > /tmp/peers.json

if [ -n "${IPFS_PEERS}" ]; then

  echo "${IPFS_PEERS}" | tr ',' '\n' | while read -r peer; do 
    peer_id="$(curl -s -X POST "http://${peer}:5001/api/v0/id" | jq -r '.ID')"
    peer_addrs="/dns4/${peer}/tcp/4001/p2p/${peer_id}"
    peering_peers="[{\"ID\":\"${peer_id}\", \"Addrs\": [\"${peer_addrs}\"]}]"
    
    cat <<< $(jq ". += ${peering_peers}" /tmp/peers.json) > /tmp/peers.json
  done

  ipfs config --json Peering.Peers "$(cat /tmp/peers.json)"

fi
