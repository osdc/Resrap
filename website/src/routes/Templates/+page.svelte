<script lang="ts">
	import { onMount } from 'svelte';
	import { loadWasm } from '$lib/wasm.js';

	import Highlight from 'svelte-highlight';
	import tokyoNight from 'svelte-highlight/styles/tokyo-night-dark';
	import abnf from 'svelte-highlight/languages/abnf';
	import c from 'svelte-highlight/languages/c';
	import cpp from 'svelte-highlight/languages/cpp';
	import go from 'svelte-highlight/languages/go';
	import haskell from 'svelte-highlight/languages/haskell';
	import xml from 'svelte-highlight/languages/xml';
	import java from 'svelte-highlight/languages/java';
	import javascript from 'svelte-highlight/languages/javascript';
	import python from 'svelte-highlight/languages/python';
	import rust from 'svelte-highlight/languages/rust';
	import sql from 'svelte-highlight/languages/sql';

	type TemplateInfo = {
		name: string;
		language: string;
		description: string;
		grammar: string;
	};

	const languages: Record<string, any> = {
		c,
		cpp,
		go,
		haskell,
		html: xml,
		java,
		js: javascript,
		python,
		rust,
		sql
	};

	let status: 'loading' | 'ready' | 'error' = 'loading';
	let errorMessage = '';
	let templates: TemplateInfo[] = [];
	let selectedName = 'cpp';
	let selectedTemplate: TemplateInfo | undefined;
	let output = '';
	let seed = 7;
	let tokens = 500;
	let isGenerating = false;
	let grammarOpen = false;

	$: selectedLanguage = languages[selectedName] ?? c;
	$: lineCount = output ? output.split('\n').length - 1 : 0;

	onMount(async () => {
		try {
			await loadWasm();
			if (!globalThis.getTemplates || !globalThis.generateTemplate) {
				throw new Error('WASM template functions are unavailable');
			}

			templates = JSON.parse(globalThis.getTemplates());
			selectedName = templates.find((template) => template.name === 'cpp')?.name ?? templates[0]?.name ?? '';
			selectTemplate();
			status = 'ready';
		} catch (error) {
			errorMessage = (error as Error).message;
			status = 'error';
		}
	});

	function selectTemplate() {
		selectedTemplate = templates.find((template) => template.name === selectedName);
		generateSelected();
	}

	function generateSelected() {
		if (!globalThis.generateTemplate || !selectedTemplate) return;

		isGenerating = true;
		output = globalThis.generateTemplate(selectedTemplate.name, 'program', seed, tokens);
		isGenerating = false;
	}

</script>

<svelte:head>
	{@html tokyoNight}
	<title>Resrap templates</title>
</svelte:head>

<svelte:window
	onkeydown={(event) => {
		if (event.key === 'Escape') grammarOpen = false;
	}}
/>

<div class="min-h-screen bg-[#06081d] p-4 text-white lg:p-8">
	<main class="mx-auto max-w-6xl">
		<a href="/" class="text-sm text-gray-400 hover:text-[#00add8]">← Back to Resrap</a>

		<header class="mt-5 mb-6">
			<h1 class="text-2xl font-bold text-[#00add8] lg:text-4xl">Built-in language templates</h1>
			<p class="mt-2 max-w-3xl text-sm text-gray-400 lg:text-base">
				These grammars are embedded in the Go module and served through WASM. Choose a language
				and each seed builds a different, complete program with the language's core constructs.
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
					Language
					<select
						bind:value={selectedName}
						onchange={selectTemplate}
						class="min-w-44 rounded-md border border-gray-600 bg-gray-950 px-3 py-2 font-mono text-sm text-gray-100 outline-none focus:border-[#00add8]"
					>
						{#each templates as template (template.name)}
							<option value={template.name}>{template.language}</option>
						{/each}
					</select>
				</label>

				<label class="flex flex-col gap-1 text-xs text-gray-400">
					Seed
					<input
						bind:value={seed}
						type="number"
						class="w-28 rounded-md border border-gray-600 bg-gray-950 px-3 py-2 font-mono text-sm text-gray-100 outline-none focus:border-[#00add8]"
					/>
				</label>

				<label class="flex flex-col gap-1 text-xs text-gray-400">
					Token budget
					<input
						bind:value={tokens}
						type="number"
						min="50"
						class="w-36 rounded-md border border-gray-600 bg-gray-950 px-3 py-2 font-mono text-sm text-gray-100 outline-none focus:border-[#00add8]"
					/>
				</label>

				<button
					type="button"
					onclick={generateSelected}
					disabled={isGenerating}
					class="rounded-md bg-[#00add8] px-4 py-2 text-sm font-semibold text-white hover:bg-[#0099c7] disabled:bg-gray-600"
				>
					{isGenerating ? 'Generating...' : 'Regenerate'}
				</button>

				<span class="ml-auto text-xs text-gray-500">
					{lineCount} lines · {output.length.toLocaleString()} characters
				</span>
			</div>

			{#if selectedTemplate}
				<section class="overflow-hidden rounded-lg border border-gray-700 bg-gray-900">
					<div class="flex flex-wrap items-center justify-between gap-2 border-b border-gray-700 px-4 py-3">
						<div>
							<h2 class="text-lg font-semibold text-gray-100">{selectedTemplate.language}</h2>
							<p class="text-xs text-gray-500">{selectedTemplate.description}</p>
						</div>
						<button
							type="button"
							onclick={() => (grammarOpen = true)}
							class="rounded border border-gray-600 px-3 py-1.5 text-xs text-gray-300 hover:border-[#00add8] hover:text-[#00add8]"
						>
							See ABNF
						</button>
					</div>

					<div class="code-theme max-h-[36rem] overflow-auto bg-[#0b1020] p-4 text-sm">
						<Highlight language={selectedLanguage} code={output} />
					</div>
				</section>
			{/if}
		{/if}
	</main>

	{#if grammarOpen && selectedTemplate}
		<div
			class="fixed inset-0 z-50 flex items-center justify-center bg-black/70 p-4 backdrop-blur-sm"
			role="presentation"
			onclick={(event) => {
				if (event.target === event.currentTarget) grammarOpen = false;
			}}
		>
			<div
				class="flex max-h-[85vh] w-full max-w-4xl flex-col overflow-hidden rounded-lg border border-gray-700 bg-gray-900 shadow-2xl"
				role="dialog"
				aria-modal="true"
				aria-labelledby="grammar-title"
			>
				<div class="flex items-center justify-between border-b border-gray-700 px-4 py-3">
					<div>
						<h2 id="grammar-title" class="text-lg font-semibold text-gray-100">
							{selectedTemplate.language} ABNF
						</h2>
						<p class="text-xs text-gray-500">The grammar behind the generated snippet</p>
					</div>
					<button
						type="button"
						onclick={() => (grammarOpen = false)}
						class="rounded border border-gray-600 px-3 py-1.5 text-xs text-gray-300 hover:border-[#00add8] hover:text-[#00add8]"
					>
						Close
					</button>
				</div>
				<div class="code-theme overflow-auto bg-[#0b1020] p-4 text-sm">
					<Highlight language={abnf} code={selectedTemplate.grammar} />
				</div>
			</div>
		</div>
	{/if}
</div>
