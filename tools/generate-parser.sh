#!/usr/bin/env sh

# run from inside container
# docker run --rm -it -v .:/app python:3.14-trixie bash

# Generates antlr Go code for theorem grammar
pip install antlr4-tools
antlr4  -Dlanguage=Go -package theorem_parser -visitor -no-listener ./pkg/parser/theorem_parser/Theorem.g4
antlr4  -Dlanguage=Go -package permission_parser -visitor -no-listener ./pkg/parser/permission_parser/PermissionExpr.g4