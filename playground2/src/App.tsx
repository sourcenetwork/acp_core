import { useEffect } from "react";
import { Moon, Sun } from "lucide-react";
import { ErrorPanel } from "@/components/ErrorPanel";
import { GenerateBar } from "@/components/GenerateBar";
import { ConstraintsTab, PolicyTab } from "@/components/SystemModelTab";
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { RequirementsTab } from "@/components/RequirementsTab";
import { WorldsView } from "@/components/WorldsView";
import { WorkspaceBadge, WorkspaceProblems } from "@/components/WorkspaceStatus";
import { useThemeState } from "@/theme";
import { useWorkspace, type MainTab } from "@/stores/workspaceStore";

export default function App() {
  const { theme, toggle } = useThemeState();
  const dark = theme === "dark";
  const error = useWorkspace((s) => s.error);
  const activeTab = useWorkspace((s) => s.activeTab);
  const setActiveTab = useWorkspace((s) => s.setActiveTab);
  const worldCount = useWorkspace((s) => s.worlds.length);
  const mode = useWorkspace((s) => s.mode);
  const initWorkspace = useWorkspace((s) => s.initWorkspace);
  useEffect(() => {
    void initWorkspace();
  }, [initWorkspace]);

  if (mode === "pending") return <p className="p-4 text-sm text-muted-foreground" data-testid="detecting">Loading...</p>;
  return (
    <div className="flex h-dvh flex-col" data-mode={mode}>
      <header className="flex items-center justify-between border-b px-4 py-2">
        <h1 className="shrink-0 text-sm font-semibold">ACP Playground 2</h1>
        <div className="mx-4 min-w-0 flex-1"><WorkspaceBadge /></div>
        <Button variant="ghost" size="iconSm" onClick={toggle} aria-label="Toggle theme">
          {dark ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
        </Button>
      </header>
      <div className="flex flex-col gap-2 border-b px-4 py-2">
        <GenerateBar />
        <WorkspaceProblems />
        {error && <ErrorPanel error={error} />}
      </div>
      <Tabs value={activeTab} onValueChange={(v) => setActiveTab(v as MainTab)} className="flex min-h-0 flex-1 flex-col">
        <TabsList className="justify-start border-b">
          <TabsTrigger value="policy">Policy</TabsTrigger>
          <TabsTrigger value="constraints">System constraints</TabsTrigger>
          <TabsTrigger value="requirements">Requirements</TabsTrigger>
          <TabsTrigger value="worlds">Worlds{worldCount > 0 ? ` (${worldCount})` : ""}</TabsTrigger>
        </TabsList>
        <TabsContent value="policy" className="mt-0 min-h-0 flex-1"><PolicyTab dark={dark} /></TabsContent>
        <TabsContent value="constraints" className="mt-0 min-h-0 flex-1"><ConstraintsTab dark={dark} /></TabsContent>
        {/* Kept mounted so the open counterexample survives a trip to the other tabs. */}
        <TabsContent value="requirements" forceMount className="mt-0 min-h-0 flex-1 data-[state=inactive]:hidden"><RequirementsTab dark={dark} /></TabsContent>
        {/* Kept mounted so a world's selection survives a trip to the other tabs. */}
        <TabsContent value="worlds" forceMount className="mt-0 min-h-0 flex-1 data-[state=inactive]:hidden"><WorldsView dark={dark} /></TabsContent>
      </Tabs>
    </div>
  );
}
