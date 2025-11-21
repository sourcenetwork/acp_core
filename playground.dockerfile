# Playground WASM Stage
FROM golang:1.23-bookworm AS playground-builder

# BUILD_COMMIT arg represents the acp_core commit from which the image was built
ARG BUILD_COMMIT

WORKDIR /app

COPY go.* /app

RUN go mod download &&\
    mkdir -p /app/build &&\
	cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" /app/build/wasm_exec.js


COPY . /app

RUN GOOS=js GOARCH=wasm go build -ldflags "-X 'github.com/sourcenetwork/acp_core/pkg/version.Commit=$BUILD_COMMIT'" -o build/playground.wasm cmd/playground_js/main.go


# UI Builder Stage
FROM node:20-alpine AS ui-builder

RUN mkdir -p /home/node/app/node_modules && chown -R node:node /home/node/app

WORKDIR /home/node/app

USER node

COPY --chown=node:node ui/package-lock.json ui/package.json ./

RUN npm ci

COPY --chown=node:node ./ui/ ./
COPY --chown=node:node --from=playground-builder /app/build/playground.wasm ./public/
COPY --chown=node:node --from=playground-builder /app/build/wasm_exec.js ./src/lib/wasm_exec.js

# Override wasm_exec.js with the correct version for the build
ENV VITE_SHARE_URL="/api/sandboxes"

RUN  npm run build

# Main entrypoint
FROM nginx:alpine
COPY --from=ui-builder /home/node/app/dist /usr/share/nginx/html
EXPOSE 80