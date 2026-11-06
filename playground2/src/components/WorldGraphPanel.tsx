import { useEffect, useMemo, useRef, useState } from "react";
import { X } from "lucide-react";
import { check, explain } from "@/lib/api";
import { runLimited, shortActor, shortError, summarize, type ActorStatus } from "@/lib/checkRunner";
import { parseConstraints } from "@/lib/constraints";
import { buildGraph } from "@/lib/graph";
import { checkableNames, defaultName, parsePolicyResources } from "@/lib/policyResources";
import { parseRelationships } from "@/lib/relationships";
import type { CheckExplainGraph, CheckRequest } from "@/lib/types";
import { useWorkspace } from "@/stores/workspaceStore";
import { ExplainGraphView } from "./ExplainGraphView";
import { RelationshipGraph, type ActorMark } from "./RelationshipGraph";

const CONCURRENCY = 4;

type ExplainState =
  | { state: "loading" }
  | { state: "error"; message: string }
  | { state: "done"; allowed: boolean; graph: CheckExplainGraph | undefined };

/**
 * A world's graph plus the interaction on top of it: pick an object and a permission or relation, see which actors
 * hold it (one Check per actor), click an actor to see why (Explain).
 */
export interface WorldGraphInitial {
  /** Object node key (`doc:doc0`). Ignored if the world has no such object. */
  object?: string;
  permission?: string;
  /** Actor whose explanation opens at once. */
  actor?: string;
}

