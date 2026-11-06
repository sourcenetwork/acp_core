import { Eye } from "lucide-react";
import { explainResult, type ContractView, type PropertyView, type RequirementView, type Summary } from "@/lib/verifyResults";
import type { VerifyStatus } from "@/lib/types";

const BADGE: Record<VerifyStatus | "NOT_RUN", string> = {
  SATISFIED: "bg-green-100 text-green-800 dark:bg-green-900/50 dark:text-green-300",
  VIOLATED: "bg-red-100 text-red-800 dark:bg-red-900/50 dark:text-red-300",
  UNVERIFIED: "bg-amber-100 text-amber-900 dark:bg-amber-900/50 dark:text-amber-200",
  DRAFT: "bg-muted text-muted-foreground",
  NOT_RUN: "border text-muted-foreground",
};

export function StatusBadge({ status, testId }: { status: VerifyStatus | undefined; testId?: string }) {
  const key = status ?? "NOT_RUN";
  return (
    <span data-testid={testId} data-status={key} className={`inline-block rounded px-1.5 py-0.5 text-[11px] font-semibold uppercase tracking-wide ${BADGE[key]}`}>
      {status ?? "not run"}
    </span>
  );
}

const ROW_TINT: Record<string, string> = {
  VIOLATED: "border-l-red-500 bg-red-50 dark:bg-red-950/40",
  UNVERIFIED: "border-l-amber-500 bg-amber-50/60 dark:bg-amber-950/20",
};

function viewLabel(p: PropertyView): string {
  if (p.type === "ASSERT" && p.status === "VIOLATED") return "View counterexample";
  if (p.type === "PRED" && p.status === "SATISFIED") return "View example";
  return "View world";
}

function PropertyRow({ p, scope, viewing, onView }: { p: PropertyView; scope: number | undefined; viewing: boolean; onView: () => void }) {
  const tint = (p.status && ROW_TINT[p.status]) ?? "border-l-transparent";
  const result = explainResult(p, scope);
  return (
    <li data-testid="property" data-status={p.status ?? "NOT_RUN"} data-type={p.type} className={`rounded-sm border-l-4 px-2 py-1.5 ${tint} ${viewing ? "ring-2 ring-sky-500" : ""}`}>
      <div className="flex flex-wrap items-center gap-2">
        <StatusBadge status={p.status} testId="property-status" />
        <span
          title={p.type === "ASSERT" ? "assert: must hold in every world" : "pred: must hold in some world"}
          className="rounded border px-1 text-[11px] font-medium uppercase tracking-wide text-muted-foreground"
          data-testid="property-type"
        >
          {p.type === "ASSERT" ? "assert" : "pred"}
        </span>
        <span className="min-w-0 flex-1 text-sm">{p.description}</span>
        {p.world && (
          <button
            type="button"
            data-testid="view-world"
            onClick={onView}
            className="inline-flex items-center gap-1 rounded border bg-background px-1.5 py-0.5 text-xs hover:bg-accent"
          >
            <Eye className="h-3 w-3" />
            {viewLabel(p)}
          </button>
        )}
      </div>
      {result && (
        <p data-testid="property-result" className={`mt-0.5 text-xs ${p.status === "UNVERIFIED" ? "font-mono text-amber-800 dark:text-amber-300" : "text-muted-foreground"}`}>
          {result}
        </p>
      )}
      {p.bindings.length > 0 && p.world && (
        <p className="mt-0.5 text-xs text-muted-foreground">
          {p.bindings.map((b) => (
            <span key={b.name} className="mr-2 font-mono">{b.name} = {b.values.join(", ")}</span>
          ))}
        </p>
      )}
      <details className="mt-0.5 text-xs">
        <summary className="cursor-pointer text-muted-foreground">formula</summary>
        <pre className="mt-0.5 overflow-x-auto rounded bg-muted p-1.5 font-mono">{p.formula.trimEnd()}</pre>
      </details>
    </li>
  );
}

const CONTRACT_EDGE: Record<string, string> = {
  SATISFIED: "border-l-green-500",
  VIOLATED: "border-l-red-500",
  UNVERIFIED: "border-l-amber-500",
  DRAFT: "border-l-border",
};

