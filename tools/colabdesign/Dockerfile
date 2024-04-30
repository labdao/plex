FROM quay.io/convexitylabs/base-images:pytorch-2.1.0-cuda11.8-cudnn8-runtime-3068995@sha256:8f75299825d3dc613d32070ac72b99439751054c06b8af994c00205954b4ee4e

WORKDIR /app

RUN apt-get update && apt-get install -y software-properties-common && \
    add-apt-repository ppa:deadsnakes/ppa && \
    apt-get update && \
    apt-get install -y python3.10 python3.10-dev

RUN apt-get update && apt-get install -y --no-install-recommends \
    wget \
    xz-utils \
    libgl-dev \
    libxi-dev \
    libxrender-dev \
    libxkbcommon-x11-0 \
    libegl1-mesa \
    libgl1-mesa-dri \
    libgomp1 \
    software-properties-common \
    libsm6 libxext6 \
    mesa-utils \
    libegl1-mesa-dev \
    libgl1-mesa-dev \
    libnvidia-egl-wayland1 \
    libnvidia-gl-460 \
    mesa-utils-extra && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Download and install pip
RUN wget -O get-pip.py 'https://bootstrap.pypa.io/get-pip.py' && \
    python get-pip.py && \
    rm get-pip.py

COPY viz-requirements.txt .
RUN pip install --upgrade pip && \
    pip install -r viz-requirements.txt

RUN wget -q https://download.blender.org/release/Blender4.1/blender-4.1.1-linux-x64.tar.xz && \
    tar -xf blender-4.1.1-linux-x64.tar.xz && \
    mv blender-4.1.1-linux-x64 blender && \
    rm blender-4.1.1-linux-x64.tar.xz && \
    echo '#!/bin/sh\n/blender/blender -noaudio -b --python-use-system-env $@' > /usr/bin/blender && \
    chmod +x /usr/bin/blender
# Copy in Files
COPY main.py /app
COPY conf /app/conf
COPY visualisers.py /app
COPY base_classes.py /app

ENV PYTHONUNBUFFERED=1
ENV HYDRA_FULL_ERROR=1

# Add so jax can find cuDnn library
ENV LD_LIBRARY_PATH=/opt/conda/lib/python3.10/site-packages/torch/lib/

# entrypoint
ENTRYPOINT ["python", "-u", "main.py"]
