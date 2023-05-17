FROM pytorch/pytorch:2.0.1-cuda11.7-cudnn8-runtime

RUN apt-get update
RUN apt-get install -y software-properties-common
RUN add-apt-repository ppa:deadsnakes/ppa
RUN apt-get update && apt-get install -y python3.10 python3.10-dev python3.10-distutils python3-pip

COPY . .
# Install apt packages
RUN tail -n +7 apt-list.txt | cut -d'/' -f1 > apt-list-clean.txt
RUN while read pkg; do apt-get install -y $pkg || true; done <apt-list-clean.txt
# Install pip packages
RUN pip install -r pip-freeze.txt

CMD ["/bin/bash"]
