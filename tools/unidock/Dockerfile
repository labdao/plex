FROM us-docker.pkg.dev/colab-images/public/runtime:latest

RUN yum update -y \
    && yum install -y wget

RUN wget https://github.com/dptech-corp/Uni-Dock/releases/download/1.0.0/unidock

RUN yum install -y wget \
    libgomp

COPY . .
