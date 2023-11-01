# FROM us-docker.pkg.dev/colab-images/public/runtime
FROM python:3.9.18-slim-bullseye 
RUN apt update -y \
    && apt install -y wget nano
RUN pip3 install hydra-core --upgrade
RUN mkdir /app
# installation
COPY install.py /app
RUN cd app && python3 install.py
COPY . /app
# entrypoint
WORKDIR /app
ENTRYPOINT []
CMD ["bash"]
