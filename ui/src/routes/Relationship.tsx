import { Button } from "@/components/ui/button";
import BaseEditor from "../components/Editor";

const Relationship = () => {

    return <>
        <div className="mb-2 flex justify-end gap-x-2">
            <Button variant="outline" asChild><a href="https://docs.source.network" target="_blank" rel="noreferrer">Guide</a></Button>
        </div>
        <BaseEditor sandboxDataType={"relationships"} />
    </>;
}

export default Relationship;

