import { useMemo } from "react";
import { buildGraph, edgePath, layoutGraph } from "@/lib/graph";
import { parseRelationships } from "@/lib/relationships";
import { usePanZoom, type View } from "@/lib/usePanZoom";

const ALLOWED = "#16a34a";
const ERROR = "#dc2626";
const ALLOWED_DARKER = (dark: boolean) => (dark ? "#86efac" : "#15803d");
const SELECTED = (dark: boolean) => (dark ? "#fbbf24" : "#d97706");

/** Stable hue per resource type (hash of the name). */
function hueOf(resource: string): number {
  let h = 0;
  for (let i = 0; i < resource.length; i++) h = (h * 31 + resource.charCodeAt(i)) >>> 0;
  return (h * 137) % 360;
}

function nodeColors(resource: string | undefined, dark: boolean) {
  if (resource === undefined) {
    return dark
      ? { fill: "#2b2b2b", stroke: "#9ca3af", text: "#f1f5f9" }
      : { fill: "#f3f4f6", stroke: "#4b5563", text: "#111827" };
  }
  const h = hueOf(resource);
  return dark
    ? { fill: `hsl(${h} 35% 24%)`, stroke: `hsl(${h} 60% 65%)`, text: "#f8fafc" }
    : { fill: `hsl(${h} 70% 90%)`, stroke: `hsl(${h} 55% 40%)`, text: "#111827" };
}

export interface ActorMark { state: "loading" | "allowed" | "denied" | "error"; message?: string }

export interface GraphSelection {
  /** Selected object node key (`doc:doc0`), if any. */
  object: string | null;
  /** Per-actor check outcome for the selected (object, permission). Empty when nothing is selected. */
  actors: Record<string, ActorMark | undefined>;
  /** Actor whose explanation is open. */
  explained: string | null;
  onSelectObject: (key: string | null) => void;
  onSelectActor: (did: string) => void;
}

const BOUND = (dark: boolean) => (dark ? "#38bdf8" : "#0369a1");

