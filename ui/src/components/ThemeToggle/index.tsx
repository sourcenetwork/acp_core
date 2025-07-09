import { useTheme } from "@/providers/ThemeProvider/useTheme";
import { LucideIcon, LucideProps, MoonIcon, SunIcon } from "lucide-react";
import { Button, ButtonProps } from "../ui/button";

interface ThemeToggleProps {
    buttonProps?: ButtonProps,
    darkIcon?: LucideIcon, // Default dark mode icon
    lightIcon?: LucideIcon, // Default light mode icon
    iconProps?: LucideProps
}

const ThemeToggle = (props: ButtonProps & ThemeToggleProps) => {
    const { darkIcon = MoonIcon, lightIcon = SunIcon, buttonProps, iconProps } = props;
    const { theme } = useTheme();
    const Icon = theme === "dark" ? darkIcon : lightIcon;

    return (
        <Button {...buttonProps} >
            <Icon size={16} {...iconProps} />
        </Button>
    )
}

export default ThemeToggle;