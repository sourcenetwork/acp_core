import { SHARE_URL } from "@/constants";
import { usePlaygroundStore } from "@/stores/playgroundStore";
import { SandboxData } from "@/types/proto-js/sourcenetwork/acp_core/sandbox";
import { useEffect, useState } from "react";
import { Button } from "../ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "../ui/dialog";

interface DialogLoadShareProps {
    shareId: string | null,
    open: boolean,
    setOpen: (state: boolean) => unknown
}

type ShareResponse = {
    state: SandboxData;
}

const fetchShare = async (shareId: string) => {
    const response = await fetch(`${SHARE_URL}/${shareId}`, {
        method: "GET",
        headers: { "Content-Type": "application/json" },
    });

    const result = await response.json() as ShareResponse;

    return result;
}

const DialogLoadShare = ({ shareId, open, setOpen }: DialogLoadShareProps) => {
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState(false);
    const [shareJson, setShareJson] = useState<ShareResponse | null>(null);

    const playgroundStatus = usePlaygroundStore((state) => state.playgroundStatus);
    const updateActiveSandbox = usePlaygroundStore((state) => state.updateActiveStoredSandbox);
    const newSandbox = usePlaygroundStore((state) => state.newSandbox);

    const isLoading = loading === true || playgroundStatus === 'loading';

    useEffect(() => {
        if (!shareId) return;
        setLoading(true);
        setError(null);

        fetchShare(shareId)
            .then((sandbox) => setShareJson(sandbox))
            .catch(error => setError((error as Error)?.message)).
            finally(() => setLoading(false));

    }, [shareId]);

    const onActionClick = (action: 'new' | 'replace' | false, data?: ShareResponse | null) => {
        setOpen(false);

        if (!data) return;

        if (action === 'replace') void updateActiveSandbox({
            data: data.state
        });

        if (action === 'new') void newSandbox({
            name: `Share: ${shareId}`,
            description: ``,
            data: data.state
        });
    }

    return <Dialog open={open} defaultOpen={true} onOpenChange={(state) => setOpen(state)}>
        <DialogContent className="sm:max-w-md">
            <DialogHeader>
                <DialogTitle>Load Share</DialogTitle>
                <DialogDescription>Load this shared policy</DialogDescription>
            </DialogHeader>


            {shareJson &&
                <div className="mb-3">
                    <div className="text-xs opacity-50 mb-1">Share: {shareId}</div>
                </div>
            }

            {error &&
                <div className="mb-3 text-sm text-destructive">Something went wrong fetching share</div>}

            {isLoading === true &&
                <div className="mb-3 text-sm">Loading ... Please wait</div>}

            <DialogFooter className="sm:justify-between">
                {isLoading === false &&
                    <>
                        <Button type="button" variant="secondary" disabled={!shareJson} onClick={() => onActionClick('replace', shareJson)}>Override</Button>
                        <Button className="" type="button" variant="default" disabled={!shareJson} onClick={() => onActionClick('new', shareJson)}>Load New</Button>
                    </>
                }
            </DialogFooter>
        </DialogContent>
    </Dialog>
}

export default DialogLoadShare;
