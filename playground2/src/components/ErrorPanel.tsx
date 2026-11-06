import { AlertCircle, WifiOff } from "lucide-react";
import type { GenerateError } from "@/stores/workspaceStore";

const TITLES: Record<GenerateError["kind"], string> = {
  constraints: "Invalid system constraints",
  requirements: "Invalid requirements",
  backend: "The backend rejected the system model",
  network: "Network error",
  protocol: "Unexpected backend response",
};

export function ErrorPanel({ error }: { error: GenerateError }) {
  const Icon = error.kind === "network" ? WifiOff : AlertCircle;
  return (
    <div role="alert" data-testid="error-panel" data-kind={error.kind}
      className="rounded-md border border-destructive/50 bg-destructive/10 p-3 text-sm">
      <div className="mb-1 flex items-center gap-2 font-medium text-destructive">
        <Icon className="h-4 w-4" />
        {TITLES[error.kind]}
        {error.kind === "backend" && error.code && (
          <code className="rounded bg-muted px-1 text-xs text-muted-foreground">{error.code}</code>
        )}
      </div>
      {error.kind === "constraints" || error.kind === "requirements" ? (
        <ul className="list-disc pl-5">
          {error.errors.map((e, i) => (
            <li key={i}>
              {error.file !== undefined && <span className="font-mono text-muted-foreground">{error.file} </span>}
              {e.line !== undefined && <span className="text-muted-foreground">line {e.line}: </span>}
              {e.message}
            </li>
          ))}
        </ul>
      ) : (
        <pre className="whitespace-pre-wrap break-words font-mono text-xs">{error.message}</pre>
      )}
    </div>
  );
}
