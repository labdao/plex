#!/bin/bash

# exit the whole script if any command returns a non-zero exit code
set -e

# IFPS Custom setup
export IPFS_PATH=$(mktemp -d)
export BACALHAU_SERVE_IPFS_PATH=${IPFS_PATH}

# Ensure there is no port conflict with other canaries
ipfs init -e
ipfs bootstrap rm --all
ipfs config Addresses.API /ip4/127.0.0.1/tcp/7001
ipfs config Addresses.Gateway /ip4/127.0.0.1/tcp/7002
ipfs config Addresses.Swarm --json '["/ip4/0.0.0.0/tcp/7000"]'
export BACALHAU_IPFS_SWARM_ADDRESSES="/dns4/bacalhau.labdao.xyz/tcp/4001/p2p/$(curl -s -X POST bacalhau.labdao.xyz:5001/api/v0/id | jq -r '.ID')"

# plex must run from the same place as tools directory
cd {{ repo_dir }}

echo "$(date) - Running Canary" | tee -a plex_equibind_out.log

plex init -t {{ repo_dir }}/tools/equibind.json -i '{"protein": ["{{ repo_dir }}/testdata/binding/abl/7n9g.pdb"], "small_molecule": ["{{ repo_dir }}/testdata/binding/abl/ZINC000003986735.sdf"]}' --scatteringMethod=dotProduct -a test -a cron --autoRun=true 2>&1 | tee -a plex_equibind_out.log
# capture the exit status of the plex call
plex_result_code=${PIPESTATUS[0]}
# exit immediately if plex exited with an error
if [ $plex_result_code -gt 0 ]; then
  exit $plex_result_code
fi

# parse the output directory from the plex stdout
result_dir=$(cat plex_equibind_out.log | grep -a 'Finished processing, results written to' | sed -n 's/^.*Finished processing, results written to //p' | sed 's/\/io.json//' | tail -n 1)

# exit if no docked files are found
cd "$result_dir/entry-0/outputs"
if [ "$(find . -name '*docked.sdf' | grep 'docked.sdf')" == "" ]; then
  echo "No docked files found"
  exit 1
else
  echo "Docked files found"
fi

rm -rf $result_dir

curl -X POST -H "Authorization: Bearer ${HEII_ON_CALL_API_KEY}" https://api.heiioncall.com./triggers/a7660c2f-5262-4392-918b-d98aee244890/checkin

echo "$(date) - Canary success" | tee -a plex_equibind_out.log

rm -rf ${IPFS_PATH}
