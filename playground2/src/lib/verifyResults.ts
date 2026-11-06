// Pure: joins the requirements that were sent with the Verify response, and summarizes it.
import type { ContractResult, PropertyResult, PropertyType, VerifyResponse, VerifyStatus, WireBinding, WireRequirement, WireWorld } from "./types";

export interface PropertyView {
  type: PropertyType;
  description: string;
  formula: string;
  /** Undefined when the response has no result for it (not run). */
  status: VerifyStatus | undefined;
  message: string;
  world: WireWorld | undefined;
  bindings: { name: string; values: string[] }[];
}

export interface ContractView {
  name: string;
  description: string;
  accepted: boolean;
  status: VerifyStatus | undefined;
  properties: PropertyView[];
}

export interface RequirementView {
  name: string;
  statement: string;
  notes: string[];
  contracts: ContractView[];
}

/** An omitted status is the enum's zero value. */
const statusOf = (s: VerifyStatus | undefined): VerifyStatus => s ?? "UNVERIFIED";

/** The world of a result, when it has relationships to show. */
export function worldOf(r: PropertyResult | undefined): WireWorld | undefined {
  const w = r?.world;
  return w && (w.relationships ?? "").trim() !== "" ? w : undefined;
}

const normBindings = (bs: WireBinding[] | undefined) => (bs ?? []).map((b) => ({ name: b.name ?? "", values: b.values ?? [] }));

/**
 * Pairs each contract with its result (by requirement and contract name, in order, so repeated names still pair up)
 * and each property with the result at the same index. Every requirement is kept, results or not.
 */
export function joinResults(requirements: WireRequirement[], response: VerifyResponse | null): RequirementView[] {
  const queues = new Map<string, ContractResult[]>();
  for (const r of response?.results ?? []) {
    const k = JSON.stringify([r.requirement ?? "", r.contract ?? ""]);
    queues.set(k, [...(queues.get(k) ?? []), r]);
  }
  return requirements.map((req) => ({
    name: req.name,
    statement: req.statement,
    notes: req.notes,
    contracts: req.contracts.map((c) => {
      const res = queues.get(JSON.stringify([req.name, c.name]))?.shift();
      const properties = c.properties.map((p, i): PropertyView => {
        const pr = res?.properties?.[i];
        return {
          type: p.type,
          description: p.description,
          formula: p.formula,
          status: res && pr ? statusOf(pr.status) : undefined,
          message: pr?.message ?? "",
          world: worldOf(pr),
          bindings: normBindings(pr?.bindings),
        };
      });
      return {
        name: c.name,
        description: c.description,
        accepted: c.status === "ACCEPTED",
        status: res ? (c.properties.length === 0 ? "DRAFT" : statusOf(res.status)) : undefined,
        properties,
      };
    }),
  }));
}

export interface Summary {
  contracts: Record<VerifyStatus, number>;
  properties: Record<"SATISFIED" | "VIOLATED" | "UNVERIFIED", number>;
  /** Contracts or properties with no result in the response. */
  notRun: number;
}

export function summarize(reqs: RequirementView[]): Summary {
  const s: Summary = {
    contracts: { SATISFIED: 0, VIOLATED: 0, UNVERIFIED: 0, DRAFT: 0 },
    properties: { SATISFIED: 0, VIOLATED: 0, UNVERIFIED: 0 },
    notRun: 0,
  };
  for (const r of reqs)
    for (const c of r.contracts) {
      if (c.status === undefined) s.notRun++;
      else s.contracts[c.status]++;
      for (const p of c.properties) {
        if (p.status === undefined) s.notRun++;
        else if (p.status !== "DRAFT") s.properties[p.status]++;
      }
    }
  return s;
}

/** "2 contracts violated, 1 satisfied" style line. */
export function summaryLine(s: Summary): string {
  const part = (n: number, w: string) => (n > 0 ? [`${n} ${w}`] : []);
  const c = [...part(s.contracts.VIOLATED, "violated"), ...part(s.contracts.UNVERIFIED, "unverified"), ...part(s.contracts.SATISFIED, "satisfied"), ...part(s.contracts.DRAFT, "draft")];
  const total = Object.values(s.contracts).reduce((a, b) => a + b, 0);
  const p = [...part(s.properties.VIOLATED, "violated"), ...part(s.properties.UNVERIFIED, "unverified"), ...part(s.properties.SATISFIED, "satisfied")];
  const out = [`${total} ${total === 1 ? "contract" : "contracts"}: ${c.join(", ") || "none"}`];
  if (p.length > 0) out.push(`properties: ${p.join(", ")}`);
  if (s.notRun > 0) out.push(`${s.notRun} not in the response`);
  return out.join(" · ");
}

/** What a property's result means, in words. */
export function explainResult(p: PropertyView, scope: number | undefined): string {
  const upTo = scope !== undefined ? ` up to scope ${scope}` : "";
  if (p.status === undefined) return "";
  if (p.status === "UNVERIFIED") return p.message || "could not be verified";
  if (p.type === "ASSERT") return p.status === "SATISFIED" ? `holds in every world${upTo}` : "fails in the counterexample";
  return p.status === "SATISFIED" ? "holds in the example" : `holds in no world${upTo}`;
}
