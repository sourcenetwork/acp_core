import { PersistedSandboxData, usePlaygroundStore } from "@/lib/playgroundStore";
import { SHARE_URL } from "@/constants";
import { useEffect, useState } from "react";
import { Button } from "../ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "../ui/dialog";

interface DialogLoadShareProps {
    shareId: string | null,
    open: boolean,
    setOpen: (state: boolean) => unknown
}

const fetchShare = async (shareId: string) => {
    const response = await fetch(`${SHARE_URL}?id=${shareId}`, {
        method: "GET",
        headers: { "Content-Type": "application/json" },
    });

    const result = await response.json() as PersistedSandboxData;
    return result;
}

const DialogLoadShare = ({ shareId, open, setOpen }: DialogLoadShareProps) => {
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState(false);
    const [shareJson, setShareJson] = useState<PersistedSandboxData | null>(null);

    const { playgroundStatus, updateActiveSandbox, newSandbox } = usePlaygroundStore((state) => ({
        playgroundStatus: state.playgroundStatus,
        updateActiveSandbox: state.updateActiveStoredSandbox,
        newSandbox: state.newSandbox
    }));

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

    const onActionClick = (action: 'new' | 'replace' | false, data?: PersistedSandboxData | null) => {
        setOpen(false);

        if (!data) return;

        if (action === 'replace') void updateActiveSandbox({
            data: data.data
        });

        if (action === 'new') void newSandbox({
            name: data.name,
            description: data.description,
            data: data.data
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
                    <div className="mb-1">{shareJson?.name}</div>
                    <div className="text-xs">{shareJson?.description}</div>
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
                        <Button className="" type="button" variant="default" disabled={!shareJson} onClick={() => onActionClick('new', shareJson)}>Load</Button>
                    </>
                }
            </DialogFooter>
        </DialogContent>
    </Dialog>
}

export default DialogLoadShare;
