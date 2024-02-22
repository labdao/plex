# FROM docker.io/pytorch/pytorch:2.1.0-cuda11.8-cudnn8-runtime
FROM rapidsai/mambaforge-cuda:cuda11.8.0-base-ubuntu22.04-py3.9

ENV DEBIAN_FRONTEND=noninteractive_autoprompt

# Install required packages
RUN apt-get update -y && \
     apt-get install --no-install-recommends -y wget aria2 git zip unzip less nano build-essential

# Install any missing CUDA libraries or tools
RUN apt-get update && apt-get install -y --no-install-recommends \
    cuda-libraries-11-8 \
    cuda-nvtx-11-8 \
    libcublas-11-8 \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Install hydra-core
RUN pip3 install --no-cache-dir hydra-core --upgrade pandas

# Set CUDA version and ColabFold version as arguments
ARG CUDA_VERSION=11.8.0
ARG COLABFOLD_VERSION=1.5.5

WORKDIR /app

# Install CUDA nvcc
RUN apt-get update && apt-get install -y cuda-nvcc-$(echo $CUDA_VERSION | cut -d'.' -f1,2 | tr '.' '-') --no-install-recommends --no-install-suggests && rm -rf /var/lib/apt/lists/*

# Install colabfold, openmm, pdbfixer, and additional Python packages
RUN mamba config --set auto_update_conda false && \
    CONDA_OVERRIDE_CUDA=$(echo $CUDA_VERSION | cut -d'.' -f1,2) mamba install -y -c conda-forge -c bioconda colabfold=$COLABFOLD_VERSION jaxlib==*=cuda* openmm pdbfixer && \
    mamba clean -afy

# Install transformers, PyTorch and other requirements
COPY requirements.txt .
RUN pip install -r requirements.txt

# # Additional installations
# RUN pip install hydra-core omegaconf pandas biopython pyyaml

# RUN git clone https://github.com/sokrypton/RFdiffusion.git
RUN git clone https://github.com/RosettaCommons/RFdiffusion.git
RUN pip -q install jedi omegaconf hydra-core icecream pyrsistent
RUN pip install dgl==1.0.2+cu116 -f https://data.dgl.ai/wheels/cu116/repo.html
RUN cd RFdiffusion/env/SE3Transformer; pip -q install --no-cache-dir -r requirements.txt; pip -q install .
# RUN cd RFdiffusion/env/SE3Transformer; pip -q install -r requirements.txt; pip -q install .
RUN wget -qnc https://files.ipd.uw.edu/krypton/ananas
RUN chmod +x ananas

# Download params, modules and schedules
RUN set -ex; \
     mkdir -p params && \
     aria2c -q -x 16 https://files.ipd.uw.edu/krypton/schedules.zip && \
     aria2c -q -x 16 http://files.ipd.uw.edu/pub/RFdiffusion/6f5902ac237024bdd0c176cb93063dc4/Base_ckpt.pt && \
     aria2c -q -x 16 http://files.ipd.uw.edu/pub/RFdiffusion/e29311f6f1bf1af907f9ef9f44b8328b/Complex_base_ckpt.pt && \
     aria2c -q -x 16 https://storage.googleapis.com/alphafold/alphafold_params_2022-12-06.tar && \
     tar -xvf alphafold_params_2022-12-06.tar -C params && \
     rm -rf alphafold_params_2022-12-06.tar && \
     mkdir -p RFdiffusion/models; mv Base_ckpt.pt Complex_base_ckpt.pt RFdiffusion/models/ && \
     unzip schedules.zip; rm schedules.zip && \
     touch params/done.txt

# Install prodigy
RUN git clone -q https://github.com/haddocking/prodigy && \
     pip install -q prodigy/

# Set environment variables
ENV PATH /usr/local/envs/colabfold/bin:$PATH
ENV MPLBACKEND Agg
VOLUME cache
ENV MPLCONFIGDIR /cache
ENV XDG_CACHE_HOME /cache
ENV DGLBACKEND="pytorch"

# download weights for colabfold
RUN python -m colabfold.download

# Cleanup
RUN apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/cache /var/lib/log /var/lib/apt/lists/*

## ProteinMPNN
RUN mkdir /app/ProteinMPNN && git clone https://github.com/dauparas/ProteinMPNN /app/ProteinMPNN

# ESM2
# Download and cache the Hugging Face model
#RUN echo "source activate evodiff" > ~/.bashrc
#ENV PATH /opt/conda/envs/base/bin:$PATH
#RUN python -c "from transformers import AutoTokenizer, EsmForMaskedLM; model_name='facebook/esm2_t33_650M_UR50D'; AutoTokenizer.from_pretrained(model_name); EsmForMaskedLM.from_pretrained(model_name)"
RUN /bin/bash -c "source activate base && \
                pip install transformers && \
                python -c \"from transformers import AutoTokenizer, EsmForMaskedLM; model_name='facebook/esm2_t33_650M_UR50D'; AutoTokenizer.from_pretrained(model_name); EsmForMaskedLM.from_pretrained(model_name)\""


# omegafold
RUN apt-get update && apt-get install -y git nano
RUN git clone https://github.com/HeliXonProtein/OmegaFold
RUN cd OmegaFold && python setup.py install
RUN pip install git+https://github.com/HeliXonProtein/OmegaFold.git

# downloading weights
RUN mkdir -p /root/.cache/omegafold_ckpt && wget https://helixon.s3.amazonaws.com/release1.pt -O /root/.cache/omegafold_ckpt/model.pt

COPY . /app

ENV PYTHONUNBUFFERED=1
ENV HYDRA_FULL_ERROR=1

# # entrypoint
ENTRYPOINT ["python", "-u", "main.py"]


### OLD

# # Use NVIDIA CUDA image as the base for CUDA and cuDNN libraries
# FROM nvidia/cuda:11.3.1-cudnn8-devel-ubuntu18.04

# # # Meta information
# # LABEL org.opencontainers.image.version="1.0.0"
# # LABEL org.opencontainers.image.authors="Gustaf Ahdritz"
# # LABEL org.opencontainers.image.source="https://github.com/aqlaboratory/openfold"
# # LABEL org.opencontainers.image.licenses="Apache License 2.0"
# # LABEL org.opencontainers.image.base.name="docker.io/nvidia/cuda:11.3.1-cudnn8-devel-ubuntu18.04"

# # Set non-interactive frontend
# ENV DEBIAN_FRONTEND=noninteractive

# # NVIDIA keys update
# RUN apt-key del 7fa2af80 && \
#     apt-key adv --fetch-keys https://developer.download.nvidia.com/compute/cuda/repos/ubuntu1804/x86_64/7fa2af80.pub && \
#     apt-key adv --fetch-keys https://developer.download.nvidia.com/compute/cuda/repos/ubuntu1804/x86_64/3bf863cc.pub

# # Install base packages, CUDA libraries, and tools
# RUN apt-get update && apt-get install -y wget libxml2 cuda-libraries-11-8 cuda-nvtx-11-8 libcublas-11-8 libcusparse-dev-11-3 libcublas-dev-11-3 libcusolver-dev-11-3 git aria2 zip unzip less nano build-essential && \
#     apt-get clean && \
#     rm -rf /var/lib/apt/lists/*

# # Install Miniforge
# RUN wget -P /tmp "https://github.com/conda-forge/miniforge/releases/download/23.3.1-1/Miniforge3-Linux-x86_64.sh" && \
#     bash /tmp/Miniforge3-Linux-x86_64.sh -b -p /opt/conda && \
#     rm /tmp/Miniforge3-Linux-x86_64.sh

# ENV PATH /opt/conda/bin:$PATH

# # Install Python packages
# COPY requirements.txt .
# RUN apt-get update && apt-get install -y zlib1g-dev
# RUN pip install --no-cache-dir -r requirements.txt && \
#     pip install --no-cache-dir hydra-core --upgrade pandas

# # Set CUDA version for compatibility with additional tools and libraries
# ARG CUDA_VERSION=11.3

# # Install colabfold, openmm, pdbfixer, additional Python packages
# RUN mamba config --set auto_update_conda false && \
#     CONDA_OVERRIDE_CUDA=$CUDA_VERSION mamba install -y -c conda-forge -c bioconda colabfold jaxlib==*=cuda* openmm pdbfixer && \
#     mamba clean -afy

# # Clone and setup RFdiffusion, install additional dependencies
# RUN git clone https://github.com/RosettaCommons/RFdiffusion.git && \
#     pip install jedi omegaconf hydra-core icecream pyrsistent && \
#     pip install dgl==1.0.2+cu116 -f https://data.dgl.ai/wheels/cu116/repo.html && \
#     cd RFdiffusion/env/SE3Transformer && pip install --no-cache-dir -r requirements.txt && pip install .

# # Clone and install prodigy
# RUN git clone https://github.com/haddocking/prodigy && \
#     pip install -e prodigy/

# # Download and prepare OmegaFold
# RUN git clone https://github.com/HeliXonProtein/OmegaFold && \
#     pip install git+https://github.com/HeliXonProtein/OmegaFold.git

# # Install OpenFold dependencies and setup OpenFold
# COPY openfold/environment.yml /opt/openfold/environment.yml
# RUN mamba env update -n base --file /opt/openfold/environment.yml && mamba clean --all
# COPY openfold/openfold /opt/openfold/openfold
# COPY openfold/scripts /opt/openfold/scripts
# COPY openfold/run_pretrained_openfold.py /opt/openfold/run_pretrained_openfold.py
# COPY openfold/train_openfold.py /opt/openfold/train_openfold.py
# COPY openfold/setup.py /opt/openfold/setup.py
# WORKDIR /opt/openfold
# RUN python3 setup.py install

# # Cleanup
# RUN apt-get clean autoclean && \
#     apt-get autoremove -y && \
#     rm -rf /var/lib/apt/lists/*

# # Set WORKDIR as per requirement
# WORKDIR /app

# # Cleanup
# RUN apt-get clean autoclean && \
#     apt-get autoremove -y && \
#     rm -rf /var/lib/{apt,dpkg,cache,log}/

# # Download weights for colabfold and setup environment variables
# # Additional setup commands can go here

# ENV PYTHONUNBUFFERED=1
# ENV HYDRA_FULL_ERROR=1
# ENV PATH /usr/local/envs/colabfold/bin:$PATH
# ENV MPLBACKEND Agg
# ENV MPLCONFIGDIR /cache
# ENV XDG_CACHE_HOME /cache
# ENV DGLBACKEND="pytorch"

# # Entry point
# ENTRYPOINT ["python", "-u", "main.py"]


### VERY OLD

# # Use the NVIDIA CUDA image as the base
# FROM nvidia/cuda:11.3.1-cudnn8-devel-ubuntu18.04

# # Set non-interactive frontend
# ENV DEBIAN_FRONTEND=noninteractive

# # NVIDIA keys update
# RUN apt-key del 7fa2af80 && \
#     apt-key adv --fetch-keys https://developer.download.nvidia.com/compute/cuda/repos/ubuntu1804/x86_64/7fa2af80.pub && \
#     apt-key adv --fetch-keys https://developer.download.nvidia.com/compute/cuda/repos/ubuntu1804/x86_64/3bf863cc.pub

# # Install required packages
# RUN apt-get update -y && \
#     apt-get install --no-install-recommends -y wget libxml2 cuda-minimal-build-11-3 libcusparse-dev-11-3 libcublas-dev-11-3 libcusolver-dev-11-3 git zip unzip less nano build-essential aria2 && \
#     apt-get clean && \
#     rm -rf /var/lib/apt/lists/*

# # Install Miniforge
# RUN wget -P /tmp "https://github.com/conda-forge/miniforge/releases/download/23.3.1-1/Miniforge3-Linux-x86_64.sh" && \
#     bash /tmp/Miniforge3-Linux-x86_64.sh -b -p /opt/conda && \
#     rm /tmp/Miniforge3-Linux-x86_64.sh

# ENV PATH /opt/conda/bin:$PATH

# # Install Python packages and CUDA nvcc
# COPY requirements.txt .
# RUN apt-get update && apt-get install -y zlib1g-dev
# RUN pip install --no-cache-dir -r requirements.txt && \
#     pip install --no-cache-dir hydra-core --upgrade pandas

# # Install colabfold, openmm, pdbfixer, additional Python packages
# ARG CUDA_VERSION=11.3
# ARG COLABFOLD_VERSION=1.5.5
# RUN mamba config --set auto_update_conda false && \
#     CONDA_OVERRIDE_CUDA=11.3 mamba install -y -c conda-forge -c bioconda colabfold=$COLABFOLD_VERSION jaxlib==*=cuda* openmm pdbfixer && \
#     mamba clean -afy

# # Before a step that fails or before running mamba env update
# RUN python --version || true
# RUN python3 --version || true

# # Install OpenFold dependencies
# # Since the environment.yml file is inside the openfold subdirectory of your working directory
# COPY openfold/environment.yml openfold/opt/openfold/environment.yml
# RUN mamba env update -n base --file openfold/opt/openfold/environment.yml && mamba clean --all
# ENV LD_LIBRARY_PATH=${CONDA_PREFIX}/lib:${LD_LIBRARY_PATH}

# # Copy the OpenFold code and related scripts from the openfold subdirectory in your working directory to the container
# COPY openfold/openfold openfold/opt/openfold

# # Copy OpenFold and related scripts
# # COPY openfold/openfold /opt/openfold/openfold/
# COPY openfold/scripts openfold/opt/openfold/scripts
# COPY openfold/run_pretrained_openfold.py openfold/opt/openfold/run_pretrained_openfold.py
# COPY openfold/train_openfold.py openfold/opt/openfold/train_openfold.py
# COPY openfold/setup.py openfold/opt/openfold/setup.py
# WORKDIR /openfold/opt/openfold
# RUN python3 setup.py install

# # Set WORKDIR as per requirement
# WORKDIR /app

# # Install other required tools and libraries
# # Additional commands for other installations can go here

# # Cleanup
# RUN apt-get clean autoclean && \
#     apt-get autoremove -y && \
#     rm -rf /var/lib/{apt,dpkg,cache,log}/

# # Download weights for colabfold and setup environment variables
# # Additional setup commands can go here

# ENV PYTHONUNBUFFERED=1
# ENV HYDRA_FULL_ERROR=1
# ENV PATH /usr/local/envs/colabfold/bin:$PATH
# ENV MPLBACKEND Agg
# ENV MPLCONFIGDIR /cache
# ENV XDG_CACHE_HOME /cache
# ENV DGLBACKEND="pytorch"

# # Entry point
# ENTRYPOINT ["python", "-u", "main.py"]
