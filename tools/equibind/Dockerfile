# FROM docker.io/pytorch/pytorch:1.11.0-cuda11.3-cudnn8-runtime
FROM python:3.8-slim 

WORKDIR /app

# Set non-interactive shell to avoid stuck installations
ENV TZ=Asia/Kolkata \
    DEBIAN_FRONTEND=noninteractive

RUN apt-get update && \
    apt-get install -y git build-essential cmake jq

# Reset the shell to interactive mode for future RUN commands and when the container is running
ENV DEBIAN_FRONTEND=interactive

RUN git clone https://github.com/HannesStark/EquiBind.git

WORKDIR /app/EquiBind

COPY requirements.txt .

RUN pip install -r requirements.txt


RUN apt-get update && apt-get install -y jq

# PLEX_JOB_INPUTS is set by plex
CMD python inference.py -o /outputs -r $(echo $PLEX_JOB_INPUTS | jq -r '.["protein"]') -l $(echo $PLEX_JOB_INPUTS | jq -r '.["small_molecule"]')
