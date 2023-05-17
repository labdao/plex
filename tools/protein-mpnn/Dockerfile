FROM nvidia/cuda:11.8.0-runtime-ubuntu20.04

RUN apt-get update

COPY . .
# Install apt packages
RUN tail -n +7 apt-list.txt | cut -d'/' -f1 > apt-list-clean.txt
RUN while read pkg; do apt-get install -y $pkg || true; done <apt-list-clean.txt
RUN apt-get install -y python3-pip

# Install pip packages
RUN pip install -r pip-freeze.txt

CMD ["/bin/bash"]
