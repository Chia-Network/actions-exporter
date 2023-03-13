FROM golang:alpine as builder

COPY . /app
WORKDIR /app
RUN apk add build-base
RUN make build

FROM alpine:latest

COPY --from=builder /app/bin/actions-exporter /actions-exporter
CMD ["/actions-exporter"]
