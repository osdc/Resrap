package templates

import (
	"strings"
	"testing"

	"github.com/osdc/resrap"
)

func TestTemplatesParseAndGenerate(t *testing.T) {
	all, err := List()
	if err != nil {
		t.Fatal(err)
	}
	if len(all) < 9 {
		t.Fatalf("expected all language templates, got %d", len(all))
	}
	for _, tmpl := range all {
		tmpl := tmpl
		t.Run(tmpl.Name, func(t *testing.T) {
			if tmpl.Grammar == "" || tmpl.Language == "" || tmpl.Description == "" {
				t.Fatalf("incomplete template metadata: %+v", tmpl)
			}
			r := resrap.NewResrap()
			if err := tmpl.Load(r); err != nil {
				t.Fatalf("parse %s: %v", tmpl.Name, err)
			}
			var first string
			for _, seed := range []uint64{1, 7, 42, 99} {
				out := r.GenerateCompleteWithSeeded(tmpl.Name, "program", seed, 60)
				if strings.TrimSpace(out) == "" {
					t.Fatalf("seed %d produced empty output", seed)
				}
				if first == "" {
					first = out
				} else if out == first {
					t.Fatalf("seed %d produced identical output", seed)
				}
			}

			// The built-in templates are demos, not five-line snippets. At a
			// realistic budget they should produce a substantial program.
			out := r.GenerateCompleteWithSeeded(tmpl.Name, "program", 7, 500)
			if lines := strings.Count(out, "\n"); lines < 20 {
				t.Fatalf("500-token output is too small: %d lines\n%s", lines, out)
			}
		})
	}
}

func TestGetIsCaseInsensitiveAndAcceptsExtension(t *testing.T) {
	a := MustGet("CPP.g4")
	b := MustGet("cpp")
	if a.Name != "cpp" || a.Grammar != b.Grammar {
		t.Fatalf("lookup normalization failed: %q vs %q", a.Name, b.Name)
	}
	if _, err := Get("not-a-language"); err == nil {
		t.Fatal("expected an error for an unknown template")
	}
}
