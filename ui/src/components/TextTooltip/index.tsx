import { Tooltip, TooltipContent, TooltipProvider, TooltipTrigger } from "../ui/tooltip";

const TextTooltip = ({ children, content }: { children: React.ReactNode, content: string }) => {
    return <TooltipProvider>
        <Tooltip>
            <TooltipTrigger asChild>
                {children}
            </TooltipTrigger>
            <TooltipContent className="TooltipContent">{content}</TooltipContent>
        </Tooltip>
    </TooltipProvider>
}

export default TextTooltip;