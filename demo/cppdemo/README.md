# Template demo

Generates source code from one of Resrap's embedded language templates.

```bash
cd demo
go run ./cppdemo -seed 10 -tokens 120
```

Flags:

* `-seed` — deterministic PRNG seed
* `-tokens` — soft token budget; the current top-level rule is always completed
* `-template` — built-in template name (`c`, `cpp`, `go`, `haskell`, `html`,
  `java`, `js`, `python`, `rust`, `sql`)
* `-out` — output directory (default `cppdemo/out`)
* `-print` — print the complete generated file instead of a preview

The demo only generates and writes source. It does **not** compile or run the
result; use your language toolchain separately if you want to validate it.

The C++ template demonstrates:

* `GenerateCompleteWithSeeded` avoiding truncation mid-function
* `@words` supplying readable strings and short identifiers
* `^` at the end of `program` repeating complete functions until the budget
