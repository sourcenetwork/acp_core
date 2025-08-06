import { usePlaygroundStore } from "@/stores/playgroundStore";
import { SandboxDataErrors } from "@acp/sandbox";
import OutputMessage from "../OutputMessage";
import { Badge } from "../ui/badge";
import { ScrollArea } from "../ui/scroll-area";

const Output = () => {
    const dataErrors = usePlaygroundStore((state) => state.setStateDataErrors);
    const setStateError = usePlaygroundStore((state) => state.setStateError);
    const verifyTheoremsError = usePlaygroundStore((state) => state.verifyTheoremsError);
    const playgroundStatus = usePlaygroundStore((state) => state.playgroundStatus);

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

    const problemCount =
        (errorMessages?.length ?? 0) +
        (setStateError ? 1 : 0) +
        (verifyTheoremsError ? 1 : 0);

    return <div className="py-2 h-full flex flex-col">
        <div className="pb-2 text-sm leading-none font-light items-center flex min-h-8">
            Problems
            {problemCount ? <Badge className="ml-2" variant="outline">{problemCount}</Badge> : null}
        </div>
        <ScrollArea className="p-4 bg-editor flex-1 rounded-md border font-mono text-[12px]">

            {playgroundStatus === 'loading' && <OutputMessage message={"Playground Loading ..."} />}

            {!!setStateError && <OutputMessage message={setStateError} />}

            {!!verifyTheoremsError && <OutputMessage message={verifyTheoremsError} />}

            {errorMessages?.map((error, index) => <OutputMessage key={index} path={error.path} prefix={error.prefix} message={error.message} />)}

        </ScrollArea>
    </div>
}

export default Output;

