package resrap

import (
	"os"
	"testing"
	"time"
)

const testGrammarMT = `program:function^;
function:'f' '(' ')';
`

func TestNewResrapMT(t *testing.T) {
	r := NewResrapMT(2, 10)
	if r == nil {
		t.Fatal("NewResrapMT returned nil")
	}
	if r.poolsize != 2 || r.waitqueuesize != 10 {
		t.Fatalf("poolsize=%d waitqueuesize=%d, want 2/10", r.poolsize, r.waitqueuesize)
	}
}

func TestResrapMTParseGrammar(t *testing.T) {
	r := NewResrapMT(2, 10)
	if err := r.ParseGrammar("Test", testGrammarMT); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResrapMTParseGrammarFile(t *testing.T) {
	dir := t.TempDir()
	path := dir + "/grammar.txt"
	if err := os.WriteFile(path, []byte(testGrammarMT), 0644); err != nil {
		t.Fatal(err)
	}

	r := NewResrapMT(2, 10)
	if err := r.ParseGrammarFile("Test", path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResrapMTParseGrammarInvalid(t *testing.T) {
	r := NewResrapMT(2, 10)
	if err := r.ParseGrammar("Test", "program:missing;"); err == nil {
		t.Fatal("expected error for invalid grammar")
	}
}

func TestResrapMTGenerateRandom(t *testing.T) {
	r := NewResrapMT(2, 10)
	if err := r.ParseGrammar("Test", testGrammarMT); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r.StartResrap()
	defer r.ShutDownResrap()

	r.GenerateRandom("job1", "Test", "program", 20)

	select {
	case res := <-r.GetCodeChannel():
		if res.Id != "job1" {
			t.Fatalf("result id = %q, want job1", res.Id)
		}
		if res.Code == "" {
			t.Fatal("generated code is empty")
		}
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for generation result")
	}
}

func TestResrapMTGenerateWithSeeded(t *testing.T) {
	r := NewResrapMT(2, 10)
	if err := r.ParseGrammar("Test", testGrammarMT); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r.StartResrap()
	defer r.ShutDownResrap()

	r.GenerateWithSeeded("job1", "Test", "program", 42, 20)

	select {
	case <-r.GetCodeChannel():
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for generation result")
	}
}

func TestResrapMTMultipleJobs(t *testing.T) {
	r := NewResrapMT(4, 100)
	if err := r.ParseGrammar("Test", testGrammarMT); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r.StartResrap()
	defer r.ShutDownResrap()

	const jobs = 10
	for i := 0; i < jobs; i++ {
		r.GenerateRandom(string(rune('A'+i)), "Test", "program", 10)
	}

	results := map[string]bool{}
	for i := 0; i < jobs; i++ {
		select {
		case res := <-r.GetCodeChannel():
			results[res.Id] = true
			if res.Code == "" {
				t.Fatal("empty code in result")
			}
		case <-time.After(5 * time.Second):
			t.Fatal("timed out waiting for generation results")
		}
	}
	for i := 0; i < jobs; i++ {
		id := string(rune('A' + i))
		if !results[id] {
			t.Fatalf("missing result for job %q", id)
		}
	}
}

func TestResrapMTSeededDeterministic(t *testing.T) {
	a := NewResrapMT(1, 10)
	b := NewResrapMT(1, 10)
	if err := a.ParseGrammar("Test", testGrammarMT); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := b.ParseGrammar("Test", testGrammarMT); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	a.StartResrap()
	b.StartResrap()
	defer a.ShutDownResrap()
	defer b.ShutDownResrap()

	a.GenerateWithSeeded("1", "Test", "program", 999, 20)
	b.GenerateWithSeeded("1", "Test", "program", 999, 20)

	var codeA, codeB string
	select {
	case res := <-a.GetCodeChannel():
		codeA = res.Code
	case <-time.After(5 * time.Second):
		t.Fatal("timed out for A")
	}
	select {
	case res := <-b.GetCodeChannel():
		codeB = res.Code
	case <-time.After(5 * time.Second):
		t.Fatal("timed out for B")
	}
	if codeA != codeB {
		t.Fatalf("seeded generation not deterministic:\n%q\nvs\n%q", codeA, codeB)
	}
}
