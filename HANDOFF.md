# Handoff: refining the policy-engineer skill

You're joining Bruno (the user) to refine `/home/lodek/source/workspace/gents-acp/policy-compiler/skills/policy-engineer/`, the skill that lets an agent design acp_core
authorization policies with a user who doesn't know acp. This document summarizes the design discussions so far. Read
it, then `/home/lodek/source/workspace/gents-acp/policy-compiler/GLOSSARY.md`, then the skill itself. Work with the user: this is a refinement conversation, not a solo task.

## The vision

acp_core is a Zanzibar-style discretionary access control engine (used by DefraDB and SourceHub). Users won't learn
it. The idea is to shift that burden to agents, and make the agents' work trustworthy with an Alloy-based
**trusted kernel**, Lean-style:

- Policies compile to Alloy models (`policy_compiler`). The agent never hand-encodes a policy.
- **Contracts are the core value.** The user states requirements in natural language; the policy engineer (the agent)
  translates them into contracts, which Alloy checks over every valid world up to a scope.
- **Contract-first.** The user is responsible for the completeness of the contracts; the agent is responsible for
  translating them faithfully and for surfacing consequences the user hasn't decided on.
- Instance generation and counterexamples are how agents and users elicit and verify intent: concrete worlds, not
  policy syntax.
- Recurring contract patterns: implied permissions (edit implies view), userset hierarchy (team members inherit),
  object hierarchy (folder contents inherit).
- Scope: layers 1–2 (relations, permissions, management/authority/influence). Layer 3, state transitions (mainly
  ownership transfer), is deferred.

## Vocabulary (agreed; see GLOSSARY.md)

| Term | Meaning |
|---|---|
| Policy | The rules acp enforces: resources, relations, permissions, management. |
| System constraints | Rules about relationship data the application enforces and acp doesn't: `disjoint`, `at_most_one`, `acyclic`. The YAML key, the CLI flag (`-a`) and `_assumptions.als` still say "assumption". |
| System domain | Policy + system constraints: every valid world, and how access works in it. |
| System specification | System domain + contracts. |
| World | One concrete state: relationships (+ optional assertions). |
| Scenario | System domain + world. |
| Assertion | A concrete expectation about one world, checked by the engine. In code/DSL: "theorem". |
| Requirement | What the user wants, in their own words, and the contracts derived from it. |
| Contract | One claim derived from a requirement, which the user can accept; verified by its properties. |
| Property | A typed Alloy formula (`assert` or `pred`) with a description, checked over all worlds up to a scope. |
| Counterexample / example | The world where an `assert` property fails / a `pred` property holds. |
| Workspace | The working set on disk: policy, constraints, contracts, worlds. |

"System model" is retired as a term (the proto message is still called `SystemModel`). Use these words with the user
and in the skill; vocabulary drift was a real problem in an earlier session.

Operations: **Verify** (contracts), **Generate** (worlds), **Evaluate** (assertions on a world, the old `Simulator`
service), **Check** and **Explain** (one question about one world, and why). There is no Expand: zanzi doesn't
implement it.

## Requirements, contracts and properties

A requirement usually derives several contracts, and a contract usually needs several properties. "Actors can read
if they're readers or writers, but not banned" derives three contracts (unbanned readers can read, unbanned writers
can read, no banned actor can read), and "unbanned readers can read" has two properties: every unbanned reader can
read (assert), and unbanned readers can exist at all (pred).

```yaml
# contracts.yaml
requirements:
- name: read-access
  statement: "actors can read if they're readers or writers, but not banned"   # the user's words, verbatim
  notes:                              # the agent's open questions, assumptions, things it didn't understand
  - "Assumed writers means the doc's writers, not its folder's."
  contracts:
  - name: unbanned-readers-can-read
    description: Readers can read the doc, unless they're banned.
    status: accepted                  # proposed (default) or accepted
    properties:
    - type: assert                    # assert: holds in every world. pred: holds in some world.
      description: every unbanned reader of a doc can read it
      formula: |                      # block scalar: `d: Doc` is invalid in a plain yaml string
        all d: Doc | d.doc_reader - d.doc_banned in d.doc_read
    - type: pred
      description: unbanned readers can exist
      formula: |
        some d: Doc | some d.doc_reader - d.doc_banned
  - name: banned-never-read
    description: No banned actor can read, whatever other role they hold.   # no properties yet: a draft
```

- **Requirement**: the user's `statement`, verbatim, and the agent's `notes`. Fidelity is judged here: do these
  contracts, together, say what the statement says? Decomposing is itself a translation step: "readers can read",
  taken literally, contradicts "no banned actor can read" for a banned reader.
