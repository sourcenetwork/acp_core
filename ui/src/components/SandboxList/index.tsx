import { usePlaygroundStore } from "@/lib/playgroundStore";
import { cn } from "@/utils/classnames";
import { Edit2, Trash } from "lucide-react";
import { MouseEvent, useState } from "react";
import EditSandboxDialog from "../EditSandboxDialog";
import TextTooltip from "../TextTooltip";
import { Button } from "../ui/button";

const SandboxList = () => {
    const [showEditSandbox, setShowEditSandbox] = useState(false);
    const [selectedSandboxId, setSelectedSandboxId] = useState<string | null>(null);

    const [sandboxes, lastActiveId, setActiveSandbox, deleteStoredSandbox] = usePlaygroundStore((state) => [
        state.sandboxes,
        state.lastActiveId,
        state.setActiveSandbox,
        state.deleteStoredSandbox
    ]);

    const handleEditSandbox = (sandboxId: string) => (event: MouseEvent) => {
        event.preventDefault();
        setSelectedSandboxId(sandboxId);
        setShowEditSandbox(true);
    };

    const handleDeleteSandbox = (id: string) => (event: MouseEvent) => {
        event.preventDefault();
        void deleteStoredSandbox(id);
    }

    const handleSandboxClick = (id: string) => () => {
        void setActiveSandbox(id);
    }

    return <>
        <EditSandboxDialog
            open={showEditSandbox}
            sandboxId={selectedSandboxId}
            setOpen={(state) => setShowEditSandbox(state)} />

        <ul className="flex flex-col gap-y-2 p-3 w-[200px] md:w-[500px]">
            {sandboxes?.map((sandbox) => (
                <div key={sandbox.id} className={cn(
                    "border-l-2 border-transparent group flex p-3 text-left hover:bg-accent hover:text-accent-foreground w-full text-xs transition-colors",
                    {
                        "border-green-500": lastActiveId === sandbox.id,
                    }
                )} onClick={handleSandboxClick(sandbox.id)}>
                    <button className="block w-full text-left">
                        <div className="text-[13px] font-medium leading-none mb-1">{sandbox.name}</div>
                        <p className="line-clamp-2 text-[12px] leading-snug text-muted-foreground">{sandbox?.description || " -"}</p>
                    </button>

                    <div className="ml-4 opacity-0 group-hover:opacity-100 whitespace-nowrap flex gap-x-1">
                        <TextTooltip content={"Edit"}>
                            <Button variant={"outline"} className="hover:border-accent-foreground" size={"iconSm"} onClick={handleEditSandbox(sandbox.id)}>
                                <Edit2 size={15} />
                            </Button>
                        </TextTooltip>
                        <TextTooltip content={"Delete"}>
                            <Button variant={"outline"} className="hover:border-destructive hover:text-destructive" size={"iconSm"} onClick={handleDeleteSandbox(sandbox.id)}>
                                <Trash size={15} />
                            </Button>
                        </TextTooltip>
                    </div>
                </div>
            ))}
        </ul>
    </>
}

export default SandboxList;