# Use the official Python base image
FROM python:3.9-slim

# Set the working directory to /app
WORKDIR /app

# Set up directories to match the expected script paths
RUN mkdir -p /inputs/fileparam
RUN mkdir -p /outputs

# Copy the Python script into the container
COPY main.py /app/

# Make the results directory writable by non-root users (useful for some Kubernetes platforms)
RUN chmod -R 777 /outputs

# Set the default command for the container to run the script with the example arguments
CMD ["python", "-u", "main.py", "exampleString", "123"]
