import { Loader2, Play } from "lucide-react";
import { MAX_CONFIG } from "@/lib/config";
import type { GenerateConfig } from "@/lib/types";
import { useWorkspace } from "@/stores/workspaceStore";
import { Button } from "@/components/ui/button";

const FIELDS: { key: keyof GenerateConfig; label: string; title: string }[] = [
  { key: "count", label: "Count", title: "Number of worlds" },
  { key: "scope", label: "Scope", title: "Entities per resource" },
  { key: "actors", label: "Actors", title: "Number of actors" },
];

export function GenerateBar() {
  const config = useWorkspace((s) => s.config);
  const setConfig = useWorkspace((s) => s.setConfig);
  const generating = useWorkspace((s) => s.generating);
  const generate = useWorkspace((s) => s.generate);

  return (
    <div className="flex flex-wrap items-center gap-3">
      {FIELDS.map((f) => (
        <label key={f.key} title={`${f.title} (1-${MAX_CONFIG[f.key]})`} className="flex items-center gap-1.5 text-sm">
          <span className="text-muted-foreground">{f.label}</span>
          <input
            type="number" min={1} max={MAX_CONFIG[f.key]} value={Number.isNaN(config[f.key]) ? "" : config[f.key]}
            aria-label={f.label}
            onChange={(e) => setConfig({ [f.key]: e.target.valueAsNumber })}
            className="h-8 w-16 rounded-md border border-input bg-background px-2 text-sm"
          />
        </label>
      ))}
      <Button size="sm" onClick={() => void generate()} disabled={generating} data-testid="generate">
        {generating ? <Loader2 className="mr-1.5 h-4 w-4 animate-spin" /> : <Play className="mr-1.5 h-4 w-4" />}
        Generate
      </Button>
    </div>
  );
}
