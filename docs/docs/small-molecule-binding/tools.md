---
title: Tools
sidebar_position: 3
---

There is an overwhelming variety of tools that have been used to predict the binding affinity between a small molecule and a protein. We have gone through the literature and open source code repositories to provide a selection of tools that are ready for you to run on PLEX. 

:::note
All models we provide are research-grade software and are provided "as-is". No model for this task has yet been demonstrated to generalise well enough to be an alternative to laboratory experiments. We make use of existing, often academic, contributions. Please give credit to the creators of open-source work. **We are standing on the shoulder of giants.**
:::

:::warning
At this point in time we are focused on docking. Stay tuned for integrated scoring functions.
:::


## Mini: Equibind
[Equibind](https://github.com/HannesStark/EquiBind) is a very fast, machine learning-based docking tool. The model is less accurate than baseline methods, but orders of magnitude faster.

````
./plex -tool equibind -input-dir testdata/binding/abl
````

:::note
Stärk, H., Ganea, O.-E., Pattanaik, L., Barzilay, R., & Jaakkola, T. (2022). EquiBind: Geometric Deep Learning for Drug Binding Structure Prediction. http://arxiv.org/abs/2202.05146
:::


## Base: Gnina (Coming Soon)
[Gnina](https://github.com/gnina/gnina) is a sampling and machine learning-based docking tool. Gnina is an implementation of [Smina](https://sourceforge.net/projects/smina/), which itself is a fork of [Vina](https://vina.scripps.edu/). These tools are considered the current open source baseline.

:::note
A McNutt, P Francoeur, R Aggarwal, T Masuda, R Meli, M Ragoza, J Sunseri, DR Koes. J. (2021). GNINA 1.0: Molecular docking with deep learning https://chemrxiv.org/engage/chemrxiv/article-details/60c753ebbb8c1a1a9d3dc142
:::

## Standard: Diffdock
[Diffdock](https://github.com/gcorso/DiffDock) is a machine learning-based docking tool. Diffdock is reportedly faster and more accurate than existing baseline tools.

````
./plex -tool diffdock -input-dir testdata/binding/abl
````

:::note
Corso, G., Stärk, H., Jing, B., Barzilay, R., & Jaakkola, T. (2022). DiffDock: Diffusion Steps, Twists, and Turns for Molecular Docking. http://arxiv.org/abs/2210.01776
:::
