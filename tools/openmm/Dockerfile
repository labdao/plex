FROM nvidia/cuda:11.6.2-base-ubuntu20.04

ENV PATH="/root/miniconda3/bin:${PATH}"
ARG PATH="/root/miniconda3/bin:${PATH}"

# Install necessary packages
RUN apt-get update && apt-get install -y wget && rm -rf /var/lib/apt/lists/*

# Install Miniconda
RUN wget https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh \
    && mkdir /root/.conda \
    && bash Miniconda3-latest-Linux-x86_64.sh -b \
    && rm -f Miniconda3-latest-Linux-x86_64.sh

# Check Conda installation
RUN conda --version

# Install OpenMM with specific CUDA version
RUN conda install -c conda-forge libstdcxx-ng openmm cudatoolkit=11.6 -y

# Copy your scripts
COPY scripts scripts
