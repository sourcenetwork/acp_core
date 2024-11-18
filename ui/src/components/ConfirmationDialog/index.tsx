import { Button } from "../ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "../ui/dialog";

interface ConfirmationDialogProps {
    title?: string,
    description?: string,
    open: boolean,
    onAction: (state: boolean) => unknown
}

const ConfirmationDialog = ({ title, description, open, onAction }: ConfirmationDialogProps) => {

    return <Dialog open={open} defaultOpen={true} onOpenChange={(state) => state === false && onAction(false)}>
        <DialogContent className="sm:max-w-md">
            <DialogHeader>
                <DialogTitle>{title}</DialogTitle>
                <DialogDescription>{description}</DialogDescription>
            </DialogHeader>
            <DialogFooter className="sm:justify-end">
                <Button type="button" variant="default" onClick={() => onAction(true)}>Confirm</Button>
                <Button type="button" variant="secondary" onClick={() => onAction(false)}>Close</Button>
            </DialogFooter>
        </DialogContent>
    </Dialog>
}

export default ConfirmationDialog;
