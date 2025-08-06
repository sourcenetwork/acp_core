import { Button } from "@/components/ui/button";
import { usePlaygroundStore } from "@/stores/playgroundStore";
import { cn } from "@/utils/classnames";
import { CircleCheck, CircleX, LoaderCircle } from "lucide-react";
import BaseEditor from "../Editor";
import TextTooltip from "../TextTooltip";

const PaneTests = () => {
    const verifyTheoremsStatus = usePlaygroundStore((state) => state.verifyTheoremsStatus);
    const verifyTheorems = usePlaygroundStore((state) => state.verifyTheorems);
    const sandboxErrorCount = usePlaygroundStore((state) => state.setStateDataErrorCount);

    const isDisabled = sandboxErrorCount > 0;
    const status = isDisabled ? 'disabled' : verifyTheoremsStatus;

    const validationStatus = {
        'passed': {
            color: "text-src-primary",
            label: "Validation Passed",
            icon: <CircleCheck className="ml-1 inline-block w-4" />
        },
        'loading': {
            color: "",
            label: "Validating",
            icon: <LoaderCircle className="ml-1 inline-block animate-spin w-4" />
        },
        'error': {
            color: "text-red-500",
            label: "Validation Failed",
            icon: <CircleX className="ml-1 inline-block w-4" />
        },
        'pending': {
            color: "",
            label: "Ready to Validate",
            icon: null
        },
        'disabled': {
            color: "",
            label: "Fix Errors to Validate",
            icon: null
        },
    }

    const currentStatus = validationStatus[status];

    const runVerification = () => {
        void verifyTheorems();
    };

    return <div className="flex flex-col h-full">
        <div className="py-2 flex items-center justify-end gap-x-2 whitespace-nowrap ">
            <span className={cn('inline-block leading-none text-xs',
                currentStatus.color,
                { 'opacity-50': isDisabled }
            )}>
                {currentStatus.label}
                {currentStatus.icon}
            </span>

            <Button variant="outline" size="xs" className="text-xs" disabled={sandboxErrorCount > 0} onClick={runVerification}>Run</Button>
            <TextTooltip content="Guide" side="right" align="center" >
                <Button variant="outline" size="xs" className="text-xs" asChild>
                    <a href="https://docs.source.network" target="_blank" rel="noreferrer" aria-label="Guide">
                        Guide
                    </a>
                </Button>
            </TextTooltip>
        </div >
        <div className="grow overflow-hidden pb-2">
            <BaseEditor sandboxDataType={"policyTheorem"} />
        </div>
    </div>;
}

export default PaneTests;
