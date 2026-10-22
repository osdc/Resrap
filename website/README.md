# Resrap website

The interactive Resrap site and playground. The browser generator is the Rust
[`resrap-rs`](https://github.com/osdc/resrap-rs) crate compiled to WebAssembly
with [`wasm-pack`](https://rustwasm.github.io/wasm-pack/).

## Prerequisites

- Node.js 20 or newer
- A stable Rust toolchain
- The WebAssembly target:

```sh
rustup target add wasm32-unknown-unknown
```

- `wasm-pack`:

```sh
cargo install wasm-pack
```

## Developing

Install the JavaScript dependencies, then start the dev server:

```sh
npm install
npm run dev
```

`npm run dev` builds the Rust WASM package first and then starts Vite. If you
only need to rebuild the WASM package:

```sh
npm run build-wasm
```

The generated package is written to `src/lib/wasm/` and is ignored by git. The
Rust bindings live in `src/srcrust/`; they depend on the tagged `resrap-rs`
release so the website can be built without a sibling checkout.

## Building

```sh
npm run build
```

This regenerates the WASM package and creates a production build. Preview it
with:

```sh
npm run preview
```

> Deployment targets may need a SvelteKit adapter other than `adapter-auto`.
> See the [SvelteKit adapter docs](https://svelte.dev/docs/kit/adapters).
