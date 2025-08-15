// TODO: Add back a button to focus the editor tab for the output message

interface OutputMessageProps {
    prefix?: string;
    path?: string;
    message: string
}

const OutputMessage = ({ prefix, message }: OutputMessageProps) => {
    return <div className="flex mb-2">
        <span className="inline-block mr-3">{`>`}</span>
        <span className="">
            <span className="underline mr-2">{prefix}</span>
            <span>{`${message}`}</span>
        </span>
    </div>
}

export default OutputMessage;