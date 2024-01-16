# Use Python 3.9 slim as the base image
FROM python:3.9-slim

# Set the working directory to /app
WORKDIR /app

# Copy the main.py file from the local directory to the /app directory in the container
COPY main.py /app

# Set the entry point to run the main.py script with unbuffered output
ENTRYPOINT ["python", "-u", "main.py"]
