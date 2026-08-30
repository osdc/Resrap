package resrap

import (
	"strings"
	"testing"
)

func TestExpandClassSimple(t *testing.T) {
	r := newRegexer()
	got := r.ExpandClass("abc")
	want := []rune{'a', 'b', 'c'}
	if len(got) != len(want) {
		t.Fatalf("ExpandClass(%q) length = %d, want %d", "abc", len(got), len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ExpandClass(%q)[%d] = %q, want %q", "abc", i, got[i], want[i])
		}
	}
}

func TestExpandClassRange(t *testing.T) {
	r := newRegexer()
	got := r.ExpandClass("a-z")
	if len(got) != 26 {
		t.Fatalf("ExpandClass(a-z) length = %d, want 26", len(got))
	}
	if got[0] != 'a' || got[25] != 'z' {
		t.Fatalf("ExpandClass(a-z) = %q, want a..z", got)
	}
}

func TestExpandClassMixedRangeAndLiteral(t *testing.T) {
	r := newRegexer()
	got := r.ExpandClass("a-c_")
	if len(got) != 4 {
		t.Fatalf("ExpandClass(a-c_) length = %d, want 4", len(got))
	}
	want := []rune{'a', 'b', 'c', '_'}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("ExpandClass(a-c_)[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestExpandClassDigitRange(t *testing.T) {
	r := newRegexer()
	got := r.ExpandClass("0-9")
	if len(got) != 10 {
		t.Fatalf("ExpandClass(0-9) length = %d, want 10", len(got))
	}
	if got[0] != '0' || got[9] != '9' {
		t.Fatalf("ExpandClass(0-9) = %q, want 0..9", got)
	}
}

func TestBias(t *testing.T) {
	r := newRegexer()
	tests := []struct {
		in   rune
		want int
	}{
		{'e', 12},
		{'a', 9},
		{'x', 1},
		{'j', 1},
		{'_', 5},
		{'0', 3},
		{'Z', 1},
		{'$', 1},
	}
	for _, tt := range tests {
		if got := r.Bias(tt.in); got != tt.want {
			t.Fatalf("Bias(%q) = %d, want %d", tt.in, got, tt.want)
		}
	}
}

func TestClosestIndex(t *testing.T) {
	cdf := []float32{0.25, 0.5, 0.75, 1.0}
	tests := []struct {
		x    float32
		want int
	}{
		{0.0, 0},
		{0.24, 0},
		{0.25, 0},
		{0.5, 1},
		{0.99, 3},
	}
	for _, tt := range tests {
		if got := closestIndex(cdf, tt.x); got != tt.want {
			t.Fatalf("closestIndex(%v) = %d, want %d", tt.x, got, tt.want)
		}
	}
}

func TestClosestIndexEmpty(t *testing.T) {
	if got := closestIndex(nil, 0.5); got != -1 {
		t.Fatalf("closestIndex(empty) = %d, want -1", got)
	}
}

func TestCacheRegex(t *testing.T) {
	r := newRegexer()
	r.CacheRegex("abc")
	state, ok := r.cachedRex["abc"]
	if !ok {
		t.Fatal("CacheRegex did not cache 'abc'")
	}
	if len(state.options) != 3 {
		t.Fatalf("cached options length = %d, want 3", len(state.options))
	}
	if len(state.cumuFreq) != 3 {
		t.Fatalf("cached cumuFreq length = %d, want 3", len(state.cumuFreq))
	}
}

func TestGenerateString(t *testing.T) {
	r := newRegexer()
	r.CacheRegex("a-z")
	p := newPRNG(7)
	out := r.GenerateString("a-z", &p)
	if out == "" {
		t.Fatal("GenerateString returned empty string")
	}
	for _, c := range out {
		if c < 'a' || c > 'z' {
			t.Fatalf("GenerateString produced rune %q outside a-z", c)
		}
	}
}

func TestGenerateStringDeterministic(t *testing.T) {
	r1 := newRegexer()
	r2 := newRegexer()
	r1.CacheRegex("a-z")
	r2.CacheRegex("a-z")
	p1 := newPRNG(42)
	p2 := newPRNG(42)
	a := r1.GenerateString("a-z", &p1)
	b := r2.GenerateString("a-z", &p2)
	if a != b {
		t.Fatalf("GenerateString not deterministic: %q vs %q", a, b)
	}
}

func TestGenerateStringUncached(t *testing.T) {
	r := newRegexer()
	// Without caching, GenerateString should not panic and returns empty.
	p := newPRNG(3)
	out := r.GenerateString("nonexistent", &p)
	if out != "" {
		t.Fatalf("GenerateString returned %q for uncached regex, want empty", out)
	}
}

func TestBiasWeightsArePositive(t *testing.T) {
	r := newRegexer()
	for c := rune(32); c < 127; c++ {
		if r.Bias(c) <= 0 {
			t.Fatalf("Bias(%q) = %d, want positive", c, r.Bias(c))
		}
	}
}

func TestGenerateStringOutputOnlyCachedChars(t *testing.T) {
	r := newRegexer()
	r.CacheRegex("xyz")
	p := newPRNG(11)
	out := r.GenerateString("xyz", &p)
	for _, c := range out {
		if !strings.ContainsRune("xyz", c) {
			t.Fatalf("GenerateString produced rune %q outside cached set", c)
		}
	}
}
