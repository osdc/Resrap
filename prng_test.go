package resrap

import "testing"

func TestPRNGDeterministic(t *testing.T) {
	seed := uint64(42)
	a := newPRNG(seed)
	b := newPRNG(seed)

	for i := 0; i < 100; i++ {
		if a.Random() != b.Random() {
			t.Fatalf("PRNG with same seed diverged at iteration %d", i)
		}
	}
}

func TestPRNGRandomInRange(t *testing.T) {
	p := newPRNG(1)
	for i := 0; i < 1000; i++ {
		v := p.Random()
		if v < 0 || v >= 1 {
			t.Fatalf("Random() returned %f, want value in [0,1)", v)
		}
	}
}

func TestPRNGRandomIntRange(t *testing.T) {
	p := newPRNG(2)
	const lo, hi = 3, 8
	for i := 0; i < 1000; i++ {
		v := p.RandomInt(lo, hi)
		if v < lo || v >= hi {
			t.Fatalf("RandomInt(%d,%d) returned %d", lo, hi, v)
		}
	}
}

func TestPRNGRandomIntInvalidRange(t *testing.T) {
	p := newPRNG(3)
	if got := p.RandomInt(5, 5); got != 5 {
		t.Fatalf("RandomInt(5,5) = %d, want 5", got)
	}
	if got := p.RandomInt(10, 2); got != 10 {
		t.Fatalf("RandomInt(10,2) = %d, want 10", got)
	}
}

func TestPRNGZeroSeedGeneratesRandom(t *testing.T) {
	a := newPRNG(0)
	b := newPRNG(0)
	// Both should produce non-zero internal numbers; they are not deterministic.
	if a.number == 0 {
		t.Fatal("PRNG with seed 0 has zero internal state")
	}
	if b.number == 0 {
		t.Fatal("PRNG with seed 0 has zero internal state")
	}
}

func TestPRNGSetSeed(t *testing.T) {
	p := prng{}
	p.setSeed(99)
	if p.seed != 99 || p.number != 99 {
		t.Fatalf("setSeed did not set seed/number correctly: %+v", p)
	}
}
