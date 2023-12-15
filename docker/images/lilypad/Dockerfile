FROM public.ecr.aws/docker/library/golang:1.20-buster as builder

# Install deps
RUN apt-get update && apt-get install -y --no-install-recommends \
  libssl-dev \
  ca-certificates \
  fuse

ARG LILYPAD_VERSION=v2.0.0-e5c38d5

WORKDIR /lilypad

# Clone Lilypad project
RUN git clone https://github.com/bacalhau-project/lilypad.git --branch ${LILYPAD_VERSION}

# Build
RUN cd lilypad && CGO_ENABLED=0 go build -o /go/bin/lilypad

FROM public.ecr.aws/docker/library/busybox:1.31.1-glibc

COPY --from=builder /go/bin/lilypad /usr/local/bin/lilypad
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

COPY run.sh /run.sh

RUN chmod +x /usr/local/bin/lilypad /run.sh

RUN mkdir -p /var/tmp/

ENV PATH="/usr/local/bin:/usr/bin"
ENTRYPOINT ["/run.sh"]
