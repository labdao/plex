FROM golang:1.20 as builder
 
COPY . /app/
RUN cd /app/ \
    && CGO_ENABLED=0 go build -o /go/bin/plex

RUN apt-get update && apt-get -y install ca-certificates

FROM alpine

COPY --from=builder /go/bin/plex /plex
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENV POSTGRES_PASSWORD=MAKE_UP_SOMETHING_RANDOM
ENV POSTGRES_USER=labdao
ENV POSTGRES_DB=labdao
ENV POSTGRES_HOST=localhost
ENV FRONTEND_URL=http://localhost:3080

EXPOSE 8080

ENTRYPOINT ["/plex"]

CMD ["web"]
