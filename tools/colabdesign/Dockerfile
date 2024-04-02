FROM quay.io/convexitylabs/base-images:pytorch-2.1.0-cuda11.8-cudnn8-runtime-f55a656@sha256:285903a49701ce19843917b9bf35e64980654cdaf1c9c165aecf3907bbc35730

# Copy in Files
COPY main.py /app
COPY conf /app/conf

ENV PYTHONUNBUFFERED=1
ENV HYDRA_FULL_ERROR=1

# Add so jax can find cuDnn library
ENV LD_LIBRARY_PATH=/opt/conda/lib/python3.10/site-packages/torch/lib/

# entrypoint
ENTRYPOINT ["python", "-u", "main.py"]
