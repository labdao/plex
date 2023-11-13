# plex 🧫×🧬→💊
⚡ **Run highly reproducible scientific applications on top of a decentralised compute and storage network.** ⚡


<p align="left">
    <a href="https://github.com/labdao/plex/blob/main/LICENSE.md" alt="License">
        <img src="https://img.shields.io/badge/license-MIT-green" />
    </a>
    <a href="https://github.com/labdao/plex/releases/" alt="Release">
        <img src="https://img.shields.io/github/v/release/labdao/plex?display_name=tag" />
    </a>
    <a href="https://github.com/labdao/plex/pulse" alt="Activity">
        <img src="https://img.shields.io/github/commit-activity/m/labdao/plex" />
    </a>
    <a href="https://img.shields.io/github/downloads/labdao/plex/total">
        <img src="https://img.shields.io/github/downloads/labdao/plex/total" alt="total download">
    </a>
    <a href="https://labdao.xyz/">
        <img alt="LabDAO website" src="https://img.shields.io/badge/website-labdao.xyz-red">
    </a>
    <a href="https://twitter.com/intent/follow?screen_name=lab_dao">
        <img src="https://img.shields.io/twitter/follow/lab_dao?style=social&logo=twitter" alt="follow on Twitter">
    </a>
    <a href="https://discord.gg/labdao" alt="Discord">
        <img src="https://dcbadge.vercel.app/api/server/labdao?compact=true&style=flat-square" />
    </a>
</p>


Plex is a simple client for distributed computation.
* 🌎 **Build once, run anywhere:** Plex is using distributed compute and storage to run containers on a public network. Need GPUs? We got you covered.
* 🔍 **Content-addressed by default:** Every file processed by plex has a deterministic address based on its content. Keep track of your files and always share the right results with other scientists.
* 🪙 **Ownernship tracking built-in** Every compute event on plex is mintable as an on-chain token that grants the holder rights over the newly generated data.
* 🔗 **Strictly composable:** Every tool in plex has declared inputs and outputs. Plugging together tools by other authors should be easy.

Plex is based on [Bacalhau](https://www.bacalhau.org/), [IPFS](https://ipfs.tech/), and inspired by the [Common Workflow Language](https://www.commonwl.org/user_guide/introduction/quick-start.html).

## 🐍 Python pip package (Python 3.8+)

1. Install plex with pip
```
pip install PlexLabExchange
```

2. Run plex example in a Python file, notebook or REPL
```
from plex import plex_run

io_json_cid, io_json_local_filepath = plex_run('QmWdKXmSz1p3zGfHmwBb5FHCS7skc4ryEA97pPVxJCT5Wx')
```
## 🚀 Plex CLI in one minute

1 . Install the client

Mac/Linux users open terminal and run
```
source <(curl -sSL https://raw.githubusercontent.com/labdao/plex/main/install.sh)
```

Windows users open terminal as an adminstrator and run
```
Invoke-Expression (Invoke-WebRequest -Uri "https://raw.githubusercontent.com/labdao/plex/main/install.ps1" -UseBasicParsing).Content
```

2. Submit an example plex job
```
./plex init -t tools/equibind.json -i '{"protein": ["testdata/binding/abl/7n9g.pdb"], "small_molecule": ["testdata/binding/abl/ZINC000003986735.sdf"]}' --scatteringMethod=dotProduct --autoRun=true
```

![Getting Started](./readme-getting-started-2x.gif)

3. [Read the docs](https://docs.labdao.xyz/) to learn how to use plex with your own data and tools

4. Request Access to our VIP Jupyter Hub Enviroment and NFT Testnet Minting.
[VIP Beta Access Form](https://try.labdao.xyz)


## 💡 Use-Cases
* 🧬 run plex to [fold proteins](https://docs.labdao.xyz/tutorials/protein-folding)
* 💊 run plex to run [small molecule docking](https://docs.labdao.xyz/tutorials/small-molecule-binding)
* 🐋 configure your containerised tool to run on plex

## 🧑‍💻 Developer Guide

### Building plex from source

```
git clone https://github.com/labdao/plex
cd plex
go build
```

### Running web app locally

# Setup

## Running the Complete Stack Locally
We have `docker-compose` files available to bring up the stack locally.

Note:
* Only `amd64` architecture is currently supported.
* New docker installation include docker compose, older installations required you install docker-compose separately and run `docker-compose up -d`

```
# Optionally, build in parallel before running
docker compose build --parallel

# Build and bring up stack
docker compose up -d --wait
```

To run `plex` cli against local environment simply set `BACALHAU_API_HOST=127.0.0.1`

# Running with private IPFS

```
docker compose -f docker-compose.yml -f docker-compose.private.yml up -d --wait
```

Once environment is up, set `BACALHAU_API_HOST=127.0.0.1` environment variable to point `plex` cli to the local environment.


## Full Stack Development
* Define necessary env variables
```
NEXT_PUBLIC_BACKEND_URL=http://localhost:8080
FRONTEND_URL=http://localhost:3000
POSTGRES_PASSWORD=labdao
POSTGRES_USER=labdao
POSTGRES_DB=backend
POSTGRES_HOST=localhost
```

Bring up the database only
```
docker compose up -d dbbackend --wait
```

# Start the Backend Go App
```
go run main.go web
```

# Start the Frontend React App
```
$ cd /frontend
$ npm install
$ npm run dev
```

## Frontend Development Only
Instead of running the backend locally, you can point the frontend to our cloud deployment.
```
$ export NEXT_PUBLIC_BACKEND_URL=https://api.prod.labdao.xyz
$ export NEXT_PUBLIC_IPFS_GATEWAY_ENDPOINT=http://bacalhau.prod.labdao.xyz:8080/ipfs/
$ cd frontend
$ npm install
$ npm run dev
```

### Running a compute node
This is a script for setting up a compute instance to run LabDAO jobs. Requires linux OS with Nvidia GPU.

Tested on Ubuntu 20.04 LTS with Nvidia T4, V100, and A10 GPUs (AWS G4, P3, and G5 instance types)

The install script sets up Docker, Nvidia Drivers, Nvidia Container Toolkit, and IPFS
```
curl -sL https://raw.githubusercontent.com/labdao/plex/main/scripts/provide-compute.sh | bash && newgrp docker
```

After the script run the following command in a separate terminal to start a Bacalhau server to accept jobs.
```
ipfs daemon
```

Once the daemon is running, configure the Bacalhau node based on the addresses used by the IPFS node.
```
ipfs id

# copy the ip4 tcp output and change port 4001 to 5001 then export
export IPFS_CONNECT=/ip4/127.0.0.1/tcp/5001/p2p/<your id goes here>

# example: export IPFS_CONNECT=/ip4/127.0.0.1/tcp/5001/p2p/12D3KooWPH1BpPfNXwkf778GMP2H5z7pwjKVQFnA5NS3DngU7pxG

LOG_LEVEL=debug bacalhau serve --job-selection-accept-networked --limit-total-gpu 1 --limit-total-memory 12gb --ipfs-connect $IPFS_CONNECT
```

To download large bacalhau results the below command may need ran
```
sudo sysctl -w net.core.rmem_max=2500000
```

## 💁 Contributing
PRs are welcome! Please consider our [Contribute Guidelines](https://docs.labdao.xyz/about-us/contributer_policy) when joining.

From time to time, we also post ```help-wanted``` bounty issues - please consider our [Bounty Policy](https://docs.labdao.xyz/about-us/bounty_policy) when engaging with LabDAO.
