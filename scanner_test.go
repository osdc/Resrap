package resrap

import "testing"

func TestExtractTokensBasics(t *testing.T) {
	input := `program:'hello'`
	tokens, errs := extracttokens(input)
	if len(errs) != 0 {
		t.Fatalf("unexpected scan errors: %v", errs)
	}
	if len(tokens) != 3 {
		t.Fatalf("got %d tokens, want 3", len(tokens))
	}
	if tokens[0].typ != identifier || tokens[0].text != "program" {
		t.Fatalf("token[0] = %+v, want identifier 'program'", tokens[0])
	}
	if tokens[1].typ != colon {
		t.Fatalf("token[1].typ = %v, want colon", tokens[1].typ)
	}
	if tokens[2].typ != character || tokens[2].text != "hello" {
		t.Fatalf("token[2] = %+v, want character 'hello'", tokens[2])
	}
}

func TestExtractTokensOperators(t *testing.T) {
	input := `a:b+c*d?|;(e)^<0.3>`
	tokens, errs := extracttokens(input)
	if len(errs) != 0 {
		t.Fatalf("unexpected scan errors: %v", errs)
	}
	wantTypes := []tokenType{
		identifier, colon, identifier, oneormore, identifier, anyno, identifier,
		maybe, option, padding, bracopen, identifier, bracclose, infinite,
		probability,
	}
	if len(tokens) != len(wantTypes) {
		t.Fatalf("got %d tokens, want %d", len(tokens), len(wantTypes))
	}
	for i, want := range wantTypes {
		if tokens[i].typ != want {
			t.Fatalf("token[%d].typ = %v, want %v", i, tokens[i].typ, want)
		}
	}
}

func TestExtractTokensProbability(t *testing.T) {
	tokens, errs := extracttokens(`a:<0.25>`)
	if len(errs) != 0 {
		t.Fatalf("unexpected scan errors: %v", errs)
	}
	if len(tokens) != 3 {
		t.Fatalf("got %d tokens, want 3", len(tokens))
	}
	if tokens[2].typ != probability || tokens[2].text != "0.25" {
		t.Fatalf("token[2] = %+v, want probability '0.25'", tokens[2])
	}
}

func TestExtractTokensRegex(t *testing.T) {
	tokens, errs := extracttokens(`a:[a-z]`)
	if len(errs) != 0 {
		t.Fatalf("unexpected scan errors: %v", errs)
	}
	if len(tokens) != 3 {
		t.Fatalf("got %d tokens, want 3", len(tokens))
	}
	if tokens[2].typ != regex || tokens[2].text != "a-z" {
		t.Fatalf("token[2] = %+v, want regex 'a-z'", tokens[2])
	}
}

func TestExtractTokensEscape(t *testing.T) {
	tokens, errs := extracttokens(`a:'\n'`)
	if len(errs) != 0 {
		t.Fatalf("unexpected scan errors: %v", errs)
	}
	if len(tokens) != 3 {
		t.Fatalf("got %d tokens, want 3", len(tokens))
	}
	if tokens[2].typ != character || tokens[2].text != `\n` {
		t.Fatalf("token[2] = %+v, want character '\\n'", tokens[2])
	}
}

func TestExtractTokensIdentifierWithUnderscore(t *testing.T) {
	tokens, errs := extracttokens(`my_ident:ok_2`)
	if len(errs) != 0 {
		t.Fatalf("unexpected scan errors: %v", errs)
	}
	if len(tokens) != 3 {
		t.Fatalf("got %d tokens, want 3", len(tokens))
	}
	if tokens[0].text != "my_ident" {
		t.Fatalf("token[0].text = %q, want 'my_ident'", tokens[0].text)
	}
	if tokens[2].text != "ok_2" {
		t.Fatalf("token[2].text = %q, want 'ok_2'", tokens[2].text)
	}
}

func TestExtractTokensUnterminatedCharacter(t *testing.T) {
	_, errs := extracttokens(`a:'oops`)
	if len(errs) != 1 {
		t.Fatalf("got %d errors, want 1", len(errs))
	}
	if errs[0].Msg != "unterminated '\\''" {
		t.Fatalf("error message = %q, want unterminated quote", errs[0].Msg)
	}
	if errs[0].Pos == 0 {
		t.Fatal("error Pos should be non-zero")
	}
}

func TestExtractTokensUnterminatedProbability(t *testing.T) {
	_, errs := extracttokens(`a:<0.5`)
	if len(errs) != 1 {
		t.Fatalf("got %d errors, want 1", len(errs))
	}
}

func TestExtractTokensEmptyInput(t *testing.T) {
	tokens, errs := extracttokens("")
	if len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	if len(tokens) != 0 {
		t.Fatalf("got %d tokens, want 0", len(tokens))
	}
}

func TestScanDelimitedHandlesNestedDelimiter(t *testing.T) {
	s := scanner{Input: `'a'b'c'`}
	s.pos = 1 // position just after the opening quote
	// Character literals stop at the first closing quote.
	val, err := s.scanDelimited('\'', '\'', false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if val != "a" {
		t.Fatalf("scanned value = %q, want 'a'", val)
	}
}
