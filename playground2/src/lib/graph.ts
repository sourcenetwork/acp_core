// Builds a graph model (nodes + labelled edges) from parsed relationships and lays it out with dagre.
import dagre from "@dagrejs/dagre";
import { objectKey, type Relationship } from "./relationships";

export interface GraphNode {
  key: string;
  kind: "object" | "actor";
  /** Resource type for objects; undefined for actors. */
  resource?: string;
  label: string;
}

export interface GraphEdge {
  id: string;
  from: string;
  to: string;
  relation: string;
  /** Set for userset subjects (`team:t1#admin`). */
  usersetRelation?: string;
  label: string;
}

export interface GraphModel {
  nodes: GraphNode[];
  edges: GraphEdge[];
}

export function buildGraph(rels: Relationship[]): GraphModel {
  const nodes = new Map<string, GraphNode>();
  const edges: GraphEdge[] = [];
  const addObject = (o: { resource: string; id: string }) => {
    const key = objectKey(o);
    if (!nodes.has(key)) nodes.set(key, { key, kind: "object", resource: o.resource, label: key });
    return key;
  };
  for (const r of rels) {
    const from = addObject(r.object);
    let to: string;
    let usersetRelation: string | undefined;
    if (r.subject.kind === "actor") {
      to = r.subject.did;
      if (!nodes.has(to)) nodes.set(to, { key: to, kind: "actor", label: to });
    } else {
      to = addObject(r.subject.object);
      if (r.subject.kind === "userset") usersetRelation = r.subject.relation;
    }
    edges.push({
      id: `e${edges.length}`,
      from,
      to,
      relation: r.relation,
      usersetRelation,
      label: usersetRelation ? `${r.relation} (#${usersetRelation})` : r.relation,
    });
  }
  return { nodes: [...nodes.values()], edges };
}

export interface Point {
  x: number;
  y: number;
}

export interface LaidNode extends GraphNode {
  x: number; // center
  y: number;
  width: number;
  height: number;
}

export interface LaidEdge extends GraphEdge {
  /** Path points from source boundary to target boundary. */
  points: Point[];
  labelPos: Point;
  labelWidth: number;
  labelHeight: number;
  selfLoop: boolean;
}

/**
 * SVG path for an edge route. Short edges (3 points or fewer) run straight into the label, whose box hides
 * the bend; longer routes are smoothed (Catmull-Rom as cubic beziers). `hull` holds every point the curve is
 * bounded by, for computing the layout extents.
 */
export function edgePath(pts: Point[]): { d: string; hull: Point[] } {
  if (pts.length < 2) return { d: "", hull: pts };
  if (pts.length <= 3) return { d: "M" + pts.map((p) => `${p.x},${p.y}`).join(" L"), hull: pts };
  let d = `M${pts[0].x},${pts[0].y}`;
  const hull: Point[] = [pts[0]];
  for (let i = 0; i < pts.length - 1; i++) {
    const p0 = pts[i - 1] ?? pts[i], p1 = pts[i], p2 = pts[i + 1], p3 = pts[i + 2] ?? p2;
    const c1 = { x: p1.x + (p2.x - p0.x) / 6, y: p1.y + (p2.y - p0.y) / 6 };
    const c2 = { x: p2.x - (p3.x - p1.x) / 6, y: p2.y - (p3.y - p1.y) / 6 };
    d += ` C${c1.x},${c1.y} ${c2.x},${c2.y} ${p2.x},${p2.y}`;
    hull.push(c1, c2, p2);
  }
  return { d, hull };
}

export interface Layout {
  nodes: LaidNode[];
  edges: LaidEdge[];
  width: number;
  height: number;
}

// Rough text metrics for a 12px sans-serif font; the SVG sizes boxes from these.
export const charWidth = 6.6;
export const textWidth = (s: string) => Math.ceil(s.length * charWidth);
const NODE_H = 32;
const LABEL_H = 18;
const LOOP_W = 36;

/**
 * dagre's route bends around the label's dummy node, which makes smoothed curves overshoot.
 * Replace the points inside the label box with the label centre so the edge passes straight through it.
 */
export function routeThroughLabel(raw: Point[], label: Point, w: number, h: number): Point[] {
  const inside = (p: Point) => Math.abs(p.x - label.x) <= w / 2 + 2 && Math.abs(p.y - label.y) <= h / 2 + 2;
  const out: Point[] = [];
  let placed = false;
  raw.forEach((p, i) => {
    const interior = i > 0 && i < raw.length - 1;
    if (interior && inside(p)) {
      if (!placed) { out.push({ ...label }); placed = true; }
    } else {
      out.push({ x: p.x, y: p.y });
    }
  });
  if (!placed) {
    // Insert the label after the last point that precedes it along the route.
    let k = 0, best = Infinity;
    out.forEach((p, i) => {
      const dist = Math.hypot(p.x - label.x, p.y - label.y);
      if (dist < best) { best = dist; k = i; }
    });
    out.splice(Math.min(k + 1, out.length - 1), 0, { ...label });
  }
  return out;
}