/** `marks` maps a node key to the names of the formula variables bound to it (a counterexample's bindings). */
export function RelationshipGraph({ text, dark, selection, fill, marks }: { text: string; dark: boolean; selection?: GraphSelection; fill?: boolean; marks?: Record<string, string[]> }) {
  const { layout, issues } = useMemo(() => {
    const parsed = parseRelationships(text);
    return { layout: layoutGraph(buildGraph(parsed.relationships)), issues: parsed.issues };
  }, [text]);

  const home: View = useMemo(
    () => ({ x: 0, y: 0, w: Math.max(layout.width, 1), h: Math.max(layout.height, 1) }),
    [layout],
  );
  const { svgRef, viewBox, reset, handlers, wasDragged } = usePanZoom(home);

  const resources = useMemo(
    () => [...new Set(layout.nodes.filter((n) => n.resource).map((n) => n.resource!))].sort(),
    [layout],
  );
  const hasActors = layout.nodes.some((n) => n.kind === "actor");
  const hasUserset = layout.edges.some((e) => e.usersetRelation);

  const fg = dark ? "#e5e7eb" : "#374151";
  const edgeColor = dark ? "#9ca3af" : "#6b7280";
  const labelBg = dark ? "#111111" : "#ffffff";
  const height = Math.round(Math.min(Math.max(layout.height + 24, 200), 720));

  return (
    <div data-testid="relationship-graph" className={fill ? "flex min-h-0 flex-1 flex-col" : undefined}>
      {issues.length > 0 && (
        <p className="mb-1 text-xs text-destructive" data-testid="graph-issues">
          {issues.length} line(s) could not be parsed:{" "}
          {issues.slice(0, 3).map((i) => `line ${i.line} (${i.reason})`).join(", ")}
          {issues.length > 3 ? ", ..." : ""}
        </p>
      )}
      {layout.nodes.length === 0 ? (
        <p className="text-xs text-muted-foreground">No relationships to draw.</p>
      ) : (
        <>
          <div className={fill ? "relative min-h-60 flex-1 rounded-md border bg-graph-bg" : "relative rounded-md border bg-graph-bg"} style={fill ? undefined : { height }}>
            <svg
              ref={svgRef}
              className="absolute inset-0 h-full w-full cursor-grab touch-none select-none active:cursor-grabbing"
              viewBox={viewBox}
              preserveAspectRatio="xMidYMid meet"
              role="img"
              aria-label="Relationship graph"
              {...handlers}
              onClick={() => { if (!wasDragged()) selection?.onSelectObject(null); }}
            >
              <defs>
                <marker id="arrow" viewBox="0 0 10 10" refX="9" refY="5" markerWidth="8" markerHeight="8" orient="auto-start-reverse" markerUnits="userSpaceOnUse">
                  <path d="M0,0 L10,5 L0,10 z" fill="context-stroke" />
                </marker>
              </defs>
              <g>
                {layout.edges.map((e) => (
                  <path
                    key={e.id}
                    d={edgePath(e.points).d}
                    fill="none"
                    stroke={edgeColor}
                    strokeWidth={1.4}
                    strokeDasharray={e.usersetRelation ? "5 4" : undefined}
                    markerEnd="url(#arrow)"
                    data-testid="graph-edge"
                  />
                ))}
                {layout.nodes.map((n) => {
                  const c = nodeColors(n.resource, dark);
                  const isActor = n.kind === "actor";
                  const selected = selection?.object === n.key;
                  const mark = isActor ? selection?.actors[n.key] : undefined;
                  const active = selection?.object != null;
                  // Dim whatever the current question rules out: denied actors, and objects other than the selected one.
                  const dim = active && (isActor ? mark?.state === "denied" : !selected);
                  const allowed = mark?.state === "allowed";
                  const explained = selection?.explained === n.key;
                  const boundNames = marks?.[n.key];
                  const clickable = !!selection && (!isActor || active);
                  const stroke = allowed ? ALLOWED : mark?.state === "error" ? ERROR : c.stroke;
                  return (
                    <g
                      key={n.key}
                      transform={`translate(${n.x - n.width / 2},${n.y - n.height / 2})`}
                      data-testid="graph-node"
                      data-kind={n.kind}
                      data-key={n.key}
                      data-selected={selected || undefined}
                      data-state={mark?.state}
                      opacity={dim ? 0.35 : 1}
                      className={clickable ? "cursor-pointer" : undefined}
                      onClick={(ev) => {
                        ev.stopPropagation();
                        if (!selection || wasDragged()) return;
                        if (isActor) selection.onSelectActor(n.key);
                        else selection.onSelectObject(selected ? null : n.key);
                      }}
                    >
                      {boundNames && (
                        <>
                          <rect x={-8} y={-8} width={n.width + 16} height={n.height + 16} rx={isActor ? n.height / 2 + 8 : 12} fill="none" stroke={BOUND(dark)} strokeWidth={2.5} strokeDasharray="3 3" data-testid="bound-ring" />
                          <g transform={`translate(${n.width / 2 - (boundNames.join(", ").length * 4.5 + 8)},-30)`} data-testid="bound-badge" data-names={boundNames.join(",")}>
                            <rect width={boundNames.join(", ").length * 9 + 16} height={21} rx={10.5} fill={BOUND(dark)} />
                            <text x={(boundNames.join(", ").length * 9 + 16) / 2} y={10.5} textAnchor="middle" dominantBaseline="central" fontSize={14} fontWeight={600} fontFamily="ui-monospace, monospace" fill={dark ? "#0c1a2b" : "#ffffff"}>{boundNames.join(", ")}</text>
                          </g>
                        </>
                      )}
                      {(selected || explained) && (
                        <rect x={-4} y={-4} width={n.width + 8} height={n.height + 8} rx={isActor ? n.height / 2 + 4 : 9} fill="none" stroke={explained ? ALLOWED_DARKER(dark) : SELECTED(dark)} strokeWidth={2.5} strokeDasharray={explained ? "4 3" : undefined} />
                      )}
                      <rect
                        width={n.width}
                        height={n.height}
                        rx={isActor ? n.height / 2 : 6}
                        fill={allowed ? (dark ? "#14532d" : "#bbf7d0") : c.fill}
                        stroke={stroke}
                        strokeWidth={allowed || selected ? 2.5 : 1.5}
                        strokeDasharray={mark?.state === "loading" ? "3 3" : undefined}
                        className={mark?.state === "loading" ? "animate-pulse" : undefined}
                      />
                      <text x={n.width / 2} y={n.height / 2} textAnchor="middle" dominantBaseline="central" fontSize={12} fontFamily="ui-monospace, monospace" fill={allowed && !dark ? "#111827" : c.text}>
                        {isActor ? (mark?.state === "allowed" ? "\u2713 " : mark?.state === "denied" ? "\u2717 " : mark?.state === "error" ? "! " : "\u{1F464} ") : ""}{n.label}
                      </text>
                      {mark?.state === "error" && <title>{mark.message}</title>}
                    </g>
                  );
                })}
                {layout.edges.map((e) => (
                  <g key={`l-${e.id}`} transform={`translate(${e.labelPos.x},${e.labelPos.y})`} data-testid="graph-edge-label">
                    <rect x={-e.labelWidth / 2} y={-e.labelHeight / 2} width={e.labelWidth} height={e.labelHeight} rx={3} fill={labelBg} fillOpacity={0.92} />
                    <text textAnchor="middle" dominantBaseline="central" fontSize={11} fontStyle={e.usersetRelation ? "italic" : undefined} fill={fg}>
                      {e.label}
                    </text>
                  </g>
                ))}
              </g>
            </svg>
            <button
              type="button"
              className="absolute right-1 top-1 rounded border bg-background/80 px-1.5 py-0.5 text-xs text-muted-foreground hover:text-foreground"
              onClick={reset}
            >
              Reset view
            </button>
          </div>
          <ul className="mt-1 flex flex-wrap items-center gap-x-3 gap-y-1 text-xs text-muted-foreground" data-testid="graph-legend">
            {resources.map((r) => {
              const c = nodeColors(r, dark);
              return (
                <li key={r} className="flex items-center gap-1">
                  <span className="inline-block h-3 w-4 rounded-sm border-2" style={{ background: c.fill, borderColor: c.stroke }} />
                  {r}
                </li>
              );
            })}
            {hasActors && (
              <li className="flex items-center gap-1">
                <span className="inline-block h-3 w-4 rounded-full border-2" style={{ background: nodeColors(undefined, dark).fill, borderColor: nodeColors(undefined, dark).stroke }} />
                actor
              </li>
            )}
            {marks && Object.keys(marks).length > 0 && (
              <li className="flex items-center gap-1">
                <span className="inline-block h-3 w-4 rounded-sm border-2 border-dotted" style={{ borderColor: BOUND(dark) }} />
                bound to a formula variable
              </li>
            )}
            <li className="flex items-center gap-1">
              <svg width="22" height="6"><line x1="0" y1="3" x2="22" y2="3" stroke={edgeColor} strokeWidth="1.5" /></svg>
              relation
            </li>
            {hasUserset && (
              <li className="flex items-center gap-1">
                <svg width="22" height="6"><line x1="0" y1="3" x2="22" y2="3" stroke={edgeColor} strokeWidth="1.5" strokeDasharray="5 4" /></svg>
                userset (#relation)
              </li>
            )}
          </ul>
        </>
      )}
    </div>
  );
}
