import { useDebounce } from "@/hooks/useDebounce";
import { useSandbox } from "@/hooks/useSandbox";
import { definePolicyTheoremTheme, POLICY_THEOREM_LANGUAGE_ID, registerPolicyTheoremLanguage } from "@/lib/languagePolicyTheorem";
import { usePlaygroundStore } from "@/lib/playgroundStore";
import { useTheme } from "@/providers/ThemeProvider/useTheme";
import { mapLocatedMessageMarkers } from "@/utils/mapLocatedMessageMarkers";
import { mapTheoremResultMarkers } from "@/utils/mapTheoremResultMarkers";
import { SandboxData, SandboxDataErrors } from "@acp/sandbox";
import { Editor, EditorProps, Monaco, OnChange, OnMount, useMonaco } from "@monaco-editor/react";
import * as monaco from 'monaco-editor';
import { useCallback, useEffect, useRef, useState } from "react";

interface BaseEditorProps {
    sandboxDataType: keyof SandboxData,
}

enum SandboxType {
    POLICY_DEFINITION = "policyDefinition",
    RELATIONSHIPS = "relationships",
    POLICY_THEOREM = "policyTheorem"
}

const SandboxDataType: Record<SandboxType, { theme: { dark: string, light: string }, language: string, errorKey: keyof SandboxDataErrors }> = {
    [SandboxType.POLICY_DEFINITION]: {
        theme: { dark: 'vs-dark', light: 'vs-light' },
        language: 'yaml',
        errorKey: "policyErrors"
    },
    [SandboxType.RELATIONSHIPS]: {
        theme: { dark: 'vs-dark', light: 'vs-light' },
        language: "yaml",
        errorKey: "relationshipsErrors"
    },
    [SandboxType.POLICY_THEOREM]: {
        theme: { dark: 'policyTheoremDark', light: 'policyTheoremLight' },
        language: POLICY_THEOREM_LANGUAGE_ID,
        errorKey: "theoremsErrors"
    }
};

const BaseEditor = (props: EditorProps & BaseEditorProps) => {
    const { sandboxDataType } = props;
    const dataType = SandboxDataType[sandboxDataType];

    const monacoRef = useMonaco();

    const editorRef = useRef<monaco.editor.IStandaloneCodeEditor | null>(null);
    const decorationsRef = useRef<monaco.editor.IEditorDecorationsCollection | null>(null);
    const [isEditorMounted, setIsEditorMounted] = useState(false);

    const { theme } = useTheme();
    const activeSandbox = useSandbox();

    const {
        dataErrors,
        annotatedPolicyTheoremResult,
        updateSandboxState,
        verifyTheorems,
        sandboxStateStatus,
        setEditorSelection,
        editorSelection
    } = usePlaygroundStore((state) => ({
        dataErrors: state.setStateDataErrors?.[dataType.errorKey],
        annotatedPolicyTheoremResult: state.verifyTheoremsResult,
        updateSandboxState: state.setPlaygroundState,
        verifyTheorems: state.verifyTheorems,
        sandboxStateStatus: state.sandboxStateStatus,
        setEditorSelection: state.setEditorSelection,
        editorSelection: state.editorSelections[sandboxDataType]
    }));

    const editorData = activeSandbox?.data?.[sandboxDataType];
    const isTheorumEditor = sandboxDataType === SandboxType.POLICY_THEOREM;

    const updateMarkers = useCallback(() => {
        const editor = editorRef.current;
        if (!editor) return;

        const errorMarkers = mapLocatedMessageMarkers(dataErrors ?? []);

        monacoRef?.editor.setModelMarkers(editor.getModel()!, 'owner', errorMarkers);

        // Update test result markers and glyphs
        if (isTheorumEditor) {

            // If there are input errors, clear the decorations as they won't align
            if (errorMarkers?.length > 0) {
                clearGlyphDecorations();
                return;
            }

            const theoremResultMarkers = mapTheoremResultMarkers(annotatedPolicyTheoremResult);

            if (!theoremResultMarkers) return;

            const { authMarkers, delegationMarkers } = theoremResultMarkers;

            monacoRef?.editor.setModelMarkers(editor.getModel()!, 'owner', [
                ...authMarkers.rejected,
                ...authMarkers.errors,
                ...delegationMarkers.rejected,
                ...delegationMarkers.errors,
            ]);

            clearGlyphDecorations();

            const decorations = [
                ...createGlyphDecorations([...authMarkers.accepted, ...delegationMarkers.accepted], 'passing'),
                ...createGlyphDecorations([...authMarkers.rejected, ...delegationMarkers.rejected], 'failing')
            ];

            if (decorations.length > 0) decorationsRef.current = editor.createDecorationsCollection(decorations);
        }
    }, [dataErrors, annotatedPolicyTheoremResult, isTheorumEditor, monacoRef]);

    const clearGlyphDecorations = () => {
        decorationsRef.current?.clear();
    };

    const createGlyphDecorations = (markers: monaco.editor.IMarkerData[] = [], type: 'failing' | 'passing'): monaco.editor.IModelDeltaDecoration[] => {
        return markers.map(marker => ({
            range: new monaco.Range(marker.startLineNumber, 1, marker.startLineNumber, 1),
            options: {
                isWholeLine: true,
                glyphMarginClassName: type === 'failing' ? 'failing-glyph' : 'passing-glyph',
                className: type === 'failing' ? 'failing-line-bg' : '',
            }
        }));
    };

    const handleEditorChange: OnChange = useDebounce((value) => {
        void updateSandboxState({ [sandboxDataType]: value });
        if (isTheorumEditor) void verifyTheorems();
    }, 500);

    const handleEditorMounted: OnMount = (editor) => {
        editorRef.current = editor;
        setIsEditorMounted(true);

        // Restore selection
        if (editorSelection) {
            editor.setSelection(editorSelection);
            editor.focus();
        }
    }

    const handleBeforeMount = (monaco: Monaco) => {
        registerPolicyTheoremLanguage(monaco);
        definePolicyTheoremTheme(monaco);
    }

    useEffect(() => {
        if (!isEditorMounted) return;

        updateMarkers();
    }, [dataErrors, annotatedPolicyTheoremResult, updateMarkers, isEditorMounted]);

    useEffect(() => {
        if (sandboxStateStatus !== "set") return;
        if (isTheorumEditor) void verifyTheorems();
    }, [sandboxStateStatus, isTheorumEditor, verifyTheorems, activeSandbox?.id]);


    // Track cursor position and selection 
    useEffect(() => {
        const editor = editorRef.current;
        if (!isEditorMounted || !editor) return;

        const selectionEvent = editor.onDidChangeCursorSelection((e) => setEditorSelection(sandboxDataType, e.selection))
        return () => {
            selectionEvent.dispose();
        };
    }, [isEditorMounted, setEditorSelection, sandboxDataType]);

    const editorLanguage = dataType.language || 'yaml';
    const editorTheme = dataType.theme[theme === 'dark' ? 'dark' : 'light'] || 'vs-dark';

    return <div className='h-full'>
        <div className='py-5 h-full rounded-md overflow-hidden bg-editor border'>
            <Editor
                height="100%"
                defaultLanguage={editorLanguage}
                defaultValue={editorData}
                value={editorData}
                onChange={handleEditorChange}
                beforeMount={handleBeforeMount}
                onMount={handleEditorMounted}
                theme={editorTheme}
                options={{
                    automaticLayout: true,
                    fixedOverflowWidgets: true,
                    glyphMargin: isTheorumEditor,
                    tabSize: 2,
                    ...props?.options
                }}
                {...props} />
        </div>
    </div>
}

export default BaseEditor;
