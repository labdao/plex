#!/bin/bash

# Initialize variables with default values
total_jobs=10
gpu=false
autoRun=true
selected_tool="equibind"

# Parse command-line options
while getopts ":n:g:a:" opt; do
  case $opt in
    n)
      total_jobs="$OPTARG"
      ;;
    g)
      gpu="$OPTARG"
      ;;
    a)
      autoRun="$OPTARG"
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      exit 1
      ;;
    :)
      echo "Option -$OPTARG requires an argument." >&2
      exit 1
      ;;
  esac
done

# Set the selected_tool and plex_command based on the options
if [ "$autoRun" = "true" ]; then
    if [ "$gpu" = "true" ]; then
      selected_tool="colabfold"
      plex_command="./plex init -t tools/colabfold-mini.json -i \"{\\\"sequence\\\": [\\\"testdata/folding/test.fasta\\\"]}\" --scatteringMethod=dotProduct --autoRun=true -a test -a debug"
    else
      plex_command="./plex init -t tools/equibind.json -i \"{\\\"protein\\\": [\\\"testdata/binding/abl/7n9g.pdb\\\"], \\\"small_molecule\\\": [\\\"testdata/binding/abl/ZINC000003986735.sdf\\\"]}\" --scatteringMethod=dotProduct --autoRun=true -a test -a debug"
    fi
else
    if [ "$gpu" = "true" ]; then
      selected_tool="colabfold"
      plex_command="./plex run -i QmZmxyfb8aTr2pY28GgWgEJomr9fSWhkhTK19x7FmsjZPn -a test -a debug"
    else
      plex_command="./plex run -i QmbbJcDA7LVCZPiTwMatU4JgSHwTnZfVjEs3SnSMK5EnG1 -a test -a debug"
    fi
fi

# Initialize the output table file based on selected_tool
output_file="multiple_${selected_tool}_runs.csv"

# Write the table header to the output file
echo "Run,Stdout,Successful run,Bacalhau job ID,Error message" > "$output_file"

# Counter for completed jobs
completed_jobs=0

# Counter for successful and failed jobs
successful_jobs=0
failed_jobs=0

echo "Running $total_jobs $selected_tool jobs..."

# Loop to run the plex command
for i in $(seq 1 $total_jobs); do
    # Run the plex command and capture its stdout
    stdout=$(eval $plex_command)

    # Initialize Bacalhau job ID and Error message as empty
    bacalhau_id=""
    error_message=""

    # Parse Bacalhau job ID
    if [[ "$stdout" =~ "Bacalhau job id: ([a-zA-Z0-9-]+)" ]]; then
        bacalhau_id=${BASH_REMATCH[1]}
    fi

    # Determine if the run was successful or failed
    if [[ "$stdout" == *"Success processing IO entry 0"* ]]; then
        success="True"
        successful_jobs=$((successful_jobs + 1))
    else
        success="False"
        failed_jobs=$((failed_jobs + 1))

        # Parse Error message
        if [[ "$stdout" =~ "error updating IO with result: (.+)$" ]]; then
            error_message=${BASH_REMATCH[1]}
        fi
    fi

    # Append the run number, the captured stdout, and other parsed data to the output table
    echo "\"$i\",\"$(echo "$stdout" | sed 's/"/\\"/g')\",\"$success\",\"$bacalhau_id\",\"$error_message\"" >> "$output_file"

    # Update completed jobs counter
    completed_jobs=$((completed_jobs + 1))

    # Calculate percentage of completion and display loading bar and success/failure counts
    percent_completed=$(awk "BEGIN { pc=100*${completed_jobs}/${total_jobs}; i=int(pc); print (pc-i<0.5)?i:i+1 }")

    # Calculate the success percentage
    success_percentage=$(awk "BEGIN { pc=100*${successful_jobs}/${completed_jobs}; i=int(pc); print (pc-i<0.5)?i:i+1 }")

    echo -ne "Processing: ${percent_completed}% | Successful: $successful_jobs | Failed: $failed_jobs | Success Percentage: ${success_percentage}%\r"
done

# Display 100% when all jobs are done
echo -ne "Processing: 100% | Successful: $successful_jobs | Failed: $failed_jobs | Success Percentage: ${success_percentage}%\n"
