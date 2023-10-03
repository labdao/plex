FROM golang:1.18 as builder
 
ADD https://truststore.pki.rds.amazonaws.com/global/global-bundle.pem /global-bundle.pem

COPY . /app/
RUN cd /app/ \
    && CGO_ENABLED=0 go build -o /go/bin/receptor

FROM scratch

COPY --from=builder /go/bin/receptor /receptor
COPY --from=builder /global-bundle.pem /global-bundle.pem

ENV PGHOST=postgres.example.com
ENV PGUSER=receptor
ENV PGDATABASE=receptor
ENV PGSSLMODE=verify-ca
ENV PGSSLROOTCERT=/global-bundle.pem

EXPOSE 8080

ENTRYPOINT ["/receptor"]
