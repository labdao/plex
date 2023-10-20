# Use an official Python runtime as a parent image
FROM continuumio/miniconda3

# The environment variable ensures that the python output is set straight
# to the terminal with out buffering it first
ENV PYTHONUNBUFFERED 1

# Install necessary libraries for rdkit
RUN apt-get update && apt-get install -y --no-install-recommends \
    cmake \
    libboost-dev \
    libboost-system-dev \
    libboost-iostreams-dev \
    git \
    unzip \
    && rm -rf /var/lib/apt/lists/*

# Create a new conda environment with python 3.7
RUN conda create -n env python=3.7

# Activate the conda environment
SHELL ["conda", "run", "-n", "env", "/bin/bash", "-c"]

# Install python libraries
RUN conda install -y pytorch==1.8.1 torchvision -c pytorch && \
    pip install scikit-learn==0.23.2 biopython==1.78 pandas==1.1.3 scipy==1.5.2 numpy==1.20.2 seaborn==0.11.0 matplotlib==3.3.2 && \
    pip install rdkit-pypi==2020.9.5.2 requests==2.26.0

# Set the working directory to /app
WORKDIR /app

# Clone the DLKcat package
RUN git clone https://github.com/SysBioChalmers/DLKcat

# Change directory to DeeplearningApproach under the DLKcat package and unzip input.zip
WORKDIR /app/DLKcat/DeeplearningApproach/Data
RUN unzip input.zip

# Set work directory to the script location
WORKDIR /app/DLKcat/DeeplearningApproach/Code/example

CMD ["conda", "run", "-n", "env", "python", "prediction_for_input.py", "/app/DLKcat/DeeplearningApproach/Code/example/input.tsv"]

