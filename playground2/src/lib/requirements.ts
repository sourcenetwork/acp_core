// Pure: parses the requirements editor text (see HANDOFF.md, `contracts.yaml`) into the wire requirements.
import { isMap, isScalar, isSeq, LineCounter, parseDocument, type Node } from "yaml";
import type { ContractStatus, PropertyType, WireContract, WireProperty, WireRequirement } from "./types";

export interface RequirementError {
  message: string;
  /** 1-based line in the requirements text, when known. */
  line?: number;
}

export type RequirementsParseResult =
  | { ok: true; requirements: WireRequirement[] }
  | { ok: false; errors: RequirementError[] };

const REQUIREMENT_KEYS = ["name", "statement", "notes", "contracts"];
const CONTRACT_KEYS = ["name", "description", "status", "properties"];
const PROPERTY_KEYS = ["type", "description", "formula"];

const isNull = (v: unknown) => v === undefined || v === null || (isScalar(v) && v.value === null);

/** Parses the requirements text. Empty text means no requirements. Errors carry the line of the offending node. */
export function parseRequirements(text: string): RequirementsParseResult {
  if (text.trim() === "") return { ok: true, requirements: [] };

  const lc = new LineCounter();
  const doc = parseDocument(text, { lineCounter: lc });
  if (doc.errors.length > 0) {
    return {
      ok: false,
      errors: doc.errors.map((e) => ({ message: e.message.split("\n")[0], line: e.linePos?.[0].line })),
    };
  }
  const lineOf = (n: unknown): number | undefined => {
    const off = (n as Node | null | undefined)?.range?.[0];
    return off === undefined ? undefined : lc.linePos(off).line;
  };

  const errors: RequirementError[] = [];
  const err = (message: string, node: unknown) => errors.push({ message, line: lineOf(node) ?? 1 });

  const root = doc.contents;
  if (root === null) return { ok: true, requirements: [] };
  if (!isMap(root)) {
    err(`expected a mapping with a "requirements" key`, root);
    return { ok: false, errors };
  }
  for (const pair of root.items) {
    const k = String((pair.key as { value?: unknown } | null)?.value ?? "");
    if (k !== "requirements") err(`unknown top-level key "${k}"`, pair.key);
  }
  const list = root.get("requirements", true);
  if (isNull(list)) return errors.length > 0 ? { ok: false, errors } : { ok: true, requirements: [] };
  if (!isSeq(list)) {
    err(`"requirements" must be a list`, list);
    return { ok: false, errors };
  }

  /** Reads a mapping's fields: reports non-mappings and unknown keys, returns the value nodes by key. */
  const fields = (node: unknown, what: string, allowed: string[]): Map<string, unknown> | null => {
    if (!isMap(node)) {
      err(`${what}: expected a mapping with ${allowed.join(", ")}`, node);
      return null;
    }
    const out = new Map<string, unknown>();
    for (const pair of node.items) {
      const k = String((pair.key as { value?: unknown } | null)?.value ?? "");
      if (!allowed.includes(k)) err(`${what}: unknown key "${k}" (expected ${allowed.join(", ")})`, pair.key);
      else out.set(k, pair.value);
    }
    return out;
  };
  const str = (f: Map<string, unknown>, key: string, what: string, required: boolean, node: unknown): string => {
    const v = f.get(key);
    if (isNull(v)) {
      if (required) err(`${what}: missing "${key}"`, node);
      return "";
    }
    if (!isScalar(v) || typeof v.value !== "string") {
      err(`${what}: "${key}" must be a string`, v);
      return "";
    }
    if (required && v.value.trim() === "") err(`${what}: "${key}" must not be empty`, v);
    return v.value;
  };
  const seq = (f: Map<string, unknown>, key: string, what: string): unknown[] => {
    const v = f.get(key);
    if (isNull(v)) return [];
    if (!isSeq(v)) {
      err(`${what}: "${key}" must be a list`, v);
      return [];
    }
    return v.items;
  };

  const requirements: WireRequirement[] = [];
  list.items.forEach((rNode, ri) => {
    const rWhat = `requirement #${ri + 1}`;
    const rf = fields(rNode, rWhat, REQUIREMENT_KEYS);
    if (!rf) return;
    const name = str(rf, "name", rWhat, true, rNode);
    const rLabel = name !== "" ? `requirement "${name}"` : rWhat;
    const statement = str(rf, "statement", rLabel, true, rNode);

    const notes: string[] = [];
    const notesNode = rf.get("notes");
    if (isScalar(notesNode) && typeof notesNode.value === "string") notes.push(notesNode.value); // a single note
    else
      seq(rf, "notes", rLabel).forEach((n) => {
        if (isScalar(n) && typeof n.value === "string") notes.push(n.value);
        else err(`${rLabel}: a note must be a string`, n);
      });

    const contracts: WireContract[] = [];
    seq(rf, "contracts", rLabel).forEach((cNode, ci) => {
      const cWhat = `${rLabel}, contract #${ci + 1}`;
      const cf = fields(cNode, cWhat, CONTRACT_KEYS);
      if (!cf) return;
      const cName = str(cf, "name", cWhat, true, cNode);
      const cLabel = cName !== "" ? `${rLabel}, contract "${cName}"` : cWhat;
      const description = str(cf, "description", cLabel, true, cNode);

      let status: ContractStatus = "PROPOSED";
      const sText = str(cf, "status", cLabel, false, cNode);
      if (sText !== "") {
        const s = sText.trim().toUpperCase();
        if (s === "PROPOSED" || s === "ACCEPTED") status = s;
        else err(`${cLabel}: status must be proposed or accepted, found "${sText}"`, cf.get("status"));
      }

      const properties: WireProperty[] = [];
      seq(cf, "properties", cLabel).forEach((pNode, pi) => {
        const pWhat = `${cLabel}, property #${pi + 1}`;
        const pf = fields(pNode, pWhat, PROPERTY_KEYS);
        if (!pf) return;
        const tText = str(pf, "type", pWhat, true, pNode);
        const t = tText.trim().toUpperCase();
        let type: PropertyType = "ASSERT";
        if (t === "ASSERT" || t === "PRED") type = t;
        else if (tText !== "") err(`${pWhat}: type must be assert or pred, found "${tText}"`, pf.get("type"));
        properties.push({
          type,
          description: str(pf, "description", pWhat, true, pNode),
          formula: str(pf, "formula", pWhat, true, pNode),
        });
      });
      contracts.push({ name: cName, description, status, properties });
    });
    requirements.push({ name, statement, notes, contracts });
  });

  return errors.length > 0 ? { ok: false, errors } : { ok: true, requirements };
}
