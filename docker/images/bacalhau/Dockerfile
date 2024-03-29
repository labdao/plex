# syntax=docker/dockerfile:1.4
# Pinned to 1.20
FROM golang:1.20 as build

# Release tag. Used to build the binary and tag the version.
ARG BACALHAU_VERSION=v1.2.0
ARG TAG

RUN apt-get update && apt-get -y --no-install-recommends install nodejs npm

RUN git clone https://github.com/bacalhau-project/bacalhau.git /work

WORKDIR /work

RUN git checkout ${BACALHAU_VERSION}
RUN make build-bacalhau
RUN find ./bin -name 'bacalhau' -exec mv -t ./bin {} +

FROM cgr.dev/chainguard/nvidia-device-plugin
COPY --from=build /work/bin/bacalhau /usr/local/bin/bacalhau
ENV PATH="/usr/local/bin:/usr/bin"
ENTRYPOINT ["bacalhau"]
LABEL org.opencontainers.image.source https://github.com/bacalhau-project/bacalhau
LABEL org.opencontainers.image.title "Bacalhau"
LABEL org.opencontainers.image.description "The Bacalhau network provices decentralised compute for compute over data. See https://bacalhau.org for more info."
LABEL org.opencontainers.image.licenses Apache-2.0
LABEL org.opencontainers.image.url https://bacalhau.org
