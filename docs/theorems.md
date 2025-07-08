---
title: Policy Theorems
---

In ACP Core Policy Theorems are a model for access control statements which can be verified for correctness.

Theorems are a general purpose tool which can be used to interactively perform batch assertions over a policy, acts as "tests" which a policy must pass in order to be mutated, one-shot simulation and more.

# Theorems

Currently, the theorem system supports two types of theorems: Authorization theorems and Delegation theorems.

## Authorization Thereom

An Authorization thereom is a statement which models an Access Request in a Policy.
It is used to verify that some access request is allowed or denied within a Policy.

Some examples of access requests are:
- bob can read file foo
- bob cannot edit file foo

## Delegation Theorem

ACP Core can be classified as a "decentralized administration" discretionary access control system.
Effectively this means that object owners are free to give access rights to whoever they want.
Furthermore, object owners are able to give delegation power to other actors, which in turn may also grant further access permissions.

Due to the dynamic nature of delegation power in a policy, delegation theorems are used to express verifications about whether some user is able to manage a certain type of relationship to verify and reason about the state.

Some examples are:
- bob can manage read relationships for file foo
- bob cannot manage write relationships for file foo

# PolicyTheoremDSL

To succinctly express theorems in a way that is both convenient for humans and machines, we've developed PolicyTheoremDSL which is a Domain Specific Language to specify theorems.
The DSL is minimal and direct, so much so that a few examples are enough to get an idea of how to use it.
For a full specification of the language, check the parser pkg in acp_core.

Examples:

Minimal Valid Policy Theorem
```
Authorizations {}
Delegations {}
```

Full example

```
// PolicyTheoremDSL supports C-style inline comments
// "//" can be used to mark that the remainder of that line will be commented out.

Authorizations {
    file:abc#reader@did:example:bob // access request verifying bob is a reader of file abc
    !file:abc#owner@did:example:bob // asserts bob is not the owner of file abc
}

Delegations {
    did:example:alice > file:abc#reader // alice can manage `reader` for `file:abc`
    !did:example:alice > file:abc#writer // alice cannot manage `writer` for `file:abc`
}
```
