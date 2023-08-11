---
title: Tools
sidebar_position: 3
sidebar_label: Tools
---

Plex is paving the way for permissionless science by ensuring that computational biology tools are not just available, but also easily accessible for open-source, early-stage drug discovery.

Plex champions open-source initiatives by incorporating state-of-the-art tools into its platform. By leveraging Docker containers, these powerful tools have been streamlined for easy accessibility. This ensures that researchers and scientists can easily tap into the most recent and effective computational biology tools without wading through intricate installation processes.

To further facilitate ease of access and the spirit of open science, all Docker images of these tools are made publicly available. This allows for transparency and easy replication, ensuring that researchers can validate, reproduce, and build upon existing work with confidence.

## Tool Configs

Plex employs **tool configs** as computation instructions. These are sent to our Bacalhau cluster, orchestrating how computations should be carried out. As demonstrated in the provided JSON example, these configs

* Specify the Docker container used
* Detail the input data format, ensuring that the data fed into the tool aligns with its expectations
* Define the output data format, allowing for standardized retrieval and further processing

This approach, reminiscent of the Common Workflow Language ([**CWL**](https://www.commonwl.org/)), ensures consistency, interoperability, and reproducibility across different tools and workflows.

### Colabfold Tool Config

```json
{
    "class": "CommandLineTool",
    "name": "colabfold-mini",
    "description": "Protein folding prediction using Colabfold (mini settings)",
    "baseCommand": ["/bin/bash", "-c"],
    "arguments": [
      "colabfold_batch --templates --max-msa 32:64 --num-recycle $(inputs.recycle.default) /inputs /outputs;"
    ],
    "dockerPull": "public.ecr.aws/p7l9w5o7/colabfold:latest",
    "gpuBool": true,
    "networkBool": true,
    "inputs": {
      "sequence": {
        "type": "File",
        "item": "",
        "glob": ["*.fasta"]
      },
      "recycle": {
        "type": "int",
        "item": "",
        "default": "1"
      }
    },
    "outputs": {
      "best_folded_protein": {
        "type": "File",
        "item": "",
        "glob": ["*rank_1*.pdb"]
      },
      "all_folded_proteins": {
        "type": "Array",
        "item": "File",
        "glob": ["*rank*.pdb"]
      }
    }
}
```