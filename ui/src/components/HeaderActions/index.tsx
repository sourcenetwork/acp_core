import { useToast } from "@/hooks/use-toast";
import { useSandbox } from "@/hooks/useSandbox";
import { usePlaygroundStore } from "@/lib/playgroundStore";
import { useTheme } from "@/providers/ThemeProvider/useTheme";
import { exportSandboxData, importSandboxData } from "@/utils/sandboxFileUtils";
import { EllipsisVertical } from "lucide-react";
import { ComponentProps, ComponentType, useState } from "react";
import DialogCopyShare from "../DialogCopyShare";
import ThemeToggle from "../ThemeToggle";
import { Button } from "../ui/button";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "../ui/dropdown-menu";

const HeaderActions = () => {
    const { theme, setTheme } = useTheme();
    const activeSandbox = useSandbox();
    const [setState] = usePlaygroundStore((state) => [state.setPlaygroundState]);
    const [showShareDialog, setShowShareDialog] = useState<boolean>(false);
    const { toast } = useToast()

    const handleShareButtonClick = () => {
        setShowShareDialog(true);
    };

    const handleExportButtonClick = async () => {
        try {
            await exportSandboxData(activeSandbox?.name, activeSandbox?.data);
        } catch (error) {
            console.error(error);
            toast({ description: "Something went wrong exporting" })
        }
    };

    const handleImportButtonClick = async () => {
        try {
            const data = await importSandboxData();
            if (!data) return;
            void setState(data);
        } catch (error) {
            console.error(error);
            toast({ description: "Something went wrong importing" })
        }
    };

    const handeThemeToggleClick = () => {
        setTheme(theme === 'dark' ? 'light' : 'dark');
    }

    const menuItems: {
        component: ComponentType<ComponentProps<typeof Button | typeof ThemeToggle>>;
        props: Record<string, unknown> & { onClick?: () => void };
        label: string;
    }[] = [
            {
                label: "Share",
                component: Button,
                props: { className: "text-xs", variant: "outline", size: "xs", onClick: handleShareButtonClick },
            },
            {
                label: "Import",
                component: Button,
                props: { className: "text-xs", variant: "outline", size: "xs", onClick: () => void handleImportButtonClick() },
            },
            {
                label: "Export",
                component: Button,
                props: { className: "text-xs", variant: "default", size: "xs", onClick: () => void handleExportButtonClick(), },
            },
            {
                label: "Theme",
                component: ThemeToggle,
                props: { buttonProps: { variant: "outline", size: "iconSm", onClick: handeThemeToggleClick }, onClick: handeThemeToggleClick }
            },
        ];

    return (
        <div className="flex items-center justify-end">

            <DialogCopyShare open={showShareDialog} setOpen={(state) => setShowShareDialog(state)} />

            <div className="hidden md:flex space-x-2">
                {menuItems.map((item, index) => {
                    const Component = item.component;
                    return <Component key={index} {...item.props}>{item.label}</Component>;
                })}
            </div>

            <div className="md:hidden">
                <DropdownMenu>
                    <DropdownMenuTrigger asChild>
                        <Button variant="ghost" size="iconSm">
                            <EllipsisVertical size={17} />
                        </Button>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent className="min-w-[170px]">
                        {menuItems.map((item, index) => (
                            <DropdownMenuItem key={index} onClick={item.props?.onClick} >
                                {item.label}
                            </DropdownMenuItem>
                        ))}
                    </DropdownMenuContent>
                </DropdownMenu>
            </div>
        </div>
    );
};

export default HeaderActions;

