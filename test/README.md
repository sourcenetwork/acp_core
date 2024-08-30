# Test

## Testing JS

In order to test the JS implementation of acp_core and the playground, the test pkg defines a wrapper over the JS wrapper.

The JS wrapper defined in the `internal/playground_js` pkg wraps a PlaygroundService such that it receives as input JS objects, transforms them into go objects and maps the output as JS objects again.

Testing in the JS environment is done by adding yet another layer, which takes Go inputs and transforms it into JS inputs.
This enables reusing the entire test suite as in, without requiring additional setup.

In the future, an interesting approach could be to employ differential testing as well, where the Go codebase can be used as a test oracle.
Any difference in output between the Wasm/JS version and the Go version indicates an error.