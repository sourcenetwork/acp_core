import { useMemo } from "react";
import { parseConstraints } from "@/lib/constraints";
import { useWorkspace } from "@/stores/workspaceStore";
import { CodeEditor } from "./CodeEditor";
import { FileNotice } from "./WorkspaceStatus";

export function PolicyTab({ dark }: { dark: boolean }) {
  const policy = useWorkspace((s) => s.policy);
  const setPolicy = useWorkspace((s) => s.setPolicy);
  const readOnly = useWorkspace((s) => s.mode === "files");
  return (
    <div className="flex h-full flex-col p-2">
      <FileNotice file="policy.yaml" />
      <div className="min-h-0 flex-1"><CodeEditor value={policy} onChange={setPolicy} readOnly={readOnly} dark={dark} testId="policy-editor" /></div>
    </div>
  );
}

export function ConstraintsTab({ dark }: { dark: boolean }) {
  const constraints = useWorkspace((s) => s.constraints);
  const setConstraints = useWorkspace((s) => s.setConstraints);
  const readOnly = useWorkspace((s) => s.mode === "files");

  // Live client-side validation markers for the constraints editor.
  const markers = useMemo(() => {
    const r = parseConstraints(constraints);
    return r.ok ? [] : r.errors;
  }, [constraints]);

  return (
    <div className="flex h-full flex-col p-2">
      <FileNotice file="constraints.yaml" />
      <div className="min-h-0 flex-1"><CodeEditor value={constraints} onChange={setConstraints} readOnly={readOnly} dark={dark} markers={markers} testId="constraints-editor" /></div>
    </div>
  );
}
