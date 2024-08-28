# Playground JS

The Playground JS executable compiles a binary containing a WASM module which exports a [PlaygroundService](../../proto/sourcenetwork/acp_core/playground.proto) implementation.

The WASM binary modifies the [global JS scope](https://pkg.go.dev/syscall/js#Global) and exports an `AcpPlayground` object which acts as a namespace with a constructor for the Playground Service.

## Building

Run `make playground:wasm_js` to build the playground the Playground as a WASM module.

## Example

With the playground loaded and the go `wasm_exec.js` runtime imported (see https://go.dev/wiki/WebAssembly), a new playground instance can be created with:

```js
let playground = await window.AcpPlayground.new();

let req = {
    data: {
        policy_definition: `
name: test
resources:
  file:
    relations:
      owner:
        types:
          - actor
      reader:
        types:
          - actor
    permissions:
      read:
        expr: owner + reader
      write:
       expr: owner
`,
        relationships: `
file:readme#owner@did:example:bob
file:readme#reader@did:example:alice
        `,
        policy_theorem: `
				Authorizations {
				  file:readme#read@did:example:bob
				  file:readme#write@did:example:bob

				  !file:readme#write@did:example:alice
				  file:readme#read@did:example:alice
				}
				Delegations {}
				ImpliedRelations {}
        `,
    }
}

let response = await playground.simulate(req)
```