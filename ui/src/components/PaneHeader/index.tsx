import { cn } from "@/utils/classnames";
import { PanelLeftOpen, PanelRightOpen } from "lucide-react";
import TextTooltip from "../TextTooltip";
import { Button } from "../ui/button";

interface PaneTitleProps {
    title: string;
    showCollapse?: boolean;
    direction?: 'left' | 'right';
    onCollapseClick?: () => void;
    className?: string;
}

const PaneHeader = ({ title, showCollapse, direction = 'left', onCollapseClick, className }: PaneTitleProps) => {
    const Icon = direction === 'left' ? PanelRightOpen : PanelLeftOpen;

    return <div className={cn('text-xs text-muted-foreground mb-4 flex justify-between items-center border-b border-divider', className)}>
        {title}

        {showCollapse &&
            <TextTooltip content={"Collapse"} side={direction} align="center" >
                <Button variant="muted" size="iconSm" onClick={onCollapseClick} >
                    <Icon size={16} />
                </Button>
            </TextTooltip>
        }
    </div>
}

export default PaneHeader;
