import './wasm_exec.js';

let wasmReady = false;
let goInstance = null;
let readyPromise = null;

async function waitForGoFunctions(timeoutMs = 10000) {
	const started = Date.now();
	while (Date.now() - started < timeoutMs) {
		if (typeof globalThis.generateText === 'function' && typeof globalThis.getTemplates === 'function') {
			return;
		}
		await new Promise((resolve) => setTimeout(resolve, 25));
	}
	throw new Error('Go WASM did not register its functions in time');
}

export async function loadWasm() {
	if (wasmReady) return;
	if (readyPromise) return readyPromise;

	readyPromise = (async () => {
		if (typeof Go === 'undefined') {
			throw new Error('Go constructor not found; wasm_exec.js did not load');
		}

		goInstance = new Go();
		const response = await fetch('/main.wasm', { cache: 'no-cache' });
		if (!response.ok) {
			throw new Error(`Could not load main.wasm: HTTP ${response.status}`);
		}

		const module = await WebAssembly.instantiate(await response.arrayBuffer(), goInstance.importObject);
		// run() starts the Go program (which blocks by design) and resolves when it exits.
		goInstance.run(module.instance);

		// The Go program registers its JS functions just after startup. Wait for
		// them so pages do not call undefined globals and fall back to stale UI.
		await waitForGoFunctions();
		wasmReady = true;
	})();

	try {
		await readyPromise;
	} catch (error) {
		readyPromise = null;
		throw error;
	}
}

export function unloadWasm() {
	if (goInstance && goInstance.exit) {
		goInstance.exit();
	}
	goInstance = null;
	wasmReady = false;
	readyPromise = null;
}
