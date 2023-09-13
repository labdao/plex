#!/bin/bash

export BACALHAU_SERVE_IPFS_PATH=/tmp/ipfs_path
export BACALHAU_IPFS_SWARM_ADDRESSES=/dns4/bacalhau.demoecs.labdao.xyz/tcp/4001/p2p/12D3KooWBeRoREGr9TD5FUDYAPznVuptieft8gfTCMiMgviTc1Vd
export IPFS_PATH=/tmp/ipfs_path
export BACALHAU_API_HOST=bacalhau.demoecs.labdao.xyz

# Doing 10 runs of RFdiffusion colabdesign with small,medium and large sizes
hyperfine --min-runs 10 --ignore-failure --show-output --export-csv /tmp/rfdiffusion-colabdesign.csv -L size small,medium,large -n "rfdiffusion-colabdesign_{size}" "BACALHAU_SERVE_IPFS_PATH=$BACALHAU_SERVE_IPFS_PATH BACALHAU_IPFS_SWARM_ADDRESSES=$BACALHAU_IPFS_SWARM_ADDRESSES IPFS_PATH=$IPFS_PATH BACALHAU_API_HOST=$BACALHAU_API_HOST plex init -t tools/colabdesign/_colabdesign-dev.json -i '{\"protein\": [\"tools/colabdesign/6vja_stripped.pdb\"], \"config\": [\"tools/colabdesign/{size}-config.yaml\"]}' --scatteringMethod=dotProduct --autoRun=true -a test -a benchmarking"

# Doing 10 runs of Equibind
hyperfine --min-runs 10 --ignore-failure --show-output --export-csv /tmp/equibind.csv -n 'equibind' "BACALHAU_SERVE_IPFS_PATH=$BACALHAU_SERVE_IPFS_PATH BACALHAU_IPFS_SWARM_ADDRESSES=$BACALHAU_IPFS_SWARM_ADDRESSES IPFS_PATH=$IPFS_PATH BACALHAU_API_HOST=$BACALHAU_API_HOST plex init -t tools/equibind.json -i '{\"protein\": [\"testdata/binding/6d08_protein_processed.pdb\"], \"small_molecule\": [\"testdata/binding/6d08_ligand.sdf\"]}' --scatteringMethod=dotProduct --autoRun=true -a test -a benchmarking"

# Doing 10 runs of diffdock
hyperfine --min-runs 10 --ignore-failure --show-output --export-csv /tmp/diffdock.csv -n 'diffdock' "BACALHAU_SERVE_IPFS_PATH=$BACALHAU_SERVE_IPFS_PATH BACALHAU_IPFS_SWARM_ADDRESSES=$BACALHAU_IPFS_SWARM_ADDRESSES IPFS_PATH=$IPFS_PATH BACALHAU_API_HOST=$BACALHAU_API_HOST plex init -t tools/diffdock.json -i '{\"protein\": [\"testdata/binding/6d08_protein_processed.pdb\"], \"small_molecule\": [\"testdata/binding/6d08_ligand.sdf\"]}' --scatteringMethod=dotProduct --autoRun=true -a test -a benchmarking"

# Doing 10 runs of each colabfold with specific input file
hyperfine --min-runs 10 --ignore-failure --show-output --export-csv /tmp/colabfold.csv -L input histatin-3,tax1-binding_protein_3,c-reactive_protein,gap_junction_protein,phospholipid_glycerol_acyltransferase_domain-containing_protein,rims-binding_protein_3b,dna_helicase,hect-type_e3_ubiquitin_transferase -n "colabfold-{input}" "BACALHAU_SERVE_IPFS_PATH=$BACALHAU_SERVE_IPFS_PATH BACALHAU_IPFS_SWARM_ADDRESSES=$BACALHAU_IPFS_SWARM_ADDRESSES IPFS_PATH=$IPFS_PATH BACALHAU_API_HOST=$BACALHAU_API_HOST plex init -t tools/colabfold-mini.json -i {\"sequence\": [\"testdata/folding/{input}.fasta\"]} --scatteringMethod=dotProduct --autoRun=true -a test -a benchmarking"


