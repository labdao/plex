# Ganglia

gan·gli·a [gang-glee-uh] : a decentralized "mini-brain" in an octopus or more generally a group of neurons in the peripheral nervous system.

This repo contains infrastructure code for running nodes in LabDAO's decentralized compute system.

## Quick install

Tested on Ubuntu 20.04 LTS with Nvidia T4, V100, and A10 GPUs (AWS G4, P3, and G5 instance types)

```
curl -sL https://raw.githubusercontent.com/labdao/ganglia/main/install.sh | bash && newgrp docker
```

### CLI example

```
export WEB3STORAGE_TOKEN=<your_token>
```