.PHONY: test
test:
	go test -coverpkg=./...  ./...

.PHONY: test\:js
test\:js:
	GOOS=js GOARCH=wasm go test -exec="$$(go env GOROOT)/misc/wasm/go_js_wasm_exec" ./...

.PHONY: proto
proto:
	docker run --rm --volume .:/app acp_core_proto

.PHONY: playground\:wasm_js
playground\:wasm_js:
	GOOS=js GOARCH=wasm go build -o build/playground.wasm cmd/playground_js/main.go

.PHONY: proto\:ts
proto\:ts:
	cd proto && buf generate --template buf.ts.gen.yaml

.PHONY: playground
playground: playground\:wasm_js
	cp build/playground.wasm cmd/playground/content/playground.wasm
	cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" cmd/playground/content/wasm_exec.js
	go build -o build/playground cmd/playground/main.go
