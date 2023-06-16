---
title: Available Tools
description: Tools available through the PLEX client
sidebar_label: Available Tools
sidebar_position: 2
---

Listed below are the tools made available through the PLEX client.

| Tool | Category | Description | Example Command |
| -------- | -------- | -------- | --- | 
| [Equibind](https://github.com/labdao/plex/blob/main/tools/equibind.json) | Small Molecule Binding | Docking of small molecules to a protein | `./plex create -t tools/equibind.json -i testdata/binding/abl --autoRun=true` |
| [Diffdock]( https://github.com/labdao/plex/blob/main/tools/diffdock.json) | Small Molecule Binding | Docking of small molecules to a protein | `./plex create -t tools/diffdock.json -i testdata/binding/abl --autoRun=true` |
| [RF Diffusion](https://github.com/labdao/plex/blob/main/tools/rfdiffusion.json) | Protein Design | Design protein binders; generally useful for conditional generation of protein backbones | `./plex create -t tools/rfdiffusion.json -i testdata/design --autoRun=true` |
| [Colabfold](https://github.com/labdao/plex/blob/main/tools/colabfold-mini.json) | Protein Folding | Protein folding prediction | `./plex create -t tools/colabfold-mini.json -i testdata/folding --autoRun=true` |
| [ODDT](https://github.com/labdao/plex/blob/main/tools/oddt.json) | Small Molecule Binding | Scoring of protein-ligand complexes | `./plex create -t tools/oddt.json -i testdata/scoring/abl --autoRun=true` |
| [bam2fastq](https://github.com/labdao/plex/blob/main/tools/bam2fastq.json) |  | Sort BAM by qname and Extract Fasta reads R1 R2 with RG using samtools |  |

If there are any additional tools you would like made available in the PLEX client, please see [how to contribute a tool](/get-involved/how-to-contribute-a-tool).