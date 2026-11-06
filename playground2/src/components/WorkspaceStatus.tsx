import { useMemo } from "react";
import { AlertCircle, Lock } from "lucide-react";
import { formatProblem, parseProblems } from "@/lib/workspaceFiles";
import { useWorkspace } from "@/stores/workspaceStore";

/** File-backed mode: the workspace directory and whether the event stream is up. */
export function WorkspaceBadge() {
  const mode = useWorkspace((s) => s.mode);
  const dir = useWorkspace((s) => s.dir);
  const connection = useWorkspace((s) => s.connection);
  if (mode !== "files") return null;
  const ok = connection === "connected";
  return (
    <div className="flex min-w-0 items-center gap-3 text-xs text-muted-foreground">
      {dir && <span className="truncate font-mono" title={dir} data-testid="workspace-dir">{dir}</span>}
      <span className="flex shrink-0 items-center gap-1.5" data-testid="connection" data-state={connection}
        title={ok ? "Receiving workspace changes from acp" : "Lost the connection to acp playground; retrying"}>
        <span className={`h-2 w-2 rounded-full ${ok ? "bg-emerald-500" : "animate-pulse bg-amber-500"}`} />
        {ok ? "connected" : connection === "connecting" ? "connecting" : "reconnecting"}
      </span>
    </div>
  );
}

/** Read-only notice above an editor showing a workspace file. */
export function FileNotice({ file }: { file: string }) {
  const mode = useWorkspace((s) => s.mode);
  const dir = useWorkspace((s) => s.dir);
  if (mode !== "files") return null;
  return (
    <p className="mb-1 flex items-center gap-1.5 text-xs text-muted-foreground" data-testid="file-notice">
      <Lock className="h-3 w-3 shrink-0" />
      <span>Read-only. Editing <code className="rounded bg-muted px-1">{dir ? `${dir}/${file}` : file}</code>? Changes appear here automatically.</span>
    </p>
  );
}

/**
 * File-backed mode: why the workspace can't be used as is, each problem naming its file. The files keep showing the
 * last good state while the workspace can't be read.
 */
export function WorkspaceProblems() {
  const mode = useWorkspace((s) => s.mode);
  const workspaceError = useWorkspace((s) => s.workspaceError);
  const constraints = useWorkspace((s) => s.constraints);
  const requirements = useWorkspace((s) => s.requirements);
  const worldProblems = useWorkspace((s) => s.worldProblems);
  const haveFiles = useWorkspace((s) => s.workspaceVersion !== null);
  const problems = useMemo(() => [...parseProblems(constraints, requirements), ...worldProblems], [constraints, requirements, worldProblems]);
  if (mode !== "files" || (!workspaceError && problems.length === 0)) return null;
  return (
    <div role="alert" data-testid="workspace-problems" className="rounded-md border border-destructive/50 bg-destructive/10 p-3 text-sm">
      {workspaceError && (
        <div data-testid="workspace-error">
          <div className="mb-1 flex items-center gap-2 font-medium text-destructive">
            <AlertCircle className="h-4 w-4" />
            Can't read the workspace
          </div>
          <pre className="whitespace-pre-wrap break-words font-mono text-xs">{workspaceError}</pre>
          <p className="mt-1 text-xs text-muted-foreground">
            {haveFiles ? "Showing the last good files; " : ""}This clears when the workspace is readable again.
          </p>
        </div>
      )}
      {problems.length > 0 && (
        <div data-testid="file-problems" className={workspaceError ? "mt-2" : ""}>
          <div className="mb-1 flex items-center gap-2 font-medium text-destructive">
            <AlertCircle className="h-4 w-4" />
            {problems.length === 1 ? "A workspace file has a problem" : "Workspace files have problems"}
          </div>
          <ul className="list-disc pl-5 font-mono text-xs">
            {problems.map((p, i) => <li key={i}>{formatProblem(p)}</li>)}
          </ul>
        </div>
      )}
    </div>
  );
}
