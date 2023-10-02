#!/bin/sh
set -euo pipefail

if [ "${BACALHAU_DEBUG}" == "true" ]; then
  set -x
fi

peer_id=""
peer_addrs=""
swarm_peers=""

# IPFS Peers
if [ "x${IPFS_PEERS}" != "x" ]; then
  for peer in ${IPFS_PEERS//,/ }; do
    peer_id="$(curl -s -X POST http://${peer}:5001/api/v0/id | jq -r '.ID')"
    peer_addrs="/dns4/${peer}/tcp/4001/p2p/${peer_id}"
    swarm_peers="${swarm_peers}${swarm_peers:+,}${peer_addrs}"
  done

  export BACALHAU_IPFS_SWARM_ADDRESSES="${swarm_peers}"
fi

# Bacalhau Peers
peer_id=""
peer_addrs=""
swarm_peers=""
if [ "x${BACALHAU_PEERS}" != "x" ]; then
  for peer in ${BACALHAU_PEERS//,/ }; do
    peer_id="$(curl -s http://${peer}:1234/node_info | jq -r '.PeerInfo.ID')"
    peer_addrs="/dns4/${peer}/tcp/1235/p2p/${peer_id}"
    swarm_peers="${swarm_peers}${swarm_peers:+,}${peer_addrs}"
  done
  export BACALHAU_PEER_CONNECT="${swarm_peers}"
fi


if [ "${BACALHAU_DEBUG}" == "true" ]; then
  echo "Dumping envs"
  env
fi

exec bacalhau $@
