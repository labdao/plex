# Ganglia

gan·gli·a [gang-glee-uh] : a decentralized "mini-brain" in an octopus or more generally a group of neurons in the peripheral nervous system.

This repo contains infrastructure code for running nodes in LabDAO's decentralized compute system.

## Quick install

Tested on Ubuntu 20.04 LTS with Nvidia T4, V100, and A10 GPUs (AWS G4, P3, and G5 instance types)

```
curl -sL https://raw.githubusercontent.com/labdao/ganglia/main/install.sh | bash && newgrp docker
```

## Development

### Run example
```
cd /home/ubuntu
git clone https://github.com/labdao/diffdock.git
git clone https://github.com/labdao/ganglia

cd ./ganglia
python3 client.py

# the json in this argument is from the client.py output
python3 process.py '{"container_id": "ghcr.io/labdao/diffdock:main", "debug_logs": true, "short_args": {"v": "/home/ubuntu/diffdock:/diffdock"}, "long_args": {"gpus": "all"}, "cmd": "/bin/bash -c \"python datasets/esm_embedding_preparation.py --protein_path test/test.pdb --out_file data/prepared_for_esm.fasta && HOME=esm/model_weights python esm/scripts/extract.py esm2_t33_650M_UR50D data/prepared_for_esm.fasta data/esm2_output --repr_layers 33 --include per_tok && python -m inference --protein_path test/test.pdb --ligand test/test.sdf --out_dir /outputs --inference_steps 20 --samples_per_complex 40 --batch_size 10 --actual_steps 18 --no_final_step_noise\""}'
```

### Run unittests
```
python3 -m unittest
```

### Lint
```
pip install black
python -m black --preview ./
```
