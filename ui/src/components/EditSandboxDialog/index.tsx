import { usePlaygroundStore } from "@/lib/playgroundStore";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { FormInputField } from "../FormInput";
import { Button } from "../ui/button";
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from "../ui/dialog";
import { Form } from "../ui/form";
import { useEffect } from "react";

interface ConfirmationDialogProps {
    sandboxId: string | null,
    open: boolean,
    setOpen: (state: boolean) => unknown
}

const FormSchema = z.object({
    name: z.string().min(2),
    description: z.string()
});

const EditSandboxDialog = ({ sandboxId, open, setOpen }: ConfirmationDialogProps) => {
    const { findSandboxById, updateStoredSandbox } = usePlaygroundStore();
    const sandbox = findSandboxById(sandboxId);

    const form = useForm<z.infer<typeof FormSchema>>({
        resolver: zodResolver(FormSchema),
        defaultValues: {
            name: sandbox?.active?.name,
            description: sandbox?.active?.description,
        },
    })

    useEffect(() => {
        form.reset({
            name: sandbox?.active?.name || "",
            description: sandbox?.active?.description || "",
        });
    }, [sandboxId]);

    const onSubmit = (data: z.infer<typeof FormSchema>) => {
        updateStoredSandbox(data);
        setOpen(false);
        form.reset();
    };

    return <Dialog open={open} defaultOpen={true} onOpenChange={(state) => state === false && setOpen(false)}>
        <DialogContent className="sm:max-w-md">
            <DialogHeader>
                <DialogTitle>Edit Sandbox</DialogTitle>
            </DialogHeader>
            <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)}>
                    <div className="grid gap-y-2 mb-4">
                        <FormInputField name="name" placeholder="Enter Name" control={form.control} />
                        <FormInputField name="description" placeholder="Enter Description" control={form.control} />
                    </div>
                    <DialogFooter className="sm:justify-end">
                        <Button type="submit" variant="default">Save</Button>
                        <Button type="button" variant="secondary" onClick={() => setOpen(false)}>Close</Button>
                    </DialogFooter>
                </form>
            </Form>
        </DialogContent>
    </Dialog>
}

export default EditSandboxDialog;
