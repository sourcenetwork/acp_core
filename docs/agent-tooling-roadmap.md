# Agent Tooling Roadmap

Working plan for making acp_core usable as a standalone service and by agents.

## Decisions

- Standalone deployment is independent of SourceHub.
- Transport: Connect (HTTP/1.1 + JSON friendly) over the existing ACPEngine / PlaygroundService protos.
- Auth (first pass): self-signed JWS from the caller's DID. The gRPC/Connect layer carries the **user** DID as the Principal.
- Concurrency (first pass): global lock serializing requests (zanzi/raccoon have no tx support).
- Two MCP servers:
  - **Ops MCP** — thin client over the deployed service. Dual identity always present: requesting user DID + agent DID. Every action logged with both. Agent identity recorded as an attribute on policy / relationship records for auditability.
  - **Policy dev MCP** — later. Embeds acp_core locally, stateless simulate/verify, plus the verification tooling.
- Verification tooling (policy -> Alloy compiler) gets its own RPC service. Orthogonal to this track.

## Phase 1 (in progress)

- [x] Standalone server `cmd/acpd` — branch `feat/acpd-server` (634d048), verified
- [x] Integrate: `feat/acpd` (cb187ec) — jws auth mode, principal+actor request logging, e2e over 4 protocols; binary smoke-tested over HTTP/JSON
- [x] JWS DID authenticator — branch `feat/jws-did-auth` (3b34267), verified; awaiting merge
- [x] Benchmarks + profiling report; metrics audit — branch `perf/benchmarks` (1b92515)

## Perf (ranked by measured impact)

- [ ] zanzi: parse permission expressions once per policy (cache); stop `fmt.Errorf` on parser backtracking (`relation_expression_parser/primitives.go` buildErr) — ~50% CPU / 55% allocs on Verify, ~45% on CreatePolicy
- [ ] zanzi: build policy lookup table once per search, plain string keys (not prototext `RelationFQN.String()`) — ~8% CPU / 16% allocs
- [ ] Fetch policy once per request (currently 2–4 fetches + double unmarshal) — ~25–30% on SetRelationship / Simulate
- [ ] zanzi: direct edge-existence lookup before successor scan — point checks O(N) → O(1) (10.9ms → ~0.3ms at 10k readers); index exact-match FilterRelationships
- [ ] raccoondb: remove stray `fmt.Printf("")` in `key.go` — 2–7%
- [ ] DefaultEventManager: unbounded, never drained, unlocked — memory leak for long-running acpd
- [ ] Metrics: real `MetricService` (counters/histograms, noop default, `WithMetricService`), timing/error decorator in service chain, KV wrapper counts

## Phase 2

- [ ] Ops MCP server (dual identity, read-only mode, guarded destructive tools, curated tool surface)
- [ ] Agent identity as record attribute on policies / relationships
- [ ] Dockerfile / deployment for `acpd`

## Backlog

- [ ] Auth: verified agent identity — user-signed delegation to agent DID; agent signs requests (today `act` is an unverified user claim)
- [ ] Auth: replay protection (jti tracking) and optional body binding
- [ ] Auth: Root principal should retain the DID that authenticated as root (audit)
- [ ] Auth: DID methods beyond did:key (rotation / revocation)
- [ ] **raccoondb iterator leak**: ObjectStore never closes iterators; MemKV deadlocks after ~64 open iterators (e.g. repeated ListPolicies). LevelDB leaks. Fix upstream in raccoondb.
- [ ] Sandbox handle allocation races under concurrent NewSandbox (affects embedded users; acpd lock covers it)
- [ ] Server: redact internal error messages returned to clients
- [ ] Server: revisit error code mapping (OPERATION_FORBIDDEN→FailedPrecondition, UNAUTHORIZED→PermissionDenied)
- [ ] Connect codegen: protoc-gen-connect-go output references APIv2 file descriptors; hand-written adapter used instead
- [ ] Pre-existing: `GOOS=js go build ./pkg/...` fails (os/user) on main

- [ ] Proper transaction isolation in the KV layer (replace global lock)
- [ ] Memoize policy loading (policies currently re-read/unmarshaled repeatedly)
- [ ] Reduce proto marshaling overhead (pending benchmark results)
- [ ] Metrics (pending audit)
- [ ] Revive playground deployment (old deployment is gone)
- [ ] Remove stale `build/playground` binary
- [ ] Policy dev MCP
- [ ] Alloy compiler RPC service
