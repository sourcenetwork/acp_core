# Playground 2 UI

Minimal UI for generating example worlds from a system model (policy + system constraints) and verifying its requirements.
Stack: React 19, TypeScript, Vite, Tailwind 4, shadcn/ui, zustand, react-resizable-panels, Monaco.

## Run

Backend (needs `alloy` and Java in `PATH`), from the repo root:

    go run ./cmd/playground2            # listens on localhost:8090 (see -h for the flag to change it)

UI, from this directory:

    npm install
    npm run dev                          # http://localhost:5173

The dev server proxies `/sourcenetwork.acp_core.*` to `http://localhost:8090`
(override with `PLAYGROUND2_BACKEND=http://host:port npm run dev`).

Two modes, detected at startup (`GET /api/workspace`: 200 means file-backed, 404 or a network error means local):

- **Local mode** (`npm run dev` against the service only): the page is the workspace. Policy, constraints, requirements
  and generate config persist in `localStorage`; everything is editable.
- **File-backed mode** (inside `acp playground <dir>`, see `docs/playground2-acp-handoff.md`): the page is a view onto a
  workspace directory. To develop against it, from the repo root:

      make playground2 && build/playground2                      # the acp service on :8090 (needs alloy)
      make acp && build/acp playground -addr localhost:7070 -ui playground2/dist -open=false <dir>

  then `PLAYGROUND2_BACKEND=http://localhost:7070 npm run dev` (Vite proxies `/api` as well as the service API to
  `acp`, so you get hot reload and the workspace). A sample `<dir>`: `policy.yaml`, `constraints.yaml`,
  `contracts.yaml` and `worlds/*.yaml` (`acp verify -o <dir>/worlds <dir>` saves counterexamples there).
  `make acp-ui` embeds the production build into `acp`.

Other scripts: `npm run build`, `npm run lint`, `npm test` (vitest).

## Notes

- Local mode: the policy, constraints, requirements and generate config persist in `localStorage`
  (`playground2.workspace`); a sample system model is preloaded. Generated worlds are ephemeral and not persisted.
- File-backed mode: the files are the source of truth. Policy, constraints and requirements come from `/api/workspace`
  (the files as written; the UI parses them, so error lines match the files) and are never read from or written to
  `localStorage`, which still keeps the UI preferences (generate config, verify scope, auto-verify; `mergePrefs`
  leaves a local-mode workspace in storage untouched). The three editors are read-only with a notice
  ("Editing `<dir>/policy.yaml`? Changes appear here automatically."). Saved worlds are `worlds/*.yaml`
  (`relationships: |`, optional `assertions: |`), shown by file name; there is no save button: a saved world shows
  its file, a generated one has "Copy as YAML" and the `acp generate -o <dir>/worlds <dir>` equivalent. Generated
  worlds stay in memory and survive reloads. The header shows the directory and a connection indicator.
- Live reload (file-backed): an `EventSource` on `/api/workspace/events`; `changed` refetches, `ready` refetches only if
  the version differs (so changes made while `acp` was down are picked up). Files replace the store and everything
  parsed from them is re-derived; Verify results are kept and marked stale (`VerifyRun.inputs`). "Auto-verify on
  change" (Requirements tab, off by default, persisted) re-verifies on change, and when switched on with stale results.
  When the workspace can't be read (500: eg. `policy.yaml` deleted) the last good files stay, the error is shown and
  clears when the files are readable again; if it is unreadable from the start the editors are empty. Vite's proxy
  keeps a browser's event stream open after `acp` dies, so a 4s heartbeat fetch also drives the indicator
  (reconnecting) and reopens the stream when `acp` is back. A world file that doesn't parse is reported on its own
  (`worlds/foo.yaml line 2: ...`) and the other worlds are unaffected; `constraints.yaml` and `contracts.yaml` problems
  are listed the same way (`contracts.yaml line 12: ...`) above the tabs.
- Network errors: a 502 `{"code":"unavailable"}` from `acp`'s proxy is shown as is; otherwise the hint names
  `localhost:8090` (local mode) or `acp playground` (file-backed).
- Pure logic is in `src/lib/workspaceFiles.ts` (mode detection, world files, sync protocol, prefs merge; tested) and
  the store (`initWorkspace`, `applyFiles`).
- Constraints are YAML, parsed client-side (`src/lib/constraints.ts`) into the wire format:

      constraints:           # a workspace's constraints.yaml uses `assumptions:`: both keys are accepted
      - disjoint: [doc.owner, doc.viewer]
      - at_most_one: doc.folder
      - acyclic: folder.parent

