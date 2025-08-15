import PanePolicy from "@/components/PanePolicy";
import PaneRelationship from "@/components/PaneRelationship";
import PaneTests from "@/components/PaneTests";
import { Pane, useLayoutStore, usePaneActions } from "@/stores/layoutStore";
import { cn } from "@/utils/classnames";
import { Columns2, icons, LucideIcon } from "lucide-react";
import { ComponentType, Fragment, useMemo } from "react";
import { useShallow } from "zustand/react/shallow";
import { DraggableTab, DroppablePane, DroppableTabGroup, DropTabInsertionIndicator } from "../DragDropComponents/DragDropComponents";
import { SandboxType } from "../Editor";
import TextTooltip from "../TextTooltip";
import { Button } from "../ui/button";
import { ResizableHandle, ResizablePanel } from "../ui/resizable";
import { Tabs, TabsList } from "../ui/tabs";

export const tabComponentMap: Record<string, {
    editor: SandboxType,
    component: ComponentType,
    icon: LucideIcon,
}> = {
    policy: {
        editor: SandboxType.POLICY_DEFINITION,
        component: PanePolicy,
        icon: icons.BookText,
    },
    relationship: {
        editor: SandboxType.RELATIONSHIPS,
        component: PaneRelationship,
        icon: icons.Waypoints,
    },
    tests: {
        editor: SandboxType.POLICY_THEOREM,
        component: PaneTests,
        icon: icons.FlaskConical,
    },
};

interface EditorPaneContainerProps {
    pane: Pane;
    index: number;
    isLast: boolean;
}

export function EditorPaneContainer({
    pane,
    index,
    isLast,
}: EditorPaneContainerProps) {
    const dropTarget = useLayoutStore((state) => state.dropTarget);
    const dropPosition = useLayoutStore((state) => state.dropPosition);
    const dropTargetIndex = useLayoutStore(useShallow((state) => state.dropTargetIndex));
    const focusedEditor = useLayoutStore((state) => state.focusedEditor);

    const { splitActivePane, setActiveTab } = usePaneActions();

    const activePaneTab = pane.activeTabKey;
    const paneTabs = pane.tabs;
    const canSplitPane = paneTabs.length > 1;

    const { component: Component } = useMemo(() => tabComponentMap[activePaneTab], [activePaneTab]);

    const handleTabClick = (tabKey: string) => {
        setActiveTab(pane.id, tabKey);
    }

    const handleSplitPane = () => {
        splitActivePane(pane.id);
    }

    return (
        <Fragment>
            <ResizablePanel
                id={pane.id}
                order={index}
                minSize={20}
                className="flex flex-col"
            >
                <DroppableTabGroup paneId={pane.id}>
                    <Tabs defaultValue={activePaneTab} className="flex justify-between items-end overflow-auto relative scrollbar-hide border-b border-divider " >
                        <TabsList className="relative flex">
                            {paneTabs.map(({ key, label }, tabIndex) => {
                                const isActive = activePaneTab === key;
                                const { editor: EditorType, icon: Icon } = tabComponentMap[key];
                                const isFocused = focusedEditor === EditorType;

                                return (
                                    <DraggableTab
                                        key={key}
                                        paneId={pane.id}
                                        tabId={key}
                                        tabIndex={tabIndex}
                                        active={isActive}
                                        onClick={() => handleTabClick(key)}
                                        dropPosition={dropPosition}
                                        dropTarget={dropTarget}
                                        isFocused={isFocused}
                                    >
                                        {Icon && <Icon size={14} className="mr-2 inline" />}
                                        {label}
                                    </DraggableTab>
                                );
                            })}

                            <DropTabInsertionIndicator
                                visible={dropTargetIndex?.paneId === pane.id}
                                targetIndex={dropTargetIndex?.index}
                            />
                        </TabsList>

                        <div className="flex items-center gap-2 px-2">
                            {canSplitPane && (
                                <TextTooltip
                                    content="Split Pane"
                                    side="left"
                                    align="center"
                                >
                                    <Button
                                        variant="muted"
                                        size="iconSm"
                                        onClick={handleSplitPane}
                                        className=" hover:text-primary text-muted-foreground hover:bg-transparent"
                                    >
                                        <Columns2 size={16} />
                                    </Button>
                                </TextTooltip>
                            )}
                        </div>
                    </Tabs>
                </DroppableTabGroup>

                <DroppablePane
                    id={pane.id}
                    position={dropPosition}
                    className={cn("pane flex flex-col grow relative overflow-hidden px-2")}
                >
                    {Component && <Component />}
                </DroppablePane>
            </ResizablePanel>
            {!isLast && <ResizableHandle />}
        </Fragment>
    );
}