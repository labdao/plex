# FROM us-docker.pkg.dev/colab-images/public/runtime
FROM python:3.10.12-slim-bullseye 
RUN apt update -y \
    && apt install --no-install-recommends -y wget unzip git aria2
RUN pip3 install --no-cache-dir hydra-core --upgrade
RUN mkdir /app
# installation
RUN pip3 install --no-cache-dir numpy cython3
RUN pip3 install --no-cache-dir matplotlib ipywidgets google-colab
COPY . /app
WORKDIR /app
RUN python3 install.py

# Cleanup
RUN apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/cache /var/lib/log /var/lib/apt/lists/*

# entrypoint
ENTRYPOINT []
CMD ["bash"]
