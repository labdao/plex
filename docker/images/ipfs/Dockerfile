ARG IPFS_VERSION=0.18.0

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

# Copy custom IPFS binary with s3ds and healthcheck plugin
COPY --from=ceramicnetwork/go-ipfs-daemon:daaf07a9dbd007a3a6378f34f2a404992ade5e93 /usr/local/bin/ipfs /usr/local/bin/ipfs

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
