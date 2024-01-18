# Use Python 3.9 slim as the base image
FROM python:3.9-slim

# Set the working directory to /app
WORKDIR /app

COPY requirements.txt .
COPY main.py /app

RUN pip install -r requirements.txt

# Set the entry point to run the main.py script with unbuffered output
ENTRYPOINT ["python", "-u", "main.py"]
