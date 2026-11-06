import { describe, expect, it } from "vitest";
import { buildGraph, layoutGraph } from "./graph";
import { parseRelationships } from "./relationships";

describe("parseRelationships", () => {
  it("parses actor, object and userset subjects", () => {
    const { relationships, issues } = parseRelationships(
      "doc:d1#viewer@did:user:alice\ndoc:d1#folder@folder:f1\nfolder:f1#viewer@team:t1#admin\n",
    );
    expect(issues).toEqual([]);
    expect(relationships).toEqual([
      { object: { resource: "doc", id: "d1" }, relation: "viewer", subject: { kind: "actor", did: "did:user:alice" } },
      { object: { resource: "doc", id: "d1" }, relation: "folder", subject: { kind: "object", object: { resource: "folder", id: "f1" } } },
      {
        object: { resource: "folder", id: "f1" },
        relation: "viewer",
        subject: { kind: "userset", object: { resource: "team", id: "t1" }, relation: "admin" },
      },
    ]);
  });

  it("tolerates blank lines and CRLF, reports bad lines with 1-based numbers", () => {
    const { relationships, issues } = parseRelationships(
      "doc:d1#viewer@did:user:a\r\n\r\n   \ngarbage\ndoc:d2@did:x\ndoc:d3#v@\nfoo#v@bar:1\n",
    );
    expect(relationships).toHaveLength(1);
    expect(issues.map((i) => i.line)).toEqual([4, 5, 6, 7]);
    expect(issues[0].text).toBe("garbage");
  });

  it("handles empty input", () => {
    expect(parseRelationships("")).toEqual({ relationships: [], issues: [] });
  });
});

describe("graph", () => {
  const text = [
    "doc:d1#viewer@did:user:alice",
    "doc:d1#viewer@did:user:bob",
    "doc:d1#folder@folder:f1",
    "doc:d1#parent@folder:f1",
    "folder:f1#parent@folder:f1",
    "folder:f1#viewer@team:t1#admin",
  ].join("\n");

  it("builds nodes and labelled edges", () => {
    const g = buildGraph(parseRelationships(text).relationships);
    expect(g.nodes.map((n) => n.key).sort()).toEqual(
      ["did:user:alice", "did:user:bob", "doc:d1", "folder:f1", "team:t1"].sort(),
    );
    expect(g.nodes.find((n) => n.key === "did:user:alice")?.kind).toBe("actor");
    const us = g.edges.find((e) => e.usersetRelation);
    expect(us).toMatchObject({ from: "folder:f1", to: "team:t1", label: "viewer (#admin)" });
  });

  it("lays out parallel edges and self loops without NaN", () => {
    const l = layoutGraph(buildGraph(parseRelationships(text).relationships));
    expect(l.edges).toHaveLength(6);
    const parallel = l.edges.filter((e) => e.from === "doc:d1" && e.to === "folder:f1");
    expect(parallel).toHaveLength(2);
    expect(parallel[0].labelPos).not.toEqual(parallel[1].labelPos);
    expect(l.edges.some((e) => e.selfLoop)).toBe(true);
    const all = [...l.nodes.map((n) => [n.x, n.y]), ...l.edges.flatMap((e) => e.points.map((p) => [p.x, p.y]))];
    expect(all.flat().every(Number.isFinite)).toBe(true);
    expect(l.width).toBeGreaterThan(0);
  });

  it("empty graph", () => {
    expect(layoutGraph(buildGraph([])).width).toBe(0);
  });
});
