package resrap

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "grammar.txt")
	content := "program : 'a';\n// this is a comment\nfunction: 'b';\n\nmore:'c';\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	statements, err := parseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(statements) != 3 {
		t.Fatalf("got %d statements, want 3", len(statements))
	}
	if statements[0] != "program : 'a';" {
		t.Fatalf("statement[0] = %q", statements[0])
	}
	if statements[1] != "function: 'b';" {
		t.Fatalf("statement[1] = %q", statements[1])
	}
	if statements[2] != "more:'c';" {
		t.Fatalf("statement[2] = %q", statements[2])
	}
}

func TestParseFileMultiLineStatement(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "grammar.txt")
	content := "program :\n  'a'\n  'b';\n"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	statements, err := parseFile(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(statements) != 1 {
		t.Fatalf("got %d statements, want 1", len(statements))
	}
	want := "program : 'a' 'b';"
	if statements[0] != want {
		t.Fatalf("statement = %q, want %q", statements[0], want)
	}
}

func TestParseFileNonexistent(t *testing.T) {
	if _, err := parseFile("/nonexistent/path/grammar.txt"); err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestLangParserFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "grammar.txt")
	if err := os.WriteFile(path, []byte("program:'a';\n"), 0644); err != nil {
		t.Fatal(err)
	}

	l := newLang()
	if err := l.ParserFile(path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.GetGraph() == nil {
		t.Fatal("graph should not be nil after parsing")
	}
}

func TestLangParserFileMissingFile(t *testing.T) {
	l := newLang()
	if err := l.ParserFile("/nonexistent/path/grammar.txt"); err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestLangParserString(t *testing.T) {
	l := newLang()
	if err := l.ParserString("program:'a';\n"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if l.GetGraph() == nil {
		t.Fatal("graph should not be nil after parsing")
	}
}

func TestLangParserStringInvalid(t *testing.T) {
	l := newLang()
	if err := l.ParserString("program:undefined;"); err == nil {
		t.Fatal("expected error for undefined symbol")
	}
}
