import { usePlaygroundStore } from "@/lib/playgroundStore";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { FormInputField } from "../FormInput";
import { Button } from "../ui/button";
import { Dialog, DialogContent, DialogFooter, DialogHeader, DialogTitle } from "../ui/dialog";
import { Form } from "../ui/form";

interface ConfirmationDialogProps {
    open: boolean,
    setOpen: (state: boolean) => unknown
}

const FormSchema = z.object({
    name: z.string().min(2),
    description: z.string()
});

const CreateSandboxDialog = ({ open, setOpen }: ConfirmationDialogProps) => {
    const { newSandbox } = usePlaygroundStore();

    const form = useForm<z.infer<typeof FormSchema>>({
        resolver: zodResolver(FormSchema),
        defaultValues: {
            name: "",
            description: "",
        },
    })

    const onSubmit = (data: any) => {
        newSandbox(data);
        setOpen(false);
        form.reset();
    };

    return <Dialog open={open} defaultOpen={true} onOpenChange={(state) => state === false && setOpen(false)}>
        <DialogContent className="sm:max-w-md">
            <DialogHeader>
                <DialogTitle>Create New Sandbox</DialogTitle>
            </DialogHeader>
            <Form {...form}>
                <form onSubmit={form.handleSubmit(onSubmit)}>
                    <div className="grid gap-y-2 mb-4">
                        <FormInputField name="name" placeholder="Enter Name" control={form.control} />
                        <FormInputField name="description" placeholder="Enter Description" control={form.control} />
                    </div>
                    <DialogFooter className="sm:justify-end">
                        <Button type="submit" variant="default">Create</Button>
                        <Button type="button" variant="secondary" onClick={() => setOpen(false)}>Close</Button>
                    </DialogFooter>
                </form>
            </Form>
        </DialogContent>
    </Dialog>
}

export default CreateSandboxDialog;
