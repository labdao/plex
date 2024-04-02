FROM quay.io/convexitylabs/base-images:pytorch-2.1.0-cuda11.8-cudnn8-runtime-4576d78@sha256:b53852474cafce3d4636c30906af8493a48669fc95908b8146382f23dfd28426

# Copy in Files
COPY main.py /app
COPY conf /app/conf

ENV PYTHONUNBUFFERED=1
ENV HYDRA_FULL_ERROR=1

# entrypoint
ENTRYPOINT ["python", "-u", "main.py"]
