# acp playground quickstart

Two pieces:

- **`acp`**, a local CLI. It reads a workspace directory, serves the playground UI for it, and runs nothing heavy
  itself: every analysis is a call to the service.
- **The service** (with the Alloy analyzer), as a Docker image, so you don't need Java or Alloy.

## Install

```
# the CLI (Go 1.24+); <commit> is a commit of the feat/policy-builder branch
go install github.com/sourcenetwork/acp_core/cmd/acp@<commit>

# the service, on localhost:8090
docker run -d --rm --name acp-service -p 8090:8090 ghcr.io/sourcenetwork/acp_playground2:feat-policy-builder
```

If the image pull is denied, the package is private: `docker login ghcr.io` with a GitHub token that has
`read:packages`. `acp` talks to `http://localhost:8090` by default; set `ACP_SERVER` (or `-server`) otherwise.

## A workspace

A workspace is a directory:

```
myapp/
  policy.yaml        the acp policy (required)
  constraints.yaml   rules the application enforces on its data (optional)
  contracts.yaml     requirements, in the user's words, and the contracts derived from them (optional)
  worlds/*.yaml      saved example states (optional)
```

## Use it

```
acp playground myapp                   # opens the playground; it follows the files as they change
acp verify myapp                       # check the contracts
acp verify -o myapp/worlds myapp       # ... and save counterexamples as worlds
acp generate -n 3 myapp                # example worlds
acp model myapp                        # the Alloy model contracts are written against
acp explain myapp <world> doc:readme view did:user:bob
```

`acp` with no arguments lists every command.

## Without Go

The image also contains `acp`. Run the service and the playground together, and the CLI with `docker exec`:

```
docker run --rm --name acp --user $(id -u):$(id -g) -p 7070:7070 -v $PWD/myapp:/workspace \
  ghcr.io/sourcenetwork/acp_playground2:feat-policy-builder playground-local
# open http://localhost:7070
docker exec acp acp verify /workspace
```
