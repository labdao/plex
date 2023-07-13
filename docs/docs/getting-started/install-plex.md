---
title: Install PLEX
description: How to install LabDAO's PLEX client
sidebar_label: Getting Started
sidebar_position: 1
---

import AsciinemaPlayer from '../../src/components/AsciinemaPlayer.js';

This tutorial will guide you through the steps required to set up PLEX, so you can easily run BioML tools from your computer's command line.

:::note

**What is PLEX?**

PLEX is software that allows you to run computational biology tools using simple commands from your computer. 

PLEX manages all the required dependencies and installations, to make the tools as easy to run as possible. 

When you run a tool, PLEX requests compute-time from members of the LabDAO network, so you don’t have to worry about hardware requirements or setting up the neccesary compute infrastructure.

:::

Let's get started!

---

**Time needed:**
- 5 minutes

**Requirements:**

---

## Install PLEX

### 1. Open your terminal
To install PLEX, first open up the Terminal application.

### 2. Type in the installation commands
Once you have your terminal open, you can download PLEX by copy and pasting this command into your terminal:

**Mac/Linux:**

```
source <(curl -sSL https://raw.githubusercontent.com/labdao/plex/main/install.sh)
```

**Windows:**

```
Invoke-Expression (Invoke-WebRequest -Uri "https://raw.githubusercontent.com/labdao/plex/main/install.ps1" -UseBasicParsing).Content
```

After you have pasted the command into the terminal, press **Enter** on your keyboard to download and install PLEX.

If the installation is successful, you will see a large LabDAO logo appear on your screen, and a confirmation that the Installation is complete. It should look something like this:

<AsciinemaPlayer 
    src="/terminal-recordings/install-plex.cast"
    rows={30}
    idleTimeLimit={3}
    preload={true}
    autoPlay={true}
    loop={false}
    speed={1.5}
/>

### 3. [Linux only] Allow download of large results

If you recieve a warning about download speeds on Linux, then you can optionally paste the following command:

```
sudo sysctl -w net.core.rmem_max=2500000
```

It may prompt you for your password when you run this command. Type your password and press **Enter**. 


**Congratulations - you've installed PLEX and are now ready to run a tool!**

## Next steps: Run a tool to check PLEX is working as expected

* Try this quick-run [small molecule binding tool](../small-molecule-binding/run-an-example.md). It's a fast algorithm, so you can **run a job** and **visualise results** in 3-5 minutes. This guide shows you how.
* Then, why not try a [protein folding tool here](../protein-folding/run-an-example.md). This guide will walk you through, step by step.

