# Use Python 3.9 slim as the base image
FROM python:3.9-slim

# Set the working directory to /app
WORKDIR /app

# Install ray serve
RUN pip3 install --no-cache-dir ray[serve]

# Copy the main.py file from the local directory to the /app directory in the container
COPY app.py /app

ENTRYPOINT ["serve", "run", "app:labray", "--host", "0.0.0.0"]
