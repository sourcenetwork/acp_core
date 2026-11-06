// Bundle Monaco with the app instead of letting @monaco-editor/react fetch it from a CDN.
// edcore.main = editor core + all editor features (find, folding, markers, ...), without
// the css/html/json/ts language services; only the YAML tokenizer is added.
import { loader } from "@monaco-editor/react";
import * as monacoApi from "monaco-editor/esm/vs/editor/editor.api.js";
import "monaco-editor/esm/vs/editor/edcore.main.js";
import "monaco-editor/esm/vs/basic-languages/yaml/yaml.contribution.js";
import EditorWorker from "monaco-editor/esm/vs/editor/editor.worker.js?worker";

self.MonacoEnvironment = {
  getWorker: () => new EditorWorker(),
};

loader.config({ monaco: monacoApi });
