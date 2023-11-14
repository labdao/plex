# Use the official Python base image
FROM python:3.9-slim

# Set up directories to match the expected script paths
RUN mkdir -p /inputs
RUN mkdir -p /outputs

# Set the working directory in the container
WORKDIR /app

# Copy the Python script into the container at /app
COPY main.py /app/

# When the container starts, it will be ready to run main.py with arguments
CMD ["python", "-u", "/app/main.py"]