export function WorldGraphPanel({ text, dark, initial, bound }: { text: string; dark: boolean; initial?: WorldGraphInitial; bound?: Record<string, string[]> }) {
  const policy = useWorkspace((s) => s.policy);
  const constraintsText = useWorkspace((s) => s.constraints);

  const resources = useMemo(() => parsePolicyResources(policy), [policy]);
  const { objects, actors } = useMemo(() => {
    const nodes = buildGraph(parseRelationships(text).relationships).nodes;
    return {
      objects: new Map(nodes.filter((n) => n.kind === "object").map((n) => [n.key, n.resource!])),
      actors: nodes.filter((n) => n.kind === "actor").map((n) => n.key),
    };
  }, [text]);

  // The scenario is the current system model (as in Generate) plus this world. Constraints play no part in a
  // Check, so an unparseable constraints text is not an error here: it is sent as no constraints.
  const scenario = useMemo(() => {
    const parsed = parseConstraints(constraintsText);
    return {
      systemModel: { policy, constraints: parsed.ok ? parsed.constraints : [] },
      world: { relationships: text },
    };
  }, [policy, constraintsText, text]);

  // `initial` only seeds the state (the panel is remounted per world): everything stays editable.
  const [picked, setPicked] = useState<string | null>(initial?.object ?? null);
  const [permission, setPermission] = useState(initial?.object ? (initial.permission ?? "") : "");
  const [draft, setDraft] = useState(initial?.object ? (initial.permission ?? "") : "");
  const [explainActor, setExplainActor] = useState<string | null>(initial?.object ? (initial.actor ?? null) : null);

  const object = picked !== null && objects.has(picked) ? picked : null;
  const resource = object ? objects.get(object)! : "";
  const id = object ? object.slice(resource.length + 1) : "";
  const names = useMemo(() => checkableNames(resources, resource), [resources, resource]);

  const selectObject = (key: string | null) => {
    setPicked(key);
    setExplainActor(null);
    const p = key ? defaultName(checkableNames(resources, objects.get(key) ?? ""), resources, objects.get(key) ?? "") : "";
    setPermission(p);
    setDraft(p);
  };
  const commit = (v: string) => {
    const t = v.trim();
    setDraft(t);
    setPermission(t);
  };

  const request = (actor: string): CheckRequest => ({ scenario, object: { resource, id }, permission, actor: { id: actor } });

  // One Check per actor; a new selection (or edit) aborts the previous run, whose results are then dropped.
  const [statuses, setStatuses] = useState<Record<string, ActorStatus>>({});
  const requestRef = useRef(request);
  requestRef.current = request;
  useEffect(() => {
    setStatuses({});
    if (!object || permission === "" || actors.length === 0) return;
    const ac = new AbortController();
    void runLimited(
      actors,
      CONCURRENCY,
      async (a, signal) => (await check(requestRef.current(a), signal)).allowed === true,
      (a, s) => setStatuses((prev) => ({ ...prev, [a]: s })),
      ac.signal,
    );
    return () => ac.abort();
  }, [object, resource, id, permission, actors, scenario]);

  // Explain for the clicked actor, same staleness handling.
  const [explained, setExplained] = useState<ExplainState>({ state: "loading" });
  const explainTarget = object && permission !== "" ? explainActor : null;
  useEffect(() => {
    if (!explainTarget) return;
    const ac = new AbortController();
    setExplained({ state: "loading" });
    explain(requestRef.current(explainTarget), ac.signal).then(
      (r) => { if (!ac.signal.aborted) setExplained({ state: "done", allowed: r.allowed === true, graph: r.graph }); },
      (e: unknown) => {
        if (ac.signal.aborted) return;
        setExplained({ state: "error", message: shortError(e instanceof Error ? e.message : String(e)) });
      },
    );
    return () => ac.abort();
  }, [explainTarget, object, resource, id, permission, scenario]);

  const marks = useMemo(() => {
    const m: Record<string, ActorMark> = {};
    for (const [a, s] of Object.entries(statuses)) m[a] = s.state === "error" ? { state: "error", message: shortError(s.message) } : { state: s.state };
    return m;
  }, [statuses]);

  const summary = summarize(actors, statuses);
  const errorMessages = [...new Set(summary.errors.map((e) => shortError(e.message)))];
  const asking = object !== null && permission !== "";
  const what = `${permission} on ${object}`;

  return (
    <div data-testid="world-graph-panel" className="flex min-h-full flex-col">
      <div className="mb-1 flex min-h-7 flex-wrap items-center gap-2 text-xs" data-testid="check-bar">
        {object === null ? (
          <span className="text-muted-foreground">Click an object to see who can act on it.</span>
        ) : (
          <>
            <span className="text-muted-foreground">Who can</span>
            <input
              list={`names-${resource}`}
              value={draft}
              aria-label="Permission or relation"
              data-testid="permission-input"
              placeholder={names.length ? "permission or relation" : "type a name"}
              className="h-6 w-40 rounded border bg-background px-1.5 font-mono"
              onChange={(e) => {
                setDraft(e.target.value);
                // picking from the list commits at once; free text commits on Enter or blur
                if (names.includes(e.target.value)) setPermission(e.target.value);
              }}
              onBlur={() => commit(draft)}
              onKeyDown={(e) => { if (e.key === "Enter") commit(draft); }}
            />
            <datalist id={`names-${resource}`}>
              {names.map((n) => <option key={n} value={n} />)}
            </datalist>
            <span className="text-muted-foreground">on</span>
            <code className="rounded border-2 border-amber-600 px-1.5 py-0.5 dark:border-amber-400" data-testid="selected-object">{object}</code>
            <button type="button" className="rounded border px-1.5 py-0.5 text-muted-foreground hover:text-foreground" onClick={() => selectObject(null)}>
              Clear
            </button>
          </>
        )}
      </div>

      <RelationshipGraph
        fill
        marks={bound}
        text={text}
        dark={dark}
        selection={{
          object,
          actors: asking ? marks : {},
          explained: asking ? explainActor : null,
          onSelectObject: selectObject,
          onSelectActor: (a) => setExplainActor((cur) => (cur === a ? null : a)),
        }}
      />

      {asking && (
        <div className="mt-1 text-xs" data-testid="check-summary" aria-live="polite">
          <p>
            <b className="font-medium">{what}:</b>{" "}
            {summary.allowed.length > 0 ? <span className="text-green-700 dark:text-green-400">{summary.allowed.map(shortActor).join(", ")} allowed</span> : errorMessages.length === 0 && summary.pending === 0 ? <span className="text-muted-foreground">no actor allowed</span> : null}
            {summary.denied.length > 0 && <span className="text-muted-foreground">{summary.allowed.length > 0 ? "; " : ""}{summary.denied.map(shortActor).join(", ")} denied</span>}
            {summary.pending > 0 && <span className="ml-2 text-muted-foreground" data-testid="check-pending">checking {summary.pending} of {actors.length}...</span>}
            {actors.length === 0 && <span className="text-muted-foreground">this world has no actors</span>}
            {actors.length > 0 && !explainActor && <span className="ml-2 text-muted-foreground">Click an actor to see why.</span>}
          </p>
          {errorMessages.length > 0 && (
            <ul className="mt-0.5 text-destructive" data-testid="check-errors">
              {errorMessages.map((m) => {
                const who = summary.errors.filter((e) => shortError(e.message) === m).map((e) => shortActor(e.actor));
                return <li key={m}>{who.length === actors.length ? "" : `${who.join(", ")}: `}{m}</li>;
              })}
            </ul>
          )}
        </div>
      )}

      {asking && explainTarget && (
        <section className="mt-2 rounded-md border p-2" data-testid="explain-panel">
          <header className="mb-1 flex flex-wrap items-center gap-2 text-xs">
            <b className="font-medium">Why:</b>
            <span>
              <code>{shortActor(explainTarget)}</code> {explained.state === "done" ? (explained.allowed ? "can" : "cannot") : "-"} <code>{permission}</code> on <code>{object}</code>
            </span>
            {explained.state === "done" && (
              <span
                data-testid="explain-result"
                className={`rounded px-1.5 py-0.5 font-medium ${explained.allowed ? "bg-green-100 text-green-800 dark:bg-green-900/50 dark:text-green-300" : "bg-red-100 text-red-800 dark:bg-red-900/50 dark:text-red-300"}`}
              >
                {explained.allowed ? "Allowed" : "Denied"}
              </span>
            )}
            {explained.state === "loading" && <span className="text-muted-foreground">loading...</span>}
            <button type="button" aria-label="Close explanation" className="ml-auto text-muted-foreground hover:text-foreground" onClick={() => setExplainActor(null)}>
              <X className="h-3.5 w-3.5" />
            </button>
          </header>
          {explained.state === "error" && <p className="text-xs text-destructive" data-testid="explain-error">{explained.message}</p>}
          {explained.state === "done" && <ExplainGraphView graph={explained.graph} dark={dark} />}
        </section>
      )}
    </div>
  );
}
