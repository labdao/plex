# FROM us-docker.pkg.dev/colab-images/public/runtime
FROM docker.io/pytorch/pytorch:2.1.0-cuda11.8-cudnn8-runtime

ENV DEBIAN_FRONTEND=noninteractive

# Install required packages
RUN apt-get update -y && \
     apt-get install --no-install-recommends -y wget aria2 git zip unzip less nano build-essential

# Install hydra-core
RUN pip3 install --no-cache-dir hydra-core boto3 --upgrade

# Install missing dependencies
RUN set -ex; \
     pip3 install ipywidgets httplib2 && \
     pip3 install -U numpy && \
     pip3 install -U 'jax[cuda11_pip]' -f https://storage.googleapis.com/jax-releases/jax_cuda_releases.html

WORKDIR /app

# RUN python3 -u install.py

# Install Google colab
RUN pip -q install --no-cache-dir git+https://github.com/googlecolab/colabtools.git@1b6462d0e66f5b17ce77c6a17d3a366e327b7096

# Install ColabDesign
RUN set -ex; \
     pip -q install --no-cache-dir git+https://github.com/sokrypton/ColabDesign.git && \
     ln -s /opt/conda/lib/python3.10/site-packages/colabdesign colabdesign

# Install ananas
ADD https://files.ipd.uw.edu/krypton/ananas ananas
RUN chmod +x ananas

# Install RFdiffusion
RUN set -ex; \
     git clone --depth 1 https://github.com/sokrypton/RFdiffusion.git && \
     pip -q install --no-cache-dir jedi omegaconf hydra-core icecream pyrsistent && \
     pip install --no-cache-dir dgl==1.0.2+cu116 -f https://data.dgl.ai/wheels/cu116/repo.html && \
     cd RFdiffusion/env/SE3Transformer; pip -q install --no-cache-dir -r requirements.txt; pip -q install --no-cache-dir .

# Download params, modules and schedules
RUN set -ex; \
     mkdir -p params && \
     aria2c -q -x 16 https://files.ipd.uw.edu/krypton/schedules.zip && \
     aria2c -q -x 16 http://files.ipd.uw.edu/pub/RFdiffusion/6f5902ac237024bdd0c176cb93063dc4/Base_ckpt.pt && \
     aria2c -q -x 16 http://files.ipd.uw.edu/pub/RFdiffusion/e29311f6f1bf1af907f9ef9f44b8328b/Complex_base_ckpt.pt && \
     aria2c -q -x 16 https://storage.googleapis.com/alphafold/alphafold_params_2022-12-06.tar && \
     tar -xvf alphafold_params_2022-12-06.tar -C params && \
     rm -rf alphafold_params_2022-12-06.tar && \
     mkdir -p RFdiffusion/models; mv Base_ckpt.pt Complex_base_ckpt.pt RFdiffusion/models/ && \
     unzip schedules.zip; rm schedules.zip && \
     touch params/done.txt

# Install prodigy
RUN git clone -q https://github.com/haddocking/prodigy && \
     pip install -q prodigy/

# Cleanup
RUN apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/cache /var/lib/log /var/lib/apt/lists/*

# Copy in Files
COPY . /app

ENV PYTHONUNBUFFERED=1
ENV HYDRA_FULL_ERROR=1


# entrypoint
ENTRYPOINT ["python", "-u", "main.py"]
