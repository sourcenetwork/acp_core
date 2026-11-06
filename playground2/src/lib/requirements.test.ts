import { describe, expect, it } from "vitest";
import { parseRequirements } from "./requirements";
import { SAMPLE_REQUIREMENTS } from "./sample";

const errs = (t: string) => {
  const r = parseRequirements(t);
  if (r.ok) throw new Error("expected failure");
  return r.errors;
};

describe("parseRequirements", () => {
  it("maps the HANDOFF example to the wire shape", () => {
    const r = parseRequirements(`requirements:
- name: read-access
  statement: "actors can read if they're readers or writers, but not banned"
  notes:
  - "Assumed writers means the doc's writers."
  contracts:
  - name: unbanned-readers-can-read
    description: Readers can read the doc, unless they're banned.
    status: accepted
    properties:
    - type: assert
      description: every unbanned reader of a doc can read it
      formula: |
        all d: Doc | d.doc_reader - d.doc_banned in d.doc_read
    - type: PRED
      description: unbanned readers can exist
      formula: 'some d: Doc | some d.doc_reader - d.doc_banned'
  - name: banned-never-read
    description: No banned actor can read.
`);
    expect(r).toEqual({
      ok: true,
      requirements: [
        {
          name: "read-access",
          statement: "actors can read if they're readers or writers, but not banned",
          notes: ["Assumed writers means the doc's writers."],
          contracts: [
            {
              name: "unbanned-readers-can-read",
              description: "Readers can read the doc, unless they're banned.",
              status: "ACCEPTED",
              properties: [
                { type: "ASSERT", description: "every unbanned reader of a doc can read it", formula: "all d: Doc | d.doc_reader - d.doc_banned in d.doc_read\n" },
                { type: "PRED", description: "unbanned readers can exist", formula: "some d: Doc | some d.doc_reader - d.doc_banned" },
              ],
            },
            { name: "banned-never-read", description: "No banned actor can read.", status: "PROPOSED", properties: [] },
          ],
        },
      ],
    });
  });

  it("parses the sample", () => {
    const r = parseRequirements(SAMPLE_REQUIREMENTS);
    expect(r.ok).toBe(true);
    if (r.ok) expect(r.requirements.map((x) => x.contracts.length)).toEqual([4, 1]);
  });

  it("treats empty text and an empty list as no requirements", () => {
    for (const t of ["", "  \n", "requirements:\n", "requirements: []", "# only a comment\n"])
      expect(parseRequirements(t)).toEqual({ ok: true, requirements: [] });
  });

  it("tolerates missing optional fields", () => {
    const r = parseRequirements("requirements:\n- name: a\n  statement: s\n");
    expect(r).toEqual({ ok: true, requirements: [{ name: "a", statement: "s", notes: [], contracts: [] }] });
    const single = parseRequirements("requirements:\n- name: a\n  statement: s\n  notes: just one\n");
    expect(single.ok && single.requirements[0].notes).toEqual(["just one"]);
  });

  it("reports yaml syntax errors with a line", () => {
    const e = errs("requirements:\n- name: a\n  statement: [oops\n");
    expect(e[0].line).toBeGreaterThan(0);
  });

  it("reports missing required fields with the line of the entry", () => {
    const e = errs(`requirements:
- name: a
  statement: s
  contracts:
  - name: c
    properties:
    - type: assert
      formula: 'all d: Doc | some d'
`);
    expect(e.map((x) => x.message)).toEqual([
      expect.stringContaining(`missing "description"`),
      expect.stringContaining(`missing "description"`),
    ]);
    expect(e.map((x) => x.line)).toEqual([5, 7]);
  });

  it("rejects bad enum values, unknown keys and wrong shapes", () => {
    expect(errs("requirements:\n- name: a\n  statement: s\n  contracts:\n  - name: c\n    description: d\n    status: maybe\n")[0]).toMatchObject({
      message: expect.stringContaining("proposed or accepted"),
      line: 7,
    });
    expect(errs("requirements:\n- name: a\n  statement: s\n  contracts:\n  - name: c\n    description: d\n    properties:\n    - type: check\n      description: x\n      formula: y\n")[0].message).toContain("assert or pred");
    expect(errs("requirements:\n- name: a\n  statement: s\n  contract: []\n")[0]).toMatchObject({ message: expect.stringContaining(`unknown key "contract"`), line: 4 });
    expect(errs("reqs: []")[0].message).toContain(`unknown top-level key "reqs"`);
    expect(errs("requirements: nope")[0].message).toContain("must be a list");
    expect(errs("- a\n")[0].message).toContain("mapping");
    expect(errs("requirements:\n- just text\n")[0].message).toContain("expected a mapping");
  });
});
