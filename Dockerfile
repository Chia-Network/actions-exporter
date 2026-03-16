FROM golang:1 AS builder

COPY . /app
WORKDIR /app
RUN make build

FROM gcr.io/distroless/static-debian13:latest

COPY --from=builder /app/bin/actions-exporter /actions-exporter
USER nonroot:nonroot
CMD ["/actions-exporter"]
