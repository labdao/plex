FROM golang:1.20-buster as builder

# Install deps
RUN apt-get update && apt-get install -y --no-install-recommends \
  libssl-dev \
  ca-certificates \
  fuse

COPY . /app/
WORKDIR /app/

# Use the race command to search for concurrency issues
# RUN CGO_ENABLED=1 go build -race -o /go/bin/plex
RUN CGO_ENABLED=0 go build -o /go/bin/plex

ARG BACALHAU_VERSION=1.2.0
ARG NEXT_PUBLIC_PRIVY_APP_ID
ARG PRIVY_PUBLIC_KEY

# For bacalhau cli
FROM ghcr.io/bacalhau-project/bacalhau:v${BACALHAU_VERSION:-1.2.0} as bacalhau

FROM busybox:1.31.1-glibc

COPY --from=builder /go/bin/plex /plex
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/gateway/migrations  /gateway/migrations

# Copy jq
COPY --from=ghcr.io/jqlang/jq /jq /usr/local/bin/jq

# This shared lib (part of glibc) doesn't seem to be included with busybox.
COPY --from=builder /lib/*-linux-gnu*/libdl.so.2 /lib/

# Copy over SSL libraries.
COPY --from=builder /usr/lib/*-linux-gnu*/libssl.so* /usr/lib/
COPY --from=builder /usr/lib/*-linux-gnu*/libcrypto.so* /usr/lib/

# COPY bacalhau cli
COPY --from=bacalhau --chmod=755 /usr/local/bin/bacalhau /usr/local/bin/bacalhau

# This creates config file needed by bacalhau golang client
RUN /usr/local/bin/bacalhau version
RUN /usr/local/bin/bacalhau config default > /root/.bacalhau/config.yaml

ENV POSTGRES_PASSWORD=MAKE_UP_SOMETHING_RANDOM
ENV POSTGRES_USER=labdao
ENV POSTGRES_DB=labdao
ENV POSTGRES_HOST=localhost
ENV FRONTEND_URL=http://localhost:3080
ENV BACALHAU_API_HOST=127.0.0.1

EXPOSE 8080

ENTRYPOINT ["/plex"]

CMD ["web"]
