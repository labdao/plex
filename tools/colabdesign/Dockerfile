FROM us-docker.pkg.dev/colab-images/public/runtime

# Remove unwanted pip packages and directories
RUN pip3 install --no-cache-dir pip-autoremove && \
     pip-autoremove tensorflow -y && \
     rm -rf /tools/

RUN apt-get update -y && \
     apt-get install --no-install-recommends -y wget aria2

# Install hydra-core
RUN pip3 install --no-cache-dir hydra-core --upgrade

# Copy in Files
COPY . /app

WORKDIR /app

# RUN python3 -u install.py

# Install ColabDesign
RUN pip -q install --no-cache-dir git+https://github.com/sokrypton/ColabDesign.git && \
     ln -s /usr/local/lib/python3.*/dist-packages/colabdesign colabdesign

# Install ananas
ADD https://files.ipd.uw.edu/krypton/ananas ananas
RUN chmod +x ananas

# Install RFdiffusion
RUN git clone --depth 1 https://github.com/sokrypton/RFdiffusion.git && \
     pip -q install --no-cache-dir jedi omegaconf hydra-core icecream pyrsistent && \
     pip install --no-cache-dir dgl==1.0.2+cu116 -f https://data.dgl.ai/wheels/cu116/repo.html && \
     cd RFdiffusion/env/SE3Transformer; pip -q install --no-cache-dir -r requirements.txt; pip -q install --no-cache-dir .

# Download params, modules and schedules
RUN mkdir -p params && \
     aria2c -q -x 16 https://files.ipd.uw.edu/krypton/schedules.zip && \
     aria2c -q -x 16 http://files.ipd.uw.edu/pub/RFdiffusion/6f5902ac237024bdd0c176cb93063dc4/Base_ckpt.pt && \
     aria2c -q -x 16 http://files.ipd.uw.edu/pub/RFdiffusion/e29311f6f1bf1af907f9ef9f44b8328b/Complex_base_ckpt.pt && \
     aria2c -q -x 16 https://storage.googleapis.com/alphafold/alphafold_params_2022-12-06.tar && \
     tar -xvf alphafold_params_2022-12-06.tar -C params && \
     rm -rf alphafold_params_2022-12-06.tar && \
     mkdir -p RFdiffusion/models; mv Base_ckpt.pt Complex_base_ckpt.pt Complex_beta_ckpt.pt RFdiffusion/models/ && \
     unzip schedules.zip; rm schedules.zip && \
     touch params/done.txt

# Cleanup
RUN apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/cache /var/lib/log /var/lib/apt/lists/*

# entrypoint
ENTRYPOINT []
CMD ["bash"]
