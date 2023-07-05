#!/bin/bash
# This script works from an index of the ZINC22 database and downloads the pdbqt files of a specified tranche
# The index is a text file with one line per file to download
# The pdbqt files are downloaded to /outputs/pdbqt
# The index of pdbqt files is written to /outputs/index.txt and can be used as input for docking tools

# First argument is the file to read from
FILE=$1

# Second argument is the pattern to match
PATTERN=$2

# Third argument is the output directory
OUTPUT=$3

# subset the index and download files
grep "$PATTERN" "$FILE" | while read -r line; do
    echo "Downloading: $line"
    eval "$line"
done

# pull out the pdbqt files 
mkdir -p $OUTPUT/tgz
#mkdir -p $OUTPUT/$PATTERN
mv zinc22/*/*/*/*/*.pdbqt.tgz $OUTPUT/tgz

# decompress pdbqt files and create an index
echo "" > $OUTPUT/index.txt

for file in $OUTPUT/tgz/*.tgz; do
  # Decompress the file
  tar -xzvf "$file" -C $OUTPUT/
done

# Find all decompressed files and append their absolute paths to the log
cd $OUTPUT && find -type f -name '*.pdbqt' -print >> $OUTPUT/index.txt

# Compress the result into a .zip file
cd $OUTPUT && zip -r $PATTERN.zip *
