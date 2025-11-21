import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "../ui/tooltip";
import * as TooltipPrimitive from "@radix-ui/react-tooltip"

type TextTooltipProps = {
    children: React.ReactNode;
    content: string;
    side?: "right" | "left";
    align?: "start" | "center" | "end";
} & React.ComponentPropsWithoutRef<typeof TooltipPrimitive.Content>

const TextTooltip = ({ children, content, ...props }: TextTooltipProps) => {
    return <TooltipProvider>
        <Tooltip delayDuration={100}>
            <TooltipTrigger asChild>
                {children}
            </TooltipTrigger>
            <TooltipContent {...props}>{content}</TooltipContent>
        </Tooltip>
    </TooltipProvider>
}

export default TextTooltip;