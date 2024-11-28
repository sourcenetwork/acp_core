import { usePlaygroundStore } from "@/lib/playgroundStore";
import { SandboxTemplate } from "@acp/sandbox";
import { useState } from "react";
import DialogLoadTemplate from "../DialogLoadTemplate";
import HeaderActions from "../HeaderActions";
import StatusIndicator from "../StatusIndicator";
import { NavigationMenu, NavigationMenuContent, NavigationMenuItem, NavigationMenuList, NavigationMenuTrigger } from "../ui/navigation-menu";

const Header = () => {
    const [playgroundStatus, samples, loadTemplate, newSandbox] = usePlaygroundStore((state) => [
        state.playgroundStatus,
        state.sandboxTemplates,
        state.loadTemplate,
        state.newSandbox
    ]);

    const [selectedSample, setSelectedSample] = useState<SandboxTemplate | null>(null);
    const [showConfirmation, setShowConfirmation] = useState(false);

    const handleSampleClick = (sample: SandboxTemplate) => () => {
        setSelectedSample(sample);
        setShowConfirmation(true);
    }

    return <div className="flex items-center justify-between px-4 py-2 border-b md:border-b-0">

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
            <span className="opacity-40 mr-3">ACP Playground</span>
            <span className="hidden md:inline-block w-[30px] h-[1px] bg-border align-middle mx-4"></span>
            <div>
                <NavigationMenu className="inline-block" delayDuration={0} >
                    <NavigationMenuList>
                        <NavigationMenuItem>
                            <NavigationMenuTrigger disabled={playgroundStatus !== 'ready'} className="text-xs">Samples</NavigationMenuTrigger>
                            <NavigationMenuContent>
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
            <StatusIndicator status={playgroundStatus} />
            <HeaderActions />
        </div>
    </div>
}

export default Header;

