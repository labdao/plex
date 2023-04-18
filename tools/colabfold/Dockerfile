# FROM ubuntu:focal
FROM athbaltzis/colabfold_proteinfold:v0.11
LABEL container="colabfold"

# installing packages required for installation
RUN echo "downloading basic packages for installation"
RUN apt-get update
RUN apt-get install -y tmux wget curl nano

# checking installation of tools
RUN gcc --version
RUN nvcc --version

# creating param dir
RUN mkdir -p /colabfold_batch/colabfold/params