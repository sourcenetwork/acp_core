#!/usr/bin/env sh

# Generates antlr Go code for theorem grammar
antlr4  -Dlanguage=Go -package parser -visitor -no-listener ./internal/parser/Theorem.g4

