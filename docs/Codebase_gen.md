# Codebase generation

V1.1 introduces a format to define a directory structure, complete with no of files, token length etc for each file and generate entire codebases, not just single files


An example structure format can be seen below:
```dsl
src/
  code/
    c[10x20 code_*.c]
    core/
      c[5x1000 *.c]
  more/
    sql[10x100 *.sql]
```
`sql[10x100 *.sql]`: access the sql grammar, generate 10 files, 100 tokens each with name [unique_identifier].sql

Then simply call GenerateCodebase
```go
	import "github.com/osdc/resrap/templates"

	r := resrap.NewResrap()

	// Built-in templates are embedded in the module.
	if err := templates.Load(r, templates.SQL); err != nil {
		panic(err)
	}
	if err := templates.Load(r, templates.C); err != nil {
		panic(err)
	}

	if err := r.GenerateCodebase("../example/format.dsl", "."); err != nil {
		panic(err)
	}
```
