# Policy patterns and the contracts that check them

DRAFT. Formulas here were run with `acp verify` against the example policy at the end of this file (scope 3), including
against deliberately wrong policies. Names like `doc_read` and `reader_group_member` come from that policy; yours
follow the same convention (see "The compiled model").

acp has a finite vocabulary: typed relations, permission expressions made of computed usersets, tuple to usersets and
`+ - &`. So a requirement usually maps to a known construct, and a construct to known contracts. Use this file to
(1) pick the construct, (2) write the contracts, (3) know what to ask the user.

Write contracts from these templates and adapt them. Don't invent a new shape when one of these fits: a known shape is a
shape whose meaning you have already checked.

## Reading a requirement: four shapes

Most requirements reduce to one of these. Name the shape before you write the contract.

| shape | the user says | contract |
|---|---|---|
| **sufficient** | "X can do Y" | assert: X's actors are within Y |
| **necessary** | "only X can do Y" | assert: Y is within X, plus the owner |
| **exclusion** | "X can never do Y" | assert: no X in Y, except the owner |
| **possible** | "Y can happen" | pred |

Every `assert` that is conditional (an "if" or "unless") gets a `pred` showing the condition can happen. An assert about
a policy where nobody can do anything passes. The pred is what stops that.

Two traps in preds:

- **The owner satisfies too much.** The owner has every permission. A pred like "someone can read this doc" is satisfied
  by its owner and proves nothing. Exclude the owner (`- d.doc_owner`) unless the owner is the point.
- **Via the thing you mean.** "A group member can read" is satisfied by a direct reader who happens to be in a group.
  Bind the formula to the relation: `g in d.reader_group_member and a in g.group_member`.

## The owner

The engine adds the owner to every permission. Consequences for contracts:

