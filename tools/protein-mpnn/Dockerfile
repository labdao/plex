FROM ubuntu:20.04

RUN apt-get update

# Install apt packages
RUN apt-get install -y $(cat apt-list.txt)

# Install pip packages
RUN pip install -r pip-freeze.txt

CMD ["/bin/bash"]
