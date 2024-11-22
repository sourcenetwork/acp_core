import { cn } from '@/utils/classnames';
import { ChevronsLeft, icons } from 'lucide-react';
import { useState } from "react";
import { Link, useLocation } from "react-router-dom";

const SideMenu = () => {
    const location = useLocation();
    const [collapsed, setCollapsed] = useState(false);

    const paths = [
        { label: "Policy", path: '/', icon: icons.BookText },
        { label: "Relationship", path: '/relationship', icon: icons.Waypoints },
        { label: "Tests", path: '/tests', icon: icons.FlaskConical },
    ];

    return <div className={cn("p-1 transition-all md:min-w-[200px] mx-2 min-w-0", {
        "min-w-0 mx-2 md:min-w-0": collapsed === true
    })}>

        <div className="mb-2 text-right">
            <button className='hidden md:inline' onClick={() => setCollapsed(!collapsed)}>
                <ChevronsLeft className={cn("ml-auto", { "rotate-180": collapsed === true })} />
            </button>
        </div>

        {paths.map(p => {
            const PathIcon = p.icon;
            const active = p.path === location.pathname;
            return <Link
                key={p.path}
                to={p.path}
                className={cn(`
                    flex mb-2 p-1 
                    items-center
                    border border-transparent rounded-md opacity-60 
                    group transition-all 
                    hover:opacity-80
                    hover:border-border`,
                    {
                        'border-border opacity-100': active
                    })}
            >
                <PathIcon className="w-5" />

                <span className={cn("ml-2 text-sm hidden md:block", {
                    "hidden md:hidden": collapsed === true
                })}>{p.label}</span>
            </Link>;
        })}

        {/* <pre>{JSON.stringify(IdHandleMap, null, 2)}</pre> */}
    </div>
}

export default SideMenu;

