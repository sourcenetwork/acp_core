import { describe, expect, it } from "vitest";
import { layoutExplain, outcomeOf } from "./explainGraph";
import type { CheckExplainGraph } from "./types";

// Shape of a real response: nodeType omitted for unions, rootNodeId not the request node, flags omitted when false.
const graph: CheckExplainGraph = {
  rootNodeId: "1",
  nodes: [
    { id: "0", nodeType: "USERSET_NODE", text: "(doc:doc0, view)", detail: "reason: Request", result: { authorized: true, explored: true, exhausted: true } },
    { id: "1", text: "+", result: { authorized: true, explored: true, exhausted: true } },
    { id: "2", nodeType: "USERSET_NODE", text: "(actor:did:user:alice)", result: { explored: true } },
    { id: "3", nodeType: "INTERSECTION_NODE", text: "&", result: {} },
  ],
  edges: [
    { sourceNodeId: "0", destNodeId: "1" },
    { sourceNodeId: "1", destNodeId: "2", message: "viewer" },
    { sourceNodeId: "1", destNodeId: "3" },
    { sourceNodeId: "1", destNodeId: "missing" },
  ],
};

describe("layoutExplain", () => {
  const l = layoutExplain(graph);
  it("lays out every node and keeps valid edges only", () => {
    expect(l.nodes).toHaveLength(4);
    expect(l.edges).toHaveLength(3);
    expect(l.dropped).toBe(1);
    expect(l.width).toBeGreaterThan(0);
    expect(l.height).toBeGreaterThan(0);
  });
  it("draws top to bottom by edge direction, ignoring rootNodeId", () => {
    const y = (id: string) => l.nodes.find((n) => n.id === id)!.y;
    expect(y("0")).toBeLessThan(y("1"));
    expect(y("1")).toBeLessThan(y("2"));
  });
  it("maps node types, labels and outcomes", () => {
    const n = (id: string) => l.nodes.find((x) => x.id === id)!;
    expect(n("1").type).toBe("UNION_NODE");
    expect(n("1").label).toBe("∪");
    expect(n("3").label).toBe("∩");
    expect(n("0").label).toBe("(doc:doc0, view)");
    expect(n("0").outcome).toBe("authorized");
    expect(n("2").outcome).toBe("explored");
    expect(n("3").outcome).toBe("unexplored");
  });
  it("gives edge messages a label box and nodes stay inside the layout", () => {
    const e = l.edges.find((x) => x.message === "viewer")!;
    expect(e.labelWidth).toBeGreaterThan(0);
    expect(e.points.length).toBeGreaterThanOrEqual(2);
    for (const n of l.nodes) {
      expect(n.x - n.width / 2).toBeGreaterThanOrEqual(0);
      expect(n.x + n.width / 2).toBeLessThanOrEqual(l.width);
    }
  });
  it("handles empty and undefined graphs", () => {
    expect(layoutExplain(undefined).nodes).toEqual([]);
    expect(layoutExplain({}).width).toBe(0);
  });
  it("tolerates cycles and duplicate ids", () => {
    const c = layoutExplain({ nodes: [{ id: "a" }, { id: "a" }, { id: "b" }], edges: [{ sourceNodeId: "a", destNodeId: "b" }, { sourceNodeId: "b", destNodeId: "a" }] });
    expect(c.nodes).toHaveLength(2);
    expect(c.edges).toHaveLength(2);
  });
});

describe("outcomeOf", () => {
  it("ranks authorized over exhausted over explored", () => {
    expect(outcomeOf({ authorized: true, exhausted: true })).toBe("authorized");
    expect(outcomeOf({ exhausted: true, explored: true })).toBe("exhausted");
    expect(outcomeOf({ explored: true })).toBe("explored");
    expect(outcomeOf(undefined)).toBe("unexplored");
  });
});
