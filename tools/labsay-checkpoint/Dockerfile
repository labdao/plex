FROM python:3.9-slim

WORKDIR /app

COPY main.py /app

RUN pip install --no-cache-dir boto3

ENTRYPOINT ["python", "-u", "main.py"]