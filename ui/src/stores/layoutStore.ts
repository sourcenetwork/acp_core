import { Rectangle } from "@dnd-kit/geometry";
import { v4 } from "uuid";
import { create } from "zustand";
import { persist, subscribeWithSelector } from "zustand/middleware";
import { useShallow } from "zustand/react/shallow";

export type DropPosition = "left" | "center" | "right";

export enum DragType {
  TabDrag = "tab-drag",
  TabDrop = "tab-drop",
  PaneDrop = "pane-drop",
  TabGroupDrop = "tab-group-drop",
}

export interface TabDragData {
  type: DragType.TabDrag;
  paneId: string;
  tabId: string;
  index: number;
}

export interface TabDropData {
  type: DragType.TabDrop;
  paneId: string;
  tabId: string;
  index: number;
  dropPosition: DropPosition | null;
}

export interface TabGroupDropData {
  type: DragType.TabGroupDrop;
  paneId: string;
}

export interface PaneDropData {
  type: DragType.PaneDrop;
  paneId: string;
  position: DropPosition | null;
}

export type DragData =
  | TabDragData
  | TabDropData
  | PaneDropData
  | TabGroupDropData;

export interface TabDefinition {
  key: string;
  label: string;
}

export interface Pane {
  id: string;
  tabs: TabDefinition[];
  activeTabKey: string;
}

export enum PaneTypes {
  Policy = "policy",
  Relationship = "relationship",
  Tests = "tests",
}

export enum SecondaryPaneTypes {
  Check = "check",
  Expand = "expand",
}

interface UIState {
  dropTarget: DragData | null;
  dropPosition: DropPosition | null;
  dropTargetIndex: { index: number; paneId: string } | null;
  secondaryPaneOpen: boolean;
  secondaryPaneType: SecondaryPaneTypes | null;
  sandboxMenuOpen: boolean;
  createSandboxDialogOpen: boolean;
  focusedEditor: string | null;
}

interface UIActions {
  // Drag state management
  setDropTarget: (target: DragData | null) => void;
  setDropPosition: (position: DropPosition | null) => void;
  clearDragState: () => void;
  setDropTargetIndex: (data: TabDropData | TabGroupDropData | null) => void;

  // Drag event handlers
  handleDragMove: (event: any) => void;
  handleDragOver: (event: any) => void;
  handleDragEnd: (event: any) => void;

  // Panel actions
  toggleSecondaryPane: () => void;
  setSecondaryPaneOpen: (collapsed: boolean) => void;
  setSecondaryPaneType: (type: SecondaryPaneTypes | null) => void;
  setSandboxMenuOpen: (collapsed: boolean) => void;
  setCreateSandboxDialogOpen: (open: boolean) => void;
  setFocusedEditor: (editor: string | null) => void;
}

const getDragData = (
  entity: { data: unknown } | null | undefined
): DragData | null => {
  return (entity?.data as DragData) || null;
};

export const editorTabs: TabDefinition[] = [
  {
    key: PaneTypes.Policy,
    label: "Policy",
  },
  {
    key: PaneTypes.Relationship,
    label: "Relationships",
  },
  {
    key: PaneTypes.Tests,
    label: "Tests",
  },
];

interface PaneState {
  panes: Pane[];
}

interface PaneActions {
  setActiveTab: (paneId: string, tabKey: string) => void;
  splitActivePane: (paneId: string) => void;
  moveTabToPane: (tabKey: string, fromPaneId: string, toPaneId: string) => void;
  moveTab: (
    sourcePaneId: string,
    targetPaneId: string,
    sourceTabKey: string,
    direction: DropPosition | null
  ) => void;
  sortTab: (
    paneId: string,
    tabKey: string,
    targetPaneId: string,
    targetIndex: number,
    targetPosition: DropPosition | null
  ) => void;
}

type LayoutStore = PaneState & PaneActions & UIState & UIActions;

const initialState: PaneState & UIState = {
  panes: [
    {
      id: v4(),
      activeTabKey: PaneTypes.Policy,
      tabs: editorTabs,
    },
  ],
  dropTarget: null,
  dropPosition: null,
  dropTargetIndex: null,
  secondaryPaneOpen: false,
  secondaryPaneType: SecondaryPaneTypes.Check,
  sandboxMenuOpen: false,
  createSandboxDialogOpen: false,
  focusedEditor: null,
};

