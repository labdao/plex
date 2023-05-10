FROM huggingface/transformers-pytorch-gpu
LABEL container="protbert"

RUN apt update

COPY requirements.txt .

RUN pip install -r requirements.txt

COPY . .

# test battery
RUN chmod +x test.sh && ./test.sh