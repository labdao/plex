# Ganglia

gan·gli·a [gang-glee-uh] : a decentralized "mini-brain" in an octopus or more generally a group of neurons in the peripheral nervous system.

This repo contains infrastructure code for running nodes in LabDAO's decentralized compute system.

## Running the client
Install script coming soon. For now git cloning the repo is required.

1) Install [GoLang](https://go.dev/doc/install)

2) Run the following commands
```
git clone https://github.com/labdao/ganglia.git
cd ganglia/plex
go build
export WEB3STORAGE_TOKEN=<your_token>
```

3) Run the diffdock example
```
./plex -app diffdock -input-dir ./testdata/pdbbind_processed_size1

# Running the outputted Bacalhau command takes ~10 min per complex, this example has 1 complex
```

4) Run diffdock with your own inputs
```
./plex -app diffdock -input-dir <path to dir on your computer>

# plex will automagically run diffdock on every protein and ligand file found in the directory
```

## Running a node
This is a script for setting up a compute instance to run LabDAO jobs. Requires linux OS with Nvidia GPU.

Tested on Ubuntu 20.04 LTS with Nvidia T4, V100, and A10 GPUs (AWS G4, P3, and G5 instance types)

The install script sets up Docker, Nvidia Drivers, Nvidia Container Toolkit, and IPFS
```
curl -sL https://raw.githubusercontent.com/labdao/ganglia/main/install.sh | bash && newgrp docker
```

After the script run the following commands to start a Bacalhau server to accept jobs
```
ipfs init

# copy the ip4 tcp output and change port 4001 to 5001 then export
export IPFS_CONNECT=/ip4/127.0.0.1/tcp/5001/p2p/<your id goes here>

# example: export IPFS_CONNECT=/ip4/127.0.0.1/tcp/5001/p2p/12D3KooWPH1BpPfNXwkf778GMP2H5z7pwjKVQFnA5NS3DngU7pxG

LOG_LEVEL=debug bacalhau serve --job-selection-accept-networked --limit-total-gpu 1 --limit-total-memory 12gb --ipfs-connect $IPFS_CONNECT
```

## Notes
To download large bacalhau results the below command may need ran 
```
sudo sysctl -w net.core.rmem_max=2500000
```
