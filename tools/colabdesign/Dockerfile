FROM us-docker.pkg.dev/colab-images/public/runtime

# Remove unwanted pip packages to slim down image
RUN pip3 install --no-cache-dir pip-autoremove
RUN pip-autoremove tensorflow -y

# Remove /tools directory
RUN rm -rf /tools/

RUN apt update -y \
    && apt install --no-install-recommends -y wget aria2
RUN pip3 install --no-cache-dir hydra-core --upgrade
RUN mkdir /app
# installation
COPY . /app
WORKDIR /app
RUN python3 install.py

# Cleanup
RUN apt-get clean autoclean \
    && apt-get autoremove -y \
    && rm -rf /var/lib/cache /var/lib/log /var/lib/apt/lists/*
# entrypoint
ENTRYPOINT []
CMD ["bash"]