- The owner has every permission and can't be locked out, whatever the policy says.
- **Exclusions never apply to the owner.** "No blocked actor can read" fails when the owner is blocked. This is not a
  bug in your policy and you can't fix it there. The counterexample has three possible resolutions, and the policy is
  not one of them:
  - the contract allows the exception: `no (d.blocked - d.doc_owner) & d.doc_read` (usually right: tell the user "the
    owner can always read their own docs"), or
  - a system constraint says the owner is never blocked: `disjoint: [doc.owner, doc.blocked]`, which is a claim about
    their application that they must confirm.
- Ask the user once, early: "Should the person who created a doc always be able to read it, whatever else happens?" The
  answer is yes in acp; the question is whether they know it.
- In `necessary` shapes ("only X can do Y"), add the owner to the right-hand side.

## Relation typing

A relation only restricts what its subjects may be: `actor`, an object (`folder`), or a userset (`group->member`).
There is no rewriting in a relation. "Group members should be readers of a doc" is `group->member` among the allowed
types of the doc's `reader`.

Ask: *can the group's owner or admins count as members?* `group->member` contains only members. `group->participant`
(a permission) also contains the owner. Point at what the user means.

Contracts:

```yaml
- type: pred          # the policy allows it (and the system constraints don't forbid it)
  description: Someone reads a doc because they are in a group it is shared with.
  formula: |
    some d: Doc, g: Group, a: Actor | g in d.reader_group_member and a in g.group_member and a in d.doc_read and a not in d.doc_owner
- type: assert        # it works: sharing with a group reaches every member
  description: Every member of a group a doc is shared with, who is not blocked, can read it.
  formula: |
    all d: Doc, g: Group, a: Actor | (g in d.reader_group_member and a in g.group_member and a not in d.blocked) implies a in d.doc_read
```

The pred says the shape is allowed; the assert says the policy does what the shape means. Write both. A separate
vacuity check isn't needed: the pred is the existence check, and it runs with the system constraints, so a constraint
that makes the shape impossible fails it.

If a relation is later removed from the policy, a contract that mentions it (`reader_group_member`) can't compile and
comes back **UNVERIFIED**, not VIOLATED. Treat that as a failure.

## Computed usersets: "A implies B"

A permission is built from computed usersets: `read: reader + writer` includes whoever holds `reader` and whoever
holds `writer`. A requirement "writing implies reading" becomes a term in the expression (`read: writer + ...`).

```yaml
- type: assert
  description: Anyone who can write a doc can read it.
  formula: |
    all d: Doc | d.doc_write in d.doc_read
- type: pred
  description: Someone other than the owner can write a doc.
  formula: |
    some d: Doc | some d.doc_write - d.doc_owner
```

The assert is the requirement. The pred stops it passing because nobody can write.

**Don't add "can read but not write" by default.** It is a contract of its own ("reading doesn't give writing") and
belongs in the list only when the user said so. Otherwise you are inventing a requirement. If the user did say both,
and the policy can't satisfy them together, the pred for the second comes back VIOLATED: the contracts conflict. Show
the user the conflicting pair in their words and ask which one they mean. Today Verify can't name the conflicting
contract for you; find it by running the pred without the other contract's asserts.

## Tuple to userset: inheritance

`folder->view` takes everyone with `view` on the objects in the `folder` relation. It is how a doc inherits from its
folder, and a folder from its parent. The permission expression already states the rule, so the assert follows from
it directly:

```yaml
- type: assert
  description: Anyone who can view a doc's folder can read the doc, unless blocked.
  formula: |
    all d: Doc, a: Actor | (a in d.folder.folder_view and a not in d.blocked) implies a in d.doc_read
- type: pred
  description: Someone other than the folder owner can view a folder that has docs.
  formula: |
    some d: Doc | some d.folder.folder_view - d.folder.folder_owner
- type: pred
  description: Reading a doc does not give access to its folder.
  formula: |
    some d: Doc, a: Actor | some d.folder and a in d.doc_read - d.doc_owner and a not in d.folder.folder_view
```

- The assert is the equivalence rule. Mind the exclusions: if `read` subtracts `blocked`, so does the assert.
- Which parent permission feeds which child permission is a question for the user: does editing a folder let you edit
  its docs? Viewing a folder, read its docs? Ask per pair.
- The last pred is the direction check: access flows down, never up. Exclude the owner, and require the doc to be in a
  folder, or the owner satisfies it.
- `d.folder` is the relation as stored (a set of folders). The generated fun `doc_folder` is empty for object
  relations; use the field.
- A recursive parent (`parent->view` inside folder `view`) makes the inheritance cover every ancestor. Ask whether
  nesting should be unlimited, and whether a folder can contain itself (that is `acyclic: folder.parent`).

## Set operators: `+`, `-`, `&`

| the user says | operator | contracts |
|---|---|---|
| "readers or writers can read" | `+` | for each term: that term alone is sufficient |
| "readers can read, unless banned" | `-` | exclusion (no banned, except owner), plus: a reader who isn't banned can |
| "must be both active and a reader" | `&` | sufficient (both → can), necessary (can → both), and each one alone can't |

```yaml
# difference
- type: assert
  description: Nobody blocked on a doc can read it, except its owner.
  formula: |
    all d: Doc | no (d.blocked - d.doc_owner) & d.doc_read
- type: pred
  description: A reader who is not blocked can read.
  formula: |
    some d: Doc | some (d.doc_reader - d.blocked - d.doc_owner) & d.doc_read

# intersection
- type: assert
  description: Active readers can review.
  formula: |
    all d: Doc | d.doc_reader & d.doc_active in d.doc_review
- type: assert
  description: Only active readers (or the owner) can review.
  formula: |
    all d: Doc | d.doc_review - d.doc_owner in d.doc_reader & d.doc_active
- type: assert
  description: An active person who is not a reader cannot review.
  formula: |
    all d: Doc | no (d.doc_active - d.doc_reader - d.doc_owner) & d.doc_review
```

Notes:

- **Evaluation is left to right, with no precedence.** `a + b - c` is `(a + b) - c`, but `a - b + c` is
  `(a - b) + c`, which doesn't exclude `b` from `c`. Always parenthesise in the policy. The exclusion contract is
  what catches the mistake: the counterexample is someone excluded who gets in through `c`.
- **Exclusions only act on what is in the expression.** `(reader + writer - blocked) + folder->view` excludes blocked
  readers but not blocked folder viewers. State the exclusion contract over the whole permission, not over the term.
- "Not X" is `no X & Y`. `X not in Y` means X is not a subset of Y, which is a different thing.

## Not covered yet

- Management, authority and influence (`manages`, who can change who has access). The first review found real problems
  here (an admin who could appoint themselves into reaching sensitive records). Contracts use `authority_*` and
  `*_influencers`. TODO.
- Ownership transfer (state transitions). Deferred.

## The compiled model

Names are generated from the policy, so a name you choose for a relation or permission becomes a name in formulas.

| in the policy | in formulas |
|---|---|
| resource `doc` | `sig Doc` (capitalised): `all d: Doc | ...` |
| actors | `Actor` |
| relation or permission `read` on `doc` | `d.doc_read`: the actors who have it on `d`, usersets expanded |
| a relation as stored | a field: `d.blocked`, `d.folder`, `d.reader_group_member` (the groups in a userset type) |
| owner | `d.doc_owner` |

- Prefer the `doc_...` funs for who has access, and the fields for object links (`d.folder`) and for the exact
  userset (`reader_group_member`).
- The comment above each generated fun shows the effective expression, with the owner on top. Read it to see what you
  are checking.
- A recursive permission becomes a closure: `(*parent).(...)`.
- Relations to objects (`folder`) have an empty fun (`doc_folder` is `none->none`). Use the field.

## Example policy

The policy these formulas were verified against.

```yaml
name: wiki
resources:
- name: group
  relations:
  - name: member
    types: [actor]
- name: folder
  relations:
  - name: parent
    types: [folder]
  - name: viewer
    types: [actor, group->member]
  - name: editor
    types: [actor]
  permissions:
  - name: edit
    expr: editor + parent->edit
  - name: view
    expr: viewer + edit + parent->view
- name: doc
  relations:
  - name: folder
    types: [folder]
  - name: reader
    types: [actor, group->member]
  - name: writer
    types: [actor]
  - name: blocked
    types: [actor]
  - name: active
    types: [actor]
  permissions:
  - name: read
    expr: ((reader + writer) + folder->view) - blocked
  - name: write
    expr: (writer + folder->edit) - blocked
  - name: review
    expr: reader & active
```

System constraints used: `acyclic: folder.parent`, `at_most_one: doc.folder`.
