import { useDragActions, useLayoutStore, usePanes, useUIActions } from "@/stores/layoutStore";
import { cn } from "@/utils/classnames";
import { DragDropProvider } from "@dnd-kit/react";
import { useEffect, useMemo, useRef, useState } from "react";
import { ImperativePanelHandle } from "react-resizable-panels";
import { useLocation, useNavigate } from "react-router-dom";
import DialogCreateSandbox from "../DialogCreateSandbox";
import DialogLoadShare from "../DialogLoadShare";
import { EditorPaneContainer } from "../EditorPaneContainer/EditorPaneContainer";
import Header from "../Header";
import Output from "../PaneOutput";
import SecondaryPaneContainer from "../SecondaryPaneContainer";
import SecondarySideMenu from "../SecondarySideMenu";
import SideMenu from "../SideMenu";
import {
  ResizableHandle,
  ResizablePanel,
  ResizablePanelGroup,
} from "../ui/resizable";
import { Toaster } from "../ui/toaster";

const RootLayout = () => {
  const showSecondaryMenu = true;

  const primaryPanelRef = useRef<ImperativePanelHandle>(null);
  const secondaryPanelRef = useRef<ImperativePanelHandle>(null);
  const location = useLocation();
  const navigate = useNavigate();
  const [shareDialogOpen, setShareDialogOpen] = useState(false);
  const queryParams = useMemo(() => new URLSearchParams(location.search), [location.search]);
  const shareId = queryParams.get("share");


  const createSandboxDialogOpen = useLayoutStore((state) => state.createSandboxDialogOpen);
  const secondaryPaneOpen = useLayoutStore((state) => state.secondaryPaneOpen);
  const { handleDragMove, handleDragOver, handleDragEnd } = useDragActions();
  const { setSecondaryPaneOpen, setCreateSandboxDialogOpen } = useUIActions();

  const panes = usePanes();

  const toggleSecondaryPane = (state: boolean) => {
    setSecondaryPaneOpen(state);
  };

  const clearShareParam = () => {
    const params = new URLSearchParams(window.location.search);
    params.delete("share");
    navigate({ pathname: window.location.pathname, search: params.toString() });
  };

  useEffect(() => {
    if (shareId) setShareDialogOpen(true);
  }, [shareId]);

  useEffect(() => {
    const pane = secondaryPanelRef.current;
    if (!pane) return;

    secondaryPaneOpen ? pane.expand() : pane.collapse();
  }, [secondaryPaneOpen]);


  const handleShareDialogOpen = (state: boolean) => {
    setShareDialogOpen(state);
    clearShareParam();
  };

  const handleCreateSandboxDialogOpen = (state: boolean) => {
    setCreateSandboxDialogOpen(state);
  };

  return (
    <div className="flex h-dvh flex-col">
      <Header />
      <div className="flex grow overflow-y-auto">
        <SideMenu />

        <div className="w-full min-w-0 ">

          {/* Primary Panel */}
          <ResizablePanelGroup
            direction="horizontal"
            autoSaveId={"primary-panel-layout"}
          >
            <ResizablePanel defaultSize={75} minSize={10} >
              <ResizablePanelGroup
                direction="vertical"
                autoSaveId={"primary-panel-layout"}
              >
                <ResizablePanel
                  ref={primaryPanelRef}
                  defaultSize={75}
                  minSize={10}
                >
                  <ResizablePanelGroup
                    direction="horizontal"
                    autoSaveId={"primary-editors"}
                  >
                    <DragDropProvider
                      onDragMove={handleDragMove}
                      onDragOver={handleDragOver}
                      onDragEnd={handleDragEnd}
                    >
                      {panes.map((pane, index) => (
                        <EditorPaneContainer
                          key={pane.id}
                          pane={pane}
                          index={index}
                          isLast={index === panes.length - 1}
                        />
                      ))}
                    </DragDropProvider>
                  </ResizablePanelGroup>
                </ResizablePanel>
                <ResizableHandle />
                <ResizablePanel defaultSize={25} minSize={10} className="px-2">
                  <Output />
                </ResizablePanel>
              </ResizablePanelGroup>
            </ResizablePanel>

            {/* Secondary Panel */}
            {showSecondaryMenu &&
              <>
                <ResizableHandle className="hidden md:flex" />
                <ResizablePanel defaultSize={40} minSize={10}
                  ref={secondaryPanelRef}
                  collapsible
                  collapsedSize={0}
                  onCollapse={() => toggleSecondaryPane(false)}
                  onExpand={() => toggleSecondaryPane(true)}
                  className={cn(
                    "md:flex flex-col h-full overflow-hidden",
                    "md:relative md:transform-none md:shadow-none md:z-auto",
                    secondaryPaneOpen && "fixed top-0 right-0 h-full w-10/12 border-l border-divider bg-background shadow-xl z-50"
                  )}>

                  <div className='relative z-10 h-full'>
                    <SecondaryPaneContainer />
                  </div>

                  <div
                    className={cn("fixed inset-0 bg-background opacity-50 z-1 md:hidden",
                      { "hidden": !secondaryPaneOpen })}
                    onClick={() => toggleSecondaryPane(false)}
                  />
                </ResizablePanel>
              </>
            }

          </ResizablePanelGroup>

        </div>

        {showSecondaryMenu && <SecondarySideMenu />}
      </div >

      <DialogLoadShare
        shareId={shareId}
        open={shareDialogOpen}
        setOpen={handleShareDialogOpen}
      />

      <DialogCreateSandbox
        open={createSandboxDialogOpen}
        setOpen={handleCreateSandboxDialogOpen}
      />

      <Toaster />
    </div >
  );
};


export default RootLayout;
