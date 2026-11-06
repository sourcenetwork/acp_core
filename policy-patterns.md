zanzibar and acp core has a finite vocabulary to build policies.

this means there are patterns which are used to solve specific user needs. 
associated with these needs, we can propose some ways to verify the contracts.


## relation typing

in acp_core, relations are by definition just relations, there is no userset rewriting expressions.

in acp core, the only configuration a relation has is restricting the possible subject sets it may have.

eg. a relation reader may have as subject `actor` or `group->member`.

a requirement such as `group members should be file readers` directly translates to adding `group->member` as one of the allowed types for the file `reader` relation.

that becomes a contract.

the contract assertion should assert that there exists some instance where an actor who is the member of a gruop is also a file reader.


TODO: i dont' think a vacuity check makes sense here, this is a binary thing, it is either possible or not

## computed usersets

in zanzibar, computed usersets are the answer to nested permissions or inhereted permissions.

a computed userset is able to express a requriement such as `writing implies read`.

in acp core, every permission is a combination of computed userset or tuple to usersets.

a permission `read` with expression `writer + reader` is two computed usersets with the relatiosn writer and reader.

a natural language requirement such as `write permission implies read permission` translates to `expr: write + ....`.

this contract should be checked two ways:

1. all subjects with write permission also have read permission (checks the statement holds)
2. there exists some subject that can read but not write (checks that the permission is not overly strong)
3. check that there exists subjects with the write permission in the first place (asserts it's possible)

note that 2 may be unsat if another contract says the opposite, `read implies write`.
that would point to a higher level problem in the contracts, which are inconsistent amongst themselves.
the unsat instance can be bubbled up to the user in the playground, which would require them to reconcile the needs

## tuple to usersets

tuple to usersets are zanzibar's answer to hiearchical objects.
the original article uses tuple to usersets to model the requirement of nested objects, such as directories and files.

a requirement such as "an actor with read permission over a directory should be able to read all files in a directory" is a natural fit for tuple to userset.

a nested object or hiearchical object contract should be checked in a few ways:

1. all subjects with the correct permission on the "parent" object should have the same permission over it's children. 
2. there exits a subject with the permission on the child resoruce but not on the parent
3. check that there exists instances where there exists an actor with the correct permission in the parent resources (checks for impossible conditions)

## union, intersection and difference

permissions are combined using set combinators.

in a permission, two terms can be joined using `+` which is equivalent to set union.
`-` denotes set subtraction
`&` denotes set intersection.

a requirement such as "actors have read permission if they are a reader, unless they are banned" means that `-` should be used to remove a banned actor from the relation sets

a require which denotes a conjunct such as "actors have read permission if they are active and reaers" implies the usage of `&` to combine the expression.

requirements which express one or the other should use `+`, such as "actors can read if they are readers or writers".

permissios should have contracts for these operators to assert the correct scenarios.

for instance, if a requirement says that no banned user can read, the correct contract should be that no actor who is banned has read permission.

TODO: refine these rules, there are probably other checks to run here.
