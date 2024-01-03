#!/bin/bash

# Absolute path to the fasta file
fasta_file="/inputs/all_seqs_to_shorten/testseq_to_shorten.fasta"

# Debug: Echo the file path and check if it exists
echo "File path: $fasta_file"
if [ -f "$fasta_file" ]; then
    echo "File exists."
else
    echo "File does not exist."
    exit 1
fi

# Initialize a variable to hold the sequence
sequence=""

# Debug: Print a message when starting to read the file
echo "Starting to read the file..."

# Read the sequences_to_shorten.fasta file line by line
while IFS= read -r line || [ -n "$line" ]; do
    # Debug: Print each line read
    echo "Read line: $line"

    # Check if the line starts with '>'
    if [[ $line == ">"* ]]; then
        # If there's a sequence in memory, write it to input_sequence.fasta
        if [ -n "$sequence" ]; then
            # Write the sequence to input_sequence.fasta with header >1
            echo -e ">1\n$sequence" > input_sequence.fasta
            # Run the python3 script
            python3 main.py
        fi
        # Reset the sequence for the next entry
        sequence=""
    else
        # Append the line to the sequence
        sequence+=$line
    fi
done < "$fasta_file"

# Process the last sequence
if [ -n "$sequence" ]; then
    echo -e ">1\n$sequence" > input_sequence.fasta
    python3 main.py
fi