#!/bin/bash

# exit the whole script if any command returns a non-zero exit code
set -e

# plex must run from the same place as tools directory
cd {{ repo_dir }}

plex init -t {{ repo_dir }}/tools/colabfold-mini.json -i '{"sequence": ["{{ repo_dir }}/testdata/folding/test.fasta"]}' --scatteringMethod=dotProduct --autoRun=true -a test -a cron | tee plex_colabfold_out.log
# capture the exit status of the plex call
plex_result_code=${PIPESTATUS[0]}
# exit immediately if plex exited with an error
if [ $plex_result_code -gt 0 ]; then
  exit $plex_result_code
fi

# parse the output directory from the plex stdout
result_dir=$(cat plex_colabfold_out.log | grep 'Finished processing, results written to' | sed -n 's/^.*Finished processing, results written to //p' | sed 's/\/io.json//')

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

echo "Colabfold mini canary success"
