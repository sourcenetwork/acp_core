#!/usr/bin/env sh

# Generates antlr Go code for theorem grammar
antlr4  -Dlanguage=Go -package theorem_parser -visitor -no-listener ./internal/parser/theorem_parser/Theorem.g4
antlr4  -Dlanguage=Go -package permission_parser -visitor -no-listener ./internal/parser/permission_parser/PermissionExpr.g4