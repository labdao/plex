# Use the official Ubuntu 20.04 LTS base image
FROM teamcsm/openbabel:py3.9-openbabel3.1.1

# Set environment variables
ENV DEBIAN_FRONTEND=noninteractive

# Update and install dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        python3 \
        python3-pip \
        wget \
        swig \
        build-essential \
        ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Install vina
RUN wget https://github.com/ccsb-scripps/AutoDock-Vina/releases/download/v1.2.3/vina_1.2.3_linux_x86_64 && \
    mv vina_1.2.3_linux_x86_64 /usr/local/bin/vina && \
    echo 'export PATH=/usr/local/bin/vina:$PATH' >> ~/.bashrc && \
    chmod a+x /usr/local/bin/vina

# Install adfr suite for protein processing
RUN wget https://ccsb.scripps.edu/adfr/download/1028/ -O adfr
RUN chmod a+x adfr
RUN ./adfr --mode silent
RUN find / -type d | grep ADFR
RUN cd /usr/local/ADFRsuite-1.0 && \
    tar zxvf ADFRsuite_x86_64Linux_1.0.tar.gz && \
    cd ADFRsuite_x86_64Linux_1.0/ && \ 
    yes Y | ./install.sh -d adfr -c 0
RUN echo 'export PATH=/ADFRsuite-1.0/ADFRsuite_x86_64Linux_1.0/adfr/bin:$PATH' >> ~/.bashrc

# Install meeko for ligand processing
RUN pip3 install --no-cache-dir -U numpy openbabel meeko

# Set the working directory
WORKDIR /app

# Run the container as an executable
CMD ["/bin/bash"]