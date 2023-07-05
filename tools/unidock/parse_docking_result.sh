#!/bin/bash

# Input directory
input_dir=$1

# Output CSV file
output_file=$2

# Print CSV header
echo "filename,ligand,model,affinity,rmsd_lb,rmsd_ub,inter,intra,unbound" > $output_file

# Scan through each PDBQT file in the directory
for pdbqt_file in $input_dir/*.pdbqt; do
    ligand=$(basename $pdbqt_file .pdbqt)
    awk -v ligand="$ligand" -v file="$pdbqt_file" '
        BEGIN {name="NA"}
        /MODEL/ {model=$2}
        /REMARK VINA RESULT/ {affinity=$4; rmsd_lb=$5; rmsd_ub=$6}
        /REMARK INTER \+ INTRA/ {inter_intra=$4}
        /REMARK INTER/ {inter=$3}
        /REMARK INTRA/ {intra=$3}
        /REMARK UNBOUND/ {unbound=$3}
        /REMARK  Name =/ {name=$4; print file","name","model","affinity","rmsd_lb","rmsd_ub","inter","intra","unbound; name="NA"}
    ' $pdbqt_file >> $output_file
done
