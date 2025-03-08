# About

Parser pkg contains an ANTLR4 grammar file which represents the DSL used to declare PolicyTheorems and Relationships within ACP Core

# Examples

The `relationship_document` production rule consumes all the input stream looking for a relationship.
Any trailing input will cause an error.
example:

```
resource:foo#owner@did:example:bob
```

The `relationship_set` production rule consumes as many relationships as it can, greedily.
Relationships are separated by whitespace characters as per the `WS` lexeme eg (`\t \n\r`)
Trailing input in the input stream causes an error.

example:
```
resource:foo#owner@did:example:bob
resource:foo#reader@group:admin#member

// relationship sets can have comments until EOL

resource:foo#year@year:2024 
```

Theorems are parsed by the `policy_theorem` production rule.
It is also a greed rule, which consumes all of the input stream or errors.

Ex:
```
Authorizations {
    file:abc#reader@did:example:bob
    !file:abc#owner@did:example:bob //coment until EOL
}

Delegations {
    did:example:alice > file:abc#reader // alice can manage `reader` for `file:abc`
}
```

# Developing

## Getting ANTLR
Installing ANTLR can be done simply using the `pip` package `antlr4-tools` as explained in the (documentation)[https://github.com/antlr/antlr4/blob/master/doc/getting-started.md#getting-started-the-easy-way-using-antlr4-tools].

[ ] TODO: A Dockerfile + script for deterministic generation 

## Parser Generation

The parser for the grammar can be generated using `tools/generate-parser.sh` using the project root as PWD.