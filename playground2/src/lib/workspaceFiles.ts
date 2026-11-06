// Pure logic of file-backed mode (inside `acp playground <dir>`): mode detection, the /api/workspace payload and how it
// maps onto the store, saved-world files, the sync protocol with /api/workspace/events. See docs/playground2-acp-handoff.md.
import { parse, stringify } from "yaml";
import { parseConstraints } from "./constraints";
import { parseRequirements } from "./requirements";
import type { World } from "./types";

export const WORKSPACE_PATH = "/api/workspace";
export const EVENTS_PATH = "/api/workspace/events";

export const POLICY_FILE = "policy.yaml";
export const CONSTRAINTS_FILE = "constraints.yaml";
export const CONTRACTS_FILE = "contracts.yaml";

/** The GET /api/workspace payload: the files as written. */
export interface WorkspaceFiles {
  dir: string;
  policy: string;
  constraints: string;
  contracts: string;
  worlds: { name: string; content: string }[];
  version: string;
}

export type Mode = "pending" | "local" | "files";

export type Detection =
  | { mode: "local" }
  | { mode: "files"; files: WorkspaceFiles }
  /** acp answered but the workspace can't be read right now (eg. policy.yaml deleted before the UI loaded). */
  | { mode: "files"; files: null; error: string };

/** One problem in one workspace file, named after the file. */
export interface FileProblem {
  file: string;
  line?: number;
  message: string;
}

export const formatProblem = (p: FileProblem) => `${p.file}${p.line !== undefined ? ` line ${p.line}` : ""}: ${p.message}`;

function isFiles(v: unknown): v is WorkspaceFiles {
  const f = v as Partial<WorkspaceFiles> | null;
  return !!f && typeof f.dir === "string" && typeof f.policy === "string" && typeof f.version === "string";
}

/** Normalises a /api/workspace body: optional files may be omitted, worlds may be null. */
export function normalizeFiles(v: WorkspaceFiles): WorkspaceFiles {
  return { ...v, constraints: v.constraints ?? "", contracts: v.contracts ?? "", worlds: v.worlds ?? [] };
}

/** Reads a /api/workspace response. 200 is the workspace; anything else is an error message. */
export async function readWorkspaceResponse(res: Response): Promise<{ files: WorkspaceFiles } | { error: string; json: boolean }> {
  const raw = await res.text();
  let body: unknown;
  let json = true;
  try {
    body = JSON.parse(raw);
  } catch {
    json = false;
  }
  if (res.ok) {
    if (isFiles(body)) return { files: normalizeFiles(body) };
    return { error: "unexpected response from /api/workspace", json };
  }
  const message = (body as { message?: unknown } | undefined)?.message;
  return { error: typeof message === "string" ? message : `HTTP ${res.status}: ${raw.slice(0, 200) || res.statusText}`, json };
}

/**
 * Local or file-backed? 200 means file-backed; 404 or a network error means local. A 500 with acp's JSON error body
 * (`{code, message}`) is a file-backed workspace that can't be read; a bare 500 is a dev proxy with nothing behind it.
 */
export async function detectMode(fetchFn: typeof fetch = fetch): Promise<Detection> {
  let res: Response;
  try {
    res = await fetchFn(WORKSPACE_PATH, { headers: { Accept: "application/json" } });
  } catch {
    return { mode: "local" };
  }
  const r = await readWorkspaceResponse(res);
  if ("files" in r) return { mode: "files", files: r.files };
  if (res.status === 500 && r.json) return { mode: "files", files: null, error: r.error };
  return { mode: "local" };
}

/**
 * Fetches the workspace again. Errors are messages so that callers keep their last good state; `reachable` says whether
 * acp answered (its own JSON error, eg. a deleted policy.yaml, is reachable; a network error or a bare 5xx from a dev
 * proxy with nothing behind it is not).
 */
