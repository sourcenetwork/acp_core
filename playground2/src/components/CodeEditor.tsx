import { Editor, useMonaco } from "@monaco-editor/react";
import { useEffect, useRef } from "react";
import type { editor } from "monaco-editor";

export interface EditorMarker {
  message: string;
  line?: number;
}

interface Props {
  value: string;
  onChange?: (v: string) => void;
  readOnly?: boolean;
  language?: string;
  dark: boolean;
  markers?: EditorMarker[];
  testId?: string;
}

/** Monaco wrapper bound to a string value; shows markers for located errors. */
export function CodeEditor({ value, onChange, readOnly, language = "yaml", dark, markers = [], testId }: Props) {
  const monaco = useMonaco();
  const ref = useRef<editor.IStandaloneCodeEditor | null>(null);

  useEffect(() => {
    const model = ref.current?.getModel();
    if (!monaco || !model) return;
    monaco.editor.setModelMarkers(
      model,
      "playground2",
      markers.map((m) => {
        const line = Math.min(m.line ?? 1, model.getLineCount());
        return {
          severity: monaco.MarkerSeverity.Error,
          message: m.message,
          startLineNumber: line,
          endLineNumber: line,
          startColumn: 1,
          endColumn: model.getLineMaxColumn(line),
        };
      }),
    );
  }, [monaco, markers]);

  return (
    <div className="h-full overflow-hidden rounded-md border bg-editor" data-testid={testId}>
      <Editor
        height="100%"
        language={language}
        value={value}
        theme={dark ? "vs-dark" : "light"}
        onChange={(v) => onChange?.(v ?? "")}
        onMount={(e) => {
          ref.current = e;
        }}
        options={{
          automaticLayout: true,
          fixedOverflowWidgets: true,
          readOnly,
          tabSize: 2,
          minimap: { enabled: false },
          scrollBeyondLastLine: false,
          padding: { top: 8 },
        }}
      />
    </div>
  );
}
