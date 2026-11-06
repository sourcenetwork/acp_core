# Domain: <application name>

The language we use with the user. Their words, not acp's. Read this first in every session; update it the moment you
learn something. Entities here become resources in `policy.yaml`, under the same name.

## Entities

### <Term>
- **Means**: one sentence, in the user's words.
- **Also called**: the user may sometimes call a <term> a <other word>. (Only after the user confirmed they are
  interchangeable.)
- **Policy name**: `<resource>` (omit until the policy exists)

**Relationships** (one line each; every answer to "how do these relate?" goes here)
- A <term> belongs to at most one <other term>. → constraint `at_most_one: <resource>.<relation>`
- A <term> can never be inside itself. → constraint `acyclic: <resource>.<relation>`
- A person is never both <role> and <role> on the same <term>. → constraint `disjoint: [...]`

**What people can do** (actions and questions, in the user's words)
- <Who> can <do what>.

## People and roles

- **<Role>**: who they are, how someone becomes one, who decides.

## Actions the application performs

What the application does to change access: share, invite, move, transfer, remove.

- **<Action>**: who may do it, to what.

## Words we are unsure about

Terms the user used that we have not settled. Ask about each before relying on it.

- "<word>": might mean <term>? Asked / not yet asked.

## Out of scope

Things the user mentioned that acp cannot express, or that they chose to leave out, and the reason.
