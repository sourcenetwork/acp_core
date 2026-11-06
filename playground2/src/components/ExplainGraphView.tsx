import { useMemo } from "react";
import { layoutExplain, type ExplainOutcome } from "@/lib/explainGraph";
import { edgePath } from "@/lib/graph";
import { usePanZoom, type View } from "@/lib/usePanZoom";
import type { CheckExplainGraph } from "@/lib/types";

function palette(outcome: ExplainOutcome, dark: boolean) {
  switch (outcome) {
    case "authorized":
      return dark ? { fill: "#14532d", stroke: "#4ade80" } : { fill: "#dcfce7", stroke: "#16a34a" };
    case "exhausted":
      return dark ? { fill: "#4c1d1d", stroke: "#f87171" } : { fill: "#fee2e2", stroke: "#dc2626" };
    case "explored":
      return dark ? { fill: "#4a3410", stroke: "#fbbf24" } : { fill: "#fef3c7", stroke: "#d97706" };
    default:
      return dark ? { fill: "#262626", stroke: "#6b7280" } : { fill: "#f3f4f6", stroke: "#9ca3af" };
  }
}

const LEGEND: { outcome: ExplainOutcome; label: string }[] = [
  { outcome: "authorized", label: "authorized" },
  { outcome: "exhausted", label: "searched, not authorized" },
  { outcome: "explored", label: "explored, search incomplete" },
  { outcome: "unexplored", label: "not explored" },
];

/** Pan/zoom SVG rendering of the engine's search graph: usersets are boxes, operators small circles. */
export function ExplainGraphView({ graph, dark }: { graph: CheckExplainGraph | undefined; dark: boolean }) {
  const layout = useMemo(() => layoutExplain(graph), [graph]);
  const home: View = useMemo(() => ({ x: 0, y: 0, w: Math.max(layout.width, 1), h: Math.max(layout.height, 1) }), [layout]);
  const { svgRef, viewBox, reset, handlers } = usePanZoom(home);

  if (layout.nodes.length === 0) return <p className="text-xs text-muted-foreground">The engine returned no search graph.</p>;

  const fg = dark ? "#e5e7eb" : "#374151";
  const edgeColor = dark ? "#9ca3af" : "#6b7280";
  const okEdge = dark ? "#4ade80" : "#16a34a";
  const labelBg = dark ? "#111111" : "#ffffff";
  const height = Math.round(Math.min(Math.max(layout.height + 24, 220), 460));

  return (
    <div data-testid="explain-graph">
      <div className="relative rounded-md border bg-graph-bg" style={{ height }}>
        <svg
          ref={svgRef}
          className="h-full w-full cursor-grab touch-none select-none active:cursor-grabbing"
          viewBox={viewBox}
          preserveAspectRatio="xMidYMid meet"
          role="img"
          aria-label="Explain graph"
          {...handlers}
        >
          <defs>
            <marker id="explain-arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="7" markerHeight="7" orient="auto-start-reverse" markerUnits="userSpaceOnUse">
              <path d="M0,0 L10,5 L0,10 z" fill="context-stroke" />
            </marker>
          </defs>
          {layout.edges.map((e) => (
            <path
              key={e.id}
              d={edgePath(e.points).d}
              fill="none"
              stroke={e.authorized ? okEdge : edgeColor}
              strokeWidth={e.authorized ? 1.8 : 1.2}
              opacity={e.authorized ? 1 : 0.7}
              markerEnd="url(#explain-arrow)"
              data-testid="explain-edge"
            />
          ))}
          {layout.nodes.map((n) => {
            const c = palette(n.outcome, dark);
            const userset = n.type === "USERSET_NODE";
            return (
              <g key={n.id} transform={`translate(${n.x},${n.y})`} data-testid="explain-node" data-type={n.type} data-outcome={n.outcome} opacity={n.outcome === "unexplored" ? 0.6 : 1}>
                {userset ? (
                  <rect x={-n.width / 2} y={-n.height / 2} width={n.width} height={n.height} rx={6} fill={c.fill} stroke={c.stroke} strokeWidth={1.5} />
                ) : (
                  <circle r={n.width / 2} fill={c.fill} stroke={c.stroke} strokeWidth={1.5} />
                )}
                <text textAnchor="middle" dominantBaseline="central" fontSize={userset ? 11 : 15} fontFamily={userset ? "ui-monospace, monospace" : undefined} fill={fg}>
                  {n.label}
                </text>
                <title>{[n.label, n.detail, n.outcome].filter(Boolean).join("\n")}</title>
              </g>
            );
          })}
          {layout.edges.filter((e) => e.message).map((e) => (
            <g key={`l-${e.id}`} transform={`translate(${e.labelPos.x},${e.labelPos.y})`} data-testid="explain-edge-label">
              <rect x={-e.labelWidth / 2} y={-e.labelHeight / 2} width={e.labelWidth} height={e.labelHeight} rx={3} fill={labelBg} fillOpacity={0.92} />
              <text textAnchor="middle" dominantBaseline="central" fontSize={10} fill={fg}>{e.message}</text>
            </g>
          ))}
        </svg>
        <button
          type="button"
          className="absolute right-1 top-1 rounded border bg-background/80 px-1.5 py-0.5 text-xs text-muted-foreground hover:text-foreground"
          onClick={reset}
        >
          Reset view
        </button>
      </div>
      <ul className="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-muted-foreground">
        {LEGEND.map(({ outcome, label }) => {
          const c = palette(outcome, dark);
          return (
            <li key={outcome} className="flex items-center gap-1">
              <span className="inline-block h-3 w-4 rounded-sm border-2" style={{ background: c.fill, borderColor: c.stroke }} />
              {label}
            </li>
          );
        })}
        <li>{"∪ union, ∩ intersection, − difference"}</li>
        {layout.dropped > 0 && <li>{layout.dropped} edge(s) to unknown nodes omitted</li>}
      </ul>
    </div>
  );
}