- **Contract**: one claim, in plain language, that the user can accept. `status` records the user's decision only:
  `accepted` means the user agreed the description states what they want, not that the policy satisfies it. A
  contract without properties is a draft. Verify results are never stored in the file: they're derived and go stale.
- **Property**: a typed formula with its own description. Every field is required except `notes` and `status`.
  Free text containing `: ` must be quoted or a block scalar, as in YAML generally.
- Verify wraps each formula, so a property can't drop the constraints or choose its own scope:
  - `assert` → `check { assumptions implies { formula } } for N`. Violated → counterexample world.
  - `pred` → `run { assumptions and { formula } } for N`. Satisfied → example world. Preds catch over-restrictive
    policies (a policy where nobody can see anything passes every assert) and replace the old "vacuity run" advice.
- Formulas use the compiled model's functions: `d.doc_edit`, `f.folder_owner`, `d.folder.folder_view`. These names
  are now effectively an API.
- Results, per property: `SATISFIED`, `VIOLATED`, or `UNVERIFIED` (formula error with a position inside the formula,
  or a timeout). Per contract: `VIOLATED` if any property is, else `UNVERIFIED` if any is, else `SATISFIED`; `DRAFT`
  without properties. Drafts don't fail `verify`. Satisfied asserts hold **up to the scope** only; say so.
- Counterexamples carry **bindings**: the values of the formula's variables (`d = doc:doc0`, `a = did:user:bob`),
  which point at the exact object and actor, ready for a Check/Explain.

### Resolving a counterexample

A counterexample is a concrete question: "in this world, bob can edit doc0; the contract says he shouldn't — which is
wrong?" There are exactly three answers, each editing a different part of the system specification:

1. **The policy is wrong** → fix the policy.
2. **The contract is wrong** (too strict, or mistranslated) → fix the contract.
3. **The world can't happen** in the application → add a system constraint.

Choosing among these with the user, in their vocabulary, is the central interaction. After a fix, the counterexample
can be kept as a world with assertions: a regression test.

Counterexamples can be large and noisy. Shrinking is low priority: at scope ~5 with 5–10 resources they're manageable,
and Check/Explain on the bound object and actor narrows attention to the relevant path.

## The core open problem: description vs formula fidelity

The user reads the description; Alloy checks the formula. A mismatch gives a green report for a contract nobody agreed
to. This is the hardest and most open part of the whole approach, and the main thing to refine. In the earlier
company-brain session, agent-written checks were the weakest link.

