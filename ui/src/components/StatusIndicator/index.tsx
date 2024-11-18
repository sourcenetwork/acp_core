// import { WasmStatus } from "@/lib/wasm";

import { PlaygroundState } from "@/lib/acpHandler";
import clsx from "clsx";
import { CircleCheck, CircleX, LoaderCircle } from "lucide-react";
import { ReactNode } from "react";

interface StatusIndicatorProps {
    status: PlaygroundState['status']
}

const StatusIndicator = ({ status }: StatusIndicatorProps) => {

    const statuses: Record<PlaygroundState['status'], { label: string, color: string, icon: ReactNode } | null> = {
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

    return <div className="flex justify-between p-4">
        {currentStatus &&
            <span className={clsx('inline-block text-[12px] leading-none ', currentStatus.color)}>
                {currentStatus.label}
                {currentStatus.icon}
            </span>
        }
    </div>
}

export default StatusIndicator;