export const useLayoutStore = create<LayoutStore>()(
  subscribeWithSelector(
    persist(
      (set, get) => ({
        ...initialState,

        setActiveTab: (paneId: string, tabKey: string) => {
          set((state) => ({
            panes: state.panes.map((pane) =>
              pane.id === paneId ? { ...pane, activeTabKey: tabKey } : pane
            ),
          }));
        },

        moveTabToPane: (
          tabKey: string,
          fromPaneId: string,
          toPaneId: string
        ) => {
          // No-op if the tab is already in the pane
          if (fromPaneId === toPaneId) return;

          set((state) => {
            const tabToMove = state.panes
              .find((p) => p.id === fromPaneId)
              ?.tabs.find((t) => t.key === tabKey);

            if (!tabToMove) return state;

            return {
              panes: state.panes
                .map((pane) => {
                  // If the tab is in the source pane, remove it
                  if (pane.id === fromPaneId) {
                    const newTabs = pane.tabs.filter((t) => t.key !== tabKey);

                    return {
                      ...pane,
                      tabs: newTabs,
                      activeTabKey:
                        pane.activeTabKey === tabKey
                          ? newTabs[0]?.key || "policy"
                          : pane.activeTabKey,
                    };
                  }

                  // If the tab is in the target pane, add it
                  if (pane.id === toPaneId) {
                    return {
                      ...pane,
                      tabs: [...pane.tabs, tabToMove],
                      activeTabKey: tabKey,
                    };
                  }

                  return pane;
                })
                .filter((pane) => pane.tabs.length > 0),
            };
          });
        },

        moveTab: (
          sourcePaneId: string,
          targetPaneId: string,
          sourceTabKey: string,
          direction: DropPosition | null
        ) => {
          const { panes } = get();
          const fromPane = panes.find((p) => p.id === sourcePaneId);
          const tabToMove = fromPane?.tabs.find((t) => t.key === sourceTabKey);

          if (!tabToMove || !fromPane) return;

          const fromPaneIndex = panes.findIndex((p) => p.id === sourcePaneId);
          const targetPaneIndex = panes.findIndex((p) => p.id === targetPaneId);
          const isSameTarget = fromPaneIndex === targetPaneIndex;
          const isLastTab = fromPane?.tabs.length === 1;

          if (direction === "center") {
            get().moveTabToPane(sourceTabKey, sourcePaneId, targetPaneId);
            return;
          }

          const sameEventualPaneTarget = isLastTab && isSameTarget;
          if (sameEventualPaneTarget) return;

          const newPaneId = v4();

          set((state) => {
            const updatedPanes = state.panes.map((pane) =>
              pane.id === sourcePaneId
                ? {
                    ...pane,
                    tabs: pane.tabs.filter((t) => t.key !== sourceTabKey),
                    activeTabKey:
                      pane.activeTabKey === sourceTabKey
                        ? pane.tabs.filter((t) => t.key !== sourceTabKey)[0]
                            ?.key || "policy"
                        : pane.activeTabKey,
                  }
                : pane
            );

            const newPane: Pane = {
              id: newPaneId,
              tabs: [tabToMove],
              activeTabKey: sourceTabKey,
            };

            const insertIndex =
              direction === "left" ? targetPaneIndex : targetPaneIndex + 1;
            updatedPanes.splice(insertIndex, 0, newPane);

            return {
              panes: updatedPanes.filter((pane) => pane.tabs.length > 0),
            };
          });
        },

        sortTab: (
          paneId: string,
          tabKey: string,
          targetPaneId: string,
          targetIndex: number,
          targetPosition: DropPosition | null
        ) => {
          set((state) => {
            const sourcePane = state.panes.find((p) => p.id === paneId);
            const tabToMove = sourcePane?.tabs.find((t) => t.key === tabKey);

            if (!tabToMove || !sourcePane) return state;

            let insertIndex = targetIndex;
            if (targetPosition === "right") {
              insertIndex = targetIndex + 1;
            }

            if (paneId === targetPaneId) {
              return {
                panes: state.panes.map((pane) => {
                  if (pane.id !== paneId) return pane;

                  const tabs = [...pane.tabs];
                  const currentIndex = tabs.findIndex((t) => t.key === tabKey);

                  if (currentIndex === -1) return pane;

                  const [movedTab] = tabs.splice(currentIndex, 1);

                  let finalInsertIndex = insertIndex;
                  if (currentIndex < insertIndex) {
                    finalInsertIndex = insertIndex - 1;
                  }

                  finalInsertIndex = Math.max(
                    0,
                    Math.min(finalInsertIndex, tabs.length)
                  );
                  tabs.splice(finalInsertIndex, 0, movedTab);

                  return {
                    ...pane,
                    tabs,
                    activeTabKey: tabKey,
                  };
                }),
              };
            }

            return {
              panes: state.panes
                .map((pane) => {
                  if (pane.id === paneId) {
                    const newTabs = pane.tabs.filter((t) => t.key !== tabKey);

                    return {
                      ...pane,
                      tabs: newTabs,
                      activeTabKey:
                        pane.activeTabKey === tabKey
                          ? newTabs[0]?.key || "policy"
                          : pane.activeTabKey,
                    };
                  }

                  if (pane.id === targetPaneId) {
                    const tabs = [...pane.tabs];
                    const finalInsertIndex = Math.max(
                      0,
                      Math.min(insertIndex, tabs.length)
                    );
                    tabs.splice(finalInsertIndex, 0, tabToMove);

                    return {
                      ...pane,
                      tabs,
                      activeTabKey: tabKey,
                    };
                  }

                  return pane;
                })
                .filter((pane) => pane.tabs.length > 0),
            };
          });
        },

        splitActivePane: (paneId: string) => {
          const state = get();
          const activePane = state.panes.find((p) => p.id === paneId);
          if (!activePane) return;

          get().moveTab(paneId, paneId, activePane.activeTabKey, "right");
        },

        setDropTarget: (target) => set({ dropTarget: target }),
        setDropPosition: (position) => set({ dropPosition: position }),

        setDropTargetIndex: (data) => {
          const { dropTarget, dropPosition } = get();

          if (!data) {
            set({ dropTargetIndex: null });
            return;
          }

          const groupDrop = data.type === DragType.TabGroupDrop;
          const isSamePane = data.paneId === dropTarget?.paneId;
          const dropIndex = (dropTarget as TabDropData)?.index || 0;

          if (!isSamePane) set({ dropTargetIndex: null });

          const targetIndex = groupDrop
            ? 100
            : dropIndex - 1 + (dropPosition === "right" ? 1 : 0);

          set({ dropTargetIndex: { index: targetIndex, paneId: data.paneId } });
        },
        clearDragState: () =>
          set({
            dropTarget: null,
            dropPosition: null,
            dropTargetIndex: null,
          }),

        handleDragMove: (event) => {
          const { setDropPosition, setDropTargetIndex } = get();

          const { to, operation } = event;
          const target = operation.target;
          const targetData = getDragData(target);
          const targetShape = target?.shape;

          switch (targetData?.type) {
            case DragType.TabDrop:
              if (targetShape && to && targetShape instanceof Rectangle) {
                const dragX = to.x;
                const targetLeft = targetShape.left;
                const targetWidth = targetShape.width;
                const leftBoundary = targetLeft + targetWidth * 0.5;
                const position = dragX < leftBoundary ? "left" : "right";
                setDropPosition(position);
                setDropTargetIndex(targetData);
              }
              break;
            case DragType.PaneDrop:
              if (target?.shape && to && target.shape instanceof Rectangle) {
                const dragX = to.x;
                const targetLeft = target.shape.left;
                const targetWidth = target.shape.width;
                const leftBoundary = targetLeft + targetWidth * 0.25;
                const rightBoundary = targetLeft + targetWidth * 0.75;

                const position =
                  dragX < leftBoundary
                    ? "left"
                    : dragX > rightBoundary
                    ? "right"
                    : "center";

                setDropPosition(position);
                setDropTargetIndex(null);
              }
              break;
            case DragType.TabGroupDrop:
              setDropTargetIndex(targetData);

              break;
          }
        },

        handleDragOver: (event) => {
          const targetData = getDragData(event.operation.target);
          get().setDropPosition(null);
          get().setDropTarget(targetData);
        },

        handleDragEnd: (event) => {
          const { operation, canceled } = event;
          const { clearDragState, sortTab, moveTab } = get();

          clearDragState();

          if (canceled) return;

          const sourceData = getDragData(operation?.source);
          const targetData = getDragData(operation?.target);

          if (
            !sourceData ||
            !targetData ||
            sourceData.type !== DragType.TabDrag
          )
            return;

          switch (targetData.type) {
            case DragType.TabDrop:
              sortTab(
                sourceData.paneId,
                sourceData.tabId,
                targetData.paneId,
                targetData.index,
                targetData.dropPosition
              );
              break;

            case DragType.PaneDrop:
              const position = targetData.position;
              const samePane = sourceData.paneId === targetData.paneId;
              if (samePane && position === "center") return;

              moveTab(
                sourceData.paneId,
                targetData.paneId,
                sourceData.tabId,
                position
              );
              break;

            case DragType.TabGroupDrop:
              const targetIndex = get().panes.find(
                (p) => p.id === targetData.paneId
              )?.tabs.length;

              if (targetIndex === undefined) return;

              sortTab(
                sourceData.paneId,
                sourceData.tabId,
                targetData.paneId,
                targetIndex,
                "left"
              );
              break;
          }
        },

        toggleSecondaryPane: () => {
          set((state) => ({ secondaryPaneOpen: !state.secondaryPaneOpen }));
        },

        setSecondaryPaneOpen: (state: boolean) => {
          set({ secondaryPaneOpen: state });
        },

        setSecondaryPaneType: (type: SecondaryPaneTypes | null) => {
          set({ secondaryPaneType: type });
        },

        setSandboxMenuOpen: (collapsed: boolean) => {
          set({ sandboxMenuOpen: collapsed });
        },

        setCreateSandboxDialogOpen: (open: boolean) => {
          set({ createSandboxDialogOpen: open });
        },

        setFocusedEditor: (editor: string | null) => {
          set({ focusedEditor: editor });
        },
      }),
      {
        name: "pane-layout-storage",
        partialize: (state) => ({
          panes: state.panes,
          secondaryPaneOpen: state.secondaryPaneOpen,
          secondaryPaneType: state.secondaryPaneType,
          sandboxMenuOpen: state.sandboxMenuOpen,
          createSandboxDialogOpen: state.createSandboxDialogOpen,
        }),
      }
    )
  )
);

