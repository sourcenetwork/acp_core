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
    state: SandboxData & { name: string };
}

const fetchShare = async (shareId: string, signal?: AbortSignal) => {
    const response = await fetch(`${SHARE_URL}/${shareId}`, {
        method: "GET",
        signal,
        headers: { "Content-Type": "application/json" },
    });

    if (!response.ok) {
        throw new Error(`Failed to fetch share (${response.status} ${response.statusText})`);
    }

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
        if (!shareId || !open) {
            setLoading(false);
            return;
        };

        const ac = new AbortController();

        setLoading(true);
        setError(null);
        setShareJson(null);

        fetchShare(shareId, ac.signal)
            .then((sandbox) => setShareJson(sandbox))
            .catch(error => {
                console.error(error);
                if ((error as any)?.name === "AbortError") return;
                setError((error as Error)?.message);
            }).
            finally(() => {
                if (!ac.signal.aborted) setLoading(false);
            });

        return () => ac.abort();
    }, [shareId, open]);

    const onActionClick = (action: 'new' | 'replace' | false, data?: ShareResponse | null) => {
        setOpen(false);

        if (!data) return;

        const { policyDefinition, relationships, policyTheorem } = data.state;

        if (action === 'replace') void updateActiveSandbox({
            data: { policyDefinition, relationships, policyTheorem }
        });

        if (action === 'new') void newSandbox({
            name: data.state.name || `Share: ${shareId}`,
            description: ``,
            data: { policyDefinition, relationships, policyTheorem }
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
                    <div className="mb-1">{shareJson.state?.name}</div>
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
