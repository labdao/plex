FROM nvcr.io/nvidia/pytorch:22.05-py3

# credit for the new version goes largely to the Tuple team 
LABEL author="Tuple, LLC <contact@tuple.xyz>"

## Environment Settings
ENV DEBIAN_FRONTEND=noninteractive \
    PIP_ROOT_USER_ACTION=ignore

USER root

## Install Basic Dependencies
RUN apt-get clean && \
    apt-get update && \
    apt-get -y install git curl wget python3.9

## Install CUDA Toolkit (and fix libcusparse.so.11 issue - /opt/conda/lib/libcusparse.so.11)
RUN conda install -vvv --yes cudatoolkit=11.3 -c pytorch && \
    conda clean --all -f -y && \
    export LD_LIBRARY_PATH="/opt/conda/lib/:$LD_LIBRARY_PATH" && \
    export PATH="/opt/conda/lib/:$PATH"

## Install DiffDock Dependencies
RUN pip install ipython-autotime && \
    pip install pyg==0.7.1 --quiet && \
    pip install pyyaml==6.0 --quiet && \
    pip install scipy==1.7.3 --quiet && \
    pip install networkx==2.6.3 --quiet && \
    pip install biopython==1.79 --quiet && \
    pip install rdkit-pypi==2022.03.5 --quiet && \
    pip install e3nn==0.5.0 --quiet && \
    pip install spyrmsd==0.5.2 --quiet && \
    pip install pandas==1.3.5 --quiet && \
    pip install biopandas==0.4.1 --quiet && \
    pip install py3Dmol==1.8.1 --quiet && \
    pip install torch==1.12.1+cu113 --extra-index-url https://download.pytorch.org/whl/cu113

### Git DiffDock
RUN mkdir /app && \
    cd /app && \
    git clone https://github.com/gcorso/DiffDock.git && \
    cd /app/DiffDock && \
    git checkout 0f9c419

### Install DiffDock-Specific Torch Dependencies
RUN pip uninstall torch-scatter torch-sparse torch-geometric torch-cluster --y && \
    pip install torch-scatter -f https://data.pyg.org/whl/torch-1.12.0%2Bcu113.html --quiet && \
    pip install torch-sparse -f https://data.pyg.org/whl/torch-1.12.0%2Bcu113.html --quiet && \
    pip install torch-cluster -f https://data.pyg.org/whl/torch-1.12.0%2Bcu113.html --quiet && \
    pip install git+https://github.com/pyg-team/pytorch_geometric.git --quiet

### Download Smina
RUN cd /app/DiffDock && \
    wget https://sourceforge.net/projects/smina/files/smina.static/download -O smina && \
    chmod +x smina

### Download Folding Models
RUN pip install "fair-esm[esmfold]"
RUN pip install 'dllogger @ git+https://github.com/NVIDIA/dllogger.git'
RUN pip install 'openfold @ git+https://github.com/aqlaboratory/openfold.git@4b41059694619831a7db195b7e0988fc4ff3a307'

### Deprecated: Install Facebook ESM
#RUN cd /app/DiffDock && \
#    git clone https://github.com/facebookresearch/esm && \
#    cd /app/DiffDock/esm && \
#    git checkout f07aed6 && \
#    pip install -e .  && \
#    cd /app/DiffDock

## Update /app permissions
RUN chmod -R 770 /app

## Set Working Directory
WORKDIR /app/DiffDock