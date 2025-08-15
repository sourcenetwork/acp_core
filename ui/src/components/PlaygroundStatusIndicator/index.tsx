import { PlaygroundState } from "@/stores/playgroundStore";
import { cn } from "@/utils/classnames";
import { CircleCheck, CircleX, LoaderCircle } from "lucide-react";
import { ReactNode, useMemo } from "react";

interface PlaygroundStatusIndicatorProps {
    status: PlaygroundState['playgroundStatus']
    syncing: boolean;
}

const PlaygroundStatusIndicator = ({ status, syncing }: PlaygroundStatusIndicatorProps) => {

    const currentStatus: { label: string, color: string, icon: ReactNode } | null = useMemo(() => {
        const statuses = {
            'ready': {
                color: "text-src-secondary",
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
            'syncing': {
                color: "text-src-secondary",
                label: "Loaded",
                icon: <LoaderCircle className="ml-1 inline-block animate-spin w-4" />
            },
            'uninitialized': null,
        }

        return statuses[syncing ? 'syncing' : status];
    }, [status, syncing]);

    if (!currentStatus) return null;

    return <div className={cn("flex items-center justify-between p-3", currentStatus.color)}>
        <span className={cn('text-[12px] leading-none')}>
            {currentStatus.label}
        </span>
        {currentStatus.icon}
    </div>
}

export default PlaygroundStatusIndicator;

