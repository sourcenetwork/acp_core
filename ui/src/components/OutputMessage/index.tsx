import { Link } from "react-router-dom";

interface OutputMessageProps {
    prefix?: string;
    path?: string;
    message: string
}

const OutputMessage = ({ path, prefix, message }: OutputMessageProps) => {
    return <div className="flex mb-2">
        <span className="inline-block mr-3">{`>`}</span>
        <span className="">
            {path && prefix ? (<Link to={path} className="underline mr-2">{prefix}</Link>) : (<span className="underline mr-2">{prefix}</span>)}
            <span>{`${message}`}</span>
        </span>
    </div>
}

export default OutputMessage;