import { Button } from "../ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "../ui/dialog";

interface LoadTemplateDialogProps {
    title?: string,
    description?: string,
    open: boolean,
    onAction: (type: 'new' | 'replace' | false) => unknown
}

const LoadTemplateDialog = ({ title, description, open, onAction }: LoadTemplateDialogProps) => {

    return <Dialog open={open} defaultOpen={true} onOpenChange={(state) => state === false && onAction(false)}>
        <DialogContent className="sm:max-w-md">
            <DialogHeader>
                <DialogTitle>{title}</DialogTitle>
                <DialogDescription>{description}</DialogDescription>
            </DialogHeader>
            <DialogFooter className="sm:justify-between">
                <Button type="button" variant="secondary" onClick={() => onAction('replace')}>Override</Button>
                <div className="flex flex-col-reverse sm:flex-row sm:space-x-2 sm:justify-between">
                    <Button className="" type="button" variant="default" onClick={() => onAction('new')}>Load</Button>
                </div>
            </DialogFooter>
        </DialogContent>
    </Dialog>
}

export default LoadTemplateDialog;
