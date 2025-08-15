import { DragData, DragType, DropPosition, PaneDropData, TabDragData, TabDropData, TabGroupDropData } from "@/stores/layoutStore";
import { cn } from "@/utils/classnames";
import { useDraggable, useDroppable } from "@dnd-kit/react";

export function DraggableTab({
    tabIndex,
    tabId,
    active,
    paneId,
    children,
    onClick,
    dropPosition,
}: {
    tabIndex: number;
    tabId: string;
    active: boolean;
    paneId: string;
    children: React.ReactNode;
    onClick: () => void;
    dropTarget: DragData | null;
    dropPosition: DropPosition | null;
    isFocused: boolean;
}) {
    const { ref: draggableRef, isDragging } = useDraggable<TabDragData>({
        id: `${paneId}-${tabId}-draggable`,
        type: DragType.TabDrag,
        data: {
            type: DragType.TabDrag,
            tabId: tabId,
            paneId: paneId,
            index: tabIndex,
        },
        feedback: "clone",
    });

    const { ref: droppableRef } = useDroppable<TabDropData>({
        id: `${paneId}-${tabId}-droppable`,
        type: DragType.TabDrop,
        collisionPriority: 5,
        data: {
            type: DragType.TabDrop,
            tabId: tabId,
            paneId: paneId,
            dropPosition: dropPosition,
            index: tabIndex,
        },
    });

    const handleClick = () => {
        if (!isDragging) onClick();
    };

    return (
        <div
            ref={droppableRef}
            style={{ order: tabIndex }}
            className={cn("text-sm border-b border-divider whitespace-nowrap inline-block relative",
                {
                    "bg-muted border-src-secondary text-foreground": active,
                }
            )}
        >
            <div ref={draggableRef} onClick={handleClick} className={cn("px-3 py-2")}>
                {children}
            </div>
        </div>
    );
}

export function DroppableTabGroup({
    children,
    paneId,
}: {
    children: React.ReactNode;
    paneId: string;
}) {
    const { ref } = useDroppable<TabGroupDropData>({
        id: `${paneId}-tab-group`,
        type: DragType.TabGroupDrop,
        collisionPriority: 1,
        data: {
            type: DragType.TabGroupDrop,
            paneId: paneId,
        },
    });

    return <div ref={ref}>{children}</div>;
}

export function DroppablePane({
    id,
    children,
    position,
    className,
}: {
    id: string;
    children: React.ReactNode;
    position: "left" | "center" | "right" | null;
    className?: string;
}) {
    const { isDropTarget, ref } = useDroppable<PaneDropData>({
        id: id,
        type: DragType.PaneDrop,
        collisionPriority: 1,
        data: {
            type: DragType.PaneDrop,
            paneId: id,
            position: position,
        },
    });

    return (
        <div ref={ref} className={className}>
            {children}
            <PaneDropIndicator visible={isDropTarget} position={position} />
        </div>
    );
}

interface PaneDropIndicatorProps {
    visible: boolean;
    position: DropPosition | null;
}

function PaneDropIndicator({ visible, position }: PaneDropIndicatorProps) {
    return (
        <div
            className={cn(
                "absolute top-0 bottom-0 z-10 bg-divider transition-all ease-out duration-300",
                {
                    "opacity-70 pointer-events-auto": visible,
                    "opacity-0 pointer-events-none": !visible,
                    "w-1/2": position === "left",
                    "w-full": position === "center",
                    "w-1/2 translate-x-full": position === "right",
                }
            )}
        />
    );
}

export function DropTabInsertionIndicator({
    visible,
    targetIndex = 0
}: {
    visible: boolean;
    targetIndex?: number;
}) {

    if (!visible) return null;

    return (
        <div
            className="self-stretch transition-all duration-200 ease-out relative z-1"
            style={{ order: targetIndex }}
        >
            <div className="absolute top-0 bottom-px w-px bg-primary  " />
        </div>
    );
}