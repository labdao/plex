# Start with a base image that has CUDA and Python
FROM rapidsai/mambaforge-cuda:cuda11.8.0-base-ubuntu22.04-py3.9
# FROM docker.io/pytorch/pytorch:2.1.0-cuda11.8-cudnn8-runtime

# Install system dependencies
RUN apt-get update && apt-get install -y wget git nano less zip unzip && rm -rf /var/lib/apt/lists/*

# Set CUDA version and ColabFold version as arguments
ARG CUDA_VERSION=11.8.0
ARG COLABFOLD_VERSION=1.5.3

# WORKDIR /app

# Install CUDA nvcc
RUN apt-get update && apt-get install -y cuda-nvcc-$(echo $CUDA_VERSION | cut -d'.' -f1,2 | tr '.' '-') --no-install-recommends --no-install-suggests && rm -rf /var/lib/apt/lists/*

# Install colabfold, openmm, pdbfixer, and additional Python packages
RUN mamba config --set auto_update_conda false && \
    CONDA_OVERRIDE_CUDA=$(echo $CUDA_VERSION | cut -d'.' -f1,2) mamba install -y -c conda-forge -c bioconda colabfold=$COLABFOLD_VERSION jaxlib==*=cuda* openmm pdbfixer && \
    mamba clean -afy

# Install transformers, PyTorch and other requirements
COPY requirements.txt .
RUN pip install -r requirements.txt

# Copy the rest of your application
COPY . .

# Run test script
RUN chmod +x test.sh && ./test.sh

# Set environment variables
ENV PATH /usr/local/envs/colabfold/bin:$PATH
ENV MPLBACKEND Agg
VOLUME cache
ENV MPLCONFIGDIR /cache
ENV XDG_CACHE_HOME /cache

# Download weights for colabfold
RUN python -m colabfold.download

# Additional installations (if needed)
RUN pip install hydra-core omegaconf pandas biopython pyyaml

# Install prodigy
RUN git clone -q https://github.com/haddocking/prodigy && \
     pip install -q prodigy/

# # Copy in Files
# COPY . /app

# ENV PYTHONUNBUFFERED=1
# ENV HYDRA_FULL_ERROR=1


# # entrypoint
# ENTRYPOINT ["python", "-u", "main.py"]

# Label the container
LABEL container="protbert_colabfold_prodigy"
