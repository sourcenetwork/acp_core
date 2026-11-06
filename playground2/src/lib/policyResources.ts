// Pure: the names a Check can ask about for each resource of a Policy (YAML), parsed client-side.
import { parse } from "yaml";

export interface ResourceNames {
  permissions: string[];
  relations: string[];
}

export type PolicyResources = Record<string, ResourceNames>;

/** Every resource implicitly has an `owner` relation, set when an object is registered (and accepted by Check). */
export const IMPLICIT_RELATION = "owner";

const names = (v: unknown): string[] => {
  const out: string[] = [];
  if (Array.isArray(v)) {
    for (const e of v) {
      const n = typeof e === "string" ? e : (e as { name?: unknown } | null)?.name;
      if (typeof n === "string" && n !== "") out.push(n);
    }
  } else if (v && typeof v === "object") {
    out.push(...Object.keys(v));
  }
  return out;
};

/** Resources of the Policy with their permission and relation names. Unparseable policies yield `{}`. */
export function parsePolicyResources(policyYaml: string): PolicyResources {
  let doc: unknown;
  try {
    doc = parse(policyYaml);
  } catch {
    return {};
  }
  const raw = (doc as { resources?: unknown } | null)?.resources;
  const entries: [string, unknown][] = Array.isArray(raw)
    ? raw.map((r) => [String((r as { name?: unknown } | null)?.name ?? ""), r])
    : raw && typeof raw === "object"
      ? Object.entries(raw)
      : [];
  const out: PolicyResources = {};
  for (const [name, r] of entries) {
    if (name === "" || !r || typeof r !== "object") continue;
    const res = r as { permissions?: unknown; relations?: unknown };
    out[name] = { permissions: names(res.permissions), relations: names(res.relations) };
  }
  return out;
}

/** Selectable names for a resource: permissions, then relations, then the implicit owner. Empty if unknown. */
export function checkableNames(resources: PolicyResources, resource: string): string[] {
  const r = resources[resource];
  if (!r) return [];
  const all = [...r.permissions, ...r.relations];
  if (!all.includes(IMPLICIT_RELATION)) all.push(IMPLICIT_RELATION);
  return [...new Set(all)];
}

/** Default permission to pre-select: `view` if present, else the first permission, else the first name. */
export function defaultName(list: string[], resources: PolicyResources, resource: string): string {
  if (list.includes("view")) return "view";
  return resources[resource]?.permissions[0] ?? list[0] ?? "";
}
