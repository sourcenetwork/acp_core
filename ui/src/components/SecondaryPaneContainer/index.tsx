import { SecondaryPaneTypes, useLayoutStore } from "@/stores/layoutStore";
import { useMemo } from "react";
import PaneCheck from "../PaneCheck";
import PaneExpand from "../PaneExpand";

const SecondaryPaneComponentMap = {
    [SecondaryPaneTypes.Check]: PaneCheck,
    [SecondaryPaneTypes.Expand]: PaneExpand,
}

const SecondaryPaneContainer = () => {
    const secondaryPaneType = useLayoutStore((state) => state.secondaryPaneType);
    const Component = useMemo(() => secondaryPaneType ? SecondaryPaneComponentMap[secondaryPaneType] : null, [secondaryPaneType]);

    return <div className='relative z-10 h-full'>
        {Component && <Component />}
    </div>
};

export default SecondaryPaneContainer;  