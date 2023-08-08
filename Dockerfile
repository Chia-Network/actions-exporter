FROM golang:alpine as builder

COPY . /app
WORKDIR /app
RUN apk add build-base
RUN make build

FROM alpine:3.18.3

COPY --from=builder /app/bin/actions-exporter /actions-exporter
CMD ["/actions-exporter"]
