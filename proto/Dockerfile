FROM golang:1.22-alpine

ARG COSMOS_GOGOPROTO_VERSION=latest
ARG PULSAR_VERSION=latest
ARG BUF_VERSION=1.27.2
ARG GRPC_GATEWAY_VERSION=latest

RUN apk update &&\
    apk add --no-cache git curl &&\
    go install github.com/cosmos/gogoproto/protoc-gen-gocosmos@${COSMOS_GOGOPROTO_VERSION} &&\
    go install github.com/cosmos/cosmos-proto/cmd/protoc-gen-go-pulsar@${PULSAR_VERSION} &&\
    go install github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@${GRPC_GATEWAY_VERSION} &&\
    curl -sSL "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-$(uname -s)-$(uname -m)" -o "/usr/local/bin/buf" &&\
    chmod +x "/usr/local/bin/buf" &&\
    adduser -h /home/proto -D proto &&\
    chown proto:proto /home/proto

USER proto
ENV HOME="/home/proto"
WORKDIR /app

ENTRYPOINT ["/usr/local/bin/buf"]