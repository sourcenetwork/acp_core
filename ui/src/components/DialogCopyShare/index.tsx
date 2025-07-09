import { useSandbox } from "@/hooks/useSandbox";
import { PersistedSandboxData } from "@/lib/playgroundStore";
import { cn } from "@/utils/classnames";
import { SHARE_URL } from "@/constants";
import { CheckIcon, ClipboardIcon } from "lucide-react";
import { useEffect, useState } from "react";
import { Button } from "../ui/button";
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from "../ui/dialog";
import { Input } from "../ui/input";

interface DialogCopyShareProps {
    open: boolean,
    setOpen: (state: boolean) => unknown
}

const postShare = async (data?: PersistedSandboxData) => {
    if (!data) return;
    const response = await fetch(SHARE_URL, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
            name: data?.name,
            description: data?.description,
            data: data?.data,
        })
    })

    const result = await response.json() as Partial<PersistedSandboxData>;
    return result;
}

const copyToClipboardWithMeta = async (value: string) => {
    await navigator.clipboard.writeText(value)
}

const DialogCopyShare = ({ open, setOpen }: DialogCopyShareProps) => {
    const baseUrl = `${window.location.origin}`;
    const activeSandbox = useSandbox();
    const [shareLink, setShareLink] = useState<string | null>(null);
    const [loading, setLoading] = useState<boolean>(false);
    const [error, setError] = useState<string | null>(null);
    const [hasCopied, setHasCopied] = useState(false);

    const link = `${baseUrl}?share=${shareLink}`;

    useEffect(() => {
        setTimeout(() => setHasCopied(false), 2000)
    }, [hasCopied])

    useEffect(() => {
        if (!open || !activeSandbox) return;

        setShareLink(null);
        setLoading(true);
        setError(null);

        postShare(activeSandbox)
            .then(result => setShareLink(result?.id ?? null))
            .catch(error => setError((error as Error)?.message))
            .finally(() => setLoading(false));

    }, [open, activeSandbox]);

    const onOpenChange = (open: boolean) => {
        setError(null);
        setOpen(open);
    }

    const onCopyButtonClick = () => {
        copyToClipboardWithMeta(link)
            .then(() => setHasCopied(true))
            .catch(error => console.error(error))
    }

    return <Dialog open={open} defaultOpen={true} onOpenChange={onOpenChange}>
        <DialogContent className="sm:max-w-md">
            <DialogHeader>
                <DialogTitle>Share</DialogTitle>
                <DialogDescription>Send share link</DialogDescription>
            </DialogHeader>

            {error &&
                <div className="mb-3 text-sm text-destructive">Something went wrong</div>}

            {loading === true &&
                <div className="mb-3 text-sm">Loading ... Please wait</div>}

            {loading === false && shareLink &&
                <div className="flex space-x-2">
                    <Input value={link} type="text" autoFocus readOnly className="text-ellipsis whitespace-nowrap" />
                    <Button
                        size="icon"
                        variant={"outline"}
                        className={cn()}
                        onClick={onCopyButtonClick}
                    >
                        <span className="sr-only">Copy</span>
                        {hasCopied ? <CheckIcon size={17} /> : <ClipboardIcon size={17} />}
                    </Button>
                </div>
            }

        </DialogContent>
    </Dialog>
}

export default DialogCopyShare;
