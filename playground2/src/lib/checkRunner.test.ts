import { describe, expect, it } from "vitest";
import { runLimited, shortError, summarize, shortActor, type ActorStatus } from "./checkRunner";

const tick = () => new Promise((r) => setTimeout(r, 1));

describe("runLimited", () => {
  it("never exceeds the concurrency limit and reports every item", async () => {
    let active = 0, peak = 0;
    const done: Record<string, ActorStatus> = {};
    await runLimited(
      ["a", "b", "c", "d", "e"], 2,
      async (x) => { active++; peak = Math.max(peak, active); await tick(); active--; return x !== "b"; },
      (x, s) => { done[x] = s; },
      new AbortController().signal,
    );
    expect(peak).toBe(2);
    expect(done.a.state).toBe("allowed");
    expect(done.b.state).toBe("denied");
    expect(Object.keys(done)).toHaveLength(5);
  });
  it("turns failures into error statuses without stopping others", async () => {
    const done: Record<string, ActorStatus> = {};
    await runLimited(["a", "b"], 1, async (x) => { if (x === "a") throw new Error("boom"); return true; }, (x, s) => { done[x] = s; }, new AbortController().signal);
    expect(done.a).toEqual({ state: "error", message: "boom" });
    expect(done.b.state).toBe("allowed");
  });
  it("stops and drops results after abort", async () => {
    const ac = new AbortController();
    const seen: string[] = [];
    await runLimited(["a", "b", "c"], 1, async () => { ac.abort(); return true; }, (x, s) => { if (s.state !== "loading") seen.push(x); }, ac.signal);
    expect(seen).toEqual([]);
  });
});

describe("summaries", () => {
  it("groups by state", () => {
    const s = summarize(["a", "b", "c", "d"], { a: { state: "allowed" }, b: { state: "denied" }, c: { state: "error", message: "x" } });
    expect(s).toEqual({ allowed: ["a"], denied: ["b"], pending: 1, errors: [{ actor: "c", message: "x" }] });
  });
  it("shortens", () => {
    expect(shortError("Check: bad; attrs={}; kind=BAD_INPUT")).toBe("Check: bad");
    expect(shortActor("did:user:alice")).toBe("alice");
    expect(shortActor("did:key:z")).toBe("did:key:z");
  });
});
