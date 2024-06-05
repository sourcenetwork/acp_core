.PHONY: wasm
wasm:
	GOOS=js GOARCH=wasm go build -o build/acp.wasm cmd/wasm/main.go


.PHONY: test
test:
	go test -coverpkg=./...  ./...


.PHONY: proto
proto:
	docker run --rm --volume .:/app acp_core_proto