export const usePanes = () => useLayoutStore((state) => state.panes);

export const usePaneActions = () =>
  useLayoutStore(
    useShallow((state) => ({
      setActiveTab: state.setActiveTab,
      moveTabToPane: state.moveTabToPane,
      moveTab: state.moveTab,
      sortTab: state.sortTab,
      splitActivePane: state.splitActivePane,
    }))
  );

export const useDragActions = () =>
  useLayoutStore(
    useShallow((state) => ({
      handleDragMove: state.handleDragMove,
      handleDragOver: state.handleDragOver,
      handleDragEnd: state.handleDragEnd,
    }))
  );

export const useUIState = () =>
  useLayoutStore(
    useShallow((state) => ({
      secondaryPaneOpen: state.secondaryPaneOpen,
      secondaryPaneType: state.secondaryPaneType,
      sandboxMenuOpen: state.sandboxMenuOpen,
      createSandboxDialogOpen: state.createSandboxDialogOpen,
      focusedEditor: state.focusedEditor,
    }))
  );

export const useUIActions = () =>
  useLayoutStore(
    useShallow((state) => ({
      toggleSecondaryPane: state.toggleSecondaryPane,
      setSecondaryPaneOpen: state.setSecondaryPaneOpen,
      setSecondaryPaneType: state.setSecondaryPaneType,
      setSandboxMenuOpen: state.setSandboxMenuOpen,
      setCreateSandboxDialogOpen: state.setCreateSandboxDialogOpen,
      setFocusedEditor: state.setFocusedEditor,
    }))
  );
