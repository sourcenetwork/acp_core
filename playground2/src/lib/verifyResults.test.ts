import { describe, expect, it } from "vitest";
import { joinResults, summarize, summaryLine } from "./verifyResults";
import type { VerifyResponse, WireRequirement } from "./types";

const prop = (type: "ASSERT" | "PRED") => ({ type, description: "d", formula: "f" });
const reqs: WireRequirement[] = [
  {
    name: "r1", statement: "s", notes: ["n"],
    contracts: [
      { name: "c1", description: "", status: "PROPOSED", properties: [prop("ASSERT"), prop("PRED")] },
      { name: "c2", description: "", status: "ACCEPTED", properties: [] },
    ],
  },
  { name: "r2", statement: "s2", notes: [], contracts: [{ name: "c1", description: "", status: "PROPOSED", properties: [prop("ASSERT")] }] },
];

describe("joinResults", () => {
  it("pairs results with every requirement, and defaults omitted statuses to UNVERIFIED", () => {
    const res: VerifyResponse = {
      scope: 3,
      results: [
        { requirement: "r1", contract: "c1", status: "VIOLATED", properties: [{ status: "VIOLATED", world: { relationships: "doc:a#owner@did:x" }, bindings: [{ name: "d", values: ["doc:a"] }] }, {}] },
        { requirement: "r1", contract: "c2", status: "DRAFT" },
        { requirement: "r2", contract: "c1", properties: [{ message: "syntax error" }] },
      ],
    };
    const v = joinResults(reqs, res);
    expect(v.map((r) => r.name)).toEqual(["r1", "r2"]);
    expect(v[0].contracts.map((c) => c.status)).toEqual(["VIOLATED", "DRAFT"]);
    expect(v[0].contracts[0].properties.map((p) => p.status)).toEqual(["VIOLATED", "UNVERIFIED"]);
    expect(v[0].contracts[0].properties[0].bindings).toEqual([{ name: "d", values: ["doc:a"] }]);
    expect(v[0].contracts[0].properties[1].world).toBeUndefined();
    expect(v[1].contracts[0].status).toBe("UNVERIFIED");
    expect(v[1].contracts[0].properties[0].message).toBe("syntax error");
    expect(v[0].contracts[1].accepted).toBe(true);
  });

  it("keeps requirements with no result as not run", () => {
    const v = joinResults(reqs, null);
    expect(v[0].contracts[0].status).toBeUndefined();
    expect(v[0].contracts[0].properties[0].status).toBeUndefined();
    expect(summarize(v).notRun).toBe(6);
  });

  it("ignores an empty world and pairs repeated names in order", () => {
    const dup: WireRequirement[] = [{ name: "r", statement: "", notes: [], contracts: [
      { name: "c", description: "", status: "PROPOSED", properties: [prop("PRED")] },
      { name: "c", description: "", status: "PROPOSED", properties: [prop("PRED")] },
    ] }];
    const v = joinResults(dup, { results: [
      { requirement: "r", contract: "c", status: "SATISFIED", properties: [{ status: "SATISFIED", world: { relationships: " \n" } }] },
      { requirement: "r", contract: "c", status: "VIOLATED", properties: [{ status: "VIOLATED" }] },
    ] });
    expect(v[0].contracts.map((c) => c.status)).toEqual(["SATISFIED", "VIOLATED"]);
    expect(v[0].contracts[0].properties[0].world).toBeUndefined();
  });
});

describe("summarize", () => {
  it("counts contracts and properties", () => {
    const v = joinResults(reqs, {
      results: [
        { requirement: "r1", contract: "c1", status: "VIOLATED", properties: [{ status: "VIOLATED" }, { status: "SATISFIED" }] },
        { requirement: "r1", contract: "c2", status: "DRAFT" },
        { requirement: "r2", contract: "c1", status: "SATISFIED", properties: [{ status: "SATISFIED" }] },
      ],
    });
    const s = summarize(v);
    expect(s.contracts).toEqual({ VIOLATED: 1, DRAFT: 1, SATISFIED: 1, UNVERIFIED: 0 });
    expect(s.properties).toEqual({ VIOLATED: 1, SATISFIED: 2, UNVERIFIED: 0 });
    expect(summaryLine(s)).toBe("3 contracts: 1 violated, 1 satisfied, 1 draft · properties: 1 violated, 2 satisfied");
  });
});
