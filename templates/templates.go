// Package templates ships ready-to-use Resrap grammars for common languages.
//
// The grammars are embedded in the binary, so they work in normal Go programs
// and in browser/WASM builds without touching the filesystem. Languages are
// available as typed values (Cpp, Rust, Python, ...) or short string names.
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

// Lang identifies one of the embedded language templates.
type Lang string

const (
	C          Lang = "c"
	Cpp        Lang = "cpp"
	Go         Lang = "go"
	Haskell    Lang = "haskell"
	HTML       Lang = "html"
	Java       Lang = "java"
	JavaScript Lang = "js"
	Python     Lang = "python"
	Rust       Lang = "rust"
	SQL        Lang = "sql"
)

// String returns the short template name, for example "cpp".
func (l Lang) String() string { return string(l) }

// Valid reports whether l is an embedded template.
func (l Lang) Valid() bool {
	_, err := Get(l)
	return err == nil
}

// Template returns the grammar and metadata for l.
func (l Lang) Template() (Template, error) { return Get(l) }

// Load parses l into a Resrap instance under its template name.
func (l Lang) Load(loader GrammarLoader) error { return Load(loader, l) }

// Langs returns every embedded language, sorted by short name.
func Langs() []Lang {
	names := Names()
	out := make([]Lang, len(names))
	for i, name := range names {
		out[i] = Lang(name)
	}
	return out
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
// ".g4" extension. It accepts both strings and Lang values.
func Get[T ~string](name T) (Template, error) {
	key := normalize(string(name))
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
func MustGet[T ~string](name T) Template {
	t, err := Get(name)
	if err != nil {
		panic(err)
	}
	return t
}

// Load parses a template into a Resrap instance under the template name.
// The name can be a string or a Lang value.
func Load[T ~string](loader GrammarLoader, name T) error {
	return LoadAs(loader, name, normalize(string(name)))
}

// Load parses this template into a Resrap instance under the template name.
func (t Template) Load(loader GrammarLoader) error {
	return loader.ParseGrammar(t.Name, t.Grammar)
}

// LoadAs parses a template into a Resrap instance under a custom grammar name.
// The template name can be a string or a Lang value.
func LoadAs[T ~string](loader GrammarLoader, name T, as string) error {
	t, err := Get(name)
	if err != nil {
		return err
	}
	return loader.ParseGrammar(as, t.Grammar)
}

func normalize(name string) string {
	return strings.ToLower(strings.TrimSuffix(strings.TrimSpace(name), ".g4"))
}
