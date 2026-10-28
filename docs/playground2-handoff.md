# Handoff: Playground 2 (minimal UI + backend)

You're building the first version of a new acp_core playground, which will eventually replace the old one in `ui/`.
This document is everything you need to start. Read it fully before writing code.

## Why a new playground

acp_core is an authorization engine (Zanzibar-style, used by DefraDB). Users write a **policy**: resources, relations,
permissions. We've built tooling that compiles a policy into an [Alloy](https://alloytools.org) model, so we can
*generate example states of a system*: concrete sets of relationships that a policy (plus some rules the application
enforces) allows. Seeing generated examples, and checking what each person can do in them, is how users will
understand and refine their policies.

The old playground runs acp_core as WASM in the browser. Generating worlds needs the Alloy analyzer (a Java program),
so the new playground talks to a Go backend instead.

## Vocabulary

Use these terms in code, UI labels and names. The old playground's names (sandbox, theorem) are deliberately not used.

| Term | Meaning |
|---|---|
| **Policy** | The authorization rules acp_core enforces, written in YAML. |
| **System constraints** | Rules about relationship data that the *application* enforces and acp_core doesn't: `disjoint`, `at_most_one`, `acyclic`. |
| **System model** | Policy + system constraints. Defines every valid state of the system. |
| **Relationships** | Relationship tuples, one per line: `doc:roadmap#viewer@did:user:bob`. |
| **World** | One concrete state of a system: relationships (+ assertions, later). |
| **Generate** | Ask the backend for worlds that are valid states of a system model. |

