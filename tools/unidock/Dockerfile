FROM us-docker.pkg.dev/colab-images/public/runtime:latest

RUN apt update -y \
    && apt install -y wget

RUN wget https://github.com/dptech-corp/Uni-Dock/releases/download/1.0.0/unidock   \
    && chmod +x unidock

COPY . .
RUN chmod +x parse_docking_result.sh