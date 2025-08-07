import { Monaco } from "@monaco-editor/react";

export const POLICY_THEOREM_LANGUAGE_ID = "policyTheorem";

let isLanguageRegistered = false;
let themeRegistered = false;

export interface ITokenThemeRule {
  token: string;
  foreground?: string;
  background?: string;
  fontStyle?: string;
}

export function registerPolicyTheoremLanguage(monaco: Monaco) {
  if (isLanguageRegistered) return;

  // Register the language
  monaco.languages.register({ id: POLICY_THEOREM_LANGUAGE_ID });
  isLanguageRegistered = true;

  // Monaco uses a global state for themes and languages, so we need to ensure there are no overlapping language tokens
  // See: https://github.com/microsoft/monaco-editor/issues/338
  const suffix = "policyTheorem";

  // Define language tokens
  monaco.languages.setMonarchTokensProvider(POLICY_THEOREM_LANGUAGE_ID, {
    tokenizer: {
      root: [
        // Comments
        [/\/\/.*$/, `comment.${suffix}`],

        // Section headers (keywords)
        [
          /\b(Authorizations|Delegations|ImpliedRelations)\b/,
          `keyword.${suffix}`,
        ],

        // Built-in operations
        [/\b(delete|create)\b/, `keyword.operation.${suffix}`],

        // Quoted strings (UTF-8 IDs)
        [/"[^"]*"/, `string.${suffix}`],

        // DID identifiers (more precise pattern matching grammar)
        [/did:[a-z0-9]+:[a-z0-9A-Z._-]+/, `string.key.${suffix}`],

        // Negation operator
        [/!/, `keyword.operator.${suffix}`],

        // Multi-character operators
        [/=>/, `operator.arrow.${suffix}`], // Implied relations
        [/->/, `operator.arrow.${suffix}`], // TTU relations
        [/>/, `operator.delegation.${suffix}`], // Delegations

        // Permission expression operators
        [/[+\-&]/, `operator.permission.${suffix}`],

        // Braces and parentheses
        [/[{}()]/, `delimiter.bracket.${suffix}`],

        // Single character operators
        [/[#@]/, `operator.${suffix}`],

        // Resource type patterns (e.g., "file:readme")
        [/\b[a-zA-Z][a-zA-Z0-9_]*:[a-zA-Z0-9_]+/, `entity.name.type.${suffix}`],

        // Relation names (after # and before @ or other operators)
        [
          /(?<=#)[a-zA-Z][a-zA-Z0-9_]*(?=[@>]|$)/,
          `entity.name.function.${suffix}`,
        ],

        // Regular identifiers
        [/\b[a-zA-Z][a-zA-Z0-9_]*/, `identifier.${suffix}`],

        // Whitespace
        [/\s+/, `white.${suffix}`],
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
  if (themeRegistered) return;

  const suffix = "policyTheorem";

  const policyTheoremRukes: {
    light: ITokenThemeRule[];
    dark: ITokenThemeRule[];
  } = {
    light: [
      { token: `comment.${suffix}`, foreground: "6A9955" },
      { token: `keyword.${suffix}`, foreground: "008080" },
      { token: `keyword.operation.${suffix}`, foreground: "0451a5" },
      { token: `delimiter.bracket.${suffix}`, foreground: "000000" },
      { token: `entity.name.type.${suffix}`, foreground: "0451a5" }, // Resource (file:readme)
      { token: `operator.${suffix}`, foreground: "0451a5" }, // # and @
      { token: `operator.arrow.${suffix}`, foreground: "D73A49" }, // => and ->
      { token: `operator.delegation.${suffix}`, foreground: "FF6F00" }, // >
      { token: `operator.permission.${suffix}`, foreground: "8E24AA" }, // +, -, &
      { token: `entity.name.function.${suffix}`, foreground: "795548" }, // Relations
      { token: `string.key.${suffix}`, foreground: "0451a5" }, // DID
      { token: `string.${suffix}`, foreground: "0451a5" }, // Quoted strings
      {
        token: `keyword.operator.${suffix}`,
        foreground: "D73A49",
        fontStyle: "bold",
      }, // !
      { token: `identifier.${suffix}`, foreground: "0451a5" }, // Default identifiers
    ],
    dark: [
      { token: `comment.${suffix}`, foreground: "6A9955" },
      { token: `keyword.${suffix}`, foreground: "3dc9b0" },
      { token: `keyword.operation.${suffix}`, foreground: "ce9178" },
      { token: `delimiter.bracket.${suffix}`, foreground: "D4D4D4" },
      { token: `entity.name.type.${suffix}`, foreground: "ce9178" }, // Resource (file:readme)
      { token: `operator.${suffix}`, foreground: "ce9178" }, // # and @
      { token: `operator.arrow.${suffix}`, foreground: "F92672" }, // => and ->
      { token: `operator.delegation.${suffix}`, foreground: "FF9800" }, // >
      { token: `operator.permission.${suffix}`, foreground: "AB47BC" }, // +, -, &
      { token: `entity.name.function.${suffix}`, foreground: "DCDCAA" }, // Relations
      { token: `string.key.${suffix}`, foreground: "ce9178" }, // DID
      { token: `string.${suffix}`, foreground: "ce9178" }, // Quoted strings
      {
        token: `keyword.operator.${suffix}`,
        foreground: "F92672",
        fontStyle: "bold",
      }, // !
      { token: `identifier.${suffix}`, foreground: "ce9178" }, // Default identifiers (permissions)
    ],
  };

  // Light theme
  monaco.editor.defineTheme("playgroundThemeLight", {
    base: "vs",
    inherit: true,
    rules: [...policyTheoremRukes.light],
    colors: {},
  });

  // Dark theme
  monaco.editor.defineTheme("playgroundThemeDark", {
    base: "vs-dark",
    inherit: true,
    rules: [...policyTheoremRukes.dark],
    colors: {},
  });

  themeRegistered = true;
}
