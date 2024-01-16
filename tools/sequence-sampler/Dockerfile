# Start with a base image that has CUDA and Python
FROM rapidsai/mambaforge-cuda:cuda11.8.0-base-ubuntu22.04-py3.9
# FROM docker.io/pytorch/pytorch:2.1.0-cuda11.8-cudnn8-runtime

# Install system dependencies
RUN apt-get update && apt-get install -y wget git nano less zip unzip && rm -rf /var/lib/apt/lists/*

# Set CUDA version and ColabFold version as arguments
ARG CUDA_VERSION=11.8.0
ARG COLABFOLD_VERSION=1.5.3

WORKDIR /app

ENV PYTHONUNBUFFERED=1
ENV HYDRA_FULL_ERROR=1

# Install CUDA nvcc
RUN apt-get update && apt-get install -y cuda-nvcc-$(echo $CUDA_VERSION | cut -d'.' -f1,2 | tr '.' '-') --no-install-recommends --no-install-suggests && rm -rf /var/lib/apt/lists/*

# Install colabfold, openmm, pdbfixer, and additional Python packages
RUN mamba config --set auto_update_conda false && \
    CONDA_OVERRIDE_CUDA=$(echo $CUDA_VERSION | cut -d'.' -f1,2) mamba install -y -c conda-forge -c bioconda colabfold=$COLABFOLD_VERSION jaxlib==*=cuda* openmm pdbfixer && \
    mamba clean -afy

# Install transformers, PyTorch and other requirements
COPY requirements.txt .
RUN pip install -r requirements.txt

# Additional installations
RUN pip install hydra-core omegaconf pandas biopython pyyaml

# Set environment variables
ENV PATH /usr/local/envs/colabfold/bin:$PATH
ENV MPLBACKEND Agg
VOLUME cache
ENV MPLCONFIGDIR /cache
ENV XDG_CACHE_HOME /cache
ENV DGLBACKEND="pytorch"

# download weights for colabfold
RUN python -m colabfold.download

# Prodigy
RUN git clone -q https://github.com/haddocking/prodigy && \
     pip install -q prodigy/

# RFdiffusion
RUN mkdir /app/RFdiffusion && git clone https://github.com/RosettaCommons/RFdiffusion.git /app/RFdiffusion

RUN pip install -q --no-cache-dir \
  dgl==1.0.2+cu118 -f https://data.dgl.ai/wheels/cu118/repo.html \
  torch==2.0.0+cu118 --extra-index-url https://download.pytorch.org/whl/cu118 \
  e3nn==0.3.3 \
  wandb==0.12.0 \
  pynvml==11.0.0 \
  git+https://github.com/NVIDIA/dllogger#egg=dllogger \
  decorator==5.1.0 \
  hydra-core==1.3.2 \
  pyrsistent==0.19.3 \
  /app/RFdiffusion/env/SE3Transformer \
  && pip install --no-cache-dir /app/RFdiffusion --no-deps

RUN pip install /app/RFdiffusion/env/SE3Transformer --no-deps && pip install --no-cache-dir /app/RFdiffusion --no-deps

# download weights
RUN mkdir -p /app/models 
RUN chmod +x /app/RFdiffusion/scripts/download_models.sh
RUN bash -x /app/RFdiffusion/scripts/download_models.sh /app/models

## ProteinMPNN
RUN mkdir /app/ProteinMPNN && git clone https://github.com/dauparas/ProteinMPNN /app/ProteinMPNN

# Copy in Files
COPY . /app

# Label the container
LABEL container="seqoptim:rfdiff-mpnn"

# entrypoint
ENTRYPOINT ["python", "-u", "main.py"]


# RUN pip install -q --no-cache-dir \
#   dgl==1.0.2+cu116 -f https://data.dgl.ai/wheels/cu116/repo.html \
#   torch==1.12.1+cu116 --extra-index-url https://download.pytorch.org/whl/cu116 \
#   e3nn==0.3.3 \
#   wandb==0.12.0 \
#   pynvml==11.0.0 \
#   git+https://github.com/NVIDIA/dllogger#egg=dllogger \
#   decorator==5.1.0 \
#   hydra-core==1.3.2 \
#   pyrsistent==0.19.3 \
#   /app/RFdiffusion/env/SE3Transformer \
#   && pip install --no-cache-dir /app/RFdiffusion --no-deps