Ideas so far, cheapest first:
1. Counterexamples catch failures for free.
2. Examples (pred contracts, or runs of an assert's antecedent) show the contract actually applies.
3. **Redundant formalization** (cf. Amazon's neurosymbolic / Automated Reasoning checks work, worth reading): two
   independent formalizations of the same description, checked for equivalence in Alloy:
   `check { assumptions implies (f1 iff f2) }`. A counterexample is a world where the readings disagree, which is
   exactly the clarifying question to ask the user.
4. Independent back-translation: another model describes the formula without seeing the description; compare.

## What exists today

In acp_core (this repository, mostly uncommitted on `feat/policy-builder`):

- `policy_compiler compile [-a constraints.yaml] <policy.yaml> <dir>`: emits `<stem>.als`,
  `<stem>_management.als`, `<stem>_assumptions.als`. Build: `make policy_compiler`.
- `policy_compiler verify [-a constraints.yaml] [-s scope] [-o dir] <policy.yaml> <contracts.yaml>`: report per
  requirement, contract and property, with bindings and worlds; exit 1 unless every non-draft contract is satisfied.
  Needs `alloy` locally.
- `policy_compiler instance [-o dir] <policy.yaml> <xml>...`: Alloy XML instance → relationships, skolems, and what
  each actor can do (use `-t xml`; Alloy's JSON drops tuples).
- `acp_simulate policy.yaml relationships.txt theorems.txt`: replays a world on the real engine (exit 0/1/2).
- `Playground2` service (Connect/gRPC): `Generate`, `Check`, `Explain`, `Verify`, all stateless and deterministic.
  Served locally by `cmd/playground2` on `localhost:8090`. A separate agent is building a web UI in `playground2/`
  (see `acp_core/docs/playground2-handoff.md`).
- `acp` (`make acp`), the thin client below. It reads a workspace directory and calls the service at `-server` /
  `$ACP_SERVER` (default `http://localhost:8090`; run `make playground2 && build/playground2` locally, which needs
  `alloy`). Exit codes: 0 success, 1 failure (contracts not satisfied, access denied), 2 error.
  - `acp verify [-s scope] [-o dir] <workspace>`: the contracts report.
  - `acp generate [-n count] [-s scope] [-a actors] [-o dir] <workspace>`: example worlds.
  - `acp model [-o dir] <workspace>`: the compiled Alloy modules (`policy.als`, `policy_management.als`,
    `policy_assumptions.als`), ie. the functions formulas can use.
  - `acp check|explain <workspace> <world> <object> <permission> <actor>`, eg.
    `acp explain myapp child-owner-sees-parent-1 doc:doc0 edit did:user:bob`: one question about a saved world;
    explain prints the engine's search as a tree, ✓ marking the path that grants access. Use it on a counterexample's
    bindings.
  - `acp playground [-addr address] [-ui dir] [-open=false] <workspace>`: serves the playground UI for the
    workspace, live-reloading on file changes, proxying API calls to the service. The UI is embedded
    (`cmd/acp/ui`, refreshed by `make acp-ui`). The UI side of file-backed mode is in progress
    (`docs/playground2-acp-handoff.md`).
  - `-o` saves counterexamples, examples or generated worlds in the `worlds/*.yaml` format, never overwriting.
- **One stateless remote service** (the Playground2 service in a container with Alloy), hosted by Source and
  self-hostable. Requests carry the whole specification; the server stores nothing.
- **A thin local binary**, `cmd/acp`, runs nothing locally: it reads the workspace directory, builds a
  `Workspace` message, and calls the service. Built: `verify`, `generate`, `model`, `check`, `explain`, and
  `playground`, which serves the UI locally, proxies API calls to the remote service (same origin; the CLI will add
  credentials), and watches the workspace so the human sees the agent's edits live.
- Access to the remote service can be ungated at first; as a product it may need OAuth/OpenID tokens or rate limits.
- Workspace layout: `policy.yaml` (required), `constraints.yaml`, `contracts.yaml`, `worlds/*.yaml` (each
  `relationships: |` and optionally `assertions: |`, named after the file).
- A remote MCP server comes later, as a thin wrapper. The CLI comes first because it reads files itself: the report
  describes exactly what's on disk, whereas over MCP the agent could verify a stale copy and report success.
- A binary-free setup is possible later but would make the service stateful (hosted workspaces); not now.

So the skill should move toward `acp verify` / `acp generate` over a workspace directory. `policy_compiler verify`
does the same with explicit file paths and local Alloy; it's likely to go away in a cleanup.

## Where the skill is out of date

The skill predates contracts and the vocabulary. Known gaps, as a starting list for the refinement:

- It has agents hand-write `checks.als` with `assert`/`check` and run `alloy exec`. Requirements, contracts and
  properties in `contracts.yaml` + `verify` should replace this.
- "Default expectations" (local control, no amplification, exclusions, inheritance direction) predate contract-first.
  An option already discussed: a per-resource contract checklist the user confirms, replacing agent-chosen defaults.
- "Vacuity runs" → `pred` properties.
- "Application assumptions" → system constraints (keep noting the YAML key is still `assumptions:`).
- The three-way counterexample resolution isn't in the skill.
- `decisions.md` overlaps with requirement statements, notes and contract status; decide what lives where.
- Reporting discipline: the earlier session reported overconfidently. Satisfied means "up to scope N", contracts are
  only as good as their translation, and anything not covered by a contract is unverified.

## Lessons from an earlier agent session (company-brain policy)

- Agent-written checks were the weak link: plausible-looking assertions that didn't capture the requirement.
- Overconfident reporting of what had been verified.
- Vocabulary drift between agent, user and tools.
- It did find a real issue: a CTO could appoint themselves team lead and reach non-engineering HR records, if teams
  were linked to the CTO. That composition consequence is the kind of finding the approach exists for.

## Engine and tooling gotchas

- Permission expressions evaluate left to right with **no precedence**: parenthesise.
- DSL ids need at least 2 characters (or quotes).
- An object's owner relationship must come before any other relationship on it; subjects needn't be registered.
- `folder->owner` is allowed in a permission; a bare `owner` is rejected.
- TTU lookup relations (left of `->`) need exactly one object type. Recursion must be a top-level `+` term.
- Engine errors are obscure (eg. "object not registered" when an owner relationship is missing). Clearer errors are
  wanted but deferred.

## Open threads (not for the skill work unless the user picks them up)

- A structured steering knob for Generate (`require: [doc.folder, folder.parent]`).
- Evaluate moving into Playground2; renaming internal "assumption" to "constraint".
- Whether `disjoint` should compare across userset kinds (undecided).
- The experiment harness; instance diversity; layer 3 transitions.

## Working with Bruno

- Prefers simple, clean code and tight scope; ask before adding scope.
- Prototyping mode: favour small, verifiable steps.
- Use the agreed vocabulary; propose new terms explicitly rather than drifting into them.
