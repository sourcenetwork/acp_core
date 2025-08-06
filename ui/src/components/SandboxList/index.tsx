import { usePlaygroundStore } from "@/stores/playgroundStore";
import { cn } from "@/utils/classnames";
import { Edit2, MoreHorizontal, Trash } from "lucide-react";
import { MouseEvent, useState } from "react";
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
    const [selectedSandboxId, setSelectedSandboxId] = useState<string | null>(null);

    const sandboxes = usePlaygroundStore((state) => state.sandboxes);
    const lastActiveId = usePlaygroundStore((state) => state.lastActiveId);
    const setActiveSandbox = usePlaygroundStore((state) => state.setActiveSandbox);
    const deleteStoredSandbox = usePlaygroundStore((state) => state.deleteStoredSandbox);

    const handleEditSandbox = (sandboxId: string) => (event: MouseEvent) => {
        event.preventDefault();
        event.stopPropagation();
        setSelectedSandboxId(sandboxId);
        setShowEditSandbox(true);
    };

    const handleDeleteSandbox = (id: string) => (event: MouseEvent) => {
        event.preventDefault();
        event.stopPropagation();
        void deleteStoredSandbox(id);
    }

    const handleSandboxClick = (id: string) => () => {
        void setActiveSandbox(id);
        if (onSandboxClick != null) onSandboxClick(id);
    }

    return <>
        <DialogEditSandbox
            open={showEditSandbox}
            sandboxId={selectedSandboxId}
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
                            <DropdownMenuItem onClick={handleEditSandbox(sandbox.id)}><Edit2 size={15} /> Edit</DropdownMenuItem>
                            <DropdownMenuItem onClick={handleDeleteSandbox(sandbox.id)}><Trash size={15} /> Delete</DropdownMenuItem>
                        </DropdownMenuContent>
                    </DropdownMenu>
                </div>
            ))}
        </ul>
    </>
}

export default SandboxList;