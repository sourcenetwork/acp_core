# Playground

Playground command builds a static HTTP server containing ACP Playground

## Building

Run `make playground` in the project root to build the embedded binary `build/playground`.

## Example

To build and execute:

```sh
make playground
build/playground
```

Once the service is running, it can be access though `http://localhost:8080`.

The index.html file should load the WASM module, which can then be interacted with through the dev tools console.