# Open Babel Tools with PLEX

Welcome to the `openbabel` directory in the PLEX project! This directory contains tools and configuration files to run Open Babel, a chemical toolbox designed to speak the many languages of chemical data, on the PLEX platform.

PLEX is a client for distributed computation that uses Bacalhau, a distributed compute network, and IPFS for decentralized storage. With PLEX, you can build highly reproducible workflows using Docker containers and execute them on a public network at any scale.

Our tools allow you to run Open Babel in Docker containers across the Bacalhau network, providing scalable, reproducible chemical computations. This README will guide you on how to utilize these tools and configurations.

## Overview of the Repository

Our Open Babel tools repository contains the following files:

1. `7n9g.pdb`: A protein structure file in PDB format.
2. `7n9g_ZINC000019632618_docked.sdf`: A file representing the docked structure of a small molecule (identified by ZINC number) with the protein `7n9g`.
3. `Dockerfile`: The Dockerfile for creating the Docker image that contains the necessary software environment to run Open Babel.
4. `ZINC000019632618.sdf`: A structure-data file for a small molecule identified by the ZINC number.
5. `pdb-to-sdf-openbabel.json`: Configuration file for converting PDB files to SDF format using Open Babel.
6. `rmsd-openbabel.json`: Configuration file for calculating the Root Mean Square Deviation (RMSD) between reference and docked structures using Open Babel.

## Docker Environment

PLEX utilizes Docker to ensure reproducibility and ease of use across different platforms. In the context of Open Babel, a Docker image is prepared with the necessary software environments required to run Open Babel applications.

For example, in our Dockerfile:

1. We use a base image that has Python 3.9 installed.
2. We install Open Babel and its development libraries.
3. We set up an `/app` directory as the working directory inside the Docker container.
4. We copy the content of the current directory into the `/app` directory in the Docker container.
5. We start the container with a Bash shell by default.

This setup ensures that no matter where you run the Docker container, the environment will always be consistent, thus maintaining the reproducibility of your scientific computation tasks. You can check out the Dockerfile in this repository.

## Configuration Files

Configuration files (e.g., `pdb-to-sdf-openbabel.json` and `rmsd-openbabel.json`) are used to define the workflows run on PLEX. These JSON files specify the input and output files, the arguments for the Open Babel command-line tool, and the Docker image to pull. These workflows can be easily run on the PLEX platform, ensuring reproducibility and scalability of your Open Babel computations.

## Conclusion

The tools in this repository provide a way for you to run scalable, reproducible Open Babel computations on the PLEX platform. By utilizing Docker and the Bacalhau network, you can focus on solving unique scientific problems without worrying about the underlying computational infrastructure.

For more information on how to use these tools, please refer to the PLEX documentation or reach out to our team. Happy computing!
