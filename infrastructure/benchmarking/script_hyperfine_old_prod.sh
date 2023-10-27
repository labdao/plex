#!/bin/bash

export BACALHAU_API_HOST=bacalhau.labdao.xyz

# Doing 10 runs of RFdiffusion colabdesign with small,medium and large sizes
hyperfine --warmup 1 --min-runs 10 --show-output --export-csv /tmp/rfdiffusion-colabdesign_${BACALHAU_API_HOST}.csv -L size small,medium,large -n "rfdiffusion-colabdesign_{size}" "BACALHAU_API_HOST=$BACALHAU_API_HOST plex_0.10.4 init -t tools/colabdesign/_colabdesign-dev.json -i '{\"protein\": [\"tools/colabdesign/6vja_stripped.pdb\"], \"config\": [\"tools/colabdesign/{size}-config.yaml\"]}' --scatteringMethod=dotProduct --autoRun=true -a test -a benchmarking"

# Doing 10 runs of Equibind
hyperfine --warmup 1 --min-runs 10 --show-output --export-csv /tmp/equibind_${BACALHAU_API_HOST}.csv -n 'equibind' "BACALHAU_API_HOST=$BACALHAU_API_HOST plex_0.10.4 init -t tools/equibind.json -i '{\"protein\": [\"testdata/binding/6d08_protein_processed.pdb\"], \"small_molecule\": [\"testdata/binding/6d08_ligand.sdf\"]}' --scatteringMethod=dotProduct --autoRun=true -a test -a benchmarking"

# Doing 10 runs of diffdock
hyperfine --warmup 1 --min-runs 10 --show-output --export-csv /tmp/diffdock_${BACALHAU_API_HOST}.csv -n 'diffdock' "BACALHAU_API_HOST=$BACALHAU_API_HOST plex_0.10.4 init -t tools/diffdock.json -i '{\"protein\": [\"testdata/binding/6d08_protein_processed.pdb\"], \"small_molecule\": [\"testdata/binding/6d08_ligand.sdf\"]}' --scatteringMethod=dotProduct --autoRun=true -a test -a benchmarking"

# Doing 10 runs of each colabfold with specific input file
hyperfine --warmup 1 --min-runs 10 --show-output --export-csv /tmp/colabfold_${BACALHAU_API_HOST}.csv -L input histatin-3,tax1-binding_protein_3,c-reactive_protein,gap_junction_protein,phospholipid_glycerol_acyltransferase_domain-containing_protein,rims-binding_protein_3b,dna_helicase,hect-type_e3_ubiquitin_transferase -n "colabfold-{input}" "BACALHAU_API_HOST=$BACALHAU_API_HOST plex init -t tools/colabfold-mini.json -i {\"sequence\": [\"testdata/folding/{input}.fasta\"]} --scatteringMethod=dotProduct --autoRun=true -a test -a benchmarking"
