
import { usePlaygroundStore } from "@/stores/playgroundStore";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { z } from "zod";
import { FormInputField } from "../FormInput";
import { Button } from "../ui/button";
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from "../ui/dialog";
import { Form } from "../ui/form";

interface DialogCreateSandboxProps {
    open: boolean,
    setOpen: (state: boolean) => unknown
}

const CreateSandboxFormSchema = z.object({
    name: z.string().min(1),
    description: z.string()
});

type CreateSandboxFormData = z.infer<typeof CreateSandboxFormSchema>;

const DialogCreateSandbox = ({ open, setOpen }: DialogCreateSandboxProps) => {

    const newSandbox = usePlaygroundStore((state) => state.newSandbox);

    const form = useForm<CreateSandboxFormData>({
        resolver: zodResolver(CreateSandboxFormSchema),
        defaultValues: {
            name: "New Sandbox",
            description: "",
        },
    })

    const onSubmit = (data: CreateSandboxFormData) => {
        void newSandbox(data);
        form.reset();
        setOpen(false);
    };

    return <Dialog open={open} defaultOpen={false} onOpenChange={(state) => state === false && setOpen(false)}>
        <DialogContent className="sm:max-w-md">
            <DialogHeader>
                <DialogTitle>Create Sandbox</DialogTitle>
                <DialogDescription>Create a new sandbox</DialogDescription>
            </DialogHeader>
            <Form {...form}>
                <form onSubmit={(event) => void form.handleSubmit(onSubmit)(event)}>
                    <div className="grid gap-y-2 mb-4">
                        <FormInputField name="name" placeholder="Enter Name" control={form.control} />
                        <FormInputField name="description" placeholder="Enter Description" control={form.control} />
                    </div>
                    <DialogFooter className="sm:justify-end ">
                        <Button type="submit" variant="default">Create</Button>
                        <Button type="button" variant="secondary" onClick={() => setOpen(false)}>Close</Button>
                    </DialogFooter>
                </form>
            </Form>
        </DialogContent>
    </Dialog>
}

export default DialogCreateSandbox;
