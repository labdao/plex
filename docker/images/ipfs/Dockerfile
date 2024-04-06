ARG IPFS_VERSION=0.18.0

# Building IPFS binary
FROM golang:1.19.1-buster AS builder

ARG IPFS_VERSION=0.18.0
ARG IPFS_DS_S3_PLUGIN_VERSION=0cc3de579103f6fd209c30e43abadf80f1024151
ARG IPFS_HEALTHCHECK_PLUGIN_VERSION=3aaefbe9799465dcfbb89d6c68d1c4bd876f3361

WORKDIR /

# Kubo build process
# See details: https://github.com/ipfs/go-ds-s3
ENV GO111MODULE on
ENV GOPROXY direct

# We clone Kubo source code
RUN git clone --depth 1 -b v${IPFS_VERSION} https://github.com/ipfs/kubo
ENV SRC_DIR /kubo

# Move to kubo folder
WORKDIR $SRC_DIR

# Install the plugin and build ipfs
RUN go get github.com/ceramicnetwork/go-ipfs-healthcheck/plugin@${IPFS_HEALTHCHECK_PLUGIN_VERSION}
RUN echo "\nhealthcheck github.com/ceramicnetwork/go-ipfs-healthcheck/plugin 0" >> plugin/loader/preload_list
RUN go get github.com/3box/go-ds-s3/plugin@${IPFS_DS_S3_PLUGIN_VERSION}
RUN echo "\ns3ds github.com/3box/go-ds-s3/plugin 0" >> plugin/loader/preload_list
RUN make build
RUN go mod tidy
RUN CGO_ENABLED=0 make build
RUN make install

# IPFS container with all the utilities
FROM ipfs/kubo:v${IPFS_VERSION} AS ipfs-with-utils

# # Copy custom IPFS binary with s3ds and healthcheck plugin
# COPY --from=ceramicnetwork/go-ipfs-daemon:daaf07a9dbd007a3a6378f34f2a404992ade5e93 /usr/local/bin/ipfs /usr/local/bin/ipfs

# Copy IPFS from locally build IPFS binary from builder
COPY --from=builder /go/bin/ipfs /usr/local/bin/ipfs

# Copy jq
COPY --from=ghcr.io/jqlang/jq /jq /usr/local/bin/jq

# Download curl
ADD https://github.com/moparisthebest/static-curl/releases/download/v8.2.1/curl-amd64 /usr/local/bin/curl
RUN chmod +x /usr/local/bin/curl

# init.d script IPFS runs before starting the daemon. Used to manipulate the IPFS config file.
COPY container-init.d container-init.d

# Set ownership
RUN chown ipfs: /container-init.d/*.sh
RUN chown ipfs: /data/ipfs

USER ipfs

FROM scratch

COPY --from=ipfs-with-utils / /

# Swarm TCP; should be exposed to the public
EXPOSE 4001
# Swarm UDP; should be exposed to the public
EXPOSE 4001/udp
# Daemon API; must not be exposed publicly but to client services under you control
EXPOSE 5001
# Web Gateway; can be exposed publicly with a proxy, e.g. as https://ipfs.example.org
EXPOSE 8080
# Swarm Websockets; must be exposed publicly when the node is listening using the websocket transport (/ipX/.../tcp/8081/ws).
EXPOSE 8081

ENV IPFS_PATH /data/ipfs

ENV IPFS_PROFILE "server"

# The default logging level
ENV IPFS_LOGGING ""

# This just makes sure that:
# 1. There's an fs-repo, and initializes one if there isn't.
# 2. The API and Gateway are accessible from outside the container.
ENTRYPOINT ["/sbin/tini", "--", "/usr/local/bin/start_ipfs"]

# Healthcheck for the container
# QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn is the CID of empty folder
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ipfs dag stat /ipfs/QmUNLLsPACCz1vLxQVkXqqLX5R1X345qqfHbsf67hvA3Nn || exit 1

# Execute the daemon subcommand by default
CMD ["daemon", "--migrate=true", "--agent-version-suffix=docker"]
