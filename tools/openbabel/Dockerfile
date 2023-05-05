# Use an official Python base image
FROM python:3.9-slim-buster

# Update the package list and install necessary packages
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
        openbabel \
        libopenbabel-dev \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# Set the working directory
WORKDIR /app

# Copy the current directory contents into the container
COPY . .

CMD [ "bash"]