<script lang="ts">
	import { onMount } from 'svelte';
	import { loadWasm } from '$lib/wasm.js';

	type TemplateInfo = {
		name: string;
		language: string;
		description: string;
		grammar: string;
	};

	let status: 'loading' | 'ready' | 'error' = 'loading';
	let errorMessage = '';
	let templates: TemplateInfo[] = [];
	let outputs: Record<string, string> = {};
	let seed = 7;
	let tokens = 500;
	let isGenerating = false;
	let totalChars = 0;

	onMount(async () => {
		try {
			await loadWasm();
			if (!globalThis.getTemplates || !globalThis.generateTemplate) {
				throw new Error('WASM template functions are unavailable');
			}
			templates = JSON.parse(globalThis.getTemplates());
			generateAll();
			status = 'ready';
		} catch (error) {
			errorMessage = (error as Error).message;
			status = 'error';
		}
	});

	function generateAll() {
		if (!globalThis.generateTemplate) return;
		isGenerating = true;
		const next: Record<string, string> = {};
		let chars = 0;
		for (const template of templates) {
			const output = globalThis.generateTemplate(template.name, 'program', seed, tokens);
			next[template.name] = output;
			chars += output.length;
		}
		outputs = next;
		totalChars = chars;
		isGenerating = false;
	}

	async function copyCode(name: string) {
		try {
			await navigator.clipboard.writeText(outputs[name] ?? '');
		} catch {
			// Clipboard access is optional.
		}
	}

	function lineCount(name: string) {
		const output = outputs[name] ?? '';
		return output ? output.split('\n').length - 1 : 0;
	}
</script>

<svelte:head>
	<title>Resrap templates</title>
</svelte:head>

<div class="min-h-screen bg-[#06081d] p-4 text-white lg:p-8">
	<main class="mx-auto max-w-6xl">
		<a href="/" class="text-sm text-gray-400 hover:text-[#00add8]">← Back to Resrap</a>

		<header class="mt-5 mb-6">
			<h1 class="text-2xl font-bold text-[#00add8] lg:text-4xl">Built-in language templates</h1>
			<p class="mt-2 max-w-3xl text-sm text-gray-400 lg:text-base">
				These grammars are embedded in the Go module and served through WASM. Each seed builds a
				different, complete program with the language's core constructs.
			</p>
		</header>

		{#if status === 'loading'}
			<div class="flex h-64 items-center justify-center gap-3 text-gray-300">
				<div class="h-7 w-7 animate-spin rounded-full border-2 border-[#00add8] border-t-transparent"></div>
				Loading templates and generating demos...
			</div>
		{:else if status === 'error'}
			<div class="rounded-lg border border-red-500/40 bg-red-500/10 p-4 text-red-300">
				Could not load the template demos: {errorMessage}
			</div>
		{:else}
			<div class="mb-6 flex flex-wrap items-end gap-3 rounded-lg border border-gray-700 bg-gray-900 p-4">
				<label class="flex flex-col gap-1 text-xs text-gray-400">
					Seed
					<input
						bind:value={seed}
						type="number"
						class="w-28 rounded-md border border-gray-600 bg-gray-950 px-3 py-2 font-mono text-sm text-gray-100 outline-none focus:border-[#00add8]"
					/>
				</label>
				<label class="flex flex-col gap-1 text-xs text-gray-400">
					Token budget per language
					<input
						bind:value={tokens}
						type="number"
						min="50"
						class="w-36 rounded-md border border-gray-600 bg-gray-950 px-3 py-2 font-mono text-sm text-gray-100 outline-none focus:border-[#00add8]"
					/>
				</label>
				<button
					type="button"
					on:click={generateAll}
					disabled={isGenerating}
					class="rounded-md bg-[#00add8] px-4 py-2 text-sm font-semibold text-white hover:bg-[#0099c7] disabled:bg-gray-600"
				>
					{isGenerating ? 'Generating...' : 'Regenerate all'}
				</button>
				<span class="ml-auto text-xs text-gray-500">
					{templates.length} languages · {totalChars.toLocaleString()} characters generated
				</span>
			</div>

			<div class="space-y-6">
				{#each templates as template (template.name)}
					<section class="overflow-hidden rounded-lg border border-gray-700 bg-gray-900">
						<div class="flex flex-wrap items-center justify-between gap-2 border-b border-gray-700 px-4 py-3">
							<div>
								<div class="flex items-center gap-2">
									<h2 class="text-lg font-semibold text-gray-100">{template.language}</h2>
									<span class="rounded bg-[#00add8]/10 px-2 py-0.5 font-mono text-[10px] uppercase text-[#00add8]">
										{template.name}
									</span>
								</div>
								<p class="text-xs text-gray-500">{template.description}</p>
							</div>
							<div class="flex items-center gap-3 text-xs text-gray-500">
								<span>{lineCount(template.name)} lines</span>
								<span>{outputs[template.name]?.length.toLocaleString() ?? 0} chars</span>
								<button
									type="button"
									on:click={() => copyCode(template.name)}
									class="rounded border border-gray-600 px-2 py-1 text-gray-300 hover:border-[#00add8] hover:text-[#00add8]"
								>
									Copy
								</button>
							</div>
						</div>
						<pre
							class="max-h-[32rem] overflow-auto bg-gray-950 p-4 font-mono text-xs leading-relaxed text-gray-200"><code>{outputs[template.name] ?? ''}</code></pre>
					</section>
				{/each}
			</div>
		{/if}
	</main>
</div>
