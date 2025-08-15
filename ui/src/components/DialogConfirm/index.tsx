import { Button } from "../ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "../ui/dialog";

interface DialogCopyShareProps {
    open: boolean,
    setOpen: (state: boolean) => unknown,
    onConfirm: (state: boolean) => void
    title: string,
    description: string
}

const DialogConfirm = ({ open, setOpen, onConfirm, title, description }: DialogCopyShareProps) => {

    const onConfirmClick = (state: boolean) => {
        onConfirm(state);
        setOpen(false);
    }

    return <Dialog open={open} defaultOpen={true} onOpenChange={(state) => state === false && setOpen(false)}>
        <DialogContent className="sm:max-w-md">
            <DialogHeader>
                <DialogTitle>{title}</DialogTitle>
                <DialogDescription>{description}</DialogDescription>
            </DialogHeader>

            <DialogFooter className="sm:justify-end ">
                <Button type="submit" variant="default" onClick={() => onConfirmClick(true)}>Confirm</Button>
                <Button type="button" variant="secondary" onClick={() => onConfirmClick(false)}>Close</Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>
}

export default DialogConfirm;
