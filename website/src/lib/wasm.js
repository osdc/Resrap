import init, {
	generate_template,
	generate_text,
	get_templates
} from './wasm/resrap_wasm.js';

let readyPromise = null;

/**
 * Treat a zero seed as "pick a fresh one", matching the playground hint.
 * The generated value is kept below Number.MAX_SAFE_INTEGER so it survives
 * the trip into Rust's u64 seed argument without precision loss.
 */
function resolveSeed(seed) {
	if (seed !== 0) return seed;

	const random = new Uint32Array(1);
	if (globalThis.crypto?.getRandomValues) {
		globalThis.crypto.getRandomValues(random);
	} else {
		random[0] = Math.floor(Math.random() * 0xffffffff);
	}
	return Date.now() * 1000 + (random[0] % 1000);
}

export function loadWasm() {
	if (!readyPromise) {
		readyPromise = init()
			.then(() => undefined)
			.catch((error) => {
				readyPromise = null;
				throw error;
			});
	}
	return readyPromise;
}

export function generateText(grammar, startpoint, seed, tokenlen) {
	return generate_text(grammar, startpoint, resolveSeed(seed), tokenlen);
}

export function generateTemplate(template, startpoint, seed, tokenlen) {
	return generate_template(template, startpoint, resolveSeed(seed), tokenlen);
}

export function getTemplates() {
	return get_templates();
}

export function unloadWasm() {
	// wasm-bindgen owns a single module instance for the page lifetime.
}
