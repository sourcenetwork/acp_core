import { useSandbox } from "@/hooks/useSandbox";
import { useDebounce } from "@/hooks/useDebounce";
import { usePlaygroundStore } from "@/lib/playgroundStore";
import { mapLocatedMessageMarkers } from "@/utils/mapLocatedMessageMarkers";
import { mapTheoremResultMarkers } from "@/utils/mapTheoremResultMarkers";
import { SandboxData, SandboxDataErrors } from "@acp/sandbox";
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
    const dataType = SandboxDataType[sandboxDataType];

    const monacoRef = useMonaco();
    const editorRef = useRef<monaco.editor.IStandaloneCodeEditor | null>(null);
    const { theme } = useTheme();
    const activeSandbox = useSandbox();
    const [dataErrors, annotatedPolicyTheoremResult, updateSandboxState] = usePlaygroundStore((state) => [state.setStateDataErrors?.[dataType.errorKey], state.verifyTheoremsResult, state.setPlaygroundState]);

    const editorData = activeSandbox?.data?.[sandboxDataType];
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
        void updateSandboxState({ [sandboxDataType]: value });
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
                defaultValue={editorData}
                value={editorData}
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