- Errors: constraint parse errors (client side, with line numbers), backend errors (`{code, message}`)
  and network errors are shown above the tabs, whatever tab is active.
- Requirements tab: a YAML editor for the requirements file (`requirements:` with `name`, `statement`, `notes`, and
  `contracts` of `name`, `description`, `status: proposed|accepted`, `properties` of `type: assert|pred`, `description`,
  `formula`; formulas are usually block scalars; format in `HANDOFF.md`). The text persists in `localStorage` with the
  Policy and constraints and a sample is preloaded. `src/lib/requirements.ts` parses it into the wire `requirements`
  (enums uppercased, optional fields tolerated, unknown keys and missing required fields reported with line numbers).
- Verify (Requirements tab): a Scope input (default 3, not capped client-side: the backend's error is shown) and a
  Verify button (Cancel while running; a newer run or a cancel ignores the stale response). It sends the current Policy,
  constraints and requirements. Results sit beside the editor (resizable): every requirement with its statement and
  notes, its contracts (rollup badge, proposed/accepted) and their properties (SATISFIED green, VIOLATED red with a tinted
  row, UNVERIFIED amber with the backend message, DRAFT neutral; ASSERT / PRED chips, with a note that satisfied asserts
  hold only up to scope N), plus a summary of counts. Results are derived: not persisted, and marked stale once the
  Policy, constraints or requirements change after the run. A run shows the requirements it was computed from.
  Joining and counting is in `src/lib/verifyResults.ts` (omitted wire fields are read as their zero values).
- Counterexamples: a property with a `world` has "View counterexample" (violated assert) or "View example" (satisfied
  pred), which opens that world in a pane below the results using the same `WorldView` as the Worlds tab (Graph /
  Relationships toggle, Check and Explain). The bindings are shown as chips and the bound nodes carry a dashed ring and a
  badge with the variable names (only values equal to a node id, so `a->b` tuples are text only). Check/Explain is
  prefilled when the bindings name one object and one actor (the permission is the first one the formula mentions as
  `resource_name`, else the default); it stays editable (`src/lib/bindings.ts`). Saving a world is not built.
- Layout: the Generate bar (count / scope / actors) and any error stay on top; below are four full-height tabs:
  Policy, System constraints, Requirements and Worlds. Worlds shows one world at a time, chosen from a dropdown
  (`World 1 (ephemeral, 16 relationships)`); a successful Generate switches to the Worlds tab and selects World 1,
  a failed one leaves the active tab alone. Switching worlds clears the object/permission selection and explain panel.
  Active tab and selected world are transient store state (`activeTab`, `selectedWorld`).
- State is shaped as a workspace: one system model plus a list of worlds (`saved` | `ephemeral`), see `src/stores/workspaceStore.ts`.
- Each world has a Graph / Relationships toggle (default Graph). The graph (`src/components/RelationshipGraph.tsx`) is
  an SVG laid out with dagre (`@dagrejs/dagre`, pinned): objects are nodes coloured per resource type, actors are grey
  pills, each tuple is an edge object -> subject labelled with the relation; userset subjects (`team:t1#admin`) point at
  the object node with a dashed edge labelled `relation (#admin)`. Wheel to zoom, drag to pan, double-click or
  "Reset view" to fit. Parsing is in `src/lib/relationships.ts` (unparseable lines are reported, not fatal), layout in
  `src/lib/graph.ts`; both are unit-tested.
- Check and Explain on a world's graph (`src/components/WorldGraphPanel.tsx`): click an object node, then pick a permission
  or relation (the input suggests the names from the Policy, parsed client-side in `src/lib/policyResources.ts`, plus the
  implicit `owner`; any other text is sent as typed and the backend error is shown). One `Check` per actor in the world
  runs (4 at a time, aborted when the selection or Policy changes): allowed actors turn green, denied ones are dimmed,
  per-actor errors are marked `!` and listed under the graph, with a summary line (`view on doc:doc0: alice, carol allowed`).
  Click the background or "Clear" to reset. Clicking an actor calls `Explain` and draws the engine's search graph below
  (`src/components/ExplainGraphView.tsx`; layout from the edges with dagre in `src/lib/explainGraph.ts`, `rootNodeId` is
  ignored): unions/intersections/differences are small circles, usersets are boxes, coloured authorized / searched /
  explored / unexplored, edge messages as labels. The scenario is the current Policy and constraints plus that world.
  Runner and summaries live in `src/lib/checkRunner.ts`; pan/zoom is shared (`src/lib/usePanZoom.ts`).
