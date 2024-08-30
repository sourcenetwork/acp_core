---
date: 2024-08-26
---

# Playground

ACP Core ships with an experimentation environment called `Playground`.
While interacting with the playground, users can make use of all the features in ACP Core in an isolated environment.
The Playground is modeled as a [protobuff service](../../proto/sourcenetwork/acp_core/playground.proto).

## Concepts

The playground system was designed to use the same entities as the primary ACP Engine.
All entities in the Engine, such as policies, relationships and theorems, apply to the Playground system.

Effectively the `Playground` system manages a collection of `Sandbox`es.
A `Sandbox` is an isolated instance of an ACP Engine where only a single `Policy` exists.

Other than a `Policy`, The `Sandbox` state is defined by a set of `Relationship`s and a `PolicyTheorem`.
All `Relationship`s added to the `Sandbox` are automatically attached to its `Policy`, likewise for the `PolicyTheorem`.

Every Sandbox contains a temporary buffer for invalid state called `Scratchpad`.
Every time an attempt to update the sandbox state is made, the new state is stored in the Scratchpad.

# Usage

Interacting with the Playground is done through some implementation of the PlaygroundService definition.

Sandboxes are created with the `NewSandbox` procedure, which returns a new and *uninitialized* Sandbox entity with a Handle.
The Handle acts as the sandbox identifier.

An unitialized sandbox contains no data and cannot do anything.
To initialize a Sandbox it's necessary to set its state through the `SetState` procedure, which loads the Sandbox's Policy, Relationships and Theorems.
Performing a `SetState` call with invalid data will not update the Sandbox's state, instead it will register those changes to the Sandbox's scratchpad.
The Scratchpad is persisted alongside the Sandbox, meaning clients are not required to store the last state they sent.
Each `SetState` call will update the Scratchpad.

Once the Sandbox has been initialized, the `EvaluateTheorems` procedure can be used to evaluate the Sandbox's theorem.

In the future, sandboxes will also support individual Check, Expand and Reverse calls.