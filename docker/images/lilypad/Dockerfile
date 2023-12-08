# Pinned to 1.20
FROM cgr.dev/chainguard/go@sha256:8ed3fdc8f6375a3fd84b4b8b696a2366c3a639931aab492d6f92ca917e726ad6 as build
ARG LILYPAD_VERSION=v2.0.0-e5c38d5

WORKDIR /work

# Clone Lilypad project
RUN git clone https://github.com/bacalhau-project/lilypad.git --branch ${LILYPAD_VERSION}

# Build
RUN cd lilypad && go build

# Base image
FROM cgr.dev/chainguard/nvidia-device-plugin

COPY --from=build /work/lilypad/lilypad /usr/local/bin/lilypad
ENV PATH="/usr/local/bin:/usr/bin"
ENTRYPOINT ["lilypad"]
