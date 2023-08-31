#!/bin/bash

# exit the whole script if any command returns a non-zero exit code
set -e

# IFPS Custom setup
export IPFS_PATH=$(mktemp -d)
export BACALHAU_SERVE_IPFS_PATH=${IPFS_PATH}

# Ensure there is no port conflict with other canaries
ipfs init -e
ipfs bootstrap rm --all
ipfs config Addresses.API /ip4/127.0.0.1/tcp/8001
ipfs config Addresses.Gateway /ip4/127.0.0.1/tcp/8002
ipfs config Addresses.Swarm --json '["/ip4/0.0.0.0/tcp/8000"]'
export BACALHAU_IPFS_SWARM_ADDRESSES="/dns4/bacalhau.labdao.xyz/tcp/4001/p2p/$(curl -s -X POST bacalhau.labdao.xyz:5001/api/v0/id | jq -r '.ID')"

# plex must run from the same place as tools directory
cd {{ repo_dir }}

echo "$(date) - Running Canary" | tee -a plex_colabfold_out.log

plex init -t {{ repo_dir }}/tools/colabfold-mini.json -i '{"sequence": ["{{ repo_dir }}/testdata/folding/test.fasta"]}' --scatteringMethod=dotProduct --autoRun=true -a test -a cron | tee -a plex_colabfold_out.log
# capture the exit status of the plex call
plex_result_code=${PIPESTATUS[0]}
# exit immediately if plex exited with an error
if [ $plex_result_code -gt 0 ]; then
  exit $plex_result_code
fi

# parse the output directory from the plex stdout
result_dir=$(cat plex_colabfold_out.log | grep -a 'Finished processing, results written to' | sed -n 's/^.*Finished processing, results written to //p' | sed 's/\/io.json//' | tail -n 1)

# exit if no docked files are found
cd "$result_dir/entry-0/outputs"
if [ "$(find . -name '*scores.json')" == "" ]; then
  echo "No scores fiels found"
  exit 1
else
  echo "Scores files found"
fi

rm -rf $result_dir

curl -X POST -H "Authorization: Bearer ${HEII_ON_CALL_API_KEY}" https://api.heiioncall.com./triggers/991a6388-5c61-422c-b8cf-202b4c4b55a6/checkin

echo "$(date) - Colabfold mini canary success" | tee -a plex_colabfold_out.log

rm -rf ${IPFS_PATH}
