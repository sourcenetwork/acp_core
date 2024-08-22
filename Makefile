.PHONY: test
test:
	go test -coverpkg=./...  ./...


.PHONY: proto
proto:
	docker run --rm --volume .:/app acp_core_proto

.PHONY: playground\:wasm
playground\:wasm:
	GOOS=js GOARCH=wasm go build -o build/playground.wasm cmd/playground_wasm/main.go
