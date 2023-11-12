# Use the official Python base image
FROM python:3.9-slim

# Set up directories to match the expected script paths
RUN mkdir /inputs /outputs

# Set the working directory in the container
WORKDIR /app

# Copy the Python script into the container at /app
COPY main.py /app/

# Use ENTRYPOINT to specify the executable (python) and the script path
ENTRYPOINT ["python", "-u", "/app/main.py"]
