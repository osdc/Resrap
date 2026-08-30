package resrap

import (
	"os"
	"path/filepath"
	"testing"
)

const testGrammar = `program:function^;
function:'f' '(' ')';
`

func TestNewResrap(t *testing.T) {
	r := NewResrap()
	if r == nil {
		t.Fatal("NewResrap returned nil")
	}
	if len(r.languageGraph) != 0 {
		t.Fatalf("languageGraph size = %d, want 0", len(r.languageGraph))
	}
}

func TestResrapParseGrammarAndGenerate(t *testing.T) {
	r := NewResrap()
	if err := r.ParseGrammar("Test", testGrammar); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := r.GenerateRandom("Test", "program", 10)
	if out == "" {
		t.Fatal("GenerateRandom returned empty output")
	}
}

func TestResrapParseGrammarInvalid(t *testing.T) {
	r := NewResrap()
	err := r.ParseGrammar("Test", "program:missing;")
	if err == nil {
		t.Fatal("expected error for invalid grammar")
	}
}

func TestResrapParseGrammarFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "grammar.txt")
	if err := os.WriteFile(path, []byte(testGrammar), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewResrap()
	if err := r.ParseGrammarFile("Test", path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := r.GenerateRandom("Test", "program", 10)
	if out == "" {
		t.Fatal("GenerateRandom returned empty output")
	}
}

func TestResrapParseGrammarFileMissing(t *testing.T) {
	r := NewResrap()
	if err := r.ParseGrammarFile("Test", "/nonexistent/file.txt"); err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestResrapGenerateWithSeededDeterministic(t *testing.T) {
	r := NewResrap()
	if err := r.ParseGrammar("Test", testGrammar); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a := r.GenerateWithSeeded("Test", "program", 1234, 50)
	b := r.GenerateWithSeeded("Test", "program", 1234, 50)
	if a != b {
		t.Fatalf("seeded generation not deterministic:\n%q\nvs\n%q", a, b)
	}
}

func TestResrapGenerateWithSeededDifferentSeeds(t *testing.T) {
	r := NewResrap()
	if err := r.ParseGrammar("Test", "program:function^;\nfunction:('a'<0.5>) | ('b'<0.5>);"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a := r.GenerateWithSeeded("Test", "program", 1, 50)
	b := r.GenerateWithSeeded("Test", "program", 2, 50)
	if a == b {
		t.Fatal("different seeds should produce different output")
	}
}

func TestResrapGenerateUnknownLanguage(t *testing.T) {
	r := NewResrap()
	// Should not panic even for unknown language.
	out := r.GenerateRandom("Unknown", "program", 10)
	if out != "" {
		t.Fatalf("GenerateRandom for unknown language = %q, want empty", out)
	}
}

func TestResrapGenerateCodebase(t *testing.T) {
	r := NewResrap()
	if err := r.ParseGrammar("c", testGrammar); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	dir := t.TempDir()
	dslPath := filepath.Join(dir, "format.dsl")
	dsl := "out/\n  c[2x10 code_*.c]\n"
	if err := os.WriteFile(dslPath, []byte(dsl), 0644); err != nil {
		t.Fatal(err)
	}

	target := filepath.Join(dir, "target")
	if err := r.GenerateCodebase(dslPath, target); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	f1 := filepath.Join(target, "out", "code_1.c")
	f2 := filepath.Join(target, "out", "code_2.c")
	for _, f := range []string{f1, f2} {
		if _, err := os.Stat(f); err != nil {
			t.Fatalf("expected file %s: %v", f, err)
		}
	}
}

func TestResrapGenerateCodebaseMissingConfig(t *testing.T) {
	r := NewResrap()
	if err := r.GenerateCodebase("/nonexistent/config.dsl", "/tmp/out"); err == nil {
		t.Fatal("expected error for missing config file")
	}
}

func TestResrapGenerateCodebaseBadDSL(t *testing.T) {
	dir := t.TempDir()
	dslPath := filepath.Join(dir, "bad.dsl")
	if err := os.WriteFile(dslPath, []byte("c[bad]"), 0644); err != nil {
		t.Fatal(err)
	}
	r := NewResrap()
	if err := r.GenerateCodebase(dslPath, filepath.Join(dir, "out")); err == nil {
		t.Fatal("expected error for bad DSL")
	}
}

func TestResrapNormalizeCalled(t *testing.T) {
	r := NewResrap()
	if err := r.ParseGrammar("Test", testGrammar); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.languageGraph["Test"].graph == nil {
		t.Fatal("graph is nil after parse")
	}
	// After Normalize, header nodes should have cumulative frequencies.
	for _, node := range r.languageGraph["Test"].graph.nodeRef {
		if node.typ == header && len(node.next) > 0 && len(node.cf) != len(node.next) {
			t.Fatalf("node %d cf length %d != next length %d", node.id, len(node.cf), len(node.next))
		}
	}
}