/** Point where the ray from the node's centre towards `toward` leaves the node's box. */
function clipToNode(n: { x: number; y: number; width: number; height: number }, toward: Point): Point {
  const dx = toward.x - n.x, dy = toward.y - n.y;
  if (dx === 0 && dy === 0) return { x: n.x, y: n.y };
  const t = Math.min(
    dx === 0 ? Infinity : n.width / 2 / Math.abs(dx),
    dy === 0 ? Infinity : n.height / 2 / Math.abs(dy),
  );
  return { x: n.x + dx * t, y: n.y + dy * t };
}

export function layoutGraph(g: GraphModel): Layout {
  const dg = new dagre.graphlib.Graph({ multigraph: true });
  dg.setGraph({ rankdir: "LR", nodesep: 24, ranksep: 90, edgesep: 24, marginx: 16, marginy: 16 });
  dg.setDefaultEdgeLabel(() => ({}));

  const sizes = new Map<string, { width: number; height: number }>();
  for (const n of g.nodes) {
    const size = { width: textWidth(n.label) + (n.kind === "actor" ? 44 : 28), height: NODE_H };
    sizes.set(n.key, size);
    dg.setNode(n.key, size);
  }
  const selfCount = new Map<string, number>();
  const selfIndex = new Map<string, number>();
  for (const e of g.edges) {
    const w = textWidth(e.label) + 10;
    if (e.from === e.to) {
      selfIndex.set(e.id, selfCount.get(e.from) ?? 0);
      selfCount.set(e.from, (selfCount.get(e.from) ?? 0) + 1);
      continue;
    }
    dg.setEdge(e.from, e.to, { width: w, height: LABEL_H, labelpos: "c" }, e.id);
  }
  dagre.layout(dg);

  const nodes: LaidNode[] = g.nodes.map((n) => {
    const p = dg.node(n.key);
    return { ...n, x: p.x, y: p.y, width: p.width, height: p.height };
  });
  const byKey = new Map(nodes.map((n) => [n.key, n]));

  let minX = Infinity, minY = Infinity, maxX = -Infinity, maxY = -Infinity;
  const grow = (x: number, y: number, hw = 0, hh = 0) => {
    minX = Math.min(minX, x - hw); maxX = Math.max(maxX, x + hw);
    minY = Math.min(minY, y - hh); maxY = Math.max(maxY, y + hh);
  };
  nodes.forEach((n) => grow(n.x, n.y, n.width / 2, n.height / 2));

  const edges: LaidEdge[] = g.edges.map((e) => {
    const lw = textWidth(e.label) + 10;
    if (e.from === e.to) {
      const n = byKey.get(e.from)!;
      const i = selfIndex.get(e.id) ?? 0;
      // Loop above the node, widening with each additional self loop.
      const top = n.y - n.height / 2;
      const rise = 26 + i * (LABEL_H + 10);
      const half = Math.min(n.width / 2 - 6, 18 + i * 4);
      const points = [
        { x: n.x - half, y: top },
        { x: n.x - half - LOOP_W / 2, y: top - rise },
        { x: n.x + half + LOOP_W / 2, y: top - rise },
        { x: n.x + half, y: top },
      ];
      const labelPos = { x: n.x, y: top - rise - LABEL_H / 2 - 2 };
      grow(labelPos.x, labelPos.y, lw / 2, LABEL_H / 2);
      return { ...e, points, labelPos, labelWidth: lw, labelHeight: LABEL_H, selfLoop: true };
    }
    const d = dg.edge(e.from, e.to, e.id);
    const labelPos = { x: d.x, y: d.y };
    const src = byKey.get(e.from)!, dst = byKey.get(e.to)!;
    // Short (adjacent-rank) edges: aim straight at the label from each node centre, clipped to the node box.
    const points =
      d.points.length <= 3
        ? [clipToNode(src, labelPos), labelPos, clipToNode(dst, labelPos)]
        : routeThroughLabel(d.points, labelPos, lw, LABEL_H);
    edgePath(points).hull.forEach((p) => grow(p.x, p.y));
    grow(labelPos.x, labelPos.y, lw / 2, LABEL_H / 2);
    return { ...e, points, labelPos, labelWidth: lw, labelHeight: LABEL_H, selfLoop: false };
  });

  if (!isFinite(minX)) return { nodes, edges, width: 0, height: 0 };
  // Normalise so the layout starts at (0,0) with a margin.
  const m = 16;
  const dx = m - minX, dy = m - minY;
  const shift = (p: Point) => { p.x += dx; p.y += dy; };
  nodes.forEach(shift);
  edges.forEach((e) => { e.points.forEach(shift); shift(e.labelPos); });
  return { nodes, edges, width: maxX - minX + 2 * m, height: maxY - minY + 2 * m };
}
