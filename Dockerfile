FROM golang:1.22 AS build

ENV CGO_ENABLED=0
ENV GOOS=linux
RUN useradd -u 10001 connect

WORKDIR /go/src/github.com/redpanda-data/redpanda-connect-plugin-example/
# Update dependencies: On unchanged dependencies, cached layer will be reused
COPY go.* /go/src/github.com/redpanda-data/redpanda-connect-plugin-example/
RUN go mod download

# Build
COPY . /go/src/github.com/redpanda-data/redpanda-connect-plugin-example/

# Tag timetzdata required for busybox base image:
# https://github.com/redpanda-data/connect/issues/897
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod go build -tags timetzdata -ldflags="-w -s" -o redpanda-connect-plugin-example

# Pack
FROM busybox AS package

LABEL org.opencontainers.image.source="https://github.com/insidegreen/redpanda-connect-plugin-example"

WORKDIR /

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /go/src/github.com/redpanda-data/redpanda-connect-plugin-example/redpanda-connect-plugin-example .
COPY ./config/example_1.yaml /connect.yaml

USER connect

EXPOSE 4195

ENTRYPOINT ["/redpanda-connect-plugin-example"]

CMD ["-c", "/connect.yaml"]
