FROM us-docker.pkg.dev/colab-images/public/runtime:latest

RUN apt update -y \
    && apt install -y wget

RUN wget https://github.com/dptech-corp/Uni-Dock/releases/download/1.0.0/unidock   \
    && chmod +x unidock

COPY . .

# run test
RUN ./unidock --receptor AF-15552_reduced.pdbqt --search_mode fast --verbosity 2 --size_x 30 --size_y 30 --size_z 30 --center_x 0 --center_y 0 --center_z 0 --ligand output.pdbqt
