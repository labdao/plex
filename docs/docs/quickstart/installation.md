---
title: Install plex
description: How to install plex
sidebar_label: Install Plex
sidebar_position: 1
---

Plex is a [Python package](https://pypi.org/project/PlexLabExchange/) developed by LabDAO that enables you to seamlessly run computational biology tools. PLEX manages all dependencies and installations and requests compute-time from the LabDAO network, ensuring an effortless experience.

:::note

**Time needed:**
- < 1 minute

**Requirements:**
- Python 3.8 or higher
- [pip](https://pip.pypa.io/en/stable/installation/)

:::

## Installation

To install [plex](https://pypi.org/project/PlexLabExchange/), run the following command:

```
pip install PlexLabExchange
```

If using a Jupyter notebook or Google Colab, you should prefix the command with an exclamation mark:

```
!pip install PlexLabExchange
```

**Congratulations.** Welcome to [DeSci](https://ethereum.org/en/desci/).

## Verification

After installation, ensure plex is working as expected by running one of the following tools:

- [Small Molecule Binding Tool](../tutorials/small-molecule-binding): A quick-run algorithm; complete a job and visualize results within 5 minutes.
- [Protein Folding Tool](../tutorials/protein-folding): Comprehensive guide provided for a step-by-step walkthrough.