Later (not in this task): **assertions** (expectations about one world), **contracts** (expectations about every world),
**Evaluate** (check a world), **workspace** (the playground's set of system model + worlds).

## Scope of this task

Build the minimal end-to-end loop:

1. **System model tab** with two editors:
   - **Policy**: a code editor where the user pastes a policy YAML (like the old playground's policy pane).
   - **Constraints**: a code editor for the system constraints (format below).
2. **Generate**: a control that calls the backend's `Generate` with the current system model and a small config
   (count, scope, actors, with defaults filled in), and shows the returned worlds.
3. **Worlds view**: shows each returned world's relationships. Readable is enough for now: a read-only editor or a
   list per world. Generated worlds are ephemeral: a new Generate replaces them.
4. **Errors**: show backend errors (invalid policy, invalid constraint, ...) clearly. The message text is meaningful;
   display it.
5. **A minimal backend** serving `Generate` over Connect (below).

Nice to have, only if trivial: persist the policy and constraints text in `localStorage`; preload a sample system model
(use the one in "Testing it" below).

Not in scope: docker, build scripts beyond the UI's own dev server, WASM, share links, explain/check graphs,
evaluating worlds, assertions, contracts, auth, multiple workspaces.

## Where to work

- UI: a new directory, `playground2/` at the repo root. Don't modify `ui/`.
- Backend: `cmd/playground2/` (a small `main.go`, plus whatever HTTP handler code you need, kept in that directory or
  a small package).

## What already exists (backend)

All on `main`, in Go:

- **API**: `proto/sourcenetwork/acp_core/playground2.proto` (service `Playground2`, rpc `Generate`) and
  `proto/sourcenetwork/acp_core/system_constraint.proto`. Generated Go types are in `pkg/types`
  (`types.GenerateRequest`, `types.SystemModel`, `types.World`, ...).
- **Implementation**: `services.NewPlayground2Service(alloyPath)` in `pkg/services/playground2.go` implements
  `types.Playground2Server`. Pass `""` to find `alloy` in `PATH`. It's stateless and deterministic, and each call runs
  Alloy (bounded by a timeout and by caps on count, scope and actors).
- **Requirements**: the `alloy` executable must be in `PATH` (it is on this machine: `/usr/local/bin/alloy`), plus Java.
- **Tests**: `go test ./internal/playground2/ ./pkg/services/` show how it's used.

You don't need to touch the compiler (`internal/policy_compiler`, `internal/alloy`).

## Backend: serve Generate over Connect

The UI talks to the backend with the [Connect protocol](https://connectrpc.com/docs/protocol). For a unary RPC,
Connect's JSON form is just an HTTP POST:

```
POST /sourcenetwork.acp_core.Playground2/Generate
Content-Type: application/json

<request as JSON>
```

The response is the response message as JSON with status 200. Errors are a non-200 status with
`{"code": "invalid_argument", "message": "..."}`.

The Go types are generated with gogoproto, not the standard Go protobuf API, so `protoc-gen-connect-go` output doesn't
work with them directly. Two options:

- **Simplest (recommended for now)**: a small `net/http` handler for that one path. Decode the body with gogo's
  `github.com/cosmos/gogoproto/jsonpb` into a `types.GenerateRequest`, call the service, and encode the response with
  `jsonpb.Marshaler{}`. Map errors to Connect error JSON: use `errors.Is(err, errors.ErrorType_BAD_INPUT)` (package
  `github.com/sourcenetwork/acp_core/pkg/errors`) for `invalid_argument` with HTTP 400, and `internal` with 500 for
  everything else.
- **Full Connect** (connect-go, also serves gRPC and gRPC-Web): the `feat/acpd` branch already serves acp_core services
  over Connect with a hand-written codec adapter for the gogo types. Read it with `git show feat/acpd:pkg/server/codec.go`
  and `git show feat/acpd:pkg/server/handlers.go`. Don't merge or depend on that branch; use it as reference only if
  you go this way.

Either way, the wire format is the same, so the UI doesn't care. Listen on `localhost:8090` by default, with a flag to
change it.

### Wire format (verified)

Request:

```json
{
  "systemModel": {
    "policy": "name: notes\nresources:\n- name: note\n  relations:\n  - name: reader\n    types: [actor]\n",
    "constraints": [
      { "disjoint":  { "resource": "note", "relations": ["owner", "reader"] } },
      { "atMostOne": { "resource": "note", "relation": "reader" } },
      { "acyclic":   { "resource": "folder", "relation": "parent" } }
    ]
  },
  "config": { "count": 2, "scope": 3, "actors": 3 }
}
```

`config` and every field in it is optional (defaults: count 1, scope 3, actors 3; maximums 50, 6 and 10).

Response:

```json
{
  "worlds": [
    { "relationships": "note:note0#owner@did:user:carol\nnote:note1#owner@did:user:carol\nnote:note0#reader@did:user:bob\nnote:note1#reader@did:user:alice\n" }
  ]
}
```

Empty fields are omitted (eg. `assertions`, and `worlds` when none were found). Fewer worlds than `count` means the
system model allows no more. `relationships` lists owner relationships first, one per line.

## UI: constraints format

The constraints editor holds YAML in this format. The UI parses it and maps it to the `constraints` JSON above:

```yaml
constraints:
- disjoint: [doc.owner, doc.viewer]   # a subject holds at most one of these relations on a doc
- at_most_one: doc.folder             # a doc has at most one folder
- acyclic: folder.parent              # following parent from a folder never leads back to it
```

- Each entry has exactly one key.
- `resource.relation` references split on the first `.`.
- `disjoint` references must all name the same resource → `{"disjoint": {"resource", "relations": [...]}}`.
- `at_most_one` and `acyclic` → `{"atMostOne": {"resource", "relation"}}` and `{"acyclic": {...}}`.

Report YAML/format errors in the UI. The backend validates the references against the policy (unknown relations,
relations that can't be disjoint, ...) and returns readable errors; show those too.

(The Go CLI's equivalent file uses the key `assumptions:`, the old name. The UI uses `constraints:`.)

## UI: stack and what to reuse from `ui/`

Use the same stack as the old playground, so the eventual replacement is natural:

- React 19 + TypeScript, Vite, Tailwind 4, shadcn/ui (Radix), lucide-react icons, zustand for state,
  `react-resizable-panels` for layout, Monaco (`@monaco-editor/react`) for editors (YAML language for policy and
  constraints).
- Worth looking at in `ui/src`:
  - `components/Editor/index.tsx` (`BaseEditor`): a Monaco wrapper bound to a store field, with debounced updates.
  - The layout structure in `components/RootLayout` and the pane/tab components, for look and feel.
  - `utils/mapLocatedMessageMarkers.ts`: maps backend error positions to Monaco markers (useful later).
  - `providers/ThemeProvider`, `utils/classnames.ts` (`cn()`), `components/ui/` (shadcn wrappers).
- Don't carry over: the WASM worker and loader (`workers/`, `playgroundWorkerProxy.ts`, `lib/wasm_exec.js`), sandbox
  handles and `SetState` flow, share/import/export, the BroadcastChannel sync, the explain-check graph.

Calling the backend: a plain `fetch` POST with hand-written TypeScript types for the few messages is fine for now.
Generated `@connectrpc/connect-web` clients can come later. Configure the Vite dev server to proxy
`/sourcenetwork.acp_core.` to the backend, so there's no CORS to deal with.

## Testing it

Backend, without the UI:

```
go run ./cmd/playground2
curl -s localhost:8090/sourcenetwork.acp_core.Playground2/Generate \
  -H 'Content-Type: application/json' -d @request.json
```

Sample system model to preload. Policy:

```yaml
name: folders
resources:
- name: folder
  relations:
  - name: parent
    types: [folder]
  - name: viewer
    types: [actor]
  permissions:
  - name: view
    expr: viewer + parent->view
- name: doc
  relations:
  - name: folder
    types: [folder]
  - name: viewer
    types: [actor]
  permissions:
  - name: view
    expr: viewer + folder->view
```

Constraints:

```yaml
constraints:
- acyclic: folder.parent
- at_most_one: doc.folder
- disjoint: [doc.owner, doc.viewer]
```

Generating 2 worlds with the defaults gives worlds with 3 actors (alice, bob, carol), folders nested without cycles,
each doc in at most one folder, and no doc owner also listed as its viewer.

More policies to try: `gents-acp/policy-compiler/assumption-examples/*/` in the user's workspace
(`/home/lodek/source/workspace/gents-acp/policy-compiler/`), which have policies and constraint files (with the old
`assumptions:` key).

Done when:
- [ ] `go run ./cmd/playground2` serves Generate; the curl above returns worlds; an invalid policy returns a Connect
      error with a readable message and a 4xx status.
- [ ] The UI loads with the sample, Generate shows worlds, and editing the policy or constraints and regenerating
      reflects the change.
- [ ] Errors (invalid YAML, invalid constraint, invalid policy) are visible in the UI.
- [ ] A short `playground2/README.md`: how to run backend and UI.

## Coming next (don't build, but don't block)

- Evaluate a world: what can each actor do in it (the engine already supports this via the `Simulator` service,
  which will move into `Playground2` under the new vocabulary).
- Save a generated world, edit worlds by hand, add assertions to them.
- A graph view of a world (entities and relationships), expandable from a node.
- Steering generation, eg. "worlds where docs are in nested folders".
- Contracts, checked across all worlds.

Keep the state model ready for this: a workspace holds one system model and a list of worlds, each either saved or
ephemeral.
