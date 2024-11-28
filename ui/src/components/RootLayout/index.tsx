import { ComponentType, useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import Policy from "../../routes/Policy";
import Relationship from "../../routes/Relationship";
import Tests from "../../routes/Tests";
import Header from "../Header";
import DialogLoadShare from "../DialogLoadShare";
import Output from "../Output";
import SideMenu from "../SideMenu";
import { ResizableHandle, ResizablePanel, ResizablePanelGroup } from "../ui/resizable";
import { Tabs, TabsList, TabsTrigger } from "../ui/tabs";
import { Toaster } from "../ui/toaster";

export interface PanelDefiniton {
    key: string,
    label: string,
    component: ComponentType
}

const RootLayout = () => {
    const location = useLocation();
    const navigate = useNavigate();
    const [shareDialogOpen, setShareDialogOpen] = useState(false);
    const queryParams = new URLSearchParams(location.search);
    const shareId = queryParams.get('share');

    const clearShareParam = () => {
        const params = new URLSearchParams(window.location.search);
        params.delete('share');
        navigate({ pathname: window.location.pathname, search: params.toString() });
    };

    useEffect(() => {
        if (shareId) setShareDialogOpen(true);
    }, [shareId]);

    const panels = {
        'policy': {
            key: 'policy',
            label: 'Policy',
            component: Policy,
        },
        'relationship': {
            key: 'relationship',
            label: 'Relationships',
            component: Relationship,
        },
        'tests': {
            key: 'tests',
            label: 'Tests',
            component: Tests,
        },
        'output': {
            key: 'output',
            label: 'Output',
            component: Output,
        },
    };

    const pathComponents: Record<string, {
        main?: PanelDefiniton[],
        secondary?: PanelDefiniton[],
        bottom?: PanelDefiniton[],
    }> = {
        '/': {
            main: [panels.policy],
            secondary: [],
            bottom: [panels.output],
        },
        '/relationship': {
            main: [panels.relationship],
            secondary: [],
            bottom: [panels.output],
        },
        '/tests': {
            main: [panels.tests],
            secondary: [],
            bottom: [panels.output],
        },
    };

    const primary = pathComponents[location?.pathname]?.main ?? [];
    const secondary = pathComponents[location?.pathname]?.secondary ?? [];

    return (
        <div className="flex h-dvh flex-col">
            <Header />
            <div className="flex flex-1 overflow-y-auto">

                <SideMenu />

                <div className="w-full min-w-0 ">
                    <ResizablePanelGroup direction="vertical" className="h-full w-full rounded-lg" >
                        <ResizablePanel defaultSize={75} className="flex h-full" >
                            <ResizablePanelGroup direction="horizontal">
                                <ResizablePanel className="py-4 pr-4" >
                                    {primary.map(({ component: Component, key }) => <Component key={key} />)}
                                </ResizablePanel>
                                {secondary?.length !== 0 && <>
                                    <ResizableHandle className="" withHandle />
                                    <ResizablePanel className="py-2 pl-4">
                                        <Tabs defaultValue={secondary[0]?.key} className="w-full">
                                            <TabsList className="mb-4">
                                                {secondary?.map(({ key, label }) => <TabsTrigger key={key} value={key}>{label}</TabsTrigger>)}
                                            </TabsList>
                                        </Tabs>
                                    </ResizablePanel>
                                </>}
                            </ResizablePanelGroup>
                        </ResizablePanel>

                        <ResizableHandle withHandle />
                        <ResizablePanel defaultSize={25}>
                            <Output />
                        </ResizablePanel>
                    </ResizablePanelGroup>
                </div>
            </div>

            <DialogLoadShare shareId={shareId} open={shareDialogOpen} setOpen={(state) => {
                setShareDialogOpen(state)
                clearShareParam();
            }} />

            <Toaster />
        </div>
    )
}

export default RootLayout;
