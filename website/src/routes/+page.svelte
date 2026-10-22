<script lang="ts">
	import ABNF from './ABNF.svelte';
	import About from './About.svelte';
	import Footer from './footer.svelte';
	import GenDel from './GenDel.svelte';
	import Header from './header.svelte';
	import Sample from './Sample.svelte';
</script>

<div class="min-h-screen w-full bg-[#06081d]">
	<section class="mx-auto flex min-h-screen w-full max-w-[1440px] flex-col px-6 py-8 lg:px-12 lg:py-10">
		<nav
			class="mb-8 flex flex-wrap items-center gap-x-6 gap-y-3 text-sm text-gray-400 lg:mb-10"
			aria-label="Page sections"
		>
			<a href="#about" class="transition-colors hover:text-[#00add8]">About</a>
			<a href="#abnf" class="transition-colors hover:text-[#00add8]">Docs</a>
			<a href="#try" class="transition-colors hover:text-[#00add8]">Playground</a>
			<a href="/Templates" class="transition-colors hover:text-[#00add8]">Templates</a>
			<a
				href="https://github.com/osdc/Resrap"
				target="_blank"
				rel="noopener noreferrer"
				class="transition-colors hover:text-[#00add8]">GitHub</a
			>
		</nav>

		<div
			class="grid w-full flex-1 items-center gap-12 lg:grid-cols-[minmax(0,0.9fr)_minmax(0,1.1fr)] lg:gap-16 xl:gap-20"
		>
			<Header />
			<GenDel />
		</div>
	</section>

	<div class="border-t border-white/[0.06]">
		<About />
	</div>
	<div class="border-t border-white/[0.06]">
		<ABNF />
	</div>
	<div class="border-t border-white/[0.06]">
		<Sample />
	</div>
	<Footer />
</div>
