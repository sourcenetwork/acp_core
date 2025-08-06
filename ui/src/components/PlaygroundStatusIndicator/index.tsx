// import { WasmStatus } from "@/lib/wasm";

import { PlaygroundState } from "@/stores/playgroundStore";
import { cn } from "@/utils/classnames";
import { CircleCheck, CircleX, LoaderCircle } from "lucide-react";
import { ReactNode } from "react";

interface PlaygroundStatusIndicatorProps {
    status: PlaygroundState['playgroundStatus']
}

const PlaygroundStatusIndicator = ({ status }: PlaygroundStatusIndicatorProps) => {

    const statuses: Record<PlaygroundState['playgroundStatus'], { label: string, color: string, icon: ReactNode } | null> = {
        'ready': {
            color: "text-src-primary",
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

export default PlaygroundStatusIndicator;