export async function fetchWorkspace(
  fetchFn: typeof fetch = fetch,
): Promise<{ files: WorkspaceFiles } | { error: string; reachable: boolean }> {
  try {
    const r = await readWorkspaceResponse(await fetchFn(WORKSPACE_PATH, { headers: { Accept: "application/json" }, cache: "no-store" }));
    return "files" in r ? r : { error: r.error, reachable: r.json };
  } catch (e) {
    return { error: `can't reach acp: ${(e as Error).message}`, reachable: false };
  }
}

// ---- worlds/*.yaml

export type WorldFileResult =
  | { ok: true; world: World }
  | { ok: false; problem: FileProblem };

export const worldFileName = (name: string) => `worlds/${name}.yaml`;
export const worldId = (name: string) => `file:${name}`;

/** Parses one worlds/<name>.yaml: `relationships` (a string, required) and `assertions` (a string, optional). */
export function parseWorldFile(name: string, content: string): WorldFileResult {
  const file = worldFileName(name);
  const fail = (message: string, line?: number): WorldFileResult => ({ ok: false, problem: { file, message, line } });
  let doc: unknown;
  try {
    doc = parse(content);
  } catch (e) {
    const err = e as { message?: string; linePos?: { line: number }[] };
    return fail(String(err.message ?? e).split("\n")[0], err.linePos?.[0]?.line);
  }
  if (typeof doc !== "object" || doc === null || Array.isArray(doc)) return fail(`expected a mapping with "relationships" (and optionally "assertions")`);
  const m = doc as Record<string, unknown>;
  const extra = Object.keys(m).find((k) => k !== "relationships" && k !== "assertions");
  if (extra !== undefined) return fail(`unknown key "${extra}" (expected "relationships" and optionally "assertions")`);
  if (typeof m.relationships !== "string") return fail(`"relationships" must be a string (a block: \`relationships: |\`)`);
  if (m.assertions !== undefined && m.assertions !== null && typeof m.assertions !== "string") return fail(`"assertions" must be a string`);
  const world: World = { id: worldId(name), kind: "saved", name, relationships: m.relationships };
  if (typeof m.assertions === "string") world.assertions = m.assertions;
  return { ok: true, world };
}

/** Parses all world files. A broken one is reported on its own and does not affect the others. */
export function parseWorldFiles(files: WorkspaceFiles["worlds"]): { worlds: World[]; problems: FileProblem[] } {
  const worlds: World[] = [];
  const problems: FileProblem[] = [];
  for (const f of files) {
    const r = parseWorldFile(f.name, f.content);
    if (r.ok) worlds.push(r.world);
    else problems.push(r.problem);
  }
  return { worlds, problems };
}

/** A world as a worlds/*.yaml file (the "copy as YAML" action). */
export function worldToYaml(w: Pick<World, "relationships" | "assertions">): string {
  const body: Record<string, string> = { relationships: w.relationships.endsWith("\n") ? w.relationships : w.relationships + "\n" };
  if (w.assertions) body.assertions = w.assertions;
  return stringify(body);
}

// ---- parse problems of the other files

/** Parse problems of constraints.yaml and contracts.yaml, named after the file (line numbers are the file's own). */
export function parseProblems(constraints: string, contracts: string): FileProblem[] {
  const out: FileProblem[] = [];
  const c = parseConstraints(constraints);
  if (!c.ok) out.push(...c.errors.map((e) => ({ file: CONSTRAINTS_FILE, line: e.line, message: e.message })));
  const r = parseRequirements(contracts);
  if (!r.ok) out.push(...r.errors.map((e) => ({ file: CONTRACTS_FILE, line: e.line, message: e.message })));
  return out;
}

// ---- keeping the store in step with the files

/** What to do with an event of the stream, given the version held now. `changed` always refetches. */
export function onEvent(event: "ready" | "changed", eventVersion: string | undefined, current: string | null): "refetch" | "none" {
  if (event === "changed") return "refetch";
  return eventVersion !== current ? "refetch" : "none";
}

