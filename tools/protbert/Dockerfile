#  Pulls the base Docker image from the Hugging Face Transformers library, specifically an image that supports PyTorch and GPU processing.
FROM huggingface/transformers-pytorch-gpu
# Labels the Docker image with the name "protbert".
LABEL container="protbert"

# Updates the package lists for upgrades and new packages from repositories.
RUN apt update

# Copies the requirements.txt file from the local system to the Docker image. This file should list the Python packages that the application requires.
COPY requirements.txt .

# Installs the Python packages listed in requirements.txt using pip.
RUN pip install -r requirements.txt

# Copies the contents of the current directory to the Docker image.
COPY . .

# test battery
RUN chmod +x test.sh && ./test.sh