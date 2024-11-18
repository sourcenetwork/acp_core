import { usePlaygroundStore } from "@/lib/acpHandler";
import { SandboxDataErrors } from "@/types/proto-js/sourcenetwork/acp_core/sandbox";
import OutputMessage from "../OutputMessage";
import { Badge } from "../ui/badge";
import { ScrollArea } from "../ui/scroll-area";

const Output = () => {
    const [dataErrors, setStateError] = usePlaygroundStore((state) => [state.sandboxErrors, state.setStateError]);

    const messageInfo: Record<keyof SandboxDataErrors, { prefix: string, path: string }> = {
        "policyErrors": {
            prefix: "Policy Error:",
            path: "/"
        },
        "theoremsErrors": {
            prefix: "Test Error:",
            path: "/tests"
        },
        "relationshipsErrors": {
            prefix: "Relationship Error:",
            path: "/relationship"
        },
    }

    const errorMessages = dataErrors &&
        Object.keys(dataErrors)
            .reduce<{ prefix: string, path: string, message: string }[]>(
                (messages, key) => {
                    const errorKey = key as keyof SandboxDataErrors;
                    const errors = dataErrors[errorKey];
                    const { prefix, path } = messageInfo[errorKey];
                    messages = [...messages, ...errors.map(({ message }) => ({ prefix, path, message }))]
                    return messages;
                }, []);

    return <div className="py-4 pr-4 h-full flex flex-col">
        <div className="pb-2 text-sm leading-none font-light">
            Problems
            {errorMessages?.length !== 0 && <Badge className="ml-2" variant="outline">{errorMessages?.length}</Badge>}
        </div>
        <ScrollArea className="p-4 bg-editor flex-1 rounded-md border font-mono text-[12px]">

            {!!setStateError && <OutputMessage message={"Something wen't wrong setting playground state"} />}

            {errorMessages?.map((error, index) => <OutputMessage key={index} path={error.path} prefix={error.prefix} message={error.message} />)}

        </ScrollArea>
    </div>
}

export default Output;

