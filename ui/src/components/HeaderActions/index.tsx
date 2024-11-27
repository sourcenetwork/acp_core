import { useSandbox } from "@/hooks/useSandbox";
import { usePlaygroundStore } from "@/lib/playgroundStore";
import { useTheme } from "@/ThemeProvider";
import { SandboxData } from "@/types/proto-js/sourcenetwork/acp_core/sandbox";
import { EllipsisVertical } from "lucide-react";
import { ComponentProps, ComponentType, useRef } from "react";
import ThemeToggle from "../ThemeToggle";
import { Button } from "../ui/button";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuTrigger } from "../ui/dropdown-menu";

const HeaderActions = () => {
    const { theme, setTheme } = useTheme();
    const activeSandbox = useSandbox();
    const [setState] = usePlaygroundStore((state) => [state.setPlaygroundState]);
    const fileInputRef = useRef<HTMLInputElement>(null);

    const computeHash = async (data: string) => {
        const encodedData = new TextEncoder().encode(data);
        const buffer = await crypto.subtle.digest("SHA-256", encodedData);
        return Array.from(new Uint8Array(buffer))
            .map((b) => b.toString(16).padStart(2, "0"))
            .join("").substring(0, 8);
    };

    const exportState = async () => {
        try {
            const activeStateData = activeSandbox?.data;
            const jsonBlob = new Blob([JSON.stringify(activeStateData, null, 2)], { type: "application/json" });
            const url = URL.createObjectURL(jsonBlob);
            const link = document.createElement("a");
            const filename = await computeHash(JSON.stringify(activeStateData));
            link.href = url;
            link.download = `acp-playground-export-${filename}.json`;
            link.click();
            URL.revokeObjectURL(url);
        } catch (error) {
            // TODO
            console.error(error);
        }
    };

    const handleImportFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        const file = event.target.files?.[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = () => {
                if (typeof reader.result !== 'string') {
                    return false;
                }

                const parsedData = JSON.parse(reader.result) as Partial<SandboxData>;
                void setState(parsedData);

                return true;
            };
            reader.readAsText(file);
        }
    };

    const handleExportButtonClick = () => {
        void exportState();
    };

    const handleImportButtonClick = () => {
        if (!fileInputRef.current) return;
        fileInputRef.current.value = "";
        fileInputRef.current.click();
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
                label: "Import",
                component: Button,
                props: { className: "text-xs", variant: "outline", size: "xs", onClick: handleImportButtonClick },
            },
            {
                label: "Export",
                component: Button,
                props: { className: "text-xs", variant: "default", size: "xs", onClick: handleExportButtonClick, },
            },
            {
                label: "Theme",
                component: ThemeToggle,
                props: { buttonProps: { variant: "outline", size: "iconSm", onClick: handeThemeToggleClick }, onClick: handeThemeToggleClick }
            },
        ];

    return (
        <div className="flex items-center justify-end">
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

            <input
                ref={fileInputRef}
                type="file"
                accept=".json"
                onChange={handleImportFileChange}
                className="hidden"
            />
        </div>
    );
};

export default HeaderActions;

