package resrap

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestMalformedGrammarsReturnErrors(t *testing.T) {
	tests := []struct {
		name    string
		grammar string
	}{
		{"missing semicolon", `program:'a'`},
		{"leading plus", `program:+'a';`},
		{"leading star", `program:*'a';`},
		{"leading question", `program:?'a';`},
		{"missing colon", `program 'a';`},
		{"missing subject", `'a';`},
		{"stray close", `program:'a');`},
		{"stray open", `program:('a';`},
		{"invalid probability", `program:'a'<nope>;`},
		{"negative probability", `program:'a'<-1>;`},
		{"nan probability", `program:'a'<NaN>;`},
		{"infinite probability", `program:'a'<Inf>;`},
		{"unknown character", `program:'a' $ 'b';`},
		{"lone probability", `program:<0.5>;`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Fatalf("ParseGrammar panicked: %v", r)
				}
			}()
			r := NewResrap()
			if err := r.ParseGrammar("x", tt.grammar); err == nil {
				t.Fatalf("expected an error for %q", tt.grammar)
			}
		})
	}
}

func TestParseErrorHasPosition(t *testing.T) {
	r := NewResrap()
	err := r.ParseGrammar("x", "program:+'a';")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "line 1:") {
		t.Fatalf("error does not contain a position: %v", err)
	}
}

func TestRegexInfiniteTerminates(t *testing.T) {
	r := NewResrap()
	if err := r.ParseGrammar("x", `program:[a-z]^;`); err != nil {
		t.Fatal(err)
	}
	done := make(chan string, 1)
	go func() {
		done <- r.GenerateCompleteWithSeeded("x", "program", 7, 5)
	}()
	select {
	case out := <-done:
		if out == "" {
			t.Fatal("expected non-empty output")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("regex-only infinity grammar hung")
	}
}

func TestZeroProbabilityDoesNotPanic(t *testing.T) {
	r := NewResrap()
	if err := r.ParseGrammar("x", `program:'a'<0>|'b'<0>;`); err != nil {
		t.Fatal(err)
	}
	defer func() {
		if rec := recover(); rec != nil {
			t.Fatalf("generation panicked: %v", rec)
		}
	}()
	out := r.GenerateWithSeeded("x", "program", 3, 10)
	if out == "" {
		t.Fatal("expected fallback output")
	}
}

func TestInfiniteLoopsStayInsideTheirRule(t *testing.T) {
	r := NewResrap()
	grammar := "program:a b;\na:'A'^;\nb:'B'^;"
	if err := r.ParseGrammar("x", grammar); err != nil {
		t.Fatal(err)
	}
	out := r.GenerateWithSeeded("x", "a", 11, 40)
	if strings.Trim(out, "A") != "" {
		t.Fatalf("rule 'a' leaked into rule 'b': %q", out)
	}
}

func TestCompleteGenerationStopsAtRuleBoundary(t *testing.T) {
	r := NewResrap()
	grammar := "program:'begin\n' item* 'end\n';\nitem:'item\n';"
	if err := r.ParseGrammar("x", grammar); err != nil {
		t.Fatal(err)
	}
	for seed := uint64(1); seed <= 50; seed++ {
		out := r.GenerateCompleteWithSeeded("x", "program", seed, 3)
		if !strings.HasPrefix(out, "begin\n") || !strings.HasSuffix(out, "end\n") {
			t.Fatalf("seed %d: incomplete output %q", seed, out)
		}
	}
}

func TestScannerEscapesAndComments(t *testing.T) {
	tokens, errs := extracttokens("# a comment\nprogram:'it\\'s';\n// trailing\n")
	if len(errs) != 0 {
		t.Fatalf("unexpected scan errors: %v", errs)
	}
	if len(tokens) != 4 {
		t.Fatalf("got %d tokens, want 4: %+v", len(tokens), tokens)
	}
	if tokens[2].typ != character || tokens[2].text != `it\'s` {
		t.Fatalf("escaped quote token = %+v", tokens[2])
	}
	// The scanner keeps the escape; generation decodes it.
	if got := unescapeString(tokens[2].text); got != "it's" {
		t.Fatalf("unescaped = %q, want %q", got, "it's")
	}
}

func TestParseFileTrailingRuleIsNotSilentlyDropped(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/grammar.g4"
	if err := writeFile(path, "program:'a';\nextra:'b'"); err != nil {
		t.Fatal(err)
	}
	l := newLang()
	err := l.ParserFile(path)
	if err == nil {
		t.Fatal("expected an error for a trailing rule without ';'")
	}
	if !strings.Contains(err.Error(), "missing ';'") {
		t.Fatalf("error should mention the missing semicolon, got: %v", err)
	}
}

func TestWordsDirective(t *testing.T) {
	grammar := "@words identifier = count, total, index, value\nprogram: 'int ' identifier ';';"
	r := NewResrap()
	if err := r.ParseGrammar("x", grammar); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	words := map[string]bool{"count": true, "total": true, "index": true, "value": true}
	for seed := uint64(1); seed <= 20; seed++ {
		out := r.GenerateCompleteWithSeeded("x", "program", seed, 10)
		found := false
		for w := range words {
			if strings.Contains(out, "int "+w+";") {
				found = true
			}
		}
		if !found {
			t.Fatalf("seed %d: expected a word from the directive, got %q", seed, out)
		}
	}
}

func TestWordsOnlyRule(t *testing.T) {
	grammar := "@keywords name = alpha, beta, gamma\nprogram:'hi ' name '!';\n"
	r := NewResrap()
	if err := r.ParseGrammar("x", grammar); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := r.GenerateCompleteWithSeeded("x", "program", 1, 10)
	if !strings.HasPrefix(out, "hi ") || !strings.HasSuffix(out, "!") {
		t.Fatalf("unexpected output %q", out)
	}
}

func TestAddWordsAPI(t *testing.T) {
	r := NewResrap()
	if err := r.ParseGrammar("x", `program:'int ' identifier ';'; identifier:[a-z];`); err != nil {
		t.Fatal(err)
	}
	if err := r.AddWords("x", "identifier", []string{"count", "total"}); err != nil {
		t.Fatal(err)
	}
	for seed := uint64(1); seed <= 20; seed++ {
		out := r.GenerateCompleteWithSeeded("x", "program", seed, 10)
		if !strings.Contains(out, "int count;") && !strings.Contains(out, "int total;") {
			t.Fatalf("seed %d: words not applied: %q", seed, out)
		}
	}
	if err := r.AddWords("missing", "identifier", []string{"x"}); err == nil {
		t.Fatal("expected error for unknown grammar")
	}
	if err := r.AddWords("x", "missing", []string{"x"}); err == nil {
		t.Fatal("expected error for unknown rule")
	}
}

func TestRegexEscapedHyphenIsLiteral(t *testing.T) {
	r := newRegexer()
	got := r.ExpandClass(`a\-z`)
	want := []rune{'a', '-', 'z'}
	if len(got) != len(want) {
		t.Fatalf("ExpandClass = %q, want %q", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ExpandClass[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestRegexLengthVaries(t *testing.T) {
	r := NewResrap()
	if err := r.ParseGrammar("x", `program:[a-z];`); err != nil {
		t.Fatal(err)
	}
	seen := map[int]bool{}
	for seed := uint64(1); seed <= 100; seed++ {
		seen[len(r.GenerateWithSeeded("x", "program", seed, 10))] = true
	}
	if len(seen) < 2 {
		t.Fatalf("regex generation always produced the same length: %v", seen)
	}
	for n := range seen {
		if n < 3 || n > 7 {
			t.Fatalf("regex length %d outside 3..7", n)
		}
	}
}

func TestResrapMTBackpressure(t *testing.T) {
	r := NewResrapMT(2, 4)
	if err := r.ParseGrammar("x", `program:'a'^;`); err != nil {
		t.Fatal(err)
	}
	r.StartResrap()

	const jobs = 40
	done := make(chan struct{})
	go func() {
		defer close(done)
		for i := 0; i < jobs; i++ {
			if err := r.GenerateRandom(string(rune('A'+i%26)), "x", "program", 5); err != nil {
				t.Errorf("submit: %v", err)
				return
			}
		}
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("submitting more jobs than the wait queue blocked (#5)")
	}

	for i := 0; i < jobs; i++ {
		select {
		case <-r.GetCodeChannel():
		case <-time.After(5 * time.Second):
			t.Fatalf("timed out reading result %d", i)
		}
	}

	shutdownDone := make(chan struct{})
	go func() {
		r.ShutDownResrap()
		close(shutdownDone)
	}()
	select {
	case <-shutdownDone:
	case <-time.After(5 * time.Second):
		t.Fatal("ShutDownResrap hung")
	}
}

func TestResrapMTRejectsJobsAfterShutdown(t *testing.T) {
	r := NewResrapMT(1, 2)
	if err := r.ParseGrammar("x", `program:'a';`); err != nil {
		t.Fatal(err)
	}
	r.StartResrap()
	r.ShutDownResrap()
	if err := r.GenerateRandom("id", "x", "program", 5); err != ErrShutDown {
		t.Fatalf("GenerateRandom after shutdown = %v, want ErrShutDown", err)
	}
	r.ShutDownResrap() // idempotent
}

func TestNoPanicOnRandomGrammarInput(t *testing.T) {
	inputs := []string{
		"", ":", ";", ")", "(", "'", "<", "[", "@words",
		"program:", "program:'a'<", "program:[a-", "program:((('a');",
		"program:a; a:'b';", "@words x = a b\nprogram:x;",
	}
	for _, in := range inputs {
		func() {
			defer func() {
				if rec := recover(); rec != nil {
					t.Fatalf("panic for %q: %v", in, rec)
				}
			}()
			r := NewResrap()
			_ = r.ParseGrammar("x", in)
		}()
	}
}

// writeFile is a tiny helper so the tests do not need to import os directly in
// every case.
func writeFile(path, content string) error {
	return os.WriteFile(path, []byte(content), 0o644)
}
