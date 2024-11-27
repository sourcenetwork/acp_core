import { usePlaygroundStore } from '@/lib/playgroundStore';
import { cn } from '@/utils/classnames';
import { Blocks, ChevronsLeft, icons, Plus } from 'lucide-react';
import { useState } from "react";
import { Link, useLocation } from "react-router-dom";
import SandboxList from '../SandboxList';
import { Button } from '../ui/button';
import { Dialog, DialogContent, DialogHeader, DialogTitle } from '../ui/dialog';
import { ScrollArea } from '../ui/scroll-area';
import { Separator } from '../ui/separator';

const SideMenu = () => {
    const location = useLocation();
    const [collapsed, setCollapsed] = useState(false);
    const [sandboxDialogOpen, setSandboxDialogOpen] = useState(false);
    const [newEmptySandbox] = usePlaygroundStore((state) => [state.newEmptySandbox]);

    const handleNewSandboxButtonClick = () => void newEmptySandbox();
    const handleMenuToggleClick = () => setCollapsed(!collapsed);
    const handleSandboxDialogClick = () => setSandboxDialogOpen(true)

    const paths = [
        { label: "Policy", path: '/', icon: icons.BookText },
        { label: "Relationship", path: '/relationship', icon: icons.Waypoints },
        { label: "Tests", path: '/tests', icon: icons.FlaskConical },
    ];

    return <ScrollArea type='always' className={cn(" transition-all min-w-0", {
        "md:min-w-[225px]": collapsed === false,
    })}
    >
        <div className={cn("transition-all w-full")}>
            <div className='grid px-3'>
                <div className={cn("mb-2 text-right", {
                    "text-center": collapsed === true
                })}>
                    <button className='hidden md:inline-block' onClick={handleMenuToggleClick}>
                        <ChevronsLeft size={20} className={cn("ml-auto", { "rotate-180": collapsed === true })} />
                    </button>
                </div>

                {paths.map(p => {
                    const PathIcon = p.icon;
                    const active = p.path === location.pathname;
                    return <Link
                        key={p.path}
                        to={p.path}
                        className={cn("sidemenu-button", { 'border-border opacity-100 bg-secondary': active })}
                    >
                        <PathIcon size={20} />
                        <span className={cn("ml-2 text-sm hidden md:block", {
                            "hidden md:hidden": collapsed === true
                        })}>{p.label}</span>
                    </Link>;
                })}

                <Separator className='my-2' />

                <button
                    onClick={handleSandboxDialogClick}
                    className={cn("sidemenu-button", "md:hidden", { "md:flex": collapsed === true })}>
                    <Blocks size={20} />
                </button>

                <div className={cn("mt-3 ", {
                    "hidden md:block": collapsed === false,
                    "hidden": collapsed === true
                })}>
                    <div className='text-xs text-primary opacity-50 mb-4 '>Sandboxes</div>
                    <Button variant={"ghost"} size={"sm"} className="border-2 text-xs w-full opacity-80 hover:opacity-100" onClick={handleNewSandboxButtonClick}>
                        <Plus size={20} className="mr-1" /> New Sandbox
                    </Button>
                    <SandboxList format='small' className='py-3' />
                </div>

                {/* Sandbox List Dialog  */}
                <Dialog open={sandboxDialogOpen} onOpenChange={(change) => setSandboxDialogOpen(change)}>
                    <DialogContent className='sm:max-w-md'>
                        <DialogHeader>
                            <DialogTitle className='text-left'>Sandboxes</DialogTitle>
                        </DialogHeader>
                        <ScrollArea type='always' className='border rounded p-2 max-h-[300px] relative -mr-[10px] pr-[10px]'>
                            <SandboxList format='expanded' className='' onSandboxClick={() => setSandboxDialogOpen(false)} />
                        </ScrollArea>
                        <div className='text-right'>
                            <Button className='' type="button" variant="secondary" onClick={handleNewSandboxButtonClick}>
                                New Sandbox
                            </Button>
                        </div>
                    </DialogContent>
                </Dialog>

            </div>
        </div>
    </ScrollArea>
}

export default SideMenu;

