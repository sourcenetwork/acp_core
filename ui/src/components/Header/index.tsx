import { useActiveSandbox } from "@/hooks/useActiveSandbox";
import { usePlaygroundStore } from "@/lib/playgroundStore";
import { SandboxData, SandboxTemplate } from "@acp/sandbox";
import { Plus } from "lucide-react";
import { useRef, useState } from "react";
import ConfirmationDialog from "../ConfirmationDialog";
import CreateSandboxDialog from "../CreateSandboxDialog";
import SandboxList from "../SandboxList";
import StatusIndicator from "../StatusIndicator";
import ThemeToggle from "../ThemeToggle";
import { Button } from "../ui/button";
import { NavigationMenu, NavigationMenuContent, NavigationMenuItem, NavigationMenuList, NavigationMenuTrigger } from "../ui/navigation-menu";

const Header = () => {
    const activeSandbox = useActiveSandbox();

    const [playgroundStatus, setState, samples, loadTemplate] = usePlaygroundStore((state) => [
        state.playgroundStatus,
        state.setPlaygroundState,
        state.sandboxTemplates,
        state.loadTemplate,
    ]);

    const [selectedSample, setSelectedSample] = useState<SandboxTemplate | null>(null);
    const [showConfirmation, setShowConfirmation] = useState(false);
    const [showCreateSandbox, setShowCreateSandbox] = useState(false);

    const fileInputRef = useRef<HTMLInputElement>(null);

    const computeHash = async (data: string) => {
        const encodedData = new TextEncoder().encode(data);
        // TODO - this isn't working
        const buffer = await crypto.subtle.digest("SHA-256", encodedData);
        return Array.from(new Uint8Array(buffer))
            .map((b) => b.toString(16).padStart(2, "0"))
            .join("").substring(0, 8);
    };

    const exportState = async () => {
        try {
            const activeStateData = activeSandbox?.data;
            const jsonBlob = new Blob([JSON.stringify(activeStateData, null, 2)], { type: "application/json" });
            const url = URL.createObjectURL(jsonBlob);
            const link = document.createElement("a");
            const filename = await computeHash(JSON.stringify(jsonBlob));
            link.href = url;
            link.download = `acp-playground-export-${filename}.json`;
            link.click();
            URL.revokeObjectURL(url);
        } catch (error) {
            // TODO
            console.error(error);
        }
    };

    const handleImportFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = () => {
                const parsedData: Partial<SandboxData> = JSON.parse(reader.result as string);
                void setState(parsedData);
            };
            reader.readAsText(file);
        }
    };

    const handleExportButtonClick = () => {
        void exportState();
    };

    const handleImportButtonClick = () => {
        if (!fileInputRef.current) return;
        fileInputRef.current.value = "";
        fileInputRef.current.click();
    };


    const handleNewSandboxButtonClick = () => {
        setShowCreateSandbox(true);
    };

    const handleSampleClick = (sample: SandboxTemplate) => () => {
        setSelectedSample(sample);
        setShowConfirmation(true);
    }

    return <div className="flex justify-between px-4 py-2">
        <div>
            <span className="opacity-40">ACP Playground</span>
            <span className="inline-block w-[30px] h-[1px] bg-border align-middle mx-4"></span>

            <CreateSandboxDialog
                open={showCreateSandbox}
                setOpen={(state) => setShowCreateSandbox(state)} />

            <ConfirmationDialog
                title={`Load Policy - ${selectedSample?.name ?? ""}`}
                description={`Are you sure you want to override current schema with ${selectedSample?.name ?? ""}`}
                open={showConfirmation}
                onAction={async (confirmed: boolean) => {
                    try {
                        if (confirmed === true && selectedSample) await loadTemplate(selectedSample);
                        setShowConfirmation(false);
                    } catch (error) {

                        console.error(error); // TODO
                    }
                }} />

            <NavigationMenu className="inline-block" delayDuration={0} >
                <NavigationMenuList>

                    <NavigationMenuItem>
                        <Button variant={"ghost"} className="w-full border-2 text-xs" onClick={handleNewSandboxButtonClick}>
                            <Plus size={20} className="mr-1" /> New Sandbox
                        </Button>
                    </NavigationMenuItem>

                    <NavigationMenuItem>
                        <NavigationMenuTrigger className="text-xs">Sandboxes</NavigationMenuTrigger>
                        <NavigationMenuContent>
                            <SandboxList />
                        </NavigationMenuContent>
                    </NavigationMenuItem>

                    <NavigationMenuItem>
                        <NavigationMenuTrigger className="text-xs">Samples</NavigationMenuTrigger>
                        <NavigationMenuContent>
                            <ul className="p-3 w-[200px] md:w-[500px]">
                                {samples?.map((sample) => (
                                    <button key={sample.name} className="p-3 rounded-md text-left hover:bg-accent hover:text-accent-foreground w-full" onClick={handleSampleClick(sample)}>
                                        <div className="text-[13px] font-medium leading-none mb-1">{sample.name}</div>
                                        <p className="line-clamp-2 text-[12px] leading-snug text-muted-foreground">{sample.description}</p>
                                    </button>
                                ))}
                            </ul>
                        </NavigationMenuContent>
                    </NavigationMenuItem>

                </NavigationMenuList>
            </NavigationMenu>
        </div>

        <div className="items-center gap-x-2 hidden md:flex">
            <StatusIndicator status={playgroundStatus} />
            <Button className="text-xs" variant="outline" size={"xs"} onClick={handleImportButtonClick}>Import</Button>

            <input
                ref={fileInputRef}
                type="file"
                accept=".json"
                onChange={handleImportFileChange}
                className="hidden"
            />

            <Button className="text-xs" variant="default" size={"xs"} onClick={handleExportButtonClick}>Export</Button>
            <ThemeToggle />
        </div>

    </div>
}

export default Header;

