---
title: Available Tools
description: Tools available through the PLEX client
sidebar_label: Available Tools
sidebar_position: 2
---

Listed below are the tools made available through the PLEX client.

| Tool | Category | Description | Example Command |
| -------- | -------- | -------- | --- | 
| [Equibind](https://github.com/labdao/plex/blob/main/tools/equibind.json) | Small Molecule Binding | Docking of small molecules to a protein | `./plex -tool equibind -input-dir testdata/binding/abl` |
| [Diffdock]( https://github.com/labdao/plex/blob/main/tools/diffdock.json) | Small Molecule Binding | Docking of small molecules to a protein | `./plex -tool diffdock -input-dir testdata/binding/abl` |
| [RF Diffusion](https://github.com/labdao/plex/blob/main/tools/rfdiffusion.json) | Protein Design | Design protein binders; generally useful for conditional generation of protein backbones | `./plex -tool rfdiffusion -input-dir testdata/design` |
| [Colabfold](https://github.com/labdao/plex/blob/main/tools/colabfold-mini.json) | Protein Folding | Protein folding prediction | `./plex -tool colabfold-mini -input-dir testdata/folding` |
| [ODDT](https://github.com/labdao/plex/blob/main/tools/oddt.json) | Small Molecule Binding | Scoring of protein-ligand complexes | `./plex -tool oddt -input-dir testdata/scoring/abl` |
| [bam2fastq](https://github.com/labdao/plex/blob/main/tools/bam2fastq.json) |  | Sort BAM by qname and Extract Fasta reads R1 R2 with RG using samtools |  |

If there are any additional tools you would like made available in the PLEX client, please see [how to contribute a tool](/get-involved/how-to-contribute-a-tool).