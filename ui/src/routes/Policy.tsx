import { Button } from "@/components/ui/button";
import BaseEditor from "../components/Editor";

const Policy = () => {

    return <>
        <div className="mb-2 flex justify-end gap-x-2">
            <Button variant="outline" asChild><a href="https://docs.source.network" target="_blank">Guide</a></Button>
        </div>
        <BaseEditor sandboxDataType={"policyDefinition"} />
    </>;
}

export default Policy;
