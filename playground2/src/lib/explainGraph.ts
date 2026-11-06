// Pure: lays out the engine's CheckExplainGraph (a debug structure whose shape is not stable) with dagre.
// Only nodes and edges are used; `rootNodeId` is ignored on purpose.
import dagre from "@dagrejs/dagre";
import { edgePath, routeThroughLabel, textWidth, type Point } from "./graph";
import type { CheckExplainGraph, ExplainNodeType } from "./types";

export type ExplainOutcome = "authorized" | "exhausted" | "explored" | "unexplored";

export interface ExplainLaidNode {
  id: string;
  type: ExplainNodeType;
  /** Text drawn in the node: the userset `(doc:doc0, view)` or the operator symbol. */
  label: string;
  detail: string;
  outcome: ExplainOutcome;
  x: number; // center
  y: number;
  width: number;
  height: number;
}

export interface ExplainLaidEdge {
  id: string;
  from: string;
  to: string;
  message: string;
  points: Point[];
  labelPos: Point;
  labelWidth: number;
  labelHeight: number;
  /** The edge leads to an authorized node. */
  authorized: boolean;
}

export interface ExplainLayout {
  nodes: ExplainLaidNode[];
  edges: ExplainLaidEdge[];
  width: number;
  height: number;
  /** Edges that referenced a node missing from the graph and were left out. */
  dropped: number;
}

/** authorized: a path was found; exhausted: searched fully without one; explored: visited, search incomplete. */
export function outcomeOf(r: { authorized?: boolean; explored?: boolean; exhausted?: boolean } | undefined): ExplainOutcome {
  if (r?.authorized) return "authorized";
  if (r?.exhausted) return "exhausted";
  if (r?.explored) return "explored";
  return "unexplored";
}

const OP_SYMBOL: Record<Exclude<ExplainNodeType, "USERSET_NODE">, string> = {
  UNION_NODE: "∪",
  DIFF_NODE: "−",
  INTERSECTION_NODE: "∩",
};
const OPERATOR_SIZE = 30;
const USERSET_H = 34;
const LABEL_H = 16;

export function layoutExplain(graph: CheckExplainGraph | undefined): ExplainLayout {
  const dg = new dagre.graphlib.Graph({ multigraph: true });
  dg.setGraph({ rankdir: "TB", nodesep: 28, ranksep: 46, edgesep: 16, marginx: 16, marginy: 16 });
  dg.setDefaultEdgeLabel(() => ({}));

  const inputs = (graph?.nodes ?? []).filter((n) => n.id !== undefined);
  const seen = new Set<string>();
  const meta = new Map<string, Omit<ExplainLaidNode, "x" | "y">>();
  for (const n of inputs) {
    if (seen.has(n.id)) continue;
    seen.add(n.id);
    const type = n.nodeType ?? "UNION_NODE";
    const userset = type === "USERSET_NODE";
    const label = userset ? n.text || n.id : OP_SYMBOL[type];
    const width = userset ? Math.min(Math.max(textWidth(label) + 24, 90), 340) : OPERATOR_SIZE;
    const height = userset ? USERSET_H : OPERATOR_SIZE;
    meta.set(n.id, { id: n.id, type, label, detail: n.detail ?? "", outcome: outcomeOf(n.result), width, height });
    dg.setNode(n.id, { width, height });
  }

  const kept: { id: string; from: string; to: string; message: string }[] = [];
  let dropped = 0;
  (graph?.edges ?? []).forEach((e, i) => {
    const from = e.sourceNodeId ?? "", to = e.destNodeId ?? "";
    if (!meta.has(from) || !meta.has(to)) { dropped++; return; }
    const message = e.message ?? "";
    const id = `e${i}`;
    kept.push({ id, from, to, message });
    dg.setEdge(from, to, message ? { width: textWidth(message) + 10, height: LABEL_H, labelpos: "c" } : {}, id);
  });
  dagre.layout(dg);

  const nodes: ExplainLaidNode[] = [...meta.values()].map((m) => {
    const p = dg.node(m.id);
    return { ...m, x: p.x, y: p.y };
  });
  const byId = new Map(nodes.map((n) => [n.id, n]));

  let minX = Infinity, minY = Infinity, maxX = -Infinity, maxY = -Infinity;
  const grow = (x: number, y: number, hw = 0, hh = 0) => {
    minX = Math.min(minX, x - hw); maxX = Math.max(maxX, x + hw);
    minY = Math.min(minY, y - hh); maxY = Math.max(maxY, y + hh);
  };
  nodes.forEach((n) => grow(n.x, n.y, n.width / 2, n.height / 2));

  const edges: ExplainLaidEdge[] = kept.map((e) => {
    const d = dg.edge(e.from, e.to, e.id);
    const lw = e.message ? textWidth(e.message) + 10 : 0;
    const labelPos = { x: d.x ?? 0, y: d.y ?? 0 };
    let points: Point[] = d.points.map((p: Point) => ({ x: p.x, y: p.y }));
    if (e.message) points = routeThroughLabel(points, labelPos, lw, LABEL_H);
    edgePath(points).hull.forEach((p) => grow(p.x, p.y));
    if (e.message) grow(labelPos.x, labelPos.y, lw / 2, LABEL_H / 2);
    return { ...e, points, labelPos, labelWidth: lw, labelHeight: e.message ? LABEL_H : 0, authorized: byId.get(e.to)!.outcome === "authorized" };
  });

  if (!isFinite(minX)) return { nodes, edges, width: 0, height: 0, dropped };
  const m = 16;
  const dx = m - minX, dy = m - minY;
  const shift = (p: Point) => { p.x += dx; p.y += dy; };
  nodes.forEach(shift);
  edges.forEach((e) => { e.points.forEach(shift); shift(e.labelPos); });
  return { nodes, edges, width: maxX - minX + 2 * m, height: maxY - minY + 2 * m, dropped };
}
