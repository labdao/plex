FROM golang:1.20-buster as builder

ARG BACALHAU_VERSION=1.1.2

# Install deps
RUN apt-get update && apt-get install -y --no-install-recommends \
  libssl-dev \
  ca-certificates \
  fuse

COPY . /app/
WORKDIR /app/
RUN CGO_ENABLED=0 go build -o /go/bin/plex

# Download bacalhau cli
ADD https://github.com/bacalhau-project/bacalhau/releases/download/v${BACALHAU_VERSION}/bacalhau_v${BACALHAU_VERSION}_linux_amd64.tar.gz /tmp/bacalhau.tgz

RUN tar -zxvf /tmp/bacalhau.tgz -C /usr/local/bin/

FROM ghcr.io/jqlang/jq as jq

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
COPY --chmod=0755 docker/images/backend/docker-entrypoint.sh /docker-entrypoint.sh

# Copy jq
COPY --from=jq /jq /usr/local/bin/jq

# This shared lib (part of glibc) doesn't seem to be included with busybox.
COPY --from=builder /lib/*-linux-gnu*/libdl.so.2 /lib/

# Copy over SSL libraries.
COPY --from=builder /usr/lib/*-linux-gnu*/libssl.so* /usr/lib/
COPY --from=builder /usr/lib/*-linux-gnu*/libcrypto.so* /usr/lib/

# COPY bacalhau cli
COPY --from=builder --chmod=755 /usr/local/bin/bacalhau /usr/local/bin/bacalhau

RUN mkdir -p /data/ipfs

# This creates config file needed by bacalhau golang client
RUN /usr/local/bin/bacalhau version

ENV POSTGRES_PASSWORD=MAKE_UP_SOMETHING_RANDOM
ENV POSTGRES_USER=labdao
ENV POSTGRES_DB=labdao
ENV POSTGRES_HOST=localhost
ENV FRONTEND_URL=http://localhost:3080
ENV BACALHAU_API_HOST=127.0.0.1
ENV IPFS_PATH=/data/ipfs
ENV IPFS_PROFILE=server
ENV BACALHAU_SERVE_IPFS_PATH=/data/ipfs

EXPOSE 8080

ENTRYPOINT ["/docker-entrypoint.sh"]

CMD ["web"]
