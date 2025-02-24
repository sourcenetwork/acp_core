---
date: 2025-02-24
title: Policy Mutation 
---

This doc proposes the scope and functionality for the Policy Mutation feature in the context of SourceHub ACP and acp core.

This is a multifaceted feature, which could be designed in many different ways to enable different use cases.
The initial version of Policy Mutation is designed to enable the following use case: *enable developers to move fast and make changes to their apps in a transparent manner to their user*.
This guiding principle was used to define the scope for this initial mutation work.

There are two properties which shall be respect during policy mutation.
Any attempt to mutate a policy which would break one of these properties is considered invalid.

1. Owners must not lose autonomy over their objects
2. A Defra compliant Policy will remain compliant after a mutation. 

If a mutated Policy maintains the previous properties, it will be accepted as the new rule set for a relation graph.

The interface for the policy mutation system will be a single Msg / GRPC Method with the following signature:
MutatePolicy(policyId: string, policy: string) -> Result PolicyRecord 

Authorization rules: Only the original creator of the Policy is allowed to mutate it.

# Alternative Approaches

## Multi-Policy Relation Graph

Source's stack addresses a challenge in local-first app which is version migration.
Through data lenses powered by LensVM, we enable local applications to seamlessly communicate with each other, even if each process runs a different application version.

We considered how the ACP engine would fit into this context and it might be possible to have relation graphs which would support multiple versions of a policy.
The rough idea would be to change the relation arity between a relation graph and a policy.
A single relation graph could be related to multiple policies (each mutations of some initial policy).
The mutations in a policy would be non destructive, meaning only adding and deprecating resources, relations and permissions would be allowed.
This approach could be compatible with Source's goal to enable transparent multi-versioned local apps.

There are some open questions wrt to security however, namely what would happen on permission expression mutation.
In theory each version of a policy could execute its own permission expression in order to define its access control, however the consequences of such approach are not obvious.
There's no clear answer as to the impacts of having multiple interpretations to the same set of relationships.

By further understanding the guarantees we desire out of the system and formalizing them, it might be possible to limit permission mutations to some desireable bounds.

## Unsafe Operations

Some operations in policy mutation are inherently unsafe, policy writes have the ability to completely change the interpretation of a relation graph.
Since users only control their relationships but not the policy, a malicious policy writer could exploit these unsafe operations to gain unwaranted access to user data.

The known unsafe operations are
- mutating permission expression
- mutating a relations management graph
- mutation a relations allows types list
- renaming relations and permissions

The guiding philosophy for the initial mutation work is to enable developers to move fast, therefore no work will be done to restrict these unsafe operations, asside from not supporting renaming relations and permissions.

With better defined requirements, it may be possible to restrict these mutations somewhat and increase user security.
Another alternative would be to mark these operations as "breaking" in a policy, requiring users to accept to a new policy ID.
This could be used in conjuction with the previously mentioned Multi-Policy Relation Graph

## Safety Guarantees

There are some security goals are still being identified and refined for policy mutation.
The goal is to make the system flexible for developers and secure for users, though that obviously requires some compromises.

Some identified goals, stated informally are: 
1. Devs must not gain acces to private data after mutation a policy
2. Collaborators should not lose access to previously shared docs
3. Mutation must be defined if model check fails 

Further work is required to understand whether all these goals are desireable or necessary, and how to enforce them.
