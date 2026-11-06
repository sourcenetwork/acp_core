import { describe, expect, it, vi } from "vitest";
import {
  connectEvents, detectMode, formatProblem, mergePrefs, mergeWorlds, onEvent, parseProblems, parseWorldFile, parseWorldFiles, reselect,
  worldToYaml, type EventSourceLike
} from "./workspaceFiles";
import type { World } from "./types";

const reply = (status: number, body: string) => (async () => new Response(body, { status })) as unknown as typeof fetch;

describe("detectMode", () => {
  it("200 is file-backed", async () => {
    const d = await detectMode(reply(200, JSON.stringify({ dir: "/w", policy: "p", version: "v" })));
    expect(d).toMatchObject({ mode: "files", files: { dir: "/w", constraints: "", contracts: "", worlds: [] } });
  });
  it("404 and network errors are local", async () => {
    expect(await detectMode(reply(404, "not found"))).toEqual({ mode: "local" });
    expect(await detectMode((async () => { throw new TypeError("failed"); }) as unknown as typeof fetch)).toEqual({ mode: "local" });
  });
  it("a bare 500 (dev proxy, nothing behind it) is local; acp's JSON 500 is a broken workspace", async () => {
    expect(await detectMode(reply(500, ""))).toEqual({ mode: "local" });
    expect(await detectMode(reply(500, "Internal Server Error"))).toEqual({ mode: "local" });
    expect(await detectMode(reply(500, JSON.stringify({ code: "internal", message: "loading workspace: no policy.yaml" })))).toEqual({
      mode: "files", files: null, error: "loading workspace: no policy.yaml",
    });
  });
  it("a 200 which isn't a workspace is local", async () => {
    expect(await detectMode(reply(200, "<html>"))).toEqual({ mode: "local" });
  });
});

describe("parseWorldFile", () => {
  it("reads relationships and assertions", () => {
    const r = parseWorldFile("w1", "relationships: |\n  doc:a#owner@did:alice\n  doc:b#owner@did:bob\nassertions: |\n  x\n");
    expect(r).toEqual({ ok: true, world: { id: "file:w1", kind: "saved", name: "w1", relationships: "doc:a#owner@did:alice\ndoc:b#owner@did:bob\n", assertions: "x\n" } });
  });
  it("names the file in its errors", () => {
    for (const bad of ["relationships: [", "- a", "relationships: 3", "relationships: |\n  x\nfoo: 1", ""]) {
      const r = parseWorldFile("foo", bad);
      expect(r.ok).toBe(false);
      if (!r.ok) expect(formatProblem(r.problem)).toMatch(/^worlds\/foo\.yaml( line \d+)?: /);
    }
  });
  it("a broken world doesn't hide the others", () => {
    const r = parseWorldFiles([{ name: "a", content: "relationships: |\n  x" }, { name: "b", content: "oops: [" }, { name: "c", content: "relationships: ''" }]);
    expect(r.worlds.map((w) => w.name)).toEqual(["a", "c"]);
    expect(r.problems).toHaveLength(1);
    expect(r.problems[0].file).toBe("worlds/b.yaml");
  });
  it("round-trips through worldToYaml", () => {
    const w = { relationships: "doc:a#owner@did:alice\ndoc:b#owner@did:bob", assertions: "a: 1" };
    const r = parseWorldFile("x", worldToYaml(w));
    expect(r.ok && r.world.relationships.trim()).toBe(w.relationships);
    expect(worldToYaml(w)).toMatch(/relationships: \|/);
  });
});

describe("parseProblems", () => {
  it("is empty for valid files, in either constraints key", () => {
    expect(parseProblems("assumptions:\n- acyclic: a.b\n", "requirements: []")).toEqual([]);
    expect(parseProblems("", "")).toEqual([]);
  });
  it("names the file and line", () => {
    const p = parseProblems("assumptions:\n- acyclic: nope\n", "requirements:\n- name: x\n  statement: s\n  bogus: 1\n");
    const out = p.map(formatProblem);
    expect(out[0]).toMatch(/^constraints\.yaml line 2: /);
    expect(out.some((m) => /^contracts\.yaml line 4: /.test(m))).toBe(true);
  });
});

describe("sync logic", () => {
  it("changed always refetches; ready only when the version differs", () => {
    expect(onEvent("changed", "v1", "v1")).toBe("refetch");
    expect(onEvent("ready", "v1", "v1")).toBe("none");
    expect(onEvent("ready", "v2", "v1")).toBe("refetch");
    expect(onEvent("ready", "v1", null)).toBe("refetch");
  });

  const w = (id: string, kind: World["kind"] = "saved"): World => ({ id, kind, relationships: "" });
  it("keeps the selected world by id, else clamps", () => {
    expect(reselect([w("a"), w("b")], 1, [w("x"), w("a"), w("b")])).toBe(2);
    expect(reselect([w("a"), w("b")], 1, [w("a")])).toBe(0);
    expect(reselect([], 0, [])).toBe(0);
  });
  it("replaces saved worlds, keeps ephemeral ones", () => {
    expect(mergeWorlds([w("old"), w("e1", "ephemeral")], [w("n1"), w("n2")]).map((x) => x.id)).toEqual(["n1", "n2", "e1"]);
  });
  it("persists only preferences in file-backed mode", () => {
    const existing = JSON.stringify({ state: { policy: "mine", config: { count: 1 }, verifyScope: 3 }, version: 1 });
    const next = JSON.stringify({ state: { policy: "from-files", config: { count: 7 }, verifyScope: 4, autoVerify: true }, version: 1 });
    expect(JSON.parse(mergePrefs(existing, next, ["config", "verifyScope", "autoVerify"])).state).toEqual({
      policy: "mine", config: { count: 7 }, verifyScope: 4, autoVerify: true,
    });
    expect(JSON.parse(mergePrefs(null, next, ["config"])).state).toEqual({ config: { count: 7 } });
    expect(JSON.parse(mergePrefs("garbage", next, ["config"])).state).toEqual({ config: { count: 7 } });
  });
});

class FakeES implements EventSourceLike {
  readyState = 0;
  onopen: ((e: Event) => unknown) | null = null;
  onerror: ((e: Event) => unknown) | null = null;
  closed = false;
  listeners: Record<string, (e: MessageEvent<string>) => void> = {};
  addEventListener(t: string, fn: (e: MessageEvent<string>) => void) { this.listeners[t] = fn; }
  close() { this.closed = true; this.readyState = 2; }
  emit(t: string, version: string) { this.listeners[t]({ data: JSON.stringify({ version }) } as MessageEvent<string>); }
}

describe("connectEvents", () => {
  it("refetches on changed, and on ready only if the version differs; reports the connection", () => {
    const es: FakeES[] = [];
    const version: string | null = "v1";
    const refetch = vi.fn();
    const conn: string[] = [];
    connectEvents({ refetch: async () => { refetch(); return true; }, version: () => version, connection: (c) => conn.push(c) }, () => { const e = new FakeES(); es.push(e); return e; });
    const s = es[0];
    s.onopen?.(new Event("x"));
    s.emit("ready", "v1");
    expect(refetch).not.toHaveBeenCalled();
    s.emit("changed", "v2");
    expect(refetch).toHaveBeenCalledTimes(1);
    s.onerror?.(new Event("x")); // acp went away: EventSource retries by itself
    expect(conn[conn.length - 1]).toBe("reconnecting");
    s.emit("ready", "v3"); // acp restarted with other files
    expect(refetch).toHaveBeenCalledTimes(2);
    expect(conn[conn.length - 1]).toBe("connected");
  });

  it("reopens a stream the browser gave up on, and stops when asked", () => {
    vi.useFakeTimers();
    const made: FakeES[] = [];
    const stop = connectEvents({ refetch: async () => true, version: () => null, connection: vi.fn() }, () => { const e = new FakeES(); made.push(e); return e; }, 1000);
    made[0].readyState = 2;
    made[0].onerror?.(new Event("x"));
    expect(made).toHaveLength(1);
    vi.advanceTimersByTime(1000);
    expect(made).toHaveLength(2);
    stop();
    expect(made[1].closed).toBe(true);
    made[1].readyState = 2;
    made[1].onerror?.(new Event("x"));
    vi.advanceTimersByTime(5000);
    expect(made).toHaveLength(2);
    vi.useRealTimers();
  });

  it("notices a proxy which kept the stream open after acp went away, and opens a new one when it is back", async () => {
    vi.useFakeTimers();
    const made: FakeES[] = [];
    let reachable = true;
    const conn: string[] = [];
    const stop = connectEvents(
      { refetch: async () => reachable, version: () => "v1", connection: (c) => conn.push(c) },
      () => { const e = new FakeES(); made.push(e); return e; }, 1000, 4000,
    );
    made[0].onopen?.(new Event("open"));
    await vi.advanceTimersByTimeAsync(4000);
    expect(conn[conn.length - 1]).toBe("connected");
    reachable = false;
    await vi.advanceTimersByTimeAsync(4000);
    expect(conn[conn.length - 1]).toBe("reconnecting");
    expect(made).toHaveLength(1);
    reachable = true;
    await vi.advanceTimersByTimeAsync(4000);
    expect(made).toHaveLength(2);
    expect(made[0].closed).toBe(true);
    stop();
    vi.useRealTimers();
  });
});

