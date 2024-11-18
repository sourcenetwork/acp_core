import { icons } from "lucide-react";
import { useTheme } from "../../ThemeProvider";
import { Button } from "../ui/button";

const ThemeToggle = () => {
    const { theme, setTheme } = useTheme();
    const Icon = theme === "dark" ? icons['Moon'] : icons['Sun'];
    return (
        <div className="flex">
            <Button variant="outline" size="icon" onClick={() => setTheme(theme === 'dark' ? 'light' : 'dark')}>
                <Icon />
            </Button>
        </div>
    )
}

export default ThemeToggle;