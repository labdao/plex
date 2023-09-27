ARG IPFS_VERSION=0.18.0
FROM golang:1.19.1-buster AS builder

ARG IPFS_VERSION=0.18.0
ARG IPFS_DS_S3_PLUGIN_VERSION=latest
ARG IPFS_HEALTHCHECK_PLUGIN_VERSION=latest

ARG IPFS_S3_ENABLED=true

WORKDIR /

RUN cat /etc/resolv.conf

# Kubo build process
# See details: https://github.com/ipfs/go-ds-s3
ENV GO111MODULE on
ENV GOPROXY direct

# We clone Kubo source code
RUN git clone https://github.com/ipfs/kubo && cd kubo && git checkout v${IPFS_VERSION}
ENV SRC_DIR /kubo

# Move to kubo folder
WORKDIR $SRC_DIR

# Install the plugin and build ipfs
RUN if [ ${IPFS_S3_ENABLED} == "true" ];then echo "Building with S3 PLUGIN Enabled";go get -v github.com/3box/go-ds-s3/plugin@${IPFS_DS_S3_PLUGIN_VERSION};fi
RUN if [ ${IPFS_S3_ENABLED} == "true" ];then echo "Building with S3 PLUGIN Enabled";echo "\ns3ds github.com/3box/go-ds-s3/plugin 0" >> plugin/loader/preload_list;fi
RUN go get -v github.com/ceramicnetwork/go-ipfs-healthcheck/plugin@${IPFS_HEALTHCHECK_PLUGIN_VERSION}
RUN echo "\nhealthcheck github.com/ceramicnetwork/go-ipfs-healthcheck/plugin 0" >> plugin/loader/preload_list
RUN make build
RUN go mod tidy
RUN CGO_ENABLED=0 make build
RUN make install

# The actual IPFS image we will use
FROM ipfs/kubo:v${IPFS_VERSION}
ENV SRC_DIR /kubo

# We copy the new binaries we built in the 'builder' stage (--from=builder)
COPY --from=builder $SRC_DIR/cmd/ipfs/ipfs /usr/local/bin/ipfs
COPY --from=builder $SRC_DIR/bin/container_daemon /usr/local/bin/start_ipfs
COPY --from=builder $SRC_DIR/bin/container_init_run /usr/local/bin/container_init_run

# Fix permissions on start_ipfs
RUN chmod 0755 /usr/local/bin/start_ipfs

# Copy jq
COPY --from=ghcr.io/jqlang/jq /jq /usr/local/bin/jq

# Download curl
ADD https://github.com/moparisthebest/static-curl/releases/download/v8.2.1/curl-amd64 /usr/local/bin/curl
RUN chmod +x /usr/local/bin/curl

# init.d script IPFS runs before starting the daemon. Used to manipulate the IPFS config file.
COPY files/001-custom-config.sh /container-init.d/001-custom-config.sh
COPY files/002-s3-keyTransformation.sh /container-init.d/002-s3-keyTransformation.sh
COPY files/003-bootstrap-rm.sh /container-init.d/003-bootstrap-rm.sh
COPY files/004-swarm-base64-key.sh /container-init.d/004-swarm-base64-key.sh
COPY files/999-debug-config.sh /container-init.d/999-debug-config.sh

# Set ownership
RUN chown ipfs: /container-init.d/*.sh
RUN chown ipfs: /data/ipfs

USER ipfs
