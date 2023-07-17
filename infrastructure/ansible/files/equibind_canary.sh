#!/bin/bash

# exit the whole script if any command returns a non-zero exit code
set -e

# plex must run from the same place as tools directory
cd {{ repo_dir }}

plex create -t {{ repo_dir }}/tools/equibind.json -i {{ repo_dir }}/testdata/binding/pdbbind_processed_size1 -a test -a cron --autoRun=true 2>&1 | tee plex_out.log
# capture the exit status of the plex call
plex_result_code=${PIPESTATUS[0]}
# exit immediately if plex exited with an error
if [ $plex_result_code -gt 0 ]; then
  exit $plex_result_code
fi

# parse the output directory from the plex stdout
result_dir=$(cat plex_out.log | grep 'Finished processing, results written to' | sed -n 's/^.*Finished processing, results written to //p' | sed 's/\/io.json//')

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

echo "Canary success"
