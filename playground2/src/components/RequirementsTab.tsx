import { useMemo, useState } from "react";
import { Loader2, ShieldCheck, X } from "lucide-react";
import { boundNodes, prefillFromBindings } from "@/lib/bindings";
import { buildGraph } from "@/lib/graph";
import { parsePolicyResources } from "@/lib/policyResources";
import { parseRelationships } from "@/lib/relationships";
import { parseRequirements } from "@/lib/requirements";
import { explainResult, joinResults, summarize, type PropertyView } from "@/lib/verifyResults";
import { useWorkspace, type VerifyRun } from "@/stores/workspaceStore";
import { Button } from "@/components/ui/button";
import { ResizableHandle, ResizablePanel, ResizablePanelGroup } from "@/components/ui/resizable";
import { CodeEditor } from "./CodeEditor";
import { RequirementList, StatusBadge, SummaryCounts } from "./VerifyResults";
import { WorldView } from "./WorldView";
import { FileNotice } from "./WorkspaceStatus";

interface Viewing {
  run: VerifyRun;
  key: string;
  title: string;
  property: PropertyView;
}

function VerifyBar({ stale, run }: { stale: boolean; run: VerifyRun | null }) {
  const scope = useWorkspace((s) => s.verifyScope);
  const setScope = useWorkspace((s) => s.setVerifyScope);
  const verifying = useWorkspace((s) => s.verifying);
  const verify = useWorkspace((s) => s.verify);
  const cancel = useWorkspace((s) => s.cancelVerify);
  const fileBacked = useWorkspace((s) => s.mode === "files");
  const autoVerify = useWorkspace((s) => s.autoVerify);
  const setAutoVerify = useWorkspace((s) => s.setAutoVerify);
  return (
    <div className="flex flex-wrap items-center gap-3 border-b px-3 py-2">
      <label className="flex items-center gap-1.5 text-sm" title="Objects per resource, and actors, in the worlds checked">
        <span className="text-muted-foreground">Scope</span>
        <input
          type="number" min={1} value={Number.isFinite(scope) ? scope : ""} aria-label="Verify scope" data-testid="verify-scope"
          onChange={(e) => setScope(e.target.valueAsNumber)}
          className="h-8 w-16 rounded-md border border-input bg-background px-2 text-sm"
        />
      </label>
      <Button size="sm" onClick={() => void verify()} disabled={verifying} data-testid="verify">
        {verifying ? <Loader2 className="mr-1.5 h-4 w-4 animate-spin" /> : <ShieldCheck className="mr-1.5 h-4 w-4" />}
        Verify
      </Button>
      {verifying && <Button size="sm" variant="outline" onClick={cancel} data-testid="verify-cancel">Cancel</Button>}
      {verifying && <span className="text-xs text-muted-foreground" data-testid="verify-loading">Verifying (Alloy can take a while)...</span>}
      {fileBacked && (
        <label className="flex items-center gap-1.5 text-sm" title="Verify again whenever the workspace files change (runs Alloy: a few seconds)">
          <input type="checkbox" checked={autoVerify} onChange={(e) => setAutoVerify(e.target.checked)} data-testid="auto-verify" />
          <span className="text-muted-foreground">Auto-verify on change</span>
        </label>
      )}
      {!verifying && run === null && <span className="text-xs text-muted-foreground">Verify checks every contract against all worlds up to the scope.</span>}
      {stale && !verifying && (
        <span data-testid="verify-stale" className="rounded bg-amber-100 px-1.5 py-0.5 text-xs font-medium text-amber-900 dark:bg-amber-900/50 dark:text-amber-200">
          {fileBacked ? "Stale: the workspace files changed since this run" : "Stale: the Policy, constraints or requirements changed since this run"}
        </span>
      )}
    </div>
  );
}

function ResultsPane({ viewingKey, onView }: { viewingKey: string | null; onView: (v: Omit<Viewing, "run">) => void }) {
  const requirements = useWorkspace((s) => s.requirements);
  const policy = useWorkspace((s) => s.policy);
  const constraints = useWorkspace((s) => s.constraints);
  const run = useWorkspace((s) => s.verifyRun);
  const verifying = useWorkspace((s) => s.verifying);

  const stale = run !== null && (run.inputs.policy !== policy || run.inputs.constraints !== constraints || run.inputs.requirements !== requirements);
  // A run shows the requirements it was computed from; before a run, the current ones.
  const shown = useMemo(() => {
    const p = parseRequirements(run ? run.inputs.requirements : requirements);
    return p.ok ? joinResults(p.requirements, run?.response ?? null) : null;
  }, [run, requirements]);
  const summary = useMemo(() => (shown && run ? summarize(shown) : null), [shown, run]);

  return (
    <div className="flex h-full flex-col" data-testid="verify-results" data-stale={stale || undefined}>
      <VerifyBar stale={stale} run={run} />
      {summary && run && (
        <div className="flex flex-wrap items-center gap-x-3 gap-y-1 border-b px-3 py-1.5" data-testid="verify-header">
          <SummaryCounts s={summary} />
          <span className="text-xs text-muted-foreground" data-testid="verify-scope-note">Satisfied asserts hold only up to scope {run.scope}.</span>
        </div>
      )}
      <div className={`min-h-0 flex-1 overflow-auto ${stale || verifying ? "opacity-60" : ""}`}>
        {shown === null ? (
          <p className="p-3 text-sm text-muted-foreground" data-testid="requirements-unparsed">Fix the requirements YAML to see the requirements here.</p>
        ) : (
          <RequirementList
            requirements={shown}
            scope={run?.scope}
            viewingKey={viewingKey}
            onView={(key, ri, ci, pi) => {
              const r = shown[ri];
              onView({ key, title: `${r.name} / ${r.contracts[ci].name}`, property: r.contracts[ci].properties[pi] });
            }}
          />
        )}
      </div>
    </div>
  );
}

