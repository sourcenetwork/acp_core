import { useUIActions } from "@/stores/layoutStore";
import { usePlaygroundStore } from "@/stores/playgroundStore";
import { cn } from "@/utils/classnames";
import { zodResolver } from "@hookform/resolvers/zod";
import { CircleCheck, CircleX, Maximize2, Minimize2 } from "lucide-react";
import { Suspense, lazy, useEffect, useMemo, useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";
import FormSelectField from "../FormSelectField";
import PaneHeader from "../PaneHeader";
import TextTooltip from "../TextTooltip";
import { Button } from "../ui/button";
import { Form } from "../ui/form";

const ExplainCheckGraphLazy = lazy(() => import('../ExplainCheckGraph'));

const ExplainCheckFormSchema = z.object({
    resourceId: z.string().nonempty(),
    resourceType: z.string().nonempty(),
    permission: z.string().nonempty(),
    actorId: z.string().nonempty(),
});

type ExplainCheckFormData = z.infer<typeof ExplainCheckFormSchema>;

const PaneCheck = () => {
    const { setSecondaryPaneOpen } = useUIActions();
    const getCatalogue = usePlaygroundStore((state) => state.getCatalogue);
    const activeHandle = usePlaygroundStore((state) => state.activeHandle);
    const catalogue = usePlaygroundStore((state) => state.catalogue);
    const explainCheck = usePlaygroundStore((state) => state.explainCheck);
    const explainCheckResult = usePlaygroundStore((state) => state.explainCheckResult);
    const explainCheckError = usePlaygroundStore((state) => state.explainCheckError);
    const explainCheckFormData = usePlaygroundStore((state) => state.explainCheckFormData);
    const setExplainCheckFormData = usePlaygroundStore((state) => state.setExplainCheckFormData);
    const [maximize, setMaximize] = useState(false);

    const form = useForm<ExplainCheckFormData>({
        resolver: zodResolver(ExplainCheckFormSchema),
        defaultValues: explainCheckFormData ?? {
            resourceId: "",
            resourceType: "",
            permission: "",
            actorId: "",
        },
    })

    useEffect(() => {
        if (!activeHandle) return;
        void getCatalogue(activeHandle);
    }, [activeHandle]);

    const onSubmit = (data: ExplainCheckFormData) => {
        if (!activeHandle) return;

        void explainCheck({
            handle: activeHandle,
            object: { id: data.resourceId, resource: data.resourceType },
            permission: data.permission,
            actor: { id: data.actorId },
        });

        setExplainCheckFormData(data);
    }

    const selectedResourceType = form.watch('resourceType');

    const selectOptions = useMemo(() => {
        if (!catalogue) {
            return {
                resourceTypes: [],
                resourceIds: [],
                permissions: [],
                actors: []
            };
        }

        const resourceTypes: { value: string; label: string }[] = [];

        let selectedResource = null;

        for (const [type, resource] of Object.entries(catalogue.resourceCatalogue)) {
            resourceTypes.push({ value: type, label: type });
            if (type === selectedResourceType) {
                selectedResource = resource;
            }
        }

        const resourceIds = selectedResource?.objectIds.map(id => ({ value: id, label: id })) ?? [];
        const permissions = selectedResource?.permissions.map(permission => ({ value: permission, label: permission })) ?? [];
        const actors = catalogue.actors.map(actor => ({ value: actor, label: actor }));

        return { resourceTypes, resourceIds, permissions, actors };
    }, [catalogue, selectedResourceType]);

    const isAuthorized = useMemo(() => {
        return explainCheckResult?.authorized;
    }, [explainCheckResult]);

    return <div className="p-3 flex flex-col gap-2 h-full">

        <div className="">
            <PaneHeader showCollapse title="Check" direction="right" onCollapseClick={() => setSecondaryPaneOpen(false)} />

            <div className={cn("border-b border-divider pb-3", {
                "fixed top-4 left-4 right-20 z-101 max-w-[600px] p-4 border border-divider rounded-lg bg-background shadow-lg/20": maximize
            })}>
                <Form {...form}>
                    <form onSubmit={(event) => void form.handleSubmit(onSubmit)(event)}>
                        <div className="grid grid-cols-2 gap-2">
                            <FormSelectField className="w-full" name="resourceType" placeholder="Enter Resource Type" control={form.control} options={selectOptions.resourceTypes} disabled={!selectOptions.resourceTypes.length} clearInvalid />
                            <FormSelectField className="w-full" name="resourceId" placeholder="Enter Resource ID" control={form.control} options={selectOptions.resourceIds} disabled={!selectOptions.resourceIds.length} clearInvalid />
                        </div>
                        <div className="grid grid-cols-2 gap-2">
                            <FormSelectField className="w-full" name="permission" placeholder="Enter Permission" control={form.control} options={selectOptions.permissions} disabled={!selectOptions.permissions.length} clearInvalid />
                            <FormSelectField className="w-full" name="actorId" placeholder="Enter Actor ID" control={form.control} options={selectOptions.actors} disabled={!selectOptions.actors.length} clearInvalid />
                        </div>

                        <div className="flex align-baseline justify-between gap-3">
                            <div className="self-end">
                                {explainCheckResult &&
                                    <div className={cn("text-sm flex items-center gap-2",
                                        isAuthorized ? "text-src-secondary" : "text-src-error"
                                    )}>
                                        {isAuthorized ? <CircleCheck className="size-4" /> : <CircleX className="size-4" />}
                                        {isAuthorized ? "Authorized" : "Not Authorized"}
                                    </div>
                                }
                            </div>

                            <Button variant="default" type="submit" size="sm" className="mt-2  justify-self-end" disabled={!form.formState.isValid}>Run Check</Button>
                        </div>
                    </form>
                </Form>

                {explainCheckError &&
                    <div className="text-sm mt-4">
                        <div className="font-medium text-src-error">Error:</div>
                        <div>{explainCheckError}</div>
                    </div>
                }
            </div>
        </div>

        {!explainCheckError && explainCheckResult?.graph && <div className={cn(
            "grow relative",
            maximize && "fixed inset-2 z-100"
        )}>
            <div className="absolute top-1 right-1 z-1">
                <TextTooltip content={maximize ? "Minimize" : "Maximize"}>
                    <Button variant="ghost" size="iconSm" onClick={() => setMaximize(!maximize)}>
                        {maximize ? <Minimize2 size={16} /> : <Maximize2 size={16} />}
                    </Button>
                </TextTooltip>
            </div>

            <Suspense fallback={<div className="py-4 text-sm text-muted-foreground">Loading...</div>}>
                <ExplainCheckGraphLazy explainGraph={explainCheckResult?.graph ?? ""} className="h-full border border-divider overflow-hidden rounded-lg" maximized={maximize} />
            </Suspense>
        </div>}
    </div>;
};

export default PaneCheck;
