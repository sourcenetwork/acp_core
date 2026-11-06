import { parseDocument } from "yaml";
import type { SystemConstraint } from "./types";

export interface ConstraintError {
  message: string;
  /** 1-based line in the constraints text, when known. */
  line?: number;
}

export type ParseResult =
  | { ok: true; constraints: SystemConstraint[] }
  | { ok: false; errors: ConstraintError[] };

const KINDS = ["disjoint", "at_most_one", "acyclic"];

/** Splits "resource.relation" on the first dot. */
export function parseRef(ref: unknown): { resource: string; relation: string } | string {
  if (typeof ref !== "string") return `reference must be a string like "resource.relation"`;
  const i = ref.indexOf(".");
  if (i <= 0 || i === ref.length - 1) return `invalid reference "${ref}": expected "resource.relation"`;
  return { resource: ref.slice(0, i), relation: ref.slice(i + 1) };
}

function lineOf(node: unknown, text: string): number | undefined {
  const range = (node as { range?: [number, number, number] } | null)?.range;
  if (!range) return undefined;
  let line = 1;
  for (let i = 0; i < range[0] && i < text.length; i++) if (text[i] === "\n") line++;
  return line;
}

/** Parses the constraints editor text (see docs/playground2-handoff.md) into wire constraints. */
export function parseConstraints(text: string): ParseResult {
  if (text.trim() === "") return { ok: true, constraints: [] };

  const doc = parseDocument(text);
  if (doc.errors.length > 0) {
    return {
      ok: false,
      errors: doc.errors.map((e) => ({ message: e.message.split("\n")[0], line: e.linePos?.[0].line })),
    };
  }

  const root = doc.contents;
  if (root === null) return { ok: true, constraints: [] };
  const top = doc.toJS() as unknown;
  if (typeof top !== "object" || top === null || Array.isArray(top)) {
    return { ok: false, errors: [{ message: `expected a mapping with a "constraints" (or "assumptions") key`, line: 1 }] };
  }
  const keys = Object.keys(top);
  // "assumptions" is the key of a workspace's constraints.yaml (and the old name): both are accepted.
  const extra = keys.filter((k) => k !== "constraints" && k !== "assumptions");
  if (extra.length > 0) {
    return { ok: false, errors: [{ message: `unknown top-level key "${extra[0]}" (expected "constraints" or "assumptions")`, line: 1 }] };
  }
  if (keys.length > 1) {
    return { ok: false, errors: [{ message: `use either "constraints" or "assumptions", not both`, line: 1 }] };
  }
  const key = keys[0];
  const list = (top as Record<string, unknown>)[key];
  if (list === null || list === undefined) return { ok: true, constraints: [] };
  if (!Array.isArray(list)) {
    return { ok: false, errors: [{ message: `"${key}" must be a list`, line: 1 }] };
  }

  const entryNodes = (doc.get(key, true) as { items?: unknown[] } | undefined)?.items ?? [];
  const errors: ConstraintError[] = [];
  const out: SystemConstraint[] = [];

  list.forEach((entry: unknown, idx) => {
    const line = lineOf(entryNodes[idx], text);
    const fail = (message: string) => errors.push({ message: `constraint #${idx + 1}: ${message}`, line });

    if (typeof entry !== "object" || entry === null || Array.isArray(entry)) {
      return fail(`expected a single-key mapping such as "acyclic: folder.parent"`);
    }
    const ks = Object.keys(entry);
    if (ks.length !== 1) return fail(`must have exactly one key, found ${ks.length}`);
    const kind = ks[0];
    const val = (entry as Record<string, unknown>)[kind];
    if (!KINDS.includes(kind)) return fail(`unknown constraint "${kind}" (expected ${KINDS.join(", ")})`);

    if (kind === "disjoint") {
      if (!Array.isArray(val) || val.length < 2) return fail(`disjoint expects a list of at least two references`);
      const refs = val.map(parseRef);
      const bad = refs.find((r) => typeof r === "string");
      if (typeof bad === "string") return fail(bad);
      const parsed = refs as { resource: string; relation: string }[];
      const resources = new Set(parsed.map((r) => r.resource));
      if (resources.size !== 1) {
        return fail(`disjoint references must all name the same resource, found ${[...resources].join(", ")}`);
      }
      out.push({ disjoint: { resource: parsed[0].resource, relations: parsed.map((r) => r.relation) } });
      return;
    }

    const ref = parseRef(val);
    if (typeof ref === "string") return fail(ref);
    out.push(kind === "at_most_one" ? { atMostOne: ref } : { acyclic: ref });
  });

  return errors.length > 0 ? { ok: false, errors } : { ok: true, constraints: out };
}
