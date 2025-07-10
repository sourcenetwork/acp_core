import { Monaco } from "@monaco-editor/react";

export const POLICY_THEOREM_LANGUAGE_ID = "policyTheorem";

export function registerPolicyTheoremLanguage(monaco: Monaco) {
  // Register the language
  monaco.languages.register({ id: POLICY_THEOREM_LANGUAGE_ID });

  // Define language tokens
  monaco.languages.setMonarchTokensProvider(POLICY_THEOREM_LANGUAGE_ID, {
    tokenizer: {
      root: [
        // Comments
        [/\/\/.*$/, "comment"],

        // Section headers (keywords)
        [/\b(Authorizations|Delegations|ImpliedRelations)\b/, "keyword"],

        // Built-in operations
        [/\b(delete|create)\b/, "keyword.operation"],

        // Quoted strings (UTF-8 IDs)
        [/"[^"]*"/, "string"],

        // DID identifiers (more precise pattern matching grammar)
        [/did:[a-z0-9]+:[a-z0-9A-Z._-]+/, "string.key"],

        // Negation operator
        [/!/, "keyword.operator"],

        // Multi-character operators
        [/=>/, "operator.arrow"], // Implied relations
        [/->/, "operator.arrow"], // TTU relations
        [/>/, "operator.delegation"], // Delegations

        // Permission expression operators
        [/[+\-&]/, "operator.permission"],

        // Braces and parentheses
        [/[{}()]/, "delimiter.bracket"],

        // Single character operators
        [/[#@]/, "operator"],

        // Resource type patterns (e.g., "file:readme")
        [/\b[a-zA-Z][a-zA-Z0-9_]*:[a-zA-Z0-9_]+/, "entity.name.type"],

        // Relation names (after # and before @ or other operators)
        [/(?<=#)[a-zA-Z][a-zA-Z0-9_]*(?=[@>]|$)/, "entity.name.function"],

        // Regular identifiers
        [/\b[a-zA-Z][a-zA-Z0-9_]*/, "identifier"],

        // Whitespace
        [/\s+/, "white"],
      ],
    },
  });

  // Configure language features
  monaco.languages.setLanguageConfiguration(POLICY_THEOREM_LANGUAGE_ID, {
    comments: {
      lineComment: "//",
    },
    brackets: [
      ["{", "}"],
      ["(", ")"],
    ],
    autoClosingPairs: [
      { open: "{", close: "}" },
      { open: "(", close: ")" },
      { open: '"', close: '"' },
    ],
    surroundingPairs: [
      { open: "{", close: "}" },
      { open: "(", close: ")" },
      { open: '"', close: '"' },
    ],
  });
}

// Define color themes for the language
export function definePolicyTheoremTheme(monaco: Monaco) {
  // Light theme
  monaco.editor.defineTheme("policyTheoremLight", {
    base: "vs",
    inherit: true,
    rules: [
      { token: "comment", foreground: "6A9955" },
      { token: "keyword", foreground: "008080" },
      { token: "keyword.operation", foreground: "0451a5" },
      { token: "delimiter.bracket", foreground: "000000" },
      { token: "entity.name.type", foreground: "0451a5" }, // Resource (file:readme)
      { token: "operator", foreground: "0451a5" }, // # and @
      { token: "operator.arrow", foreground: "D73A49" }, // => and ->
      { token: "operator.delegation", foreground: "FF6F00" }, // >
      { token: "operator.permission", foreground: "8E24AA" }, // +, -, &
      { token: "entity.name.function", foreground: "795548" }, // Relations
      { token: "string.key", foreground: "0451a5" }, // DID
      { token: "string", foreground: "0451a5" }, // Quoted strings
      { token: "keyword.operator", foreground: "D73A49", fontStyle: "bold" }, // !
      { token: "identifier", foreground: "0451a5" }, // Default identifiers
    ],
    colors: {},
  });

  // Dark theme
  monaco.editor.defineTheme("policyTheoremDark", {
    base: "vs-dark",
    inherit: true,
    rules: [
      { token: "comment", foreground: "6A9955" },
      { token: "keyword", foreground: "3dc9b0" },
      { token: "keyword.operation", foreground: "ce9178" },
      { token: "delimiter.bracket", foreground: "D4D4D4" },
      { token: "entity.name.type", foreground: "ce9178" }, // Resource (file:readme)
      { token: "operator", foreground: "ce9178" }, // # and @
      { token: "operator.arrow", foreground: "F92672" }, // => and ->
      { token: "operator.delegation", foreground: "FF9800" }, // >
      { token: "operator.permission", foreground: "AB47BC" }, // +, -, &
      { token: "entity.name.function", foreground: "DCDCAA" }, // Relations
      { token: "string.key", foreground: "ce9178" }, // DID
      { token: "string", foreground: "ce9178" }, // Quoted strings
      { token: "keyword.operator", foreground: "F92672", fontStyle: "bold" }, // !
      { token: "identifier", foreground: "ce9178" }, // Default identifiers (permissions)
    ],
    colors: {},
  });
}
