FROM continuumio/miniconda3:22.11.1

# Install system dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    g++ \
    python3-dev \
    libeigen3-dev \
    wget \
    cmake \
    less \
    unzip

# Installing OpenBabel
RUN wget https://github.com/openbabel/openbabel/releases/download/openbabel-3-1-1/openbabel-3.1.1-source.tar.bz2
RUN tar -xjf openbabel-3.1.1-source.tar.bz2
RUN cd openbabel-3.1.1
RUN cmake ../openbabel-3.1.1 \
    -DPYTHON_BINDINGS=ON \
    -DRUN_SWIG=ON \
    -DCMAKE_INSTALL_PREFIX=/opt/conda \
    -DPYTHON_INCLUDE_DIR=/opt/conda/include/python3.10 \
    -DCMAKE_LIBRARY_PATH=/opt/conda/lib \
    -DSWIG_DIR=/opt/conda/share/swig/4.0.2 \
    -DSWIG_EXECUTABLE=/opt/conda/bin/swig \
    -DPYTHON_LIBRARY=/opt/conda/lib/libpython3.10.so \
    -DCMAKE_BUILD_TYPE=DEBUG
RUN make -j4
RUN make install

# download and install rf-score VS binary
#RUN wget http://wojcikowski.pl/travis/rf-score-vs_v1.0_linux_2.7.zip
#RUN unzip rf-score-vs_v1.0_linux_2.7.zip -d rf-score-vs

# Set the working directory
WORKDIR /app

# Copy the requirements file and install the dependencies
#COPY requirements.txt .
#RUN pip3 install -r requirements.txt

# Install the application
RUN pip3 install oddt

# Copy additional files
COPY . .
RUN oddt_cli /app/testdata/6d08_ligand.sdf \
    --receptor /app/testdata/6d08_protein_processed.pdb \
    #--score autodock_vina \
    #--score rfscore \
    --score rfscore_v1 \
    --score rfscore_v2 \
    --score rfscore_v3 \
    --score nnscore \
    #--score pleclinear \
    #--score plecnn \
    #--score plecrf \
    -O /app/testdata/6d08_ligand_scored.sdf
