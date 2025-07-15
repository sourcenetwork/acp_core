FROM golang:1.23-bookworm AS playground-builder

WORKDIR /app

COPY go.* /app

RUN go mod download &&\
    mkdir -p /app/build &&\
	cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" /app/build/wasm_exec.js


COPY . /app

RUN GOOS=js GOARCH=wasm go build -o build/playground.wasm cmd/playground_js/main.go

FROM node:20-alpine AS ui-builder

RUN mkdir -p /home/node/app/node_modules && chown -R node:node /home/node/app

WORKDIR /home/node/app

USER node

COPY --chown=node:node ui/package-lock.json ui/package.json ./

RUN npm ci

COPY --chown=node:node --from=playground-builder /app/build/playground.wasm ./public/
COPY --chown=node:node ./ui/ ./
# Override wasm_exec.js with the correct version for the build
COPY --chown=node:node --from=playground-builder /app/build/wasm_exec.js ./src/lib/wasm_exec.js

ENV VITE_SHARE_URL="https://yvc9ygu76d.execute-api.us-east-1.amazonaws.com/dev"
RUN  npm run build


FROM nginx:alpine

COPY --from=ui-builder /home/node/app/dist /usr/share/nginx/html
