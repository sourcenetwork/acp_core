import { describe, expect, it } from "vitest";
import { checkableNames, defaultName, parsePolicyResources } from "./policyResources";
import { SAMPLE_POLICY } from "./sample";

describe("parsePolicyResources", () => {
  it("reads permissions and relations per resource", () => {
    const r = parsePolicyResources(SAMPLE_POLICY);
    expect(Object.keys(r)).toEqual(["folder", "doc"]);
    expect(r.doc).toEqual({ permissions: ["view"], relations: ["folder", "viewer"] });
  });
  it("tolerates map-shaped resources and missing sections", () => {
    const r = parsePolicyResources("resources:\n  file:\n    permissions:\n      read: {expr: owner}\n");
    expect(r.file).toEqual({ permissions: ["read"], relations: [] });
  });
  it("yields nothing for invalid or empty yaml", () => {
    expect(parsePolicyResources("a: [")).toEqual({});
    expect(parsePolicyResources("")).toEqual({});
    expect(parsePolicyResources("resources: 3")).toEqual({});
  });
});

describe("checkableNames", () => {
  const r = parsePolicyResources(SAMPLE_POLICY);
  it("lists permissions, relations then the implicit owner", () => {
    expect(checkableNames(r, "doc")).toEqual(["view", "folder", "viewer", "owner"]);
  });
  it("does not duplicate a declared owner", () => {
    const p = parsePolicyResources("resources:\n- name: x\n  relations:\n  - name: owner\n");
    expect(checkableNames(p, "x")).toEqual(["owner"]);
  });
  it("is empty for unknown resources", () => {
    expect(checkableNames(r, "nope")).toEqual([]);
  });
  it("picks view by default", () => {
    expect(defaultName(checkableNames(r, "doc"), r, "doc")).toBe("view");
    const p = parsePolicyResources("resources:\n- name: x\n  permissions:\n  - name: read\n");
    expect(defaultName(checkableNames(p, "x"), p, "x")).toBe("read");
  });
});
