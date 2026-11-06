// Pure: turns a counterexample's bindings into marks on the graph and a prefilled Check/Explain question.
import { checkableNames, defaultName, type PolicyResources } from "./policyResources";

export interface BindingInput {
  name: string;
  values: string[];
}

/** Node key -> the variables bound to it. Only values that equal a node key are marked (tuples like `a->b` are not). */
export function boundNodes(bindings: BindingInput[], nodeKeys: Iterable<string>): Record<string, string[]> {
  const keys = new Set(nodeKeys);
  const out: Record<string, string[]> = {};
  for (const b of bindings)
    for (const v of b.values)
      if (keys.has(v) && !(out[v] ?? []).includes(b.name)) out[v] = [...(out[v] ?? []), b.name];
  return out;
}

export interface Prefill {
  object?: string;
  permission?: string;
  actor?: string;
}

const uniq = (xs: string[]) => [...new Set(xs)];

/**
 * Best effort: the object and actor when the bindings name exactly one of each (variables may repeat a value),
 * and, for the object, the first permission of its resource that the formula mentions as `resource_name`
 * (the compiled model's naming), else the default one.
 */
export function prefillFromBindings(
  bindings: BindingInput[],
  objects: Map<string, string>,
  actors: Iterable<string>,
  formula: string,
  resources: PolicyResources,
): Prefill {
  const actorSet = new Set(actors);
  const values = bindings.flatMap((b) => b.values);
  const objs = uniq(values.filter((v) => objects.has(v)));
  const acts = uniq(values.filter((v) => actorSet.has(v)));
  const out: Prefill = {};
  if (acts.length === 1) out.actor = acts[0];
  if (objs.length === 1) {
    out.object = objs[0];
    const resource = objects.get(objs[0])!;
    const names = checkableNames(resources, resource);
    const perms = resources[resource]?.permissions ?? [];
    const mentioned = names
      .map((n) => ({ n, at: formula.search(new RegExp(`\\b${resource}_${n}\\b`)) }))
      .filter((m) => m.at >= 0)
      .sort((a, b) => a.at - b.at);
    out.permission = (mentioned.find((m) => perms.includes(m.n)) ?? mentioned[0])?.n ?? defaultName(names, resources, resource);
  }
  return out;
}
