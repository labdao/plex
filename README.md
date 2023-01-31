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
Set-up directory and run server
```
git clone https://github.com/labdao/ganglia
cd ./ganglia
./ipfs.sh  # TODO move to install script
pip install -r requirements.txt 
cd ./server
python3 server.py
```

Run client in new tab
````
python3

import asyncio
from ganglia import generate_diffdock_instructions, run_with_socket

asyncio.run(run_with_socket(generate_diffdock_instructions()))
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
