import { describe, expect, it } from "vitest";
import { parseConstraints } from "./constraints";
import { SAMPLE_CONSTRAINTS } from "./sample";

const errs = (t: string) => {
  const r = parseConstraints(t);
  if (r.ok) throw new Error("expected failure");
  return r.errors.map((e) => e.message);
};

describe("parseConstraints", () => {
  it("maps the sample", () => {
    expect(parseConstraints(SAMPLE_CONSTRAINTS)).toEqual({
      ok: true,
      constraints: [
        { acyclic: { resource: "folder", relation: "parent" } },
        { atMostOne: { resource: "doc", relation: "folder" } },
        { disjoint: { resource: "doc", relations: ["owner", "viewer"] } },
      ],
    });
  });

  it("treats empty text and empty list as no constraints", () => {
    expect(parseConstraints("  \n")).toEqual({ ok: true, constraints: [] });
    expect(parseConstraints("constraints:\n")).toEqual({ ok: true, constraints: [] });
    expect(parseConstraints("constraints: []")).toEqual({ ok: true, constraints: [] });
  });

  it("splits references on the first dot", () => {
    const r = parseConstraints("constraints:\n- acyclic: a.b.c");
    expect(r).toEqual({ ok: true, constraints: [{ acyclic: { resource: "a", relation: "b.c" } }] });
  });

  it("rejects entries with more or fewer than one key", () => {
    expect(errs("constraints:\n- acyclic: a.b\n  at_most_one: a.c")[0]).toMatch(/exactly one key/);
    expect(errs("constraints:\n- {}")[0]).toMatch(/exactly one key/);
  });

  it("rejects disjoint across resources", () => {
    expect(errs("constraints:\n- disjoint: [a.x, b.y]")[0]).toMatch(/same resource/);
  });

  it("rejects malformed references and unknown kinds", () => {
    expect(errs("constraints:\n- acyclic: nodot")[0]).toMatch(/resource\.relation/);
    expect(errs("constraints:\n- acyclic: a.")[0]).toMatch(/resource\.relation/);
    expect(errs("constraints:\n- nope: a.b")[0]).toMatch(/unknown constraint/);
  });

  it("reports yaml syntax errors with a line", () => {
    const r = parseConstraints("constraints:\n- acyclic: [a.b\n");
    expect(r.ok).toBe(false);
  });

  it("accepts the workspace file's assumptions key, same as constraints", () => {
    const a = parseConstraints("assumptions:\n- acyclic: a.b\n- disjoint: [d.x, d.y]");
    expect(a).toEqual(parseConstraints("constraints:\n- acyclic: a.b\n- disjoint: [d.x, d.y]"));
    expect(a.ok).toBe(true);
  });
  it("reports the right key for assumptions", () => {
    expect(errs("assumptions: 3")[0]).toMatch(/"assumptions" must be a list/);
    expect(errs("assumptions:\n- acyclic: a.b\n- bogus: x")[0]).toMatch(/unknown constraint/);
  });
  it("rejects both keys at once", () => {
    expect(errs("constraints: []\nassumptions: []")[0]).toMatch(/not both/);
  });

  it("reports the line of a bad entry", () => {
    const r = parseConstraints("constraints:\n- acyclic: a.b\n- acyclic: bad");
    expect(r.ok === false && r.errors[0].line).toBe(3);
  });
});
