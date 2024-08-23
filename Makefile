.PHONY: test
test:
	go test -coverpkg=./...  ./...


.PHONY: proto
proto:
	docker run --rm --volume .:/app acp_core_proto

.PHONY: playground\:wasm
playground\:wasm:
	mkdir -p build/playground
	GOOS=js GOARCH=wasm go build -o build/playground/playground.wasm cmd/playground_wasm/main.go
	cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" build/playground/
	cp static/playground-index.html build/playground/index.html
	cd proto && buf generate --template buf.ts.gen.yaml