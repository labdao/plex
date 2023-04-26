---
title: Tools
sidebar_position: 2
---

Since the release of Alphafold there has been a variety of related models released for protein folding. We have gone through the literature and open source code repositories to provide a selection of tools that are ready for you to run on PLEX. 

:::note
All models we provide are research-grade software and are provided "as-is". We make use of existing, often academic, contributions. Please give credit to the creators of open-source work. **We are standing on the shoulder of giants.**
:::

:::info
We have prepared a set of configurations for you that we list below. We are working on ways to have more control over individual parameters.
:::


## Mini: Colabfold
[Colabfold](https://github.com/sokrypton/ColabFold) is an implementation of Alphafold that uses a multiple sequence alignment (MSA) Server, MMSeq2, instead of a local database to make using [Alphafold](https://github.com/deepmind/alphafold) more lightweight. The "mini" configuration runs a shallow MSA, performs one recycling and uses available templates to make a prediction. This is best used if you want to predict a protein structure very fast. 

````
./plex -tool colabfold-mini -input-dir testdata/folding
````

:::note
Mirdita M, Schütze K, Moriwaki Y, Heo L, Ovchinnikov S and Steinegger M. ColabFold: Making protein folding accessible to all.
Nature Methods (2022) doi: 10.1038/s41592-022-01488-1
:::

## Standard: Colabfold
The "standard" configuration runs a full MSA, performs three recycling rounds and uses available templates to make a prediction. It runs this prediction 5 times with different randomness seeds. This is best used if you want to predict a state of the art structure and draw from a distribution of potential conformational substates.

````
./plex -tool colabfold-standard -input-dir testdata/folding
````

:::note
Mirdita M, Schütze K, Moriwaki Y, Heo L, Ovchinnikov S and Steinegger M. ColabFold: Making protein folding accessible to all.
Nature Methods (2022) doi: 10.1038/s41592-022-01488-1
:::

## Large: Colabfold
The "large" configuration runs just like the standard configuration, but includes a GPU-accelerated relaxation step using Amber. It returns 25 predictions. This is best used when a lot of ressources are available and you want to predict a state of the art structure while drawing from a larger distribution of potential conformational substates.

````
./plex -tool colabfold-large -input-dir testdata/folding
````

:::note
Mirdita M, Schütze K, Moriwaki Y, Heo L, Ovchinnikov S and Steinegger M. ColabFold: Making protein folding accessible to all.
Nature Methods (2022) doi: 10.1038/s41592-022-01488-1
:::

