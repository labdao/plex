import torch
import subprocess
import sys

def format_pytorch_version(version):
  return version.split('+')[0]

def format_cuda_version(version):
  return 'cu' + version.replace('.', '')

TORCH_version = torch.__version__
TORCH = format_pytorch_version(TORCH_version)

CUDA_version = torch.version.cuda
CUDA = format_cuda_version(CUDA_version)

def install_package(package):
    subprocess.check_call([sys.executable, "-m", "pip", "install", package])

packages = [
    f"torch-scatter -f https://pytorch-geometric.com/whl/torch-{TORCH}+{CUDA}.html",
    f"torch-sparse -f https://pytorch-geometric.com/whl/torch-{TORCH}+{CUDA}.html",
    f"torch-cluster -f https://pytorch-geometric.com/whl/torch-{TORCH}+{CUDA}.html",
    f"torch-spline-conv -f https://pytorch-geometric.com/whl/torch-{TORCH}+{CUDA}.html",
    "torch-geometric"
]

for package in packages:
    install_package(package)
