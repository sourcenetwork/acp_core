import { useSandbox } from "@/hooks/useSandbox";
import { usePlaygroundStore } from "@/stores/playgroundStore";
import { zodResolver } from "@hookform/resolvers/zod";
import { useEffect, useMemo } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { FormInputField } from "../FormInput";
import { Button } from "../ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "../ui/dialog";
import { Form } from "../ui/form";

interface DialogEditSandboxProps {
    sandboxId?: string,
    open: boolean,
    setOpen: (state: boolean) => unknown
}

const EditSandboxFormSchema = z.object({
    name: z.string().min(1),
    description: z.string()
});

type EditSandboxFormData = z.infer<typeof EditSandboxFormSchema>;

const DialogEditSandbox = ({ sandboxId, open, setOpen }: DialogEditSandboxProps) => {

    const sandbox = useSandbox(sandboxId);
    const updateStoredSandbox = usePlaygroundStore((state) => state.updateStoredSandbox);
    const sandboxData = useMemo(() => ({ name: sandbox?.name ?? "", description: sandbox?.description ?? "", }), [sandbox?.name, sandbox?.description]);

    const form = useForm<EditSandboxFormData>({
        resolver: zodResolver(EditSandboxFormSchema),
        defaultValues: {
            name: sandbox?.name,
            description: sandbox?.description,
        },
    })

    useEffect(() => {
        form.reset(sandboxData);
    }, [sandboxId, sandboxData, form]);

    const onSubmit = (data: EditSandboxFormData) => {
        if (sandboxId) updateStoredSandbox(data, sandboxId);
        setOpen(false);
        form.reset();
    };

    return <Dialog open={open} defaultOpen={true} onOpenChange={(state) => state === false && setOpen(false)}>
        <DialogContent className="sm:max-w-md">
            <DialogHeader>
                <DialogTitle>Edit Sandbox</DialogTitle>
                <DialogDescription>Update sandbox details</DialogDescription>
            </DialogHeader>
            <Form {...form}>
                <form onSubmit={(event) => void form.handleSubmit(onSubmit)(event)}>
                    <div className="grid gap-y-2 mb-4">
                        <FormInputField name="name" placeholder="Enter Name" control={form.control} />
                        <FormInputField name="description" placeholder="Enter Description" control={form.control} />
                    </div>
                    <DialogFooter className="sm:justify-end ">
                        <Button type="submit" variant="default">Save</Button>
                        <Button type="button" variant="secondary" onClick={() => setOpen(false)}>Close</Button>
                    </DialogFooter>
                </form>
            </Form>
        </DialogContent>
    </Dialog>
}

export default DialogEditSandbox;
