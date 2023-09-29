ARG IPFS_VERSION=0.18.0
FROM golang:1.19.1-buster AS builder

ARG IPFS_VERSION=0.18.0
ARG IPFS_DS_S3_PLUGIN_VERSION=latest
ARG IPFS_HEALTHCHECK_PLUGIN_VERSION=latest

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

# Final IPFS container with custom IPFS binary and utils
FROM ipfs-with-utils
ENV SRC_DIR /kubo

# We copy the new binaries we built in the 'builder' stage (--from=builder)
COPY --from=builder $SRC_DIR/cmd/ipfs/ipfs /usr/local/bin/ipfs

USER ipfs
