---
title: Small molecule binding with Equibind
sidebar_label: Small Molecule Binding
sidebar_position: 1
---

import OpenInColab from '../../src/components/OpenInColab.js';

<OpenInColab link="https://colab.research.google.com/drive/15nZrm5k9fMdAHfzpR1g_8TPIz9qgRoys?usp=sharing"></OpenInColab>

## Small molecule docking with plex

In this tutorial we perform small molecule docking with **plex**.

There are multiple reasons we believe plex is a new standard for computational biology 🧫:
1. with a simple python interface, running containerised tools with your data is only a few commands away
2. the infrastructure of the compute network is fully open source - use the public network or work with us to set up your own node
3. every event on the compute network is tracked - no more results are lost in an interactive compute session. You can base your decisions and publications on fully reproducible results.
4. we made adding new tools to the network as easy as possible - moving your favorite tool to PLEX is one JSON document away.

In the following tutorial, we illustrate how plex can be used to conduct small molecule binding studies to explore potential drug interactions with proteins. We demonstrate this with [Equibind](https://hannes-stark.com/assets/EquiBind.pdf).

We will also walk through the process of minting a ProofOfScience NFT. These tokens represent on-chain, verifiable records of the compute job and its input/output data. This enables reproducible scientific results.

![docking-graphic](../../static/img/small-molecule-binding-graphic.png)

## Install plex


```python
!pip install PlexLabExchange
```

    Collecting PlexLabExchange
      Downloading PlexLabExchange-0.8.20-py3-none-manylinux2014_x86_64.whl (26.9 MB)
    [2K     [90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━[0m [32m26.9/26.9 MB[0m [31m20.1 MB/s[0m eta [36m0:00:00[0m
    [?25hInstalling collected packages: PlexLabExchange
    Successfully installed PlexLabExchange-0.8.20


Then, create a directory where we can save our project files.


```python
import os

cwd = os.getcwd()
!mkdir project

dir_path = f"{cwd}/project"
```

## Download small molecule and protein data

We'll download the small molecule `.sdf` and protein `.pdb` we want to dock with Equibind.


```python
# small molecule
!wget https://raw.githubusercontent.com/labdao/plex/main/testdata/binding/abl/ZINC000003986735.sdf -O {dir_path}/ZINC000003986735.sdf
# protein
!wget https://raw.githubusercontent.com/labdao/plex/main/testdata/binding/abl/7n9g.pdb -O {dir_path}/7n9g.pdb
```

    --2023-08-08 18:56:14--  https://raw.githubusercontent.com/labdao/plex/main/testdata/binding/abl/ZINC000003986735.sdf
    Resolving raw.githubusercontent.com (raw.githubusercontent.com)... 185.199.108.133, 185.199.109.133, 185.199.110.133, ...
    Connecting to raw.githubusercontent.com (raw.githubusercontent.com)|185.199.108.133|:443... connected.
    HTTP request sent, awaiting response... 200 OK
    Length: 2967 (2.9K) [text/plain]
    Saving to: ‘/content/project/ZINC000003986735.sdf’
    
    /content/project/ZI 100%[===================>]   2.90K  --.-KB/s    in 0s      
    
    2023-08-08 18:56:14 (47.2 MB/s) - ‘/content/project/ZINC000003986735.sdf’ saved [2967/2967]
    
    --2023-08-08 18:56:14--  https://raw.githubusercontent.com/labdao/plex/main/testdata/binding/abl/7n9g.pdb
    Resolving raw.githubusercontent.com (raw.githubusercontent.com)... 185.199.111.133, 185.199.109.133, 185.199.108.133, ...
    Connecting to raw.githubusercontent.com (raw.githubusercontent.com)|185.199.111.133|:443... connected.
    HTTP request sent, awaiting response... 200 OK
    Length: 580284 (567K) [text/plain]
    Saving to: ‘/content/project/7n9g.pdb’
    
    /content/project/7n 100%[===================>] 566.68K  --.-KB/s    in 0.05s   
    
    2023-08-08 18:56:14 (12.1 MB/s) - ‘/content/project/7n9g.pdb’ saved [580284/580284]
    


## Small molecule docking

With the small molecule and protein files downloaded, we can now use Equibind to run a docking simulation.


```python
from plex import CoreTools, plex_init

protein_path = [f"{dir_path}/7n9g.pdb"]
small_molecule_path = [f"{dir_path}/ZINC000003986735.sdf"]

initial_io_cid = plex_init(
    CoreTools.EQUIBIND.value,
    protein=protein_path,
    small_molecule=small_molecule_path,
)
```

    plex init -t QmZ2HarAgwZGjc3LBx9mWNwAQkPWiHMignqKup1ckp8NhB -i {"protein": ["/content/project/7n9g.pdb"], "small_molecule": ["/content/project/ZINC000003986735.sdf"]} --scatteringMethod=dotProduct
    Plex version (v0.8.4) up to date.
    Pinned IO JSON CID: QmShD7ApeDBUqqy98RuuKdyv8AdmBsvyZqqxSLAEvB9EKP


This code initiates the docking process. We'll need to run it to complete the operation.


```python
from plex import plex_run

completed_io_cid, io_local_filepath = plex_run(initial_io_cid, dir_path)
```

    Plex version (v0.8.4) up to date.
    Created working directory:  /content/project/2e3a8afd-928d-4fb7-a381-fff63c7d51de
    Initialized IO file at:  /content/project/2e3a8afd-928d-4fb7-a381-fff63c7d51de/io.json
    Processing IO Entries
    Starting to process IO entry 0 
    Job running...
    Bacalhau job id: 892bf30d-7f6d-4cc7-a490-c1fa17d82171 
    
    Computing default go-libp2p Resource Manager limits based on:
        - 'Swarm.ResourceMgr.MaxMemory': "6.8 GB"
        - 'Swarm.ResourceMgr.MaxFileDescriptors': 524288
    
    Applying any user-supplied overrides on top.
    Run 'ipfs swarm limit all' to see the resulting limits.
    
    Success processing IO entry 0 
    Finished processing, results written to /content/project/2e3a8afd-928d-4fb7-a381-fff63c7d51de/io.json
    Completed IO JSON CID: QmVG4mT2kkPSb6wzT5QxYZndB5VbKLU8nH2dErZW2zxae6
    2023/08/08 18:56:21 failed to sufficiently increase receive buffer size (was: 208 kiB, wanted: 2048 kiB, got: 416 kiB). See https://github.com/quic-go/quic-go/wiki/UDP-Receive-Buffer-Size for details.


After the job is complete, we can retrieve and view the results. The state of each object is written in a JSON object. Every file has a unique content-address.


```python
import json

with open(io_local_filepath, 'r') as f:
  data = json.load(f)
  pretty_data = json.dumps(data, indent=4, sort_keys=True)
  print(pretty_data)
```

    [
        {
            "errMsg": "",
            "inputs": {
                "protein": {
                    "class": "File",
                    "filepath": "7n9g.pdb",
                    "ipfs": "QmUWCBTqbRaKkPXQ3M14NkUuM4TEwfhVfrqLNoBB7syyyd"
                },
                "small_molecule": {
                    "class": "File",
                    "filepath": "ZINC000003986735.sdf",
                    "ipfs": "QmV6qVzdQLNM6SyEDB3rJ5R5BYJsQwQTn1fjmPzvCCkCYz"
                }
            },
            "outputs": {
                "best_docked_small_molecule": {
                    "class": "File",
                    "filepath": "7n9g_ZINC000003986735_docked.sdf",
                    "ipfs": "QmZdoaKEGtESnLoHFMb9bvqdwXjyUuRK6DbEoYz8PYpZ8W"
                },
                "protein": {
                    "class": "File",
                    "filepath": "7n9g.pdb",
                    "ipfs": "QmUWCBTqbRaKkPXQ3M14NkUuM4TEwfhVfrqLNoBB7syyyd"
                }
            },
            "state": "completed",
            "tool": {
                "ipfs": "QmZ2HarAgwZGjc3LBx9mWNwAQkPWiHMignqKup1ckp8NhB",
                "name": "equibind"
            }
        }
    ]


This output provides us with key information about the small molecule-protein interaction. The "best_docked_small_molecule" represents the most likely interaction between the protein and the small molecule, which can inform subsequent analysis and experiments.

The results can also be viewed using an IPFS gateway. Below, the state of the IO JSON is read using the ipfs.io gateway.

**Note:** Depending on how long it takes for the results to propagate to the ipfs.io nodes, the data may not be available immediately. The results can also be viewed on IPFS Desktop or by accessing IPFS through the Brave browser (ipfs://completed_io_cid)


```python
print(f"View this result on IPFS: https://ipfs.io/ipfs/{completed_io_cid}")
```

    View this result on IPFS: https://ipfs.io/ipfs/QmVG4mT2kkPSb6wzT5QxYZndB5VbKLU8nH2dErZW2zxae6

## Visualization and NFT minting

For visualization and NFT minting steps, please visit the Colab notebook below.

<OpenInColab link="https://colab.research.google.com/drive/15nZrm5k9fMdAHfzpR1g_8TPIz9qgRoys?usp=sharing"></OpenInColab>
