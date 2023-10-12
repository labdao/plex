FROM golang:1.20 as builder
 
COPY . /app/
RUN cd /app/ \
    && CGO_ENABLED=0 go build -o /go/bin/plex

RUN apt-get update && apt-get -y install ca-certificates

FROM alpine

COPY --from=builder /go/bin/plex /plex
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /app/gateway/migrations  /gateway/migrations

# Copy custom IPFS binary with s3ds and healthcheck plugin
COPY --from=quay.io/labdao/ipfs@sha256:461646b6ea97dffc86b1816380360be3d38d5a2c6c7c86352a2e3b0a5a4ccca5 /usr/local/bin/ipfs /usr/local/bin/ipfs

# Copy init script from ipfs image
COPY --from=quay.io/labdao/ipfs@sha256:461646b6ea97dffc86b1816380360be3d38d5a2c6c7c86352a2e3b0a5a4ccca5 /usr/local/bin/container_init_run /usr/local/bin/container_init_run

# Copy init script from ipfs image
COPY --from=quay.io/labdao/ipfs@sha256:461646b6ea97dffc86b1816380360be3d38d5a2c6c7c86352a2e3b0a5a4ccca5 /container-init.d /container-init.d

# init.d script IPFS runs before starting the daemon. Used to manipulate the IPFS config file.
COPY docker/images/backend/docker-entrypoint.sh /docker-entrypoint.sh

ENV POSTGRES_PASSWORD=MAKE_UP_SOMETHING_RANDOM
ENV POSTGRES_USER=labdao
ENV POSTGRES_DB=labdao
ENV POSTGRES_HOST=localhost
ENV FRONTEND_URL=http://localhost:3080
ENV BACALHAU_API_HOST=127.0.0.1
ENV IPFS_API_HOST=127.0.0.1
ENV IPFS_PATH=/data/ipfs

EXPOSE 8080

ENTRYPOINT ["/docker-entrypoint.sh"]

CMD ["web"]
