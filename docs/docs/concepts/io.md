---
title: Input / Output (IO)
sidebar_position: 4
sidebar_label: Input / Output (IO)
---



### Initialized IO

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

### Completed IO

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