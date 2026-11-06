import { useState } from "react";
import { Check, Copy } from "lucide-react";
import { worldFileName, worldToYaml } from "@/lib/workspaceFiles";
import { clampWorldIndex, worldLabel } from "@/lib/worlds";
import { useWorkspace } from "@/stores/workspaceStore";
import type { World } from "@/lib/types";
import { WorldView } from "./WorldView";

/** File-backed mode: worlds are files, so instead of a save button, where the world lives or how to put it there. */
function WorldSource({ world, dir }: { world: World; dir: string | null }) {
  const [copied, setCopied] = useState(false);
  const copy = async () => {
    try {
      await navigator.clipboard.writeText(worldToYaml(world));
      setCopied(true);
      setTimeout(() => setCopied(false), 1500);
    } catch {
      /* clipboard unavailable (insecure context or denied) */
    }
  };
  if (world.kind === "saved") {
    return <span className="font-mono text-xs text-muted-foreground" data-testid="world-source">{dir}/{worldFileName(world.name ?? "")}</span>;
  }
  return (
    <span className="flex min-w-0 flex-wrap items-center gap-x-2 gap-y-1 text-xs text-muted-foreground" data-testid="world-source">
      <button type="button" onClick={() => void copy()} data-testid="copy-world-yaml"
        className="flex shrink-0 items-center gap-1 whitespace-nowrap rounded border px-1.5 py-0.5 hover:text-foreground">
        {copied ? <Check className="h-3 w-3" /> : <Copy className="h-3 w-3" />}
        {copied ? "Copied" : "Copy as YAML"}
      </button>
      <span className="min-w-0 break-all">
        To save it: put it in <code className="rounded bg-muted px-1">{dir}/worlds/</code>, or run{" "}
        <code className="rounded bg-muted px-1">acp generate -o {dir}/worlds {dir}</code>
      </span>
    </span>
  );
}

export function WorldsView({ dark }: { dark: boolean }) {
  const worlds = useWorkspace((s) => s.worlds);
  const generating = useWorkspace((s) => s.generating);
  const selected = useWorkspace((s) => s.selectedWorld);
  const setSelected = useWorkspace((s) => s.setSelectedWorld);
  const fileBacked = useWorkspace((s) => s.mode === "files");
  const dir = useWorkspace((s) => s.dir);

  if (worlds.length === 0) {
    return (
      <p className="p-4 text-sm text-muted-foreground" data-testid="no-worlds">
        {generating ? "Generating..." : fileBacked ? `Generate to see worlds, or add files to ${dir}/worlds/.` : "Generate to see worlds."}
      </p>
    );
  }
  const index = clampWorldIndex(selected, worlds.length);
  const w = worlds[index];
  return (
    <div className="flex h-full flex-col gap-2 p-2" data-testid="worlds">
      <div className="flex flex-wrap items-center gap-2 text-sm">
        <label htmlFor="world-select" className="text-muted-foreground">World</label>
        <select
          id="world-select"
          data-testid="world-select"
          value={index}
          onChange={(e) => setSelected(Number(e.target.value))}
          className="h-8 rounded-md border border-input bg-background px-2 text-sm"
        >
          {worlds.map((x, i) => (
            <option key={x.id} value={i}>{worldLabel(i, x.kind, x.relationships, x.name)}</option>
          ))}
        </select>
        <span className="text-xs text-muted-foreground">{worlds.length} {worlds.length === 1 ? "world" : "worlds"}</span>
        {fileBacked && <WorldSource world={w} dir={dir} />}
      </div>
      {/* keyed by world: switching worlds drops the object/permission selection and the explain panel */}
      <section key={w.id} data-testid="world" className="min-h-0 flex-1">
        <WorldView text={w.relationships} dark={dark} />
      </section>
    </div>
  );
}
