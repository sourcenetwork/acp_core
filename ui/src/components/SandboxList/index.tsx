import { PersistedSandboxData, usePlaygroundStore } from "@/stores/playgroundStore";
import { cn } from "@/utils/classnames";
import { Edit2, MoreHorizontal, Trash } from "lucide-react";
import { MouseEvent, useState } from "react";
import DialogConfirm from "../DialogConfirm";
import DialogEditSandbox from "../DialogEditSandbox";
import { Button } from "../ui/button";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "../ui/dropdown-menu";

interface SandboxListProps {
    className?: string;
    collpased?: boolean;
    format: 'small' | 'expanded';
    onSandboxClick?: (id: string) => unknown;
}

const SandboxList = (props: SandboxListProps) => {
    const { className, format, onSandboxClick } = props;
    const [showEditSandbox, setShowEditSandbox] = useState(false);
    const [selectedSandbox, setSelectedSandbox] = useState<PersistedSandboxData | null>(null);
    const [showConfirmDelete, setShowConfirmDelete] = useState(false);

    const sandboxes = usePlaygroundStore((state) => state.sandboxes);
    const lastActiveId = usePlaygroundStore((state) => state.lastActiveId);
    const setActiveSandbox = usePlaygroundStore((state) => state.setActiveSandbox);
    const deleteStoredSandbox = usePlaygroundStore((state) => state.deleteStoredSandbox);

    const handleEditSandbox = (sandbox: PersistedSandboxData) => (event: MouseEvent) => {
        event.preventDefault();
        event.stopPropagation();

        setSelectedSandbox(sandbox);
        setShowEditSandbox(true);
    };

    const handleDeleteSandbox = (sandbox: PersistedSandboxData) => (event: MouseEvent) => {
        event.preventDefault();
        event.stopPropagation();

        setSelectedSandbox(sandbox);
        setShowConfirmDelete(true);
    }

    const handleSandboxClick = (id: string) => () => {
        void setActiveSandbox(id);
        if (onSandboxClick != null) onSandboxClick(id);
    }

    const handleConfirmDelete = (confirmed: boolean) => {
        if (confirmed && selectedSandbox) void deleteStoredSandbox(selectedSandbox.id);
        setShowConfirmDelete(false);
        setSelectedSandbox(null);
    }

    return <>
        <DialogEditSandbox
            open={showEditSandbox}
            sandboxId={selectedSandbox?.id}
            setOpen={(state) => setShowEditSandbox(state)} />

        <ul className={cn("space-y-1 w-full", className)}>

            {sandboxes?.map((sandbox) => (
                <div key={sandbox.id}
                    className={cn(`flex px-2 py-1 group 
                    border-l-2 border-transparent text-left 
                    hover:bg-accent hover:text-accent-foreground text-xs transition-colors`,
                        {
                            "border-src-secondary bg-secondary/50": lastActiveId === sandbox.id,
                            "py-3": format === 'expanded'
                        }
                    )}
                    onClick={handleSandboxClick(sandbox.id)}>

                    <button className="block w-full text-left">
                        <div className={
                            cn("text-xs leading-none overflow-hidden whitespace-nowrap text-ellipsis max-w-[180px] mr-2", {
                                "text-sm leading-normal max-w-full": format === 'expanded',
                            })
                        }>{sandbox.name}</div>

                        {format === 'expanded' &&
                            <div className="text-[12px] opacity-60 leading-relaxed mr-2">{sandbox.description || '-'}</div>}
                    </button>

                    <DropdownMenu>
                        <DropdownMenuTrigger asChild>
                            <Button variant={"outline"} size={"iconXs"} className="opacity-100 group-hover:opacity-100 md:opacity-0" >
                                <MoreHorizontal size={15} />
                            </Button>
                        </DropdownMenuTrigger>
                        <DropdownMenuContent>
                            <DropdownMenuItem onClick={handleEditSandbox(sandbox)}><Edit2 size={15} /> Edit</DropdownMenuItem>
                            <DropdownMenuItem onClick={handleDeleteSandbox(sandbox)}><Trash size={15} /> Delete</DropdownMenuItem>
                        </DropdownMenuContent>
                    </DropdownMenu>
                </div>
            ))}
        </ul>

        <DialogConfirm
            title={`Delete Sandbox`}
            description={`Are you sure you want to delete sandbox: (${selectedSandbox?.name})? This action cannot be undone.`}
            open={showConfirmDelete}
            setOpen={setShowConfirmDelete}
            onConfirm={handleConfirmDelete}
        />
    </>
}

export default SandboxList;