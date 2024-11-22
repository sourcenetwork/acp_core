// import { WasmStatus } from "@/lib/wasm";

import { PlaygroundState } from "@/lib/playgroundStore";
import { cn } from "@/utils/classnames";
import { CircleCheck, CircleX, LoaderCircle } from "lucide-react";
import { ReactNode } from "react";

interface StatusIndicatorProps {
    status: PlaygroundState['playgroundStatus']
}

const StatusIndicator = ({ status }: StatusIndicatorProps) => {

    const statuses: Record<PlaygroundState['playgroundStatus'], { label: string, color: string, icon: ReactNode } | null> = {
        'ready': {
            color: "text-green-500",
            label: "Loaded",
            icon: <CircleCheck className="ml-1 inline-block w-4" />
        },
        'loading': {
            color: "",
            label: "Loading",
            icon: <LoaderCircle className="ml-1 inline-block animate-spin w-4" />
        },
        'error': {
            color: "text-red-500",
            label: "Error Loading",
            icon: <CircleX className="ml-1 inline-block w-4" />
        },
        'uninitialized': null,
    }

    const currentStatus = statuses[status];

    return <div className="flex justify-between p-3">
        {currentStatus &&
            <span className={cn('inline-block text-[12px] leading-none ', currentStatus.color)}>
                {currentStatus.label}
                {currentStatus.icon}
            </span>
        }
    </div>
}

export default StatusIndicator;

