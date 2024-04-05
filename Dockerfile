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

ENV POSTGRES_PASSWORD=MAKE_UP_SOMETHING_RANDOM
ENV POSTGRES_USER=labdao
ENV POSTGRES_DB=labdao
ENV POSTGRES_HOST=localhost
ENV FRONTEND_URL=http://localhost:3080
ENV IPFS_API_HOST=ipfs
ENV MAX_QUEUE_TIME_SECONDS=259200 # default 72 hours
ENV MAX_COMPUTE_TIME_SECONDS=259200 # default 72 hours
ENV STRIPE_PRODUCT_SLUG=price_1OlwVPE7xzGf7nZbaccQCnHv
# ENV STRIPE_SECRET_KEY: ${STRIPE_SECRET_KEY}
# ENV STRIPE_WEBHOOK_SECRET_KEY: ${STRIPE_WEBHOOK_SECRET_KEY}
# ENV NEXT_PUBLIC_PRIVY_APP_ID: ${NEXT_PUBLIC_PRIVY_APP_ID}
# ENV PRIVY_PUBLIC_KEY: ${PRIVY_PUBLIC_KEY}
# ENV AWS_ACCESS_KEY_ID: ${AWS_ACCESS_KEY_ID}
# ENV AWS_SECRET_ACCESS_KEY : ${AWS_SECRET_ACCESS_KEY}


EXPOSE 8080

ENTRYPOINT ["/plex"]

CMD ["web"]
