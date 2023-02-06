FROM nvidia/cuda:11.0.3-base-ubuntu20.04

ENV PATH="/root/miniconda3/bin:${PATH}"
ARG PATH="/root/miniconda3/bin:${PATH}"

RUN apt-get update
RUN apt install -y git
RUN apt-get install -y wget && rm -rf /var/lib/apt/lists/*
RUN wget \
    https://repo.anaconda.com/miniconda/Miniconda3-latest-Linux-x86_64.sh \
    && mkdir /root/.conda \
    && bash Miniconda3-latest-Linux-x86_64.sh -b \
    && rm -f Miniconda3-latest-Linux-x86_64.sh

RUN conda create -n diffdock python=3.8

# Make RUN commands use the new environment
SHELL ["conda", "run", "--no-capture-output", "-n", "diffdock", "/bin/bash", "-c"]

RUN conda install pytorch pytorch-cuda=11.7 -c pytorch -c nvidia
RUN pip install torch-scatter torch-sparse torch-cluster torch-spline-conv torch-geometric -f https://data.pyg.org/whl/torch-1.13.0+cu117.html
RUN python -m pip install PyYAML scipy "networkx[default]" biopython rdkit-pypi e3nn spyrmsd pandas biopandas

RUN git clone https://github.com/gcorso/DiffDock.git

WORKDIR DiffDock

RUN git clone https://github.com/facebookresearch/esm
RUN cd esm && pip install -e .

COPY test.sh test.sh

CMD ["./test.sh"]
