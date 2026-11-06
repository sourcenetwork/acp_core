import { useState } from "react";
import { CodeEditor } from "./CodeEditor";
import { WorldGraphPanel, type WorldGraphInitial } from "./WorldGraphPanel";

type Mode = "graph" | "text";

/**
 * One world: the Graph / Relationships toggle over the graph with Check and Explain. Shared by the Worlds tab and
 * the Requirements tab (a counterexample or example is just another world). Key it by world to reset it.
 */
export function WorldView({ text, dark, initial, marks }: { text: string; dark: boolean; initial?: WorldGraphInitial; marks?: Record<string, string[]> }) {
  const [mode, setMode] = useState<Mode>("graph");
  const tab = (m: Mode, label: string) => (
    <button
      type="button"
      key={m}
      aria-pressed={mode === m}
      data-testid={`toggle-${m}`}
      onClick={() => setMode(m)}
      className={`px-2 py-0.5 text-xs ${mode === m ? "bg-primary text-primary-foreground" : "text-muted-foreground hover:text-foreground"}`}
    >
      {label}
    </button>
  );
  return (
    <div className="flex h-full flex-col">
      <div className="mb-1 inline-flex w-fit overflow-hidden rounded border" role="group" aria-label="World view">
        {tab("graph", "Graph")}
        {tab("text", "Relationships")}
      </div>
      {/* The graph stays mounted while hidden so a selection survives switching views. */}
      <div className="min-h-0 flex-1 overflow-auto" hidden={mode !== "graph"}>
        <WorldGraphPanel text={text} dark={dark} initial={initial} bound={marks} />
      </div>
      {mode === "text" && (
        <div className="min-h-0 flex-1">
          <CodeEditor value={text} readOnly language="plaintext" dark={dark} />
        </div>
      )}
    </div>
  );
}
