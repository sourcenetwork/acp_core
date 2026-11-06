import { describe, expect, it } from "vitest";
import { boundNodes, prefillFromBindings } from "./bindings";
import { parsePolicyResources } from "./policyResources";
import { SAMPLE_POLICY } from "./sample";

const resources = parsePolicyResources(SAMPLE_POLICY);
const objects = new Map([["doc:doc0", "doc"], ["doc:doc1", "doc"], ["folder:folder0", "folder"]]);
const actors = ["did:user:alice", "did:user:bob"];

describe("boundNodes", () => {
  it("marks only values that are nodes, collecting the variables", () => {
    const m = boundNodes(
      [{ name: "d", values: ["doc:doc0"] }, { name: "x", values: ["doc:doc0", "nope:1"] }, { name: "r", values: ["doc:doc0->did:user:bob"] }, { name: "a", values: ["did:user:bob"] }],
      [...objects.keys(), ...actors],
    );
    expect(m).toEqual({ "doc:doc0": ["d", "x"], "did:user:bob": ["a"] });
  });
});

describe("prefillFromBindings", () => {
  it("takes the unambiguous object and actor, and the permission named in the formula", () => {
    const p = prefillFromBindings(
      [{ name: "a", values: ["did:user:bob"] }, { name: "d", values: ["doc:doc0"] }],
      objects, actors, "all d: Doc, a: Actor | a in d.doc_viewer implies a in d.doc_view", resources,
    );
    expect(p).toEqual({ object: "doc:doc0", actor: "did:user:bob", permission: "view" });
  });

  it("prefers a permission over a relation mentioned earlier, and falls back to the default", () => {
    const b = [{ name: "d", values: ["doc:doc0"] }];
    expect(prefillFromBindings(b, objects, actors, "d.doc_viewer in d.doc_view", resources).permission).toBe("view");
    expect(prefillFromBindings(b, objects, actors, "d.doc_viewer = none", resources).permission).toBe("viewer");
    expect(prefillFromBindings(b, objects, actors, "true", resources).permission).toBe("view");
  });

  it("leaves things out when ambiguous", () => {
    const p = prefillFromBindings(
      [{ name: "d", values: ["doc:doc0", "doc:doc1"] }, { name: "a", values: ["did:user:alice", "did:user:bob"] }],
      objects, actors, "", resources,
    );
    expect(p).toEqual({});
  });
});
