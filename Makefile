GIT_HEAD_COMMIT=$(shell git rev-parse HEAD)
BUILD_FLAGS=-ldflags "-X 'github.com/sourcenetwork/acp_core/version.Commit=$(GIT_HEAD_COMMIT)'"

.PHONY: test
test:
	go test -coverpkg=./... ./...

.PHONY: test\:bench
test\:bench:
	go test ./... -bench .
	# To profile use -cpuprofile cpu.out

.PHONY: test\:js
test\:js:
	GOOS=js GOARCH=wasm go test -exec="$$(go env GOROOT)/misc/wasm/go_js_wasm_exec" ./...

.PHONY: proto
proto:
	docker image build --file proto/Dockerfile --tag acp_core_proto:latest .
	docker run --rm --volume=".:/app" --workdir="/app/proto" --env="HOST_USER=$$(id -u)" --entrypoint sh acp_core_proto ./docker-generate.sh
	mv github.com/sourcenetwork/acp_core/pkg/types/* pkg/types/
	mv github.com/sourcenetwork/acp_core/pkg/errors/* pkg/errors/
	rm -r github.com

.PHONY: fmt
fmt:
	docker run --rm --volume=".:/app" --workdir="/app/proto" --user="$$(id -u)" acp_core_proto format -w
	gofmt -w .

.PHONY: playground\:wasm_js
playground\:wasm_js:
	GOOS=js GOARCH=wasm go build $(BUILD_FLAGS) -o build/playground.wasm cmd/playground_js/main.go

.PHONY: playground
playground: playground\:wasm_js
	cp build/playground.wasm cmd/playground/content/playground.wasm
	cp "$$(go env GOROOT)/misc/wasm/wasm_exec.js" cmd/playground/content/wasm_exec.js
	go build -o build/playground cmd/playground/main.go
