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

        // Section headers
        [/\b(Authorizations|Delegations)\b/, "keyword"],

        // Braces
        [/[{}]/, "delimiter.bracket"],

        // Negation operator at start of line
        [/^\s*!/, "keyword.operator"],

        // Identifiers
        [/did:[a-zA-Z0-9_]+:[a-zA-Z0-9_]+/, "string.key"],

        // Resource type (e.g., "file:readme")
        [/\b[a-zA-Z][a-zA-Z0-9_]*:[a-zA-Z0-9_]+/, "entity.name.type"],

        // Hash separator
        [/#/, "operator"],

        // At separator
        [/@/, "operator"],

        // Action names (words that come after # and before @)
        [/\b[a-zA-Z][a-zA-Z0-9_]*(?=@)/, "entity.name.function"],

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
    brackets: [["{", "}"]],
    autoClosingPairs: [{ open: "{", close: "}" }],
    surroundingPairs: [{ open: "{", close: "}" }],
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
      { token: "keyword", foreground: "0000FF" },
      { token: "delimiter.bracket", foreground: "000000" },
      { token: "entity.name.type", foreground: "267F99" }, // Resource (file:readme)
      { token: "operator", foreground: "267F99" }, // # and @
      { token: "entity.name.function", foreground: "267F99" }, // Action (read, write)
      { token: "string.key", foreground: "267F99" }, // DID (did:user:bob)
      { token: "keyword.operator", foreground: "D73A49" }, // !
      { token: "identifier", foreground: "000000" }, // Default identifiers
    ],
    colors: {},
  });

  // Dark theme
  monaco.editor.defineTheme("policyTheoremDark", {
    base: "vs-dark",
    inherit: true,
    rules: [
      { token: "comment", foreground: "6A9955" },
      { token: "keyword", foreground: "10CBFF" },
      { token: "delimiter.bracket", foreground: "D4D4D4" },
      { token: "entity.name.type", foreground: "0DE09E" }, // Resource (file:readme)
      { token: "operator", foreground: "0DE09E" }, // # and @
      { token: "entity.name.function", foreground: "0DE09E" }, // Action (read, write)
      { token: "string.key", foreground: "0DE09E" }, // DID (did:user:bob)
      { token: "keyword.operator", foreground: "F92672" }, // !
      { token: "identifier", foreground: "D4D4D4" }, // Default identifiers
    ],
    colors: {},
  });
}
