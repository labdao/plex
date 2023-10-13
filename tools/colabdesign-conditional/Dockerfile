# Use the official Ubuntu 20.04 base image
FROM niklastr/dev:latest

# Set environment variable to non-interactive (this prevents some prompts)
ENV DEBIAN_FRONTEND=noninteractive

# Update the package list and install the specified packages
RUN apt-get update

WORKDIR /app
RUN mkdir -p /app/conf
COPY install.py /app
RUN python install.py
COPY requirements.txt /app
RUN pip3 install -r requirements.txt
COPY conf* /app/conf
COPY main.py /app

# Command to run when starting the container
CMD ["/bin/bash"]