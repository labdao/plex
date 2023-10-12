#!/bin/sh
set -e

repo="$IPFS_PATH"

if [ -e "$repo/config" ]; then
  echo "Found IPFS fs-repo at $repo"
else
  ipfs init ${IPFS_PROFILE:+"--profile=$IPFS_PROFILE"}
  ipfs config Addresses.API /ip4/0.0.0.0/tcp/5001
  ipfs config Addresses.Gateway /ip4/0.0.0.0/tcp/8080
fi


find /container-init.d -maxdepth 1 -type f -iname '*.sh' -print0 | sort -z | xargs -n 1 -0 -r container_init_run

exec /plex "$@"
