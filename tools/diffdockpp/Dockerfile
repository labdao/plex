# FROM nvcr.io/nvidia/pytorch:22.05-py3
FROM nvcr.io/nvidia/pytorch:22.08-py3
# FROM nvcr.io/nvidia/pytorch:23.04-py3

LABEL author="Tuple, LLC <contact@tuple.xyz>"

## Environment Settings
ENV DEBIAN_FRONTEND=noninteractive \
    PIP_ROOT_USER_ACTION=ignore

USER root

RUN mkdir /DiffDock-PP
WORKDIR /DiffDock-PP

## Install Basic Dependencies
RUN apt-get clean && \
    apt-get update && \
    apt-get -y install git curl wget python3.10

## Install Dependencies
# RUN conda install pytorch=1.13.0 pytorch-cuda=11.6 -c pytorch -c nvidia

# RUN pip uninstall -y torch

# RUN pip install --no-cache-dir \
#         torch-scatter==2.0.9 \
#         torch-sparse==0.6.15 \
#         torch-scatter \
#         torch-sparse \
#         torch-cluster \
#         torch-spline-conv \
#         torch-geometric \
#         -f https://data.pyg.org/whl/torch-1.13.0+cu116.html

RUN pip install --no-cache-dir \
        torch-scatter \
        torch-sparse \
        torch-cluster \
        torch-spline-conv \
        torch-geometric

RUN pip install \
        numpy \
        dill \
        tqdm \
        pyyaml \
        pandas \
        biopandas \
        scikit-learn \
        biopython \
        e3nn \
        wandb \
        tensorboard \
        tensorboardX \
        matplotlib

## Get DiffDock-PP Scripts
# RUN git clone https://github.com/ketatam/DiffDock-PP.git /DiffDock-PP
COPY . /DiffDock-PP/.

## Get DIPS dataset
# RUN curl -L -o /DiffDock-PP/datasets/DIPS/dips.zip https://www.dropbox.com/s/sqknqofy58nlosh/DIPS.zip?dl=0 &&\
#    unzip /DiffDock-PP/datasets/DIPS/dips.zip -d /DiffDock-PP/datasets/DIPS/pairs_pruned

## Update /diffdock-pp permissions
# RUN chmod -R 770 /diffdock-pp