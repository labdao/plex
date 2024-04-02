FROM quay.io/convexitylabs/base-images:pytorch-2.1.0-cuda11.8-cudnn8-runtime-af9f8b5@sha256:5828a2b0f27db1fcfe8b39475acf44bb4a3c93b8aad444d4f70b837a75cc0b3e

# Copy in Files
COPY main.py /app
COPY conf /app/conf

ENV PYTHONUNBUFFERED=1
ENV HYDRA_FULL_ERROR=1

# Add so jax can find cuDnn library
ENV LD_LIBRARY_PATH=/opt/conda/lib/python3.10/site-packages/torch/lib/

WORKDIR /app

# entrypoint
ENTRYPOINT ["python", "-u", "main.py"]
