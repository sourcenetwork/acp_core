FROM golang:1-alpine

ARG GOCOSMOS_VERSION=latest
ARG BUF_VERSION=1.27.2

RUN apk update &&\
    apk add --no-cache git curl &&\
    go install github.com/cosmos/gogoproto/protoc-gen-gocosmos@${GOCOSMOS_VERSION} &&\
    curl -sSL "https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-$(uname -s)-$(uname -m)" -o "/usr/local/bin/buf" &&\
    chmod +x "/usr/local/bin/buf"

COPY . /app

WORKDIR /app/proto

CMD ["/usr/local/bin/buf", "generate"]