import os, time, signal
import sys, random, string, re


os.system("apt-get install aria2")
os.system("mkdir params")
# send param download into background
os.system("(\
aria2c -q -x 16 https://files.ipd.uw.edu/krypton/schedules.zip; \
aria2c -q -x 16 http://files.ipd.uw.edu/pub/RFdiffusion/6f5902ac237024bdd0c176cb93063dc4/Base_ckpt.pt; \
aria2c -q -x 16 http://files.ipd.uw.edu/pub/RFdiffusion/e29311f6f1bf1af907f9ef9f44b8328b/Complex_base_ckpt.pt; \
aria2c -q -x 16 http://files.ipd.uw.edu/pub/RFdiffusion/f572d396fae9206628714fb2ce00f72e/Complex_beta_ckpt.pt; \
aria2c -q -x 16 https://storage.googleapis.com/alphafold/alphafold_params_2022-12-06.tar; \
tar -xf alphafold_params_2022-12-06.tar -C params; \
touch params/done.txt) &")


print("installing RFdiffusion...")
os.system("git clone https://github.com/sokrypton/RFdiffusion.git")
os.system("pip -q install jedi omegaconf hydra-core icecream pyrsistent")
os.system("pip install dgl==1.0.2+cu116 -f https://data.dgl.ai/wheels/cu116/repo.html")
os.system("cd RFdiffusion/env/SE3Transformer; pip -q install --no-cache-dir -r requirements.txt; pip -q install .")
os.system("wget -qnc https://files.ipd.uw.edu/krypton/ananas")
os.system("chmod +x ananas")


print("installing ColabDesign...")
os.system("pip -q install git+https://github.com/sokrypton/ColabDesign.git")
os.system("ln -s /usr/local/lib/python3.*/dist-packages/colabdesign colabdesign")


print("downloading RFdiffusion params...")
os.system("mkdir RFdiffusion/models")
models = ["Base_ckpt.pt","Complex_base_ckpt.pt","Complex_beta_ckpt.pt"]
for m in models:
  while os.path.isfile(f"{m}.aria2"):
    time.sleep(5)
os.system(f"mv {' '.join(models)} RFdiffusion/models")
os.system("unzip schedules.zip; rm schedules.zip")

#if 'RFdiffusion' not in sys.path:
os.environ["DGLBACKEND"] = "pytorch"
sys.path.append('RFdiffusion')