/** Selected world index after the list changed: follow the world by id, else clamp. */
export function reselect(prev: World[], index: number, next: World[]): number {
  const id = prev[index]?.id;
  const found = id === undefined ? -1 : next.findIndex((w) => w.id === id);
  if (found >= 0) return found;
  return Math.max(0, Math.min(index, next.length - 1));
}

/** The saved worlds come from the files; ephemeral ones are kept (after them). */
export function mergeWorlds(current: World[], saved: World[]): World[] {
  return [...saved, ...current.filter((w) => w.kind === "ephemeral")];
}

/**
 * Persistence in file-backed mode: only UI preferences are written; the local-mode workspace (policy, constraints,
 * requirements, saved worlds) already in storage is left alone, so going back to local mode finds it intact.
 */
export function mergePrefs(existingRaw: string | null, nextRaw: string, prefs: string[]): string {
  let existing: { state?: Record<string, unknown>; version?: number } = {};
  try {
    existing = existingRaw ? JSON.parse(existingRaw) : {};
  } catch {
    existing = {};
  }
  const next = JSON.parse(nextRaw) as { state: Record<string, unknown>; version?: number };
  const state = { ...(existing.state ?? {}) };
  for (const k of prefs) if (k in next.state) state[k] = next.state[k];
  return JSON.stringify({ ...next, state });
}

// ---- the event stream

export type Connection = "connecting" | "connected" | "reconnecting";

export interface EventSourceLike {
  readyState: number;
  onopen: ((e: Event) => unknown) | null;
  onerror: ((e: Event) => unknown) | null;
  addEventListener(type: string, fn: (e: MessageEvent<string>) => void): void;
  close(): void;
}

export interface SyncHandlers {
  /** Fetches the files again and applies them; resolves to whether acp was reachable. */
  refetch: () => Promise<boolean>;
  connection: (c: Connection) => void;
  version: () => string | null;
}

const CLOSED = 2;

/**
 * Subscribes to the workspace events. EventSource reconnects by itself when acp restarts; when the browser gives up
 * (the stream closed, eg. a non-200 from a proxy) it is opened again after `retryMs`. Some proxies (Vite's) keep the
 * browser's stream open after acp is gone, so every `heartbeatMs` the workspace is also fetched: that notices acp
 * going away (indicator: reconnecting), and when it is back, opens a new stream. Returns the stop function.
 */
export function connectEvents(
  h: SyncHandlers,
  make: (url: string) => EventSourceLike = (u) => new EventSource(u),
  retryMs = 2000,
  heartbeatMs = 4000,
): () => void {
  let es: EventSourceLike | null = null;
  let timer: ReturnType<typeof setTimeout> | null = null;
  let stopped = false;
  let down = false;

  const open = () => {
    if (stopped) return;
    es?.close();
    const cur = (es = make(EVENTS_PATH));
    h.connection("connecting");
    cur.onopen = () => h.connection("connected");
    cur.onerror = () => {
      h.connection("reconnecting");
      if (cur.readyState === CLOSED && !stopped) {
        cur.close();
        timer = setTimeout(open, retryMs);
      }
    };
    for (const event of ["ready", "changed"] as const) {
      cur.addEventListener(event, (e) => {
        let version: string | undefined;
        try {
          version = (JSON.parse(e.data) as { version?: string }).version;
        } catch {
          /* a malformed event still means "look again" */
        }
        h.connection("connected");
        if (onEvent(event, version, h.version()) === "refetch") void h.refetch();
      });
    }
  };
  open();
  const beat = setInterval(() => {
    void h.refetch().then((reachable) => {
      if (stopped) return;
      if (!reachable) {
        down = true;
        h.connection("reconnecting");
      } else if (down) {
        down = false;
        if (timer) clearTimeout(timer);
        open();
      }
    });
  }, heartbeatMs);
  return () => {
    stopped = true;
    clearInterval(beat);
    if (timer) clearTimeout(timer);
    es?.close();
  };
}
