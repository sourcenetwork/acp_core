import { usePlaygroundStore } from "@/lib/acpHandler";
import { AcpSample, allSamples } from "@/lib/acpSamples";
import { usePlaygroundStorageStore } from "@/lib/acpStorage";
import { useRef, useState } from "react";
import ConfirmationDialog from "../ConfirmationDialog";
import StatusIndicator from "../StatusIndicator";
import ThemeToggle from "../ThemeToggle";
import { Button } from "../ui/button";
import { NavigationMenu, NavigationMenuContent, NavigationMenuItem, NavigationMenuList, NavigationMenuTrigger } from "../ui/navigation-menu";

const Header = () => {
    const [loadSample, data] = usePlaygroundStorageStore((state) => [state.loadSample, state.data]);
    const [status, setState] = usePlaygroundStore((state) => [state.status, state.setState, state.active]);
    const [selectedSample, setSelectedSample] = useState<AcpSample | null>(null);
    const [showConfirmation, setShowConfirmation] = useState(false);
    const fileInputRef = useRef<HTMLInputElement>(null);

    const computeHash = async (data: string) => {
        const encodedData = new TextEncoder().encode(data);
        const buffer = await crypto.subtle.digest("SHA-256", encodedData);
        return Array.from(new Uint8Array(buffer))
            .map((b) => b.toString(16).padStart(2, "0"))
            .join("").substring(0, 8);
    };

    const exportState = async () => {
        const activeStateData = data;
        const jsonBlob = new Blob([JSON.stringify(activeStateData, null, 2)], { type: "application/json" });
        const url = URL.createObjectURL(jsonBlob);
        const link = document.createElement("a");
        const filename = await computeHash(JSON.stringify(jsonBlob));
        link.href = url;
        link.download = `acp-playground-export-${filename}.json`;
        link.click();
        URL.revokeObjectURL(url);
    };

    const handleImportFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = () => {
                try {
                    const parsedData = JSON.parse(reader.result as string);
                    setState(parsedData);
                } catch (error) {
                    console.error("Invalid JSON file", error);
                }
            };
            reader.readAsText(file);
        }
    };

    const handleImportButtonClick = () => {
        if (!fileInputRef.current) return;
        fileInputRef.current.value = "";
        fileInputRef.current.click();
    };

    return <div className="flex justify-between px-4 py-2">
        <div>
            <span className="opacity-40">ACP Playground</span>
            <span className="inline-block w-[50px] h-[1px] bg-border align-middle mx-4"></span>

            <ConfirmationDialog
                title={`Load Policy - ${selectedSample?.title ?? ""}`}
                description={`Are you sure you want to replace with ${selectedSample?.title ?? ""}`}
                open={showConfirmation}
                onAction={(confirmed: boolean) => {
                    if (confirmed === true && selectedSample) {
                        loadSample(selectedSample.id);
                        setState({});
                    }
                    setShowConfirmation(false);
                }} />

            <NavigationMenu className="inline-block" delayDuration={0} >
                <NavigationMenuList>
                    <NavigationMenuItem>
                        <NavigationMenuTrigger>Sample Schemas</NavigationMenuTrigger>
                        <NavigationMenuContent>
                            <ul className="p-3 w-[200px] md:w-[500px]">
                                {allSamples.map((sample) => (
                                    <button key={sample.id}
                                        className="p-3 rounded-md text-left hover:bg-accent hover:text-accent-foreground"
                                        onClick={() => {
                                            setSelectedSample(sample);
                                            setShowConfirmation(true);
                                        }}>
                                        <div className="text-[13px] font-medium leading-none mb-1">{sample.title}</div>
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
            <StatusIndicator status={status} />

            <Button variant="outline" onClick={handleImportButtonClick}>Import</Button>

            <input
                ref={fileInputRef}
                type="file"
                accept=".json"
                onChange={handleImportFileChange}
                className="hidden"
            />

            <Button variant="default" onClick={exportState}>Export</Button>
            <ThemeToggle />
        </div>
    </div>
}

export default Header;

