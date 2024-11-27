import { Button } from "@/components/ui/button";
import { usePlaygroundStore } from "@/lib/playgroundStore";
import { cn } from "@/utils/classnames";
import { CircleCheck, CircleX, LoaderCircle } from "lucide-react";
import BaseEditor from "../components/Editor";

const Tests = () => {
    const [verifyTheoremsStatus, verifyTheorems, sandboxErrorCount] = usePlaygroundStore((state) => [state.verifyTheoremsStatus, state.verifyTheorems, state.setStateDataErrorCount]);

    const isDisabled = sandboxErrorCount > 0;
    const status = isDisabled ? 'disabled' : verifyTheoremsStatus;

    const validationStatus = {
        'passed': {
            color: "text-green-500",
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

    return <>
        <div className="mb-2 flex justify-end gap-x-2">
            <div className="my-auto font-bold text-[12px]">
                <span className={cn('inline-block text-[12px] leading-none ',
                    currentStatus.color,
                    { 'opacity-50': isDisabled }
                )}>
                    {currentStatus.label}
                    {currentStatus.icon}
                </span>
            </div>

            <Button variant="outline" disabled={sandboxErrorCount > 0} onClick={runVerification}>Run</Button>
            <Button variant="outline" asChild><a href="https://docs.source.network" target="_blank" rel="noreferrer">Guide</a></Button>
        </div >
        <BaseEditor sandboxDataType={"policyTheorem"} />
    </>;
}

export default Tests;
