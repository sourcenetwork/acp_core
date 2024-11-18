import { useDebounce } from "@/hooks/useDebounce";
import { usePlaygroundStore } from "@/lib/acpHandler";
import { usePlaygroundStorageStore } from "@/lib/acpStorage";
import { SandboxData, SandboxDataErrors } from "@/types/proto-js/sourcenetwork/acp_core/sandbox";
import { mapLocatedMessageMarkers } from "@/utils/mapLocatedMessageMarkers";
import { mapTheoremResultMarkers } from "@/utils/mapTheoremResultMarkers";
import { Editor, EditorProps, OnChange, OnMount, useMonaco } from "@monaco-editor/react";
import * as monaco from 'monaco-editor';
import { useCallback, useEffect, useRef } from "react";
import { useTheme } from "../../ThemeProvider";

interface BaseEditorProps {
    sandboxDataType: keyof SandboxData,
}

const SandboxDataType: Record<keyof SandboxData, { errorKey: keyof SandboxDataErrors }> = {
    "policyDefinition": {
        errorKey: "policyErrors"
    },
    "relationships": {
        errorKey: "relationshipsErrors"
    },
    "policyTheorem": {
        errorKey: "theoremsErrors"
    }
};

const BaseEditor = (props: EditorProps & BaseEditorProps) => {
    const { sandboxDataType } = props;
    const monacoRef = useMonaco();
    const editorRef = useRef<monaco.editor.IStandaloneCodeEditor | null>(null);
    const { theme } = useTheme();
    const dataType = SandboxDataType[sandboxDataType];
    const [dataContent] = usePlaygroundStorageStore((state) => [state.data[sandboxDataType], state.updateStore]);
    const [dataErrors, annotatedPolicyTheoremResult, setState] = usePlaygroundStore((state) => [state.sandboxErrors?.[dataType.errorKey], state.annotatedPolicyTheoremResult, state.setState]);

    const isTestEditor = sandboxDataType === 'policyTheorem';

    const updateMarkers = useCallback(() => {
        const editor = editorRef.current;
        if (!editor) return;

        const errorMarkers = mapLocatedMessageMarkers(dataErrors ?? []);
        const theoremResultMarkers = isTestEditor ? mapTheoremResultMarkers(annotatedPolicyTheoremResult) : [];

        monacoRef?.editor.setModelMarkers(editor.getModel()!, 'owner', [
            ...errorMarkers,
            ...theoremResultMarkers
        ]);
    }, [dataErrors, annotatedPolicyTheoremResult, isTestEditor]);

    useEffect(() => {
        updateMarkers();
    }, [dataErrors, annotatedPolicyTheoremResult, updateMarkers])

    const handleEditorChange: OnChange = useDebounce((value) => {
        setState({ [sandboxDataType]: value });
    }, 500);

    const handleEditorMounted: OnMount = (editor) => {
        editorRef.current = editor;
        updateMarkers();
    }

    return <div className='h-full'>
        <div className='py-5 h-full rounded-md overflow-hidden bg-editor border'>
            <Editor
                height="100%"
                defaultLanguage="yaml"
                defaultValue={dataContent}
                value={dataContent}
                onChange={handleEditorChange}
                onMount={handleEditorMounted}
                theme={theme === 'dark' ? 'vs-dark' : 'vs-light'}
                options={{
                    automaticLayout: true,
                    fixedOverflowWidgets: true,
                    ...props?.options
                }}
                {...props} />
        </div>
    </div>
}

export default BaseEditor;

