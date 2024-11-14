import React from 'react';
import ReactDOM from 'react-dom';
import { Editor } from './components/editor/Editor';
import './userWorker';
import './runtime/wasm_exec.js';
import { loadGoWASM } from './utils/wasm';


const module = await loadGoWASM("playground.wasm");


ReactDOM.render(
	<React.StrictMode>
		<Editor/>
	</React.StrictMode>,
	document.getElementById('root')
);