function ContractCard({ c, scope, viewingKey, keyOf, onView }: { c: ContractView; scope: number | undefined; viewingKey: string | null; keyOf: (pi: number) => string; onView: (pi: number) => void }) {
  const violated = c.properties.filter((p) => p.status === "VIOLATED").length;
  return (
    <section data-testid="contract" data-status={c.status ?? "NOT_RUN"} className={`rounded border border-l-4 p-2 ${CONTRACT_EDGE[c.status ?? ""] ?? "border-l-border"}`}>
      <header className="flex flex-wrap items-center gap-2">
        <StatusBadge status={c.status} testId="contract-status" />
        <b className="text-sm font-medium">{c.name}</b>
        <span className="rounded border px-1 text-[11px] uppercase tracking-wide text-muted-foreground" title="the user's decision, not a verification result">
          {c.accepted ? "accepted" : "proposed"}
        </span>
        {c.properties.length > 0 && c.status !== undefined && (
          <span className="text-xs text-muted-foreground">
            {violated > 0 ? `${violated} of ${c.properties.length} properties violated` : `${c.properties.length} ${c.properties.length === 1 ? "property" : "properties"}`}
          </span>
        )}
      </header>
      <p className="mt-0.5 text-sm text-muted-foreground">{c.description}</p>
      {c.properties.length === 0 ? (
        <p className="mt-1 text-xs italic text-muted-foreground">No properties yet: a draft, nothing is verified.</p>
      ) : (
        <ul className="mt-1 flex flex-col gap-1">
          {c.properties.map((p, pi) => (
            <PropertyRow key={pi} p={p} scope={scope} viewing={viewingKey === keyOf(pi)} onView={() => onView(pi)} />
          ))}
        </ul>
      )}
    </section>
  );
}

export function RequirementList({ requirements, scope, viewingKey, onView }: { requirements: RequirementView[]; scope: number | undefined; viewingKey: string | null; onView: (key: string, ri: number, ci: number, pi: number) => void }) {
  if (requirements.length === 0) {
    return <p className="p-3 text-sm text-muted-foreground" data-testid="no-requirements">No requirements.</p>;
  }
  return (
    <div className="flex flex-col gap-3 p-3" data-testid="requirement-list">
      {requirements.map((r, ri) => (
        <article key={ri} data-testid="requirement" className="rounded-md border p-2">
          <h3 className="text-sm font-semibold">{r.name}</h3>
          <blockquote className="mt-0.5 border-l-2 pl-2 text-sm italic">{r.statement}</blockquote>
          {r.notes.length > 0 && (
            <ul className="mt-1 list-disc pl-5 text-xs text-muted-foreground" aria-label="Notes">
              {r.notes.map((n, i) => <li key={i}>{n}</li>)}
            </ul>
          )}
          <div className="mt-2 flex flex-col gap-2">
            {r.contracts.length === 0 && <p className="text-xs italic text-muted-foreground">No contracts yet.</p>}
            {r.contracts.map((c, ci) => (
              <ContractCard
                key={ci}
                c={c}
                scope={scope}
                viewingKey={viewingKey}
                keyOf={(pi) => `${ri}.${ci}.${pi}`}
                onView={(pi) => onView(`${ri}.${ci}.${pi}`, ri, ci, pi)}
              />
            ))}
          </div>
        </article>
      ))}
    </div>
  );
}

export function SummaryCounts({ s }: { s: Summary }) {
  const chip = (n: number, status: VerifyStatus, label: string) =>
    n > 0 ? <span key={status} className={`rounded px-1.5 py-0.5 text-xs font-medium ${BADGE[status]}`}>{n} {label}</span> : null;
  const total = s.contracts.SATISFIED + s.contracts.VIOLATED + s.contracts.UNVERIFIED + s.contracts.DRAFT;
  return (
    <span className="flex flex-wrap items-center gap-1.5" data-testid="verify-summary">
      <span className="text-xs text-muted-foreground">{total} {total === 1 ? "contract" : "contracts"}:</span>
      {chip(s.contracts.VIOLATED, "VIOLATED", "violated")}
      {chip(s.contracts.UNVERIFIED, "UNVERIFIED", "unverified")}
      {chip(s.contracts.SATISFIED, "SATISFIED", "satisfied")}
      {chip(s.contracts.DRAFT, "DRAFT", "draft")}
      <span className="text-xs text-muted-foreground">
        ({s.properties.VIOLATED} violated, {s.properties.UNVERIFIED} unverified, {s.properties.SATISFIED} satisfied properties)
      </span>
    </span>
  );
}
