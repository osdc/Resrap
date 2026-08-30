package resrap

import (
	"math/rand"
)

type prng struct {
	seed   uint64
	number uint64
}

func (p *prng) setSeed(seed uint64) {
	p.number = seed
	p.seed = seed
}
func (p *prng) generateSeed() {
	seed := rand.Uint64()
	p.number = seed
	p.seed = seed
}
func (p *prng) nextPRN() uint64 {
	// Using the XOR shift method for PRN generation
	p.number ^= p.number << 13
	p.number ^= p.number >> 7
	p.number ^= p.number << 17
	return p.number
}

// random returns a float64 in [0,1)
func (p *prng) Random() float64 {
	// Take the next 53 random bits (same precision as math/rand.Float64)
	v := p.nextPRN() >> 11        // keep top 53 bits
	return float64(v) / (1 << 53) // normalize to [0,1)
}

// randomInt returns an int in [low, high)
func (p *prng) RandomInt(low, high int) int {
	if high <= low {
		return low // avoid division by zero or negative range
	}
	r := p.Random()
	return low + int(r*float64(high-low))
}

func newPRNG(seed uint64) prng {
	p := prng{}
	if seed == 0 {
		p.generateSeed()
	} else {
		p.setSeed(seed)
	}
	return p
}
