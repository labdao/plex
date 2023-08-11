---
title: Input / Output (IO)
sidebar_position: 4
sidebar_label: Input / Output (IO)
---

Plex employs a streamlined approach to input and output data management, facilitating consistency and transparency throughout the computation process.

Plex begins its IO process with [`plex_init`](../reference/python.md), which creates an `io.json` file. This file serves as the cornerstone of instruction for the [Bacalhau](https://docs.bacalhau.org/) compute cluster, dictating the parameters and expected outputs for each computational job.

Key components of the initialized `io.json`

* **Input Data:** lists the provided input files, detailing their filename and corresponding CID
* **Output Data Placeholder:** lays out the expected outputs, as defined by the tool config
* **Tool Information:** indicates the computational tool to be used, along with the CID of its config
* **Job State:** initially set to `created`, it tracks the job's progression
* **Bacalhau Job ID Placeholder:** reserved for the unique job identifier once submitted to the Bacalhau compute cluster

## Initialized `io.json`

```json
[
  {
    "outputs": {
      "best_docked_small_molecule": {
        "class": "File",
        "filepath": "",
        "ipfs": ""  
      },
      "protein": {
        "class": "File",
        "filepath": "",
        "ipfs": ""
      }
    },
    "tool": {
      "name": "equibind",
      "ipfs": "QmZ2HarAgwZGjc3LBx9mWNwAQkPWiHMignqKup1ckp8NhB"
    },
    "inputs": {
      "protein": {
        "class": "File",
        "filepath": "6d08_protein_processed.pdb",
        "ipfs": "QmeTreLhxMmBaRqHemJcStvdyHZThdzi4gTmvTyY1igeCk"
      },
      "small_molecule": {
        "class": "File",
        "filepath": "6d08_ligand.sdf",
        "ipfs": "QmPErdymxLwpXcEHnWXYqEVHvRBVnh7kr3Uu5DNt2Y8wMR"
      }
    },
    "state": "created",
    "errMsg": "",
    "bacalhauJobId": ""
  }
]
```

## Execution with `plex_run`

The action commences with [`plex_run`](../reference/python.md). Upon its call, the computational job(s) outlined in the `io.json` are dispatched to the Bacalhau cluster for processing.

As the computations unfold and conclude, the `io.json` undergoes real-time updates

* **Output Data:** once a job completes, the `io.json` populates with the resultant data and its CID
* **Bacalhau Job ID:** the unique identifier for the job is added, facilitating traceability; useful in cases when a job fails to run
* **Updated Job State:** reflects the final status of the job, transitioning to `completed` if successful

## Completed `io.json`

```json
[
  {
    "outputs": {
      "best_docked_small_molecule": {
        "class": "File",
        "filepath": "6d08_protein_processed_6d08_ligand_docked.sdf",
        "ipfs": "QmWdzgrt5wtUJPyCrcKycU3voKGmT59FZXMasuaa1XCbkk"
      },
      "protein": {
        "class": "File",
        "filepath": "6d08_protein_processed.pdb",
        "ipfs": "QmeTreLhxMmBaRqHemJcStvdyHZThdzi4gTmvTyY1igeCk"
      }
    },
    "tool": {
      "name": "equibind",
      "ipfs": "QmZ2HarAgwZGjc3LBx9mWNwAQkPWiHMignqKup1ckp8NhB"
    },
    "inputs": {
      "protein": {
        "class": "File",
        "filepath": "6d08_protein_processed.pdb",
        "ipfs": "QmeTreLhxMmBaRqHemJcStvdyHZThdzi4gTmvTyY1igeCk"
      },
      "small_molecule": {
        "class": "File",
        "filepath": "6d08_ligand.sdf",
        "ipfs": "QmPErdymxLwpXcEHnWXYqEVHvRBVnh7kr3Uu5DNt2Y8wMR"
      }
    },
    "state": "completed",
    "errMsg": "",
    "bacalhauJobId": "7a01e92a-877e-4d1b-ba91-9effec6f170e"
  }
]
```