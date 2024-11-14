import { VFC, useRef, useState, useEffect } from 'react';
import * as monaco from 'monaco-editor/esm/vs/editor/editor.api';
import styles from './Editor.module.css';
import { PlaygroundService } from '../../types/sourcenetwork/acp_core/playground';
import { SandboxDataErrors } from '../../types/sourcenetwork/acp_core/sandbox';
import { LocatedMessage } from '../../types/sourcenetwork/acp_core/parser_message';
import { AnnotatedAuthorizationTheoremResult, AnnotatedDelegationTheoremResult, ResultStatus } from '../../types/sourcenetwork/acp_core/theorem';

const policy = `name: test
resources:
  file:
    relations:
	  owner:
	    types:
		  - actor
	  reader:
	    types:
		  - actor
	permissions:
	  read:
	    expr: owner + reader
	  write:
	    expr: owner
`
const relationships = `file:readme#owner@did:user:bob
`

const theorems = `Authorizations {
  file:readme#read@did:user:bob
  file:readme#read@did:user:alice
}

Delegations {   
}
`

function highlightSelected(menu: HTMLElement, newIdx: number, setSelected) {
	setSelected((oldIdx: number) => {
		menu.children[oldIdx].className = ""
		menu.children[newIdx].className = styles.Selected;
		return newIdx
	})
}

export const Editor: VFC = () => {
	const [editor, setEditor] = useState<monaco.editor.IStandaloneCodeEditor | null>(null);
	const policyModel = useRef(monaco.editor.createModel(policy, "yaml"));
	const relationshipsModel = useRef(monaco.editor.createModel(relationships));
	const theoremsModel = useRef(monaco.editor.createModel(theorems));
	const monacoEl = useRef(null);

	const [playground, setPlayground] = useState<PlaygroundService | null>(null);

	const menuEl = useRef<HTMLElement | null>(null);
	const [_, setSelected] = useState(0);

	useEffect(() => {
		if (monacoEl) {
			setEditor((editor) => {
				if (editor) return editor;
				highlightSelected(menuEl.current!, 0, setSelected)
				return monaco.editor.create(monacoEl.current!, {
					model: policyModel.current,
				});
			});
		}

		(async () => {
			let playground = await window.AcpPlayground.new();
			setPlayground((p) => {
				if (p) return p;
				(async () => {
					let response = await playground.NewSandbox({name: "test"});
					console.log("sandbox created");
					console.log(response);
				})();
				return playground;
			})
			console.log("playground created");
		})();

		return () => editor?.dispose();
	}, [monacoEl.current]);

	function changeModel(model: monaco.editor.ITextModel, idx: number) {
		editor?.setModel(model);
		highlightSelected(menuEl.current!, idx, setSelected)
		editor?.focus();
	}

	function setPlaygroundState() {
		(async () => { 
			let setResult = await playground!.SetState({
			handle: 1,
			data: {
				policy_definition: policyModel.current.getValue(),
				relationships: relationshipsModel.current.getValue(),
				policy_theorem: theoremsModel.current.getValue(),
			}
		 });

		 console.log(setResult);

		let policyMarkers = [];
		let relationshipMarkers = [];
		let theoremMarkers = [];

		if (setResult && !setResult.ok) {
			policyMarkers = setResult.errors.policyErrors.map(mapLocatedMessage);
			relationshipMarkers = setResult.errors.relationshipsErrors.map(mapLocatedMessage);
			theoremMarkers =  setResult.errors.theoremsErrors.map(mapLocatedMessage);
		}
		monaco.editor.setModelMarkers(policyModel.current, "owner", policyMarkers);
		monaco.editor.setModelMarkers(relationshipsModel.current, "owner", relationshipMarkers);
		monaco.editor.setModelMarkers(theoremsModel.current, "owner", theoremMarkers);
		})();
	}


	function verifyTheorems() {
		setPlaygroundState();
		(async () => {
		 let verifyResult = await playground!.VerifyTheorems({
			handle: 1
		 });

		console.log(verifyResult)

		 if (verifyResult) {
			verifyResult.result?.authorizationTheoremsResult
			let authMarkers = verifyResult.result?.authorizationTheoremsResult.filter((result: AnnotatedAuthorizationTheoremResult) => result.result?.result?.status != "Accept")
			.map((result: AnnotatedAuthorizationTheoremResult): monaco.editor.IMarkerData => {
					return {
						message: result.result?.result?.status.toString(),
						startLineNumber: Number(result.interval?.start?.line!),
						endLineNumber: Number(result.interval?.end?.line!),
						startColumn: Number(result.interval?.start?.column!),
						endColumn: Number(result.interval?.end?.column!),
						severity: monaco.MarkerSeverity.Error,
					}
			});

			let delegationMarkers = verifyResult.result?.delegationTheoremsResult.filter((result: AnnotatedDelegationTheoremResult) => result.result?.result?.status != "Accept")
			.map((result: AnnotatedDelegationTheoremResult): monaco.editor.IMarkerData => {
					return {
						message: result.result?.result?.status.toString(),
						startLineNumber: Number(result.interval?.start?.line!),
						endLineNumber: Number(result.interval?.end?.line!),
						startColumn: Number(result.interval?.start?.column!),
						endColumn: Number(result.interval?.end?.column!),
						severity: monaco.MarkerSeverity.Error,
					}
			});

			let markers: monaco.editor.IMarkerData[] = [];
			markers.push(...authMarkers);
			markers.push(...delegationMarkers);
			console.log(markers);

			monaco.editor.setModelMarkers(theoremsModel.current, "owner", markers);
		 }
		})();
	}


	const focusPolicy = () => changeModel(policyModel.current, 0);
	const focusRelationships = () => changeModel(relationshipsModel.current, 1);
	const focusTheorems = () => changeModel(theoremsModel.current, 2);

	policyModel.current.onDidChangeContent((e: monaco.editor.IModelContentChangedEvent) => setPlaygroundState());
	relationshipsModel.current.onDidChangeContent((e: monaco.editor.IModelContentChangedEvent) => setPlaygroundState());
	theoremsModel.current.onDidChangeContent((e: monaco.editor.IModelContentChangedEvent) => setPlaygroundState());

	return (
	<div>
		<div>
			<ul id="menu" className="flex-container" ref={menuEl}>
				<li onClick={focusPolicy} >Policy </li>
				<li onClick={focusRelationships} >Relationships</li>
				<li onClick={focusTheorems} >Theorems</li>
			</ul>
		</div>
		<div className={styles.Editor} ref={monacoEl}>
	</div>
	<div onClick={verifyTheorems}>Run</div>
	</div>
	);
};

function mapLocatedMessage(msg: LocatedMessage): monaco.editor.IMarkerData {
	console.log(typeof(msg.interval?.start.column))
	return {
		message: msg.message,
		startLineNumber: Number(msg.interval?.start?.line!),
		endLineNumber: Number(msg.interval?.end?.line!),
		startColumn: Number(msg.interval?.start?.column!),
		endColumn: Number(msg.interval?.end?.column!),
		severity: monaco.MarkerSeverity.Error,
	}
}