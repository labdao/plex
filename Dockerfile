FROM golang:1.20-buster as builder
 
# Install deps
RUN apt-get update && apt-get install -y \
  libssl-dev \
  ca-certificates \
  fuse

COPY . /app/
RUN cd /app/ \
    && CGO_ENABLED=0 go build -o /go/bin/plex

RUN apt-get update && apt-get -y install ca-certificates

FROM busybox:1.31.1-glibc

COPY --from=builder /go/bin/plex /plex
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/gateway/migrations  /gateway/migrations

# Copy custom IPFS binary with s3ds and healthcheck plugin
COPY --from=quay.io/labdao/ipfs@sha256:461646b6ea97dffc86b1816380360be3d38d5a2c6c7c86352a2e3b0a5a4ccca5 /usr/local/bin/ipfs /usr/local/bin/ipfs

# Copy init script from ipfs image
COPY --from=quay.io/labdao/ipfs@sha256:461646b6ea97dffc86b1816380360be3d38d5a2c6c7c86352a2e3b0a5a4ccca5 /usr/local/bin/container_init_run /usr/local/bin/container_init_run

# Copy container-init
COPY docker/images/ipfs/container-init.d /container-init.d

# init.d script IPFS runs before starting the daemon. Used to manipulate the IPFS config file.
COPY docker/images/backend/docker-entrypoint.sh /docker-entrypoint.sh

# Copy jq
COPY --from=ghcr.io/jqlang/jq /jq /usr/local/bin/jq

# This shared lib (part of glibc) doesn't seem to be included with busybox.
COPY --from=builder /lib/*-linux-gnu*/libdl.so.2 /lib/

# Copy over SSL libraries.
COPY --from=builder /usr/lib/*-linux-gnu*/libssl.so* /usr/lib/
COPY --from=builder /usr/lib/*-linux-gnu*/libcrypto.so* /usr/lib/

RUN chmod +x /docker-entrypoint.sh

RUN mkdir -p /data/ipfs

ENV POSTGRES_PASSWORD=MAKE_UP_SOMETHING_RANDOM
ENV POSTGRES_USER=labdao
ENV POSTGRES_DB=labdao
ENV POSTGRES_HOST=localhost
ENV FRONTEND_URL=http://localhost:3080
ENV BACALHAU_API_HOST=127.0.0.1
ENV IPFS_PATH=/data/ipfs
ENV IPFS_PROFILE=server
ENV BACALHAU_SERVE_IPFS_PATH=/data/ipfs

# Needed until we figure out a better way to set the bacalhau config file
# and a better way to stream logs through the bacalhau Go pkg
# Update the package list and install curl and bash
RUN apt-get update && apt-get install -y \
    curl \
    bash \
    && rm -rf /var/lib/apt/lists/*
RUN curl -sL https://get.bacalhau.org/install.sh | bash

EXPOSE 8080

ENTRYPOINT ["/docker-entrypoint.sh"]

CMD ["web"]
