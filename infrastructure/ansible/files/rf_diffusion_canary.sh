#!/bin/bash

# exit the whole script if any command returns a non-zero exit code
set -e

# IFPS Custom setup
export IPFS_PATH=$(mktemp -d)
export BACALHAU_SERVE_IPFS_PATH=${IPFS_PATH}
export LOGFILE={{ repo_dir }}/plex_new_canary_out.log

# Ensure there is no port conflict with other canaries
ipfs init -e
ipfs bootstrap rm --all
ipfs config Addresses.API /ip4/127.0.0.1/tcp/9001
ipfs config Addresses.Gateway /ip4/127.0.0.1/tcp/9002
ipfs config Addresses.Swarm --json '["/ip4/0.0.0.0/tcp/9000"]'
export BACALHAU_IPFS_SWARM_ADDRESSES="/dns4/bacalhau.labdao.xyz/tcp/4001/p2p/$(curl -s -X POST bacalhau.labdao.xyz:5001/api/v0/id | jq -r '.ID')"

# plex must run from the same place as tools directory
cd {{ repo_dir }}
touch ${LOGFILE}
echo "$(date) - Running Canary" | tee -a ${LOGFILE}

# Plex command to be provided



# plex init ...

# capture the exit status of the plex call
plex_result_code=${PIPESTATUS[0]}
# exit immediately if plex exited with an error
if [ $plex_result_code -gt 0 ]; then
  exit $plex_result_code
fi

# parse the output directory from the plex stdout
result_dir=$(cat ${LOGFILE} | grep -a 'Finished processing, results written to' | sed -n 's/^.*Finished processing, results written to //p' | sed 's/\/io.json//' | tail -n 1)

# exit if no output files are found
cd "$result_dir/entry-0/outputs"
# Replace '*scores.json' with the actual output file you expect






if [ "$(find . -name '*scores.json')" == "" ]; then
  echo "No output files found"
  exit 1
else
  echo "Output files found"
fi

rm -rf $result_dir

# Replace the URL and API key with the actual ones for the new canary
curl -X POST -H "Authorization: Bearer ${HEII_ON_CALL_API_KEY}" https://api.heiioncall.com./triggers/.../checkin

echo "$(date) - New canary success" | tee -a ${LOGFILE}

rm -rf ${IPFS_PATH}