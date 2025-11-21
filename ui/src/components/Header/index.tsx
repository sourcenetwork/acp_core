import { useSandbox } from "@/hooks/useSandbox";
import { useLayoutStore } from "@/stores/layoutStore";
import { usePlaygroundStore } from "@/stores/playgroundStore";
import { SandboxTemplate } from "@acp/sandbox";
import { Box } from "lucide-react";
import { useState } from "react";
import DialogLoadTemplate from "../DialogLoadTemplate";
import HeaderActions from "../HeaderActions";
import PlaygroundStatusIndicator from "../PlaygroundStatusIndicator";
import { NavigationMenu, NavigationMenuContent, NavigationMenuItem, NavigationMenuList, NavigationMenuTrigger } from "../ui/navigation-menu";

const Header = () => {

    const activeSandbox = useSandbox();
    const playgroundStatus = usePlaygroundStore((state) => state.playgroundStatus);
    const samples = usePlaygroundStore((state) => state.sandboxTemplates);
    const loadTemplate = usePlaygroundStore((state) => state.loadTemplate);
    const newSandbox = usePlaygroundStore((state) => state.newSandbox);
    const playgroundSyncing = usePlaygroundStore((state) => state.playgroundSyncing);
    const setSandboxMenuOpen = useLayoutStore((state) => state.setSandboxMenuOpen);



    const [selectedSample, setSelectedSample] = useState<SandboxTemplate | null>(null);
    const [showConfirmation, setShowConfirmation] = useState(false);

    const handleSampleClick = (sample: SandboxTemplate) => () => {
        setSelectedSample(sample);
        setShowConfirmation(true);
    }

    const handleSandboxClick = () => {
        setSandboxMenuOpen(true);
    }

    return <div className="flex items-center justify-between px-4 py-2 border-b border-divider relative bg-background z-1">

        <DialogLoadTemplate
            title={`Load Policy - ${selectedSample?.name ?? ""}`}
            description={`Load this policy template to a new sandbox or override the current schema with ${selectedSample?.name ?? ""}`}
            open={showConfirmation}
            onAction={(action) => {
                setShowConfirmation(false);
                if (!selectedSample) return;
                if (action === 'replace') void loadTemplate(selectedSample);
                if (action === 'new') void newSandbox(selectedSample);
            }} />

        <div className="flex items-center">
            <span className="opacity-60 mr-3 text-sm">ACP Playground</span>
            <span className="hidden md:inline-block w-[30px] h-px bg-border align-middle mx-4"></span>
            <div className="flex items-center gap-x-2">
                {activeSandbox &&
                    <div className="text-xs text-muted-foreground hover:text-accent-foreground items-center gap-x-2 hidden md:flex border-r border-divider pr-6 cursor-pointer" onClick={handleSandboxClick}>
                        <Box size={16} />
                        <span>{activeSandbox?.name}</span>
                    </div>
                }

                <NavigationMenu className="inline-block" delayDuration={0} >
                    <NavigationMenuList>
                        <NavigationMenuItem>
                            <NavigationMenuTrigger disabled={playgroundStatus !== 'ready'} className="text-xs">Samples</NavigationMenuTrigger>
                            <NavigationMenuContent >
                                <ul className="p-3 w-[300px] md:w-[400px]">
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
        </div>

        <div className="self-start items-center gap-x-1 flex">
            <PlaygroundStatusIndicator status={playgroundStatus} syncing={playgroundSyncing} />
            <HeaderActions />
        </div>
    </div>
}

export default Header;

