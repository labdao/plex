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

# Additional installations
RUN pip install hydra-core omegaconf pandas biopython pyyaml

# RUN git clone https://github.com/sokrypton/RFdiffusion.git
RUN git clone https://github.com/RosettaCommons/RFdiffusion.git
RUN pip -q install jedi omegaconf hydra-core icecream pyrsistent
RUN pip install dgl==1.0.2+cu116 -f https://data.dgl.ai/wheels/cu116/repo.html
RUN cd RFdiffusion/env/SE3Transformer; pip -q install --no-cache-dir -r requirements.txt; pip -q install .
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

# colabfold here
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

# Copy in Files
COPY . /app

ENV PYTHONUNBUFFERED=1
ENV HYDRA_FULL_ERROR=1


# entrypoint
ENTRYPOINT ["python", "-u", "main.py"]
