---
title: Adding a service
description: How to add a new service to the lab-exchange.
sidebar_position: 3
---

## Adding a service

The lab-exchange is currently in version 0.1 and only supports running simple services on our centralized compute infrastructure. We are currently preparing the release of the complete backend infrastructure to allow more nodes to join the network. 

In order to add a tool to the lab-exchange, you have to go through multiple steps. We start with the basics because it is important to us that community members from all experience-levels know how to contribute tools to the ecosystem. You can find a very basic application repository for [reverse-complement](https://github.com/openlab-apps/lab-reverse_complement) generation on our GitHub.  

0. Make sure your tool can be launched from the command line and does not need any manual change to variables within your scripts to run. You can define variables for [python scripts](https://www.tutorialspoint.com/python/python_command_line_arguments.htm) and [R scripts](https://www.r-bloggers.com/2015/09/passing-arguments-to-an-r-script-from-command-lines/) easily from the command line.
1. Make sure your tool is on github, ideally in a *public* repository (this will make it easier for community members to help you).
2. Dockerize your tool. [Docker](https://www.docker.com/) is a popular container framework that enables high reproducibility for running code. You can think of it as a lightweight virtual machine that makes sharing code very easy (no more missing packages/libraries!). The key element of every docker container is the dockerfile. Tutorials on how to dockerize [python](https://medium.com/swlh/dockerize-your-python-command-line-program-6a273f5c5544) and [R](https://www.r-bloggers.com/2019/02/running-your-r-script-in-docker/) command line tools exist. It is good practice to add docker-related files to your project's git reposority. 
3. Wrap your tool in a [nextflow](https://www.nextflow.io/docs/latest/getstarted.html), or even [NF-core](https://nf-co.re/tools/#creating-a-new-pipeline) workflow. The community can support you during this step.
4. Integrate with the lab-exchange API (Documentation pending). 
5. Deploy the application on your own hardware or work with someone in the community to host your service. Once hosted, everybody in the community can use your tool to run jobs and mint NFTs.
