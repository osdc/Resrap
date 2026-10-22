<section id="abnf" class="px-6 py-16 text-white lg:px-12 lg:py-24">
	<div class="mx-auto max-w-5xl">
		<h1 class="mb-4 text-2xl lg:text-4xl font-bold text-[#00add8]">ABNF – Awesome BNF</h1>
		<p class="mb-8 lg:mb-12 text-base lg:text-lg leading-relaxed text-gray-300">
			ABNF (Awesome BNF) is a lightweight, custom grammar format designed for the Resrap code
			generation tool. It extends standard EBNF/BNF with probabilities, infinite generation, and
			readable word lists.
		</p>

		<h2 class="mb-4 lg:mb-6 text-xl lg:text-2xl font-bold text-[#00add8]">Quick Revision: EBNF and Grammar</h2>
		<p class="mb-4 text-gray-300 text-sm lg:text-base">ABNF inherits the basics of EBNF:</p>
		<pre class="mb-6 lg:mb-8 rounded bg-gray-900 p-3 font-mono text-xs lg:text-sm text-gray-100 overflow-x-auto">function : rules ;</pre>

		<h3 class="mb-4 text-lg lg:text-xl font-semibold text-gray-200">ABNF Operators</h3>
		<div class="mb-8 lg:mb-12 space-y-2 lg:space-y-3">
			<p class="font-mono text-sm lg:text-base">
				<span class="font-bold text-[#00add8]">+</span>
				<span class="text-gray-300">One or more repetitions</span>
			</p>
			<p class="font-mono text-sm lg:text-base">
				<span class="font-bold text-[#00add8]">*</span>
				<span class="text-gray-300">Zero or more repetitions</span>
			</p>
			<p class="font-mono text-sm lg:text-base">
				<span class="font-bold text-[#00add8]">?</span>
				<span class="text-gray-300">Optional element</span>
			</p>
			<p class="font-mono text-sm lg:text-base">
				<span class="font-bold text-[#00add8]">()</span> <span class="text-gray-300">Grouping</span>
			</p>
			<p class="font-mono text-sm lg:text-base">
				<span class="font-bold text-[#00add8]">^</span>
				<span class="text-gray-300">Infinite generation - loops back infinitely</span>
			</p>
		</div>

		<div class="mb-6 lg:mb-8">
			<p class="mb-3 text-gray-300 text-sm lg:text-base">Example with infinite generation:</p>
			<pre class="rounded bg-gray-900 p-3 font-mono text-xs lg:text-sm text-gray-100 overflow-x-auto">program : function^;</pre>
			<p class="mt-2 text-xs lg:text-sm text-gray-400">
				Generates functions endlessly until token limit is reached
			</p>
		</div>

		<h2 class="mb-4 lg:mb-6 text-xl lg:text-2xl font-bold text-[#00add8]">Probabilities (&lt;...&gt;)</h2>
		<p class="mb-4 text-gray-300 text-sm lg:text-base">
			ABNF allows specifying weighted probabilities for finer control:
		</p>
		<pre
			class="mb-4 rounded bg-gray-900 p-3 font-mono text-xs lg:text-sm text-gray-100 overflow-x-auto">char : a&lt;0.2&gt; | b&lt;0.8&gt;;
a    : 'A';
b    : 'B';</pre>
		<p class="mb-4 text-gray-300 text-sm lg:text-base">
			<span class="text-[#00add8]">a</span> is chosen 20% of the time,
			<span class="text-[#00add8]">b</span> 80%. Probabilities are normalized automatically.
		</p>

		<div class="mb-8 lg:mb-12">
			<p class="mb-3 text-gray-300 text-sm lg:text-base">Probabilities work with all operators:</p>
			<pre
				class="rounded bg-gray-900 p-3 font-mono text-xs lg:text-sm text-gray-100 overflow-x-auto">statement : loop+&lt;0.3&gt; | condition?&lt;0.7&gt;;</pre>
		</div>

		<h2 class="mb-4 lg:mb-6 text-xl lg:text-2xl font-bold text-[#00add8]">Word Lists (<code>@words</code>)</h2>
		<p class="mb-4 text-gray-300 text-sm lg:text-base">
			Random character classes produce random-looking identifiers. A word list lets a rule emit
			readable words instead of expanding its grammar:
		</p>
		<pre
			class="mb-4 rounded bg-gray-900 p-3 font-mono text-xs lg:text-sm text-gray-100 overflow-x-auto">@words identifier = count, total, index, value, result
@words type       = int, long, double

program: type ' ' identifier ';\n';</pre>
		<p class="mb-4 text-gray-300 text-sm lg:text-base">
			Every reference to <code class="text-[#00add8]">identifier</code> now emits one of those words.
			Word-only rules do not need a separate grammar definition, and
			<code class="text-[#00add8]">@keywords</code> is accepted as an alias for
			<code class="text-[#00add8]">@words</code>.
		</p>
		<p class="mb-8 lg:mb-12 text-gray-300 text-sm lg:text-base">
			Words can also be attached at runtime with <code class="text-[#00add8]">Resrap.AddWords</code>.
			This is useful when keywords come from another source or need to change per job.
		</p>

		<h2 class="mb-4 lg:mb-6 text-xl lg:text-2xl font-bold text-[#00add8]">Comments and escapes</h2>
		<div class="mb-6 grid gap-4 sm:grid-cols-2">
			<div class="rounded bg-gray-900 p-4">
				<h3 class="mb-2 text-sm font-semibold text-gray-200">Comments</h3>
				<p class="text-sm leading-6 text-gray-400">
					Start a line comment with <code class="text-[#00add8]">#</code> or
					<code class="text-[#00add8]">//</code>. Both run to the end of the line.
				</p>
			</div>
			<div class="rounded bg-gray-900 p-4">
				<h3 class="mb-2 text-sm font-semibold text-gray-200">Escapes</h3>
				<p class="text-sm leading-6 text-gray-400">
					Single-quoted terminals support <code class="text-[#00add8]">\n</code>,
					<code class="text-[#00add8]">\t</code>, <code class="text-[#00add8]">\r</code>,
					<code class="text-[#00add8]">\\</code>, <code class="text-[#00add8]">\'</code>, and
					<code class="text-[#00add8]">\"</code>.
				</p>
			</div>
		</div>
		<p class="mb-4 text-gray-300 text-sm lg:text-base">
			Character classes support ranges such as <code class="text-[#00add8]">[a-z0-9_]</code> plus
			the common escapes.
		</p>
	</div>
</section>