/** A property's world with its bindings marked on the graph, and Check/Explain prefilled from them. */
function PropertyWorld({ property, title, dark, onClose }: { property: PropertyView; title: string; dark: boolean; onClose: () => void }) {
  const policy = useWorkspace((s) => s.policy);
  const text = property.world?.relationships ?? "";
  const { marks, prefill } = useMemo(() => {
    const nodes = buildGraph(parseRelationships(text).relationships).nodes;
    const objects = new Map(nodes.filter((n) => n.kind === "object").map((n) => [n.key, n.resource!]));
    const actors = nodes.filter((n) => n.kind === "actor").map((n) => n.key);
    return {
      marks: boundNodes(property.bindings, nodes.map((n) => n.key)),
      prefill: prefillFromBindings(property.bindings, objects, actors, property.formula, parsePolicyResources(policy)),
    };
  }, [text, property, policy]);

  return (
    <div className="flex h-full flex-col p-2" data-testid="property-world">
      <header className="mb-1 flex flex-wrap items-center gap-2 text-sm">
        <StatusBadge status={property.status} />
        <b className="font-medium">{property.type === "ASSERT" ? "Counterexample" : "Example"}</b>
        <span className="text-muted-foreground">{title}: {property.description}</span>
        <button type="button" aria-label="Close world" data-testid="close-world" className="ml-auto text-muted-foreground hover:text-foreground" onClick={onClose}>
          <X className="h-4 w-4" />
        </button>
      </header>
      <p className="mb-1 text-xs text-muted-foreground">{explainResult(property, undefined)}</p>
      {property.bindings.length > 0 && (
        <ul className="mb-1 flex flex-wrap gap-1.5 text-xs" data-testid="bindings" aria-label="Bindings">
          {property.bindings.map((b) => (
            <li key={b.name} className="rounded border border-sky-600 px-1.5 py-0.5 font-mono dark:border-sky-400">
              <span className="font-semibold">{b.name}</span> = {b.values.join(", ")}
            </li>
          ))}
        </ul>
      )}
      <div className="min-h-0 flex-1">
        <WorldView text={text} dark={dark} initial={prefill} marks={marks} />
      </div>
    </div>
  );
}

export function RequirementsTab({ dark }: { dark: boolean }) {
  const requirements = useWorkspace((s) => s.requirements);
  const setRequirements = useWorkspace((s) => s.setRequirements);
  const readOnly = useWorkspace((s) => s.mode === "files");
  const run = useWorkspace((s) => s.verifyRun);
  const [viewing, setViewing] = useState<Viewing | null>(null);
  // A world belongs to its run: a new run drops it.
  const open = viewing && viewing.run === run ? viewing : null;

  const markers = useMemo(() => {
    const r = parseRequirements(requirements);
    return r.ok ? [] : r.errors;
  }, [requirements]);

  return (
    <ResizablePanelGroup direction="horizontal" className="p-2" autoSaveId="playground2.requirements.h">
      <ResizablePanel id="req-editor" order={1} defaultSize={36} minSize={15}>
        <div className="flex h-full flex-col pr-1">
          <FileNotice file="contracts.yaml" />
          <div className="min-h-0 flex-1"><CodeEditor value={requirements} onChange={setRequirements} readOnly={readOnly} dark={dark} markers={markers} testId="requirements-editor" /></div>
        </div>
      </ResizablePanel>
      <ResizableHandle withHandle />
      <ResizablePanel id="req-right" order={2} defaultSize={64} minSize={25}>
        <ResizablePanelGroup direction="vertical" autoSaveId="playground2.requirements.v">
          <ResizablePanel id="req-results" order={1} defaultSize={45} minSize={15}>
            <ResultsPane viewingKey={open?.key ?? null} onView={(v) => run && setViewing({ ...v, run })} />
          </ResizablePanel>
          {open && (
            <>
              <ResizableHandle withHandle />
              <ResizablePanel id="req-world" order={2} defaultSize={55} minSize={20}>
                <PropertyWorld key={open.key} property={open.property} title={open.title} dark={dark} onClose={() => setViewing(null)} />
              </ResizablePanel>
            </>
          )}
        </ResizablePanelGroup>
      </ResizablePanel>
    </ResizablePanelGroup>
  );
}
