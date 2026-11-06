# Handoff: file-backed workspaces in Playground 2 (`acp playground`)

This follows `docs/playground2-handoff.md`. Read that first if you haven't; vocabulary is in `HANDOFF.md`
("Vocabulary") and in the glossary.

## Why

The playground is moving from "a page you edit" to "a view onto a workspace on disk". An agent (or the user) edits a
workspace directory, and the human watches the playground update live:

```
myapp/
  policy.yaml        required
  constraints.yaml   optional (the `assumptions:` YAML the constraints editor already parses)
  contracts.yaml     optional (the `requirements:` YAML the Requirements tab already parses)
  worlds/*.yaml      optional saved worlds, one per file: `relationships: |` and optionally `assertions: |`
```

Users won't have Alloy, so Generate/Verify run on a remote acp service. The local `acp` binary bridges the two:

```
browser ──► acp playground <dir>   (local, eg. http://127.0.0.1:7070)
              ├─ /                                     the UI (this app's dist/)
              ├─ /sourcenetwork.acp_core.Playground2/* proxied to the acp service (same API as today)
              ├─ GET /api/workspace                    the workspace's files
              └─ GET /api/workspace/events             server-sent events on file changes
```

The UI keeps calling the API with relative paths, exactly as today: in the playground, `acp` proxies them; in
`npm run dev`, Vite does. Nothing about the API calls changes. The backend side is done (`cmd/acp/playground.go`).

## The local endpoints

### `GET /api/workspace`

The files **as written**, not parsed: the UI keeps parsing them itself, so error line numbers match the files.

```json
{
  "dir": "/home/bruno/myapp",
  "policy": "name: folders\nresources:\n...",
  "constraints": "assumptions:\n- acyclic: folder.parent\n",
  "contracts": "requirements:\n- name: doc-edit\n...",
  "worlds": [
    { "name": "child-owner-sees-parent-1", "content": "relationships: |\n  doc:doc0#owner@did:user:bob\n  ..." }
  ],
  "version": "0e4dc2e7d1fcb2e0"
}
```

- Missing optional files are `""`; no saved worlds is `[]`. `version` changes whenever any file does.
- If the workspace can't be read (eg. `policy.yaml` deleted, or mid-save), the response is HTTP 500 with the usual
  `{"code": "internal", "message": "..."}`. Show it and keep the last good state; the next change event will retry.
- A world's `content` is YAML (`relationships`, `assertions`), parse it with the `yaml` package already in use. A
  world that doesn't parse should show its own error, not break the workspace.

### `GET /api/workspace/events`

Server-sent events (use `EventSource`):

```
event: ready
data: {"version":"0e4dc2e7d1fcb2e0"}

event: changed
data: {"version":"8e0b16d90dc2326a"}
```

- `ready` on connect, `changed` after any file changes (polled every 500ms, so edits arrive within ~half a second).
  Comment lines (`: keepalive`) arrive every 25s.
- On `changed`, refetch `/api/workspace`. On `ready` after a reconnect, refetch too if the version differs from what
  you have. `EventSource` reconnects by itself if `acp` restarts.

## What changes in the UI

### 1. Two modes, detected at startup

- `GET /api/workspace` answers 200 → **file-backed mode** (inside `acp playground`).
- 404 (eg. plain `npm run dev` without `acp`), or a network error → **local mode**: today's behavior, unchanged
  (`localStorage`, sample workspace, editable).

Keep local mode working: it's how the UI is developed and demoed without a workspace.

### 2. File-backed mode

- **Files are the source of truth.** Policy, constraints and requirements come from `/api/workspace`, not from
  `localStorage` (`workspaceStore.ts` persists them today; in file-backed mode don't read or write them there, but
  keep persisting UI preferences such as the verify scope and generate config).
- **Editors are read-only**, with a short notice: "Editing `<dir>/policy.yaml`? Changes appear here automatically."
  Writing files from the UI is undecided; don't build it.
- **Saved worlds come from `worlds/*.yaml`**: `kind: "saved"`, named after the file (add a `name` to the `World`
  type; today saved worlds have only an `id`). Generated worlds stay ephemeral and in memory, as now. There's no
  "save" button in this mode; it's fine to show the CLI equivalent (`acp generate -o <dir>/worlds <dir>`) or a
  "copy as YAML" action.
- **Live reload**: on a change, replace the files in the store and re-derive everything parsed from them. The
  existing `VerifyRun.inputs` staleness check already tells you when results are out of date: show results as stale
  rather than clearing them. Add an **"auto-verify on change"** toggle, off by default (Verify runs Alloy once per
  property, a few seconds in total).
- **Header**: show the workspace `dir` and a connection indicator (connected / reconnecting) for the event stream.
- **Errors**: per-file parse errors should name the file (`contracts.yaml line 12: ...`, `worlds/foo.yaml: ...`).

### 3. Small fixes

- `api.ts`'s network error says "Is it running on localhost:8090?". In file-backed mode the backend is `acp`; a 502
  from the proxy already carries a useful message (`{"code":"unavailable","message":"can't reach the acp service at
  ..."}`): show that instead.
- `vite.config.ts`: also proxy `/api` so you can develop against a real `acp playground` (below).

## Developing against `acp playground`

```
make playground2 && build/playground2                       # the acp service on :8090 (needs alloy)
make acp && build/acp playground -addr localhost:7070 -ui playground2/dist -open=false <workspace dir>
PLAYGROUND2_BACKEND=http://localhost:7070 npm run dev       # once vite proxies /api as well
```

`acp playground` proxies the service, so pointing Vite at it gives you both the API and `/api/workspace` with hot
reload. A sample workspace: `policy.yaml` and `assumptions.yaml` (rename to `constraints.yaml`) from
`/home/lodek/source/workspace/gents-acp/policy-compiler/assumption-examples/01-folders`, plus a `contracts.yaml` in
the format the Requirements tab uses. `acp verify -o <dir>/worlds <dir>` saves counterexamples as worlds.

Production build: `make acp-ui` runs `npm run build` and copies `dist/` into `cmd/acp/ui/`, which is committed and
embedded in every `acp` build (so `go install` includes the UI). Run it after UI changes. The UI must keep working
when served from `/` with the default Vite `base` (it does today).

## Done when

- [ ] Inside `acp playground`, the UI shows the workspace's files, read-only, with the directory in the header.
- [ ] Editing `policy.yaml`, `constraints.yaml` or `contracts.yaml` on disk updates the UI within a second, and marks
      Verify results stale; with auto-verify on, it re-verifies.
- [ ] Saved worlds from `worlds/*.yaml` appear by name; adding or deleting a file adds or removes the world; a broken
      world file shows its own error.
- [ ] `npm run dev` without `acp` still works as before (local mode).
- [ ] A broken workspace (deleted `policy.yaml`) shows an error and recovers when the file comes back.
- [ ] README updated with both modes.

## Not now

- Writing files from the UI (saving worlds, editing policy): undecided.
- Auth: the remote service is ungated for now; when it isn't, `acp` will add credentials to proxied calls, so the UI
  won't change.
- A model view: `POST /sourcenetwork.acp_core.Playground2/Compile` returns the compiled Alloy modules if you want a
  read-only "model" tab next to Requirements, but it's optional.
