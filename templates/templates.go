// Package templates ships ready-to-use Resrap grammars for common languages.
//
// The grammars are embedded in the binary, so they work in normal Go programs
// and in browser/WASM builds without touching the filesystem.
package templates

import (
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strings"
)

//go:embed *.g4
var grammars embed.FS

// Template describes one embedded grammar.
type Template struct {
	// Name is the short identifier used by the API, for example "cpp".
	Name string
	// Language is a human-readable language name, for example "C++".
	Language string
	// Description is a one-line summary suitable for a picker.
	Description string
	// Grammar is the complete ABNF grammar text.
	Grammar string
}

// GrammarLoader is implemented by *resrap.Resrap and *resrap.ResrapMT.
type GrammarLoader interface {
	ParseGrammar(name, grammar string) error
}

type info struct {
	language    string
	description string
}

var infoByName = map[string]info{
	"c":       {"C", "Headers, loops, conditionals and printf output"},
	"cpp":     {"C++", "Classes, helpers, lambdas, exceptions and complete functions"},
	"go":      {"Go", "Package main with imports, loops and fmt output"},
	"haskell": {"Haskell", "Module Main with do blocks and list operations"},
	"html":    {"HTML", "Document structure, headings, paragraphs and lists"},
	"java":    {"Java", "A Main class with a static main method"},
	"js":      {"JavaScript", "Browser-ready functions, loops and console output"},
	"python":  {"Python", "A main function with loops and print statements"},
	"rust":    {"Rust", "A main function with loops, ranges and println!"},
	"sql":     {"SQL", "SELECT, INSERT and CREATE TABLE statements"},
}

// Names returns the short names of every embedded template, sorted.
func Names() []string {
	entries, err := fs.ReadDir(grammars, ".")
	if err != nil {
		return nil
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".g4") {
			continue
		}
		names = append(names, strings.TrimSuffix(entry.Name(), ".g4"))
	}
	sort.Strings(names)
	return names
}

// List returns every embedded template with its grammar and metadata.
func List() ([]Template, error) {
	names := Names()
	out := make([]Template, 0, len(names))
	for _, name := range names {
		t, err := Get(name)
		if err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}

// Get returns one template. The name is case-insensitive and may include the
// ".g4" extension.
func Get(name string) (Template, error) {
	key := normalize(name)
	body, err := grammars.ReadFile(key + ".g4")
	if err != nil {
		return Template{}, fmt.Errorf("templates: unknown template %q (available: %s)", name, strings.Join(Names(), ", "))
	}
	meta, ok := infoByName[key]
	if !ok {
		meta = info{strings.ToUpper(key[:1]) + key[1:], "Embedded Resrap grammar"}
	}
	return Template{
		Name:        key,
		Language:    meta.language,
		Description: meta.description,
		Grammar:     string(body),
	}, nil
}

// MustGet is Get for callers that treat a missing built-in template as a
// programming error.
func MustGet(name string) Template {
	t, err := Get(name)
	if err != nil {
		panic(err)
	}
	return t
}

// Load parses a template into a Resrap instance under the template name.
func Load(loader GrammarLoader, name string) error {
	return LoadAs(loader, name, normalize(name))
}

// Load parses this template into a Resrap instance under the template name.
func (t Template) Load(loader GrammarLoader) error {
	return loader.ParseGrammar(t.Name, t.Grammar)
}

// LoadAs parses a template into a Resrap instance under a custom grammar name.
func LoadAs(loader GrammarLoader, name, as string) error {
	t, err := Get(name)
	if err != nil {
		return err
	}
	return loader.ParseGrammar(as, t.Grammar)
}

func normalize(name string) string {
	return strings.ToLower(strings.TrimSuffix(strings.TrimSpace(name), ".g4"))
}
