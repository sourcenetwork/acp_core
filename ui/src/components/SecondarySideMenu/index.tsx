import { SecondaryPaneTypes, useLayoutStore, useUIActions } from '@/stores/layoutStore';
import { cn } from '@/utils/classnames';
import { SquareCheckBig } from 'lucide-react';
import TextTooltip from '../TextTooltip';
import { Button } from '../ui/button';
import { ScrollArea } from '../ui/scroll-area';

const SecondarySideMenu = () => {
    const { setSecondaryPaneType, setSecondaryPaneOpen } = useUIActions();
    const secondaryPaneOpen = useLayoutStore((state) => state.secondaryPaneOpen);
    const secondaryPaneType = useLayoutStore((state) => state.secondaryPaneType);

    const onHandleClick = (type: SecondaryPaneTypes) => {
        const isActivePane = secondaryPaneType === type;
        setSecondaryPaneType(type);
        setSecondaryPaneOpen(isActivePane && secondaryPaneOpen ? false : true);
    }

    return <ScrollArea type='always' className={cn("shrink-0 border-l border-divider ",
        secondaryPaneOpen ? "md:border-l" : "md:border-none"
    )}>
        <div className="flex flex-col gap-2 p-2 border-divider shrink-0">
            <TextTooltip content="Check" side="right" align="center" >
                <Button
                    size="icon"
                    variant="muted"
                    onClick={() => onHandleClick(SecondaryPaneTypes.Check)}
                    className={cn({ 'opacity-100 text-src-primary': secondaryPaneOpen && secondaryPaneType === SecondaryPaneTypes.Check })}>
                    <SquareCheckBig size={20} />
                </Button>
            </TextTooltip>

            {/* <TextTooltip content="Expand" side="left" align="center" >
                <Button
                    size="icon"
                    variant="muted"
                    onClick={() => onHandleClick(SecondaryPaneTypes.Expand)}
                    className={cn({ 'opacity-100 text-src-primary': secondaryPaneOpen && secondaryPaneType === SecondaryPaneTypes.Expand })}>
                    <GanttChart size={20} />
                </Button>
            </TextTooltip> */}
        </div>
    </ScrollArea >
}

export default SecondarySideMenu;
