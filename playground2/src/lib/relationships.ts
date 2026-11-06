// Pure parser for a world's relationships text: one `resource:id#relation@subject` tuple per line.
// The subject is an actor (`did:user:alice`), an object (`folder:folder1`) or a userset (`team:t1#admin`).

export interface ObjectRef {
  resource: string;
  id: string;
}

export type Subject =
  | { kind: "actor"; did: string }
  | { kind: "object"; object: ObjectRef }
  | { kind: "userset"; object: ObjectRef; relation: string };

export interface Relationship {
  object: ObjectRef;
  relation: string;
  subject: Subject;
}

export interface ParseIssue {
  line: number; // 1-based
  text: string;
  reason: string;
}

export interface ParsedRelationships {
  relationships: Relationship[];
  issues: ParseIssue[];
}

const OBJECT = /^([^\s:#@]+):([^\s#@]+)$/;
const NAME = /^[^\s:#@]+$/;

function parseObject(s: string): ObjectRef | null {
  const m = OBJECT.exec(s);
  return m ? { resource: m[1], id: m[2] } : null;
}

function parseSubject(s: string): Subject | string {
  if (s.startsWith("did:")) {
    return s.length > 4 && !/[\s#@]/.test(s) ? { kind: "actor", did: s } : "invalid actor subject";
  }
  const hash = s.indexOf("#");
  const objText = hash < 0 ? s : s.slice(0, hash);
  const object = parseObject(objText);
  if (!object) return "invalid subject";
  if (hash < 0) return { kind: "object", object };
  const relation = s.slice(hash + 1);
  if (!NAME.test(relation)) return "invalid subject relation";
  return { kind: "userset", object, relation };
}

export function parseRelationship(line: string): Relationship | string {
  const text = line.trim();
  const at = text.indexOf("@");
  const hash = text.indexOf("#");
  if (at < 0) return "missing '@subject'";
  if (hash < 0 || hash > at) return "missing '#relation'";
  const object = parseObject(text.slice(0, hash));
  if (!object) return "invalid object, expected resource:id";
  const relation = text.slice(hash + 1, at);
  if (!NAME.test(relation)) return "invalid relation";
  const subject = parseSubject(text.slice(at + 1));
  if (typeof subject === "string") return subject;
  return { object, relation, subject };
}

export function parseRelationships(text: string): ParsedRelationships {
  const relationships: Relationship[] = [];
  const issues: ParseIssue[] = [];
  text.split(/\r?\n/).forEach((raw, i) => {
    if (raw.trim() === "") return;
    const r = parseRelationship(raw);
    if (typeof r === "string") issues.push({ line: i + 1, text: raw.trim(), reason: r });
    else relationships.push(r);
  });
  return { relationships, issues };
}

export const objectKey = (o: ObjectRef) => `${o.resource}:${o.id}`;
