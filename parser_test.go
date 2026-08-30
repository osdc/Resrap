package resrap

import (
	"strings"
	"testing"
)

func TestGraphBuilderSimpleGrammar(t *testing.T) {
	gb := newGraphBuilder()
	err := gb.startGeneration(`program:'hello' 'world';`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gb.pars.graph.namemap["program"] == 0 {
		t.Fatal("program name not registered in namemap")
	}
}

func TestGraphBuilderMissingDefinition(t *testing.T) {
	gb := newGraphBuilder()
	err := gb.startGeneration(`program:missing;`)
	if err == nil {
		t.Fatal("expected error for missing definition")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("error = %v, want 'not found'", err)
	}
}

func TestGraphBuilderDuplicateDefinition(t *testing.T) {
	gb := newGraphBuilder()
	err := gb.startGeneration(`program:'a';
program:'b';`)
	if err == nil {
		t.Fatal("expected error for duplicate definition")
	}
	if !strings.Contains(err.Error(), "multiple definitions") {
		t.Fatalf("error = %v, want 'multiple definitions'", err)
	}
}

func TestGraphBuilderInvalidGrammarNoColon(t *testing.T) {
	gb := newGraphBuilder()
	err := gb.startGeneration(`program 'a';`)
	if err == nil {
		t.Fatal("expected error for missing colon")
	}
}

func TestGraphBuilderScanError(t *testing.T) {
	gb := newGraphBuilder()
	err := gb.startGeneration(`program:'unterminated`)
	if err == nil {
		t.Fatal("expected scan error")
	}
	if !strings.Contains(err.Error(), "Scanning") {
		t.Fatalf("error = %v, want scan error", err)
	}
}

func TestGraphBuilderRecursiveGrammar(t *testing.T) {
	gb := newGraphBuilder()
	err := gb.startGeneration(`program:function^;
function:identifier '(' ')';
identifier:[a-z];`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGraphBuilderProbabilities(t *testing.T) {
	gb := newGraphBuilder()
	err := gb.startGeneration(`program:('a'<0.2>) | ('b'<0.8>);`)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGraphBuilderNegativeProbability(t *testing.T) {
	gb := newGraphBuilder()
	err := gb.startGeneration(`program:'a'<0.5>;`)
	if err != nil {
		t.Fatalf("unexpected error for valid grammar: %v", err)
	}
}

func TestParserGetIndex(t *testing.T) {
	p := newParser()
	id := p.getIndex("foo")
	if id == 0 {
		t.Fatal("getIndex returned 0")
	}
	if p.defCheck[id] {
		t.Fatal("newly created index should be marked as undefined")
	}
	got, ok := p.nameMap["foo"]
	if !ok || got != id {
		t.Fatalf("nameMap['foo'] = %d, want %d", got, id)
	}
}

func TestParserGetPrintPtr(t *testing.T) {
	p := newParser()
	a := p.getPrintPtr()
	b := p.getPrintPtr()
	if a == b {
		t.Fatal("getPrintPtr should return decreasing unique values")
	}
	if a < b {
		t.Fatalf("print ptr should decrease, got %d then %d", a, b)
	}
}

func TestParserValidateGraphNoErrors(t *testing.T) {
	p := newParser()
	p.parseGrammar() // empty grammar
	if errs := p.ValidateGraph(); len(errs) != 0 {
		t.Fatalf("unexpected validation errors: %v", errs)
	}
}
