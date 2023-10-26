# Use the official Ubuntu 20.04 base image
FROM python:3.9.18-slim

# Set environment variable to non-interactive (this prevents some prompts)
ENV DEBIAN_FRONTEND=noninteractive

COPY conf /app/conf
COPY install.py requirements.txt main.py /app
COPY conf/inputs /app/conf/inputs
COPY conf/outputs /app/conf/inputs

RUN apt-get update && apt-get -y install git

WORKDIR /app

RUN mkdir -p /app/conf/params /app/conf/inputs /app/conf/outputs
ADD https://files.ipd.uw.edu/krypton/schedules.zip /tmp/
ADD http://files.ipd.uw.edu/pub/RFdiffusion/60f09a193fb5e5ccdc4980417708dbab/Complex_Fold_base_ckpt.pt /tmp/
ADD https://storage.googleapis.com/alphafold/alphafold_params_2022-12-06.tar /tmp/params.tar
RUN tar -xvf /tmp/params.tar -C /app/conf/params/
RUN git clone https://github.com/sokrypton/RFdiffusion.git \
     && pip -q install jedi omegaconf hydra-core icecream pyrsistent \
     && pip install dgl==1.0.2+cu116 -f https://data.dgl.ai/wheels/cu116/repo.html \
     && cd RFdiffusion/env/SE3Transformer; pip -q install --no-cache-dir -r requirements.txt; pip -q install . \
     && pip -q install py3Dmol pydssp

ADD https://raw.githubusercontent.com/sokrypton/ColabDesign/v1.1.1/colabdesign/rf/blueprint.js
ADD https://raw.githubusercontent.com/sokrypton/ColabDesign/v1.1.1/colabdesign/rf/blueprint.css

if not os.path.isdir("RFdiffusion/models"):
  print("downloading RFdiffusion params...")
  os.system("mkdir RFdiffusion/models")
  models = ["Complex_Fold_base_ckpt.pt"]
  for m in models:
    while os.path.isfile(f"{m}.aria2"):
      time.sleep(5)
  os.system(f"mv {' '.join(models)} RFdiffusion/models")
  os.system("unzip schedules.zip; rm schedules.zip")
  print("----------------------------------")



# # SETUP RFDIFFUSION
# if not os.path.isdir("RFdiffusion"):
#   print("installing RFdiffusion...")
#   os.system("apt-get install aria2")
#   os.system("mkdir params")
#   os.system("(\
#     aria2c -q -x 16 https://files.ipd.uw.edu/krypton/schedules.zip; \
#     aria2c -q -x 16 http://files.ipd.uw.edu/pub/RFdiffusion/60f09a193fb5e5ccdc4980417708dbab/Complex_Fold_base_ckpt.pt; \
#     aria2c -q -x 16 https://storage.googleapis.com/alphafold/alphafold_params_2022-12-06.tar; \
#     tar -xf alphafold_params_2022-12-06.tar -C params; \
#     touch params/done.txt) &")
#
#   os.system("git clone https://github.com/sokrypton/RFdiffusion.git")
#   os.system("pip -q install jedi omegaconf hydra-core icecream pyrsistent")
#   os.system("pip install dgl==1.0.2+cu116 -f https://data.dgl.ai/wheels/cu116/repo.html")
#   os.system("cd RFdiffusion/env/SE3Transformer; pip -q install --no-cache-dir -r requirements.txt; pip -q install .")
#
#   os.system("pip -q install py3Dmol pydssp")
#   os.system("wget -qnc https://raw.githubusercontent.com/sokrypton/ColabDesign/v1.1.1/colabdesign/rf/blueprint.js")
#   os.system("wget -qnc https://raw.githubusercontent.com/sokrypton/ColabDesign/v1.1.1/colabdesign/rf/blueprint.css")
#
# if not os.path.isdir("RFdiffusion/models"):
#   print("downloading RFdiffusion params...")
#   os.system("mkdir RFdiffusion/models")
#   models = ["Complex_Fold_base_ckpt.pt"]
#   for m in models:
#     while os.path.isfile(f"{m}.aria2"):
#       time.sleep(5)
#   os.system(f"mv {' '.join(models)} RFdiffusion/models")
#   os.system("unzip schedules.zip; rm schedules.zip")
#   print("----------------------------------")
#
# if 'RFdiffusion' not in sys.path:
#   os.environ["DGLBACKEND"] = "pytorch"
#   sys.path.append('RFdiffusion')
#
# if not os.path.isdir("colabdesign"):
#   print("installing ColabDesign...")
#   os.system("pip -q install git+https://github.com/sokrypton/ColabDesign.git")
#   os.system("ln -s /usr/local/lib/python3.*/dist-packages/colabdesign colabdesign")
#
# #install prodigy
# os.system("git clone -q https://github.com/haddocking/prodigy")
# # os.system("pip install -q /content/prodigy/")
# os.system("pip install -q prodigy/")

#
# # Command to run when starting the container
# CMD ["/bin/bash"]
