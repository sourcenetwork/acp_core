import { Button } from "@/components/ui/button";
import BaseEditor from "../Editor";
import TextTooltip from "../TextTooltip";

const PanePolicy = () => {

    return <div className="flex flex-col h-full">
        <div className="py-2 flex justify-end gap-x-2">
            <TextTooltip content="Guide" side="right" align="center" >
                <Button variant="outline" size="xs" className="text-xs" asChild>
                    <a href="https://docs.source.network" target="_blank" rel="noreferrer" aria-label="Guide">
                        Guide
                    </a>
                </Button>
            </TextTooltip>
        </div>
        <div className="grow overflow-hidden pb-2">
            <BaseEditor sandboxDataType={"policyDefinition"} />
        </div>
    </div>;
}

export default PanePolicy;
