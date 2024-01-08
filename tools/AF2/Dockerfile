ARG CUDA_VERSION=11.8.0
ARG COLABFOLD_VERSION=1.5.3
FROM rapidsai/mambaforge-cuda:cuda${CUDA_VERSION}-base-ubuntu22.04-py3.9

# colabfold
RUN apt-get update && apt-get install -y wget cuda-nvcc-$(echo $CUDA_VERSION | cut -d'.' -f1,2 | tr '.' '-') --no-install-recommends --no-install-suggests && rm -rf /var/lib/apt/lists/* && \
    mamba config --set auto_update_conda false && \
    CONDA_OVERRIDE_CUDA=$(echo $CUDA_VERSION | cut -d'.' -f1,2) mamba install -y -c conda-forge -c bioconda colabfold=$COLABFOLD_VERSION jaxlib==*=cuda* && \
    mamba clean -afy

ENV PATH /usr/local/envs/colabfold/bin:$PATH
ENV MPLBACKEND Agg
VOLUME cache
ENV MPLCONFIGDIR /cache
ENV XDG_CACHE_HOME /cache

# download weights
RUN python -m colabfold.download

# openmm, pdbfixer, and additional Python packages
RUN mamba install -y openmm pdbfixer && \
    mamba clean -afy
RUN pip install hydra-core omegaconf pandas biopython pyyaml

# omegafold
RUN apt-get update && apt-get install -y git nano
RUN git clone https://github.com/HeliXonProtein/OmegaFold
RUN cd OmegaFold && python setup.py install
RUN pip install git+https://github.com/HeliXonProtein/OmegaFold.git

# downloading weights
RUN wget https://helixon.s3.amazonaws.com/release1.pt -O /root/.cache/omegafold_ckpt/model.pt 
