#!/bin/bash
# First argument is the file to read from
FILE=$1

# Second argument is the pattern to match
PATTERN=$2

# Third argument is the output zip name
OUTPUT=$3


# Subset the index and download files
egrep "$PATTERN" "$FILE" | while read -r line; do
    echo "Downloading: $line"
    eval "$line"
done

# Decompress tgz files directly in place
find zinc22 -type f -name '*.pdbqt.tgz' -exec tar -xzvf {} -C zinc22/zinc-22d \;

# Search for all '.pdbqt' files and add them to a zip file
find zinc22 -type f -name '*.pdbqt' -print | zip $OUTPUT.zip -@

# cleaning up
# rm -r zinc22