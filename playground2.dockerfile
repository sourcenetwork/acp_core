# syntax=docker/dockerfile:1
#
# The playground service (Playground2: Generate, Check, Explain, Verify, Compile) with Alloy,
# and the acp CLI with the playground UI embedded.
#
#   docker run -p 8090:8090 <image>                                          # the service, for a local acp
#   docker run -p 7070:7070 -v $PWD/myapp:/workspace <image> playground-local # service + playground UI in one
#
# Build stages run on the build platform and cross-compile, so multi-platform builds need no emulation.

ARG ALLOY_VERSION=6.2.0
ARG ALLOY_SHA256=6b8c1cb5bc93bedfc7c61435c4e1ab6e688a242dc702a394628d9a9801edb78d

# Playground UI
FROM --platform=$BUILDPLATFORM node:20-alpine AS ui
WORKDIR /ui
COPY playground2/package.json playground2/package-lock.json ./
RUN npm ci
COPY playground2/ ./
RUN npm run build

# Go binaries
FROM --platform=$BUILDPLATFORM golang:1.24-bookworm AS go
ARG TARGETOS TARGETARCH BUILD_COMMIT
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# the committed cmd/acp/ui may lag behind playground2/: embed a fresh build
COPY --from=ui /ui/dist cmd/acp/ui
ENV CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH
RUN LDFLAGS="-X 'github.com/sourcenetwork/acp_core/pkg/version.Commit=$BUILD_COMMIT'" && \
    go build -ldflags "$LDFLAGS" -o /out/playground2 ./cmd/playground2 && \
    go build -ldflags "$LDFLAGS" -o /out/acp ./cmd/acp

# Alloy: the platform independent jar, checked against a known checksum
FROM --platform=$BUILDPLATFORM debian:bookworm-slim AS alloy
ARG ALLOY_VERSION ALLOY_SHA256
ADD --checksum=sha256:${ALLOY_SHA256} \
    https://github.com/AlloyTools/org.alloytools.alloy/releases/download/v${ALLOY_VERSION}/org.alloytools.alloy.dist.jar \
    /opt/alloy/alloy.jar

FROM eclipse-temurin:21-jre
# --chmod applies to the directories COPY creates too, which must stay traversable
COPY --from=alloy --chmod=755 /opt/alloy /opt/alloy
COPY --chmod=755 docker/alloy docker/playground-local /usr/local/bin/
COPY --from=go /out/playground2 /out/acp /usr/local/bin/
# alloy only needs a writable home and temp dir
ENV HOME=/tmp
USER 65532:65532
EXPOSE 8090 7070
CMD ["playground2", "-addr", "0.0.0.0:8090"]
