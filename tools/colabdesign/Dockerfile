FROM quay.io/convexitylabs/base-images:pytorch-2.1.0-cuda11.8-cudnn8-runtime-3068995@sha256:8f75299825d3dc613d32070ac72b99439751054c06b8af994c00205954b4ee4e

# Copy in Files
COPY main.py /app
COPY conf /app/conf

ENV PYTHONUNBUFFERED=1
ENV HYDRA_FULL_ERROR=1

# Add so jax can find cuDnn library
ENV LD_LIBRARY_PATH=/opt/conda/lib/python3.10/site-packages/torch/lib/

# entrypoint
ENTRYPOINT ["python", "-u", "main.py"]
