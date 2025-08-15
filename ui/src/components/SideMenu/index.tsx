import { useLayoutStore, useUIActions } from '@/stores/layoutStore';
import { cn } from '@/utils/classnames';
import { Box, Plus } from 'lucide-react';
import PaneHeader from '../PaneHeader';
import SandboxList from '../SandboxList';
import TextTooltip from '../TextTooltip';
import { Button } from '../ui/button';
import { ScrollArea } from '../ui/scroll-area';
import { Separator } from '../ui/separator';

const SideMenu = () => {
    const sandboxMenuOpen = useLayoutStore((state) => state.sandboxMenuOpen);
    const { setSandboxMenuOpen, setCreateSandboxDialogOpen } = useUIActions();

    const handleNewSandboxButtonClick = () => setCreateSandboxDialogOpen(true);
    const handleMenuToggleClick = () => setSandboxMenuOpen(!sandboxMenuOpen);

    return <div className="flex shrink-0">
        <div className='min-w-0 shrink-0 flex flex-col gap-2 p-2 border-r border-divider '>
            <TextTooltip content="New Sandbox" side="right" align="center"  >
                <Button variant="ghost" size="icon" onClick={handleNewSandboxButtonClick} >
                    <Plus size={20} />
                </Button>
            </TextTooltip>
            <Separator />
            <TextTooltip content="Sandboxes" side="right" align="center"  >
                <Button variant="ghost" size="icon" onClick={handleMenuToggleClick} className={cn({ "text-src-primary": sandboxMenuOpen })}>
                    <Box size={20} />
                </Button>
            </TextTooltip>
        </div>

        {sandboxMenuOpen &&
            <div className='fixed md:static md:flex flex-col inset-0 shrink-0 z-11 md:z-auto '>
                <ScrollArea type='always' className={cn("h-full min-w-0 p-2 bg-background border-r border-divider w-[240px] relative z-2")}>
                    <PaneHeader title="Sandboxes" showCollapse onCollapseClick={handleMenuToggleClick} />
                    <Button variant="default" size={"sm"} className="border-2 text-xs w-full opacity-80 hover:opacity-100" onClick={handleNewSandboxButtonClick}>
                        <Plus size={16} className="mr-1" /> New Sandbox
                    </Button>
                    <SandboxList format='small' className='py-3' />
                </ScrollArea>
                <div className='fixed inset-0 bg-background opacity-50 md:hidden' onClick={handleMenuToggleClick} />
            </div>
        }
    </div>
}

export default SideMenu;

