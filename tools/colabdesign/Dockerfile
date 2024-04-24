FROM quay.io/convexitylabs/base-images:pytorch-2.1.0-cuda11.8-cudnn8-runtime-3068995@sha256:8f75299825d3dc613d32070ac72b99439751054c06b8af994c00205954b4ee4e

# Copy in Files
COPY main.py /app
COPY conf /app/conf
COPY visualisers.py /app
COPY base_classes.py /app

COPY viz-requirements.txt .

RUN pip install --no-cache-dir -r viz-requirements.txt

# RUN apt-get update && apt-get install -y \
#     libegl1-mesa \
#     libgl1-mesa-dri \
#     wget xz-utils libgl-dev libxi-dev libxrender-dev libxkbcommon-x11-0 &&\
#     pip install -U pip wheel setuptools && pip install numpy requests &&\
#     wget -q https://download.blender.org/release/Blender3.6/blender-3.6.0-linux-x64.tar.xz &&\
#     tar -xf blender-3.6.0-linux-x64.tar.xz &&\
#     rm blender-3.6.0-linux-x64.tar.xz &&\
#     mv blender-3.6.0-linux-x64 blender &&\
#     rm -rf /blender/3.6/python &&\
#     rm -rf /var/lib/apt/lists/* &&\
#     echo '#!/bin/sh\n/blender/blender -noaudio -b --python-use-system-env $@' > /usr/bin/blender &&\
#     chmod +x /usr/bin/blender

RUN apt-get update && apt-get install -y libegl1-mesa libgl1-mesa-dri libsm6 libxext6 libxrender-dev libxi6 wget xz-utils libgl-dev libxi-dev libxrender-dev libxkbcommon-x11-0 && \
    pip install -U pip wheel setuptools && \
    pip install numpy requests && \
    wget -q https://download.blender.org/release/Blender3.6/blender-3.6.0-linux-x64.tar.xz && \
    tar -xf blender-3.6.0-linux-x64.tar.xz && \
    rm blender-3.6.0-linux-x64.tar.xz && \
    mv blender-3.6.0-linux-x64 blender && \
    rm -rf /blender/3.6/python && \
    rm -rf /var/lib/apt/lists/* && \
    echo '#!/bin/sh\n/blender/blender -noaudio -b --python-use-system-env $@' > /usr/bin/blender && \
    chmod +x /usr/bin/blender

ENV PYTHONUNBUFFERED=1
ENV HYDRA_FULL_ERROR=1

# Add so jax can find cuDnn library
ENV LD_LIBRARY_PATH=/opt/conda/lib/python3.10/site-packages/torch/lib/

# entrypoint
ENTRYPOINT ["python", "-u", "main.py"]
