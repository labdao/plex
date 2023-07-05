# Run Gnina on PLEX 

Gnina is a powerful tool for protein-ligand docking, and now you can run it seamlessly on PLEX! This guide will take you through the steps of running Gnina in a highly reproducible manner on PLEX's decentralized compute network.

## Prerequisites

Ensure you have installed PLEX. If not, follow the [PLEX installation guide](https://docs.labdao.xyz/getting-started/install-plex).

## Available Tools

In this directory, we have two different tool configurations for Gnina:

1. **gnina-redocking.json**: Used for protein-ligand redocking. This configuration takes as input a protein file and a best-docked small molecule file.

2. **gnina.json**: Used for protein-ligand docking. This configuration takes as input a protein file and a small molecule file.

Each tool has several parameters that can be customized, including `exhaustiveness` and `cnn_scoring`.

## How to Run

1. Ensure you're in the correct directory:

   ```
   cd tools/gnina
   ```

2. Submit a PLEX job. Here's an example of running the `gnina.json` tool:

   ```
   plex create -t gnina.json -i {input_directory} --autoRun=True
   ```
   
   Replace `{input_directory}` with the path to your input directory. The input directory should contain the necessary input files (for `gnina.json`, a protein file in .pdb format and a small molecule file in .sdf or .mol2 format).

The PLEX client will take care of the rest: fetching the correct Docker image, running the job on the decentralized network, and storing the results. You can check the status of your jobs and retrieve results using the PLEX client.

## Outputs

The output of each tool is a docked and scored .sdf file, which is stored and can be retrieved using the PLEX client.

## Getting Help

For more information on using PLEX, visit the [PLEX documentation](https://docs.labdao.xyz/). If you encounter any issues, feel free to raise an issue on this repository or reach out to us through our [Discord](https://discord.gg/labdao).

Happy Docking!
