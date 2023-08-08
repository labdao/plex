---
title: Folding proteins with ColabFold
sidebar_label: Protein Folding
sidebar_position: 2
---

import OpenInColab from '../../src/components/OpenInColab.js';

<OpenInColab link="https://colab.research.google.com/drive/1312M2VOx_YpTFgy60ZYChgR9h3a7aorr?usp=sharing"></OpenInColab>

## Protein folding in silico

In this tutorial we perform protein folding with **plex**.

There are multiple reasons we believe plex is a new standard for computational biology 🧫:
1. with a simple python interface, running containerised tools with your data is only a few commands away
2. the infrastructure of the compute network is fully open source - use the public network or work with us to set up your own node
3. every event on the compute network is tracked - no more results are lost in an interactive compute session. You can base your decisions and publications on fully reproducible results.
4. we made adding new tools to the network as easy as possible - moving your favorite tool to plex is one JSON document away.

In this tutorial, we'll walk through an example of how to use plex to predict a protein's 3D structure using [ColabFold](https://www.nature.com/articles/s41592-022-01488-1). We will use the sequence of the Streptavidin protein for this demo.

We will also walk through the process of minting a ProofOfScience NFT. These tokens represent on-chain, verifiable records of the compute job and its input/output data. This enables reproducible scientific results.

![protein-folding-graphic](../../static/img/protein-folding-graphic.png)

## Install plex


```python
!pip install PlexLabExchange
```

    Collecting PlexLabExchange
      Downloading PlexLabExchange-0.8.20-py3-none-manylinux2014_x86_64.whl (26.9 MB)
    [2K     [90m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━[0m [32m26.9/26.9 MB[0m [31m16.6 MB/s[0m eta [36m0:00:00[0m
    [?25hInstalling collected packages: PlexLabExchange
    Successfully installed PlexLabExchange-0.8.20


Then, create a directory where we can save our project files.


```python
import os

cwd = os.getcwd()
!mkdir project

dir_path = f"{cwd}/project"
```

## Download protein sequence

We'll download a `.fasta` file containing the sequence of the protein we want to fold. Here, we're using the sequence of Streptavidin.




```python
!wget https://rest.uniprot.org/uniprotkb/P22629.fasta -O {dir_path}/P22629.fasta # Streptavidin
```

    --2023-08-08 18:49:21--  https://rest.uniprot.org/uniprotkb/P22629.fasta
    Resolving rest.uniprot.org (rest.uniprot.org)... 193.62.193.81
    Connecting to rest.uniprot.org (rest.uniprot.org)|193.62.193.81|:443... connected.
    HTTP request sent, awaiting response... 200 OK
    Length: unspecified [text/plain]
    Saving to: ‘/content/project/P22629.fasta’
    
    /content/project/P2     [ <=>                ]     264  --.-KB/s    in 0s      
    
    2023-08-08 18:49:21 (157 MB/s) - ‘/content/project/P22629.fasta’ saved [264]
    


## Fold the protein

With the sequence downloaded, we can now use ColabFold to fold the protein.




```python
from plex import CoreTools, plex_init

fasta_local_filepaths = [f"{dir_path}/P22629.fasta"]

initial_io_cid = plex_init(
    CoreTools.COLABFOLD_MINI.value,
    sequence=fasta_local_filepaths
)
```

    plex init -t QmcRH74qfqDBJFku3mEDGxkAf6CSpaHTpdbe1pMkHnbcZD -i {"sequence": ["/content/project/P22629.fasta"]} --scatteringMethod=dotProduct
    Plex version (v0.8.4) up to date.
    Pinned IO JSON CID: QmZgLQypfjvK9kTsqLXwbNRiFifEU5CC7eduWWPbminybi


This code initiates the folding process. We'll need to run it to complete the operation.


```python
from plex import plex_run

completed_io_cid, completed_io_filepath = plex_run(initial_io_cid, dir_path)
```

    Plex version (v0.8.4) up to date.
    Created working directory:  /content/project/9102a179-ac65-4823-9a03-93766ea32671
    Initialized IO file at:  /content/project/9102a179-ac65-4823-9a03-93766ea32671/io.json
    Processing IO Entries
    Starting to process IO entry 0 
    Job running...
    Bacalhau job id: 271f4b64-cb2d-4be6-86af-ed16186e69e0 
    
    Computing default go-libp2p Resource Manager limits based on:
        - 'Swarm.ResourceMgr.MaxMemory': "6.8 GB"
        - 'Swarm.ResourceMgr.MaxFileDescriptors': 524288
    
    Applying any user-supplied overrides on top.
    Run 'ipfs swarm limit all' to see the resulting limits.
    
    Success processing IO entry 0 
    Finished processing, results written to /content/project/9102a179-ac65-4823-9a03-93766ea32671/io.json
    Completed IO JSON CID: QmdnjMsUar6nTqGwgjCwN1Fyjaan4i3zyht9SE9L235YRm
    2023/08/08 18:51:17 failed to sufficiently increase receive buffer size (was: 208 kiB, wanted: 2048 kiB, got: 416 kiB). See https://github.com/quic-go/quic-go/wiki/UDP-Receive-Buffer-Size for details.


After the job is complete, we can retrieve and view the results. The state of each object is written in a JSON object. Every file has a unique content-address.




```python
import json

with open(completed_io_filepath, 'r') as f:
  data = json.load(f)
  pretty_data = json.dumps(data, indent=4, sort_keys=True)
  print(pretty_data)
```

    [
        {
            "errMsg": "",
            "inputs": {
                "sequence": {
                    "class": "File",
                    "filepath": "P22629.fasta",
                    "ipfs": "QmR3TRtG1EWszHJTpZWZut6VFqzBPWT5KYVJvaMdXFLWXn"
                }
            },
            "outputs": {
                "all_folded_proteins": {
                    "class": "Array",
                    "files": [
                        {
                            "class": "File",
                            "filepath": "P22629_unrelaxed_rank_1_model_1.pdb",
                            "ipfs": "QmXZHhB7qP1tnJNyR2TeH7m4gB1R5UF84SzvK94eYB9qdL"
                        },
                        {
                            "class": "File",
                            "filepath": "P22629_unrelaxed_rank_2_model_4.pdb",
                            "ipfs": "QmPWGR36mbm5qptniHxd5KjUQKVn8EFMc57DMJzwcetNnU"
                        },
                        {
                            "class": "File",
                            "filepath": "P22629_unrelaxed_rank_3_model_3.pdb",
                            "ipfs": "QmXQ1F8xD3TP1qDvU1HDhpuR5JDZvxv1G2udJSdTsimKvH"
                        },
                        {
                            "class": "File",
                            "filepath": "P22629_unrelaxed_rank_4_model_2.pdb",
                            "ipfs": "QmV4TZJyWbu4CcmLTvD6nKM8YpzDK4fBsiiA3KQkHjW1RG"
                        },
                        {
                            "class": "File",
                            "filepath": "P22629_unrelaxed_rank_5_model_5.pdb",
                            "ipfs": "QmVHT7nQzmNkxDJsRTJPAFqwqhqEgmD3QBGZpUPneogVqX"
                        }
                    ]
                },
                "best_folded_protein": {
                    "class": "File",
                    "filepath": "P22629_unrelaxed_rank_1_model_1.pdb",
                    "ipfs": "QmTxVHTSUr8kLa9W8yM7KUNth2pNn8m3x6M18x8yiaV2SU"
                }
            },
            "state": "completed",
            "tool": {
                "ipfs": "QmcRH74qfqDBJFku3mEDGxkAf6CSpaHTpdbe1pMkHnbcZD",
                "name": "colabfold-mini"
            }
        }
    ]


The results can also be viewed using an IPFS gateway. Below, the state of the IO JSON is read using the ipfs.io gateway.

**Note:** Depending on how long it takes for the results to propagate to the ipfs.io nodes, the data may not be available immediately. The results can also be viewed on IPFS Desktop or by accessing IPFS through the Brave browser (ipfs://completed_io_cid)


```python
print(f"View this result on IPFS: https://ipfs.io/ipfs/{completed_io_cid}")
```

    View this result on IPFS: https://ipfs.io/ipfs/QmdnjMsUar6nTqGwgjCwN1Fyjaan4i3zyht9SE9L235YRm


## Visualization and NFT minting

For visualization and NFT minting steps, please visit the Colab notebook below.

<OpenInColab link="https://colab.research.google.com/drive/1312M2VOx_YpTFgy60ZYChgR9h3a7aorr?usp=sharing"></OpenInColab>
