// Command cppdemo generates source code from one of the built-in templates.
//
// It is the reference demo for complete generation and @words/keyword support.
// Generated output is written to disk but is deliberately not compiled or run;
// Resrap is a generator, not a build system.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/osdc/resrap"
	"github.com/osdc/resrap/templates"
)

var extensions = map[string]string{
	"c":       ".c",
	"cpp":     ".cpp",
	"go":      ".go",
	"haskell": ".hs",
	"html":    ".html",
	"java":    ".java",
	"js":      ".js",
	"python":  ".py",
	"rust":    ".rs",
	"sql":     ".sql",
}

func main() {
	seed := flag.Uint64("seed", 2, "PRNG seed")
	tokens := flag.Int("tokens", 5000, "soft token budget; the current top-level construct is always finished")
	name := flag.String("template", "cpp", "built-in template: "+strings.Join(templates.Names(), ", "))
	outDir := flag.String("out", "cppdemo/out", "directory for generated source")
	printAll := flag.Bool("print", false, "print the complete generated file to stdout")
	flag.Parse()

	if flag.NArg() > 0 {
		fmt.Fprintf(os.Stderr, "unexpected arguments: %v\n", flag.Args())
		fmt.Fprintln(os.Stderr, "flag syntax is -seed 10 -tokens 50 (no space after the dash)")
		os.Exit(2)
	}
	if *tokens <= 0 {
		fmt.Fprintln(os.Stderr, "-tokens must be greater than zero")
		os.Exit(2)
	}

	tmpl, err := templates.Get(*name)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	r := resrap.NewResrap()
	if err := tmpl.Load(r); err != nil {
		fmt.Fprintf(os.Stderr, "parse template %q: %v\n", tmpl.Name, err)
		os.Exit(1)
	}

	code := r.GenerateCompleteWithSeeded(tmpl.Name, "program", *seed, *tokens)
	if strings.TrimSpace(code) == "" {
		fmt.Fprintln(os.Stderr, "generation produced no output")
		os.Exit(1)
	}

	ext, ok := extensions[tmpl.Name]
	if !ok {
		ext = ".txt"
	}
	if err := os.MkdirAll(*outDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "create output dir: %v\n", err)
		os.Exit(1)
	}
	src := filepath.Join(*outDir, tmpl.Name+ext)
	if err := os.WriteFile(src, []byte(code), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write source: %v\n", err)
		os.Exit(1)
	}

	lines := strings.Count(code, "\n")
	if !strings.HasSuffix(code, "\n") {
		lines++
	}
	fmt.Printf("template: %s (%s)\n", tmpl.Name, tmpl.Language)
	fmt.Printf("seed:     %d\n", *seed)
	fmt.Printf("tokens:   %d (soft budget, top-level rule finished)\n", *tokens)
	fmt.Printf("output:   %d bytes, %d lines\n", len(code), lines)
	fmt.Printf("wrote:    %s\n", src)

	if *printAll {
		fmt.Println("\n--- generated source ---")
		fmt.Print(code)
		return
	}
	fmt.Println("\n--- preview (first 20 lines; use -print for the full file) ---")
	preview := strings.Split(code, "\n")
	if len(preview) > 20 {
		preview = preview[:20]
	}
	fmt.Println(strings.Join(preview, "\n"))
}
