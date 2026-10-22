package resrap

import (
	"strings"
	"unicode"
)

type regexer struct {
	cachedRex map[string]cacheRexState
}

func newRegexer() regexer {
	return regexer{cachedRex: make(map[string]cacheRexState)}
}

type cacheRexState struct {
	cumuFreq []float32
	options  []rune
}

func (r *regexer) GenerateString(regex string, prn *prng) string {
	state, ok := r.cachedRex[regex]
	if !ok || len(state.cumuFreq) == 0 || len(state.options) == 0 {
		return ""
	}
	size := prn.RandomInt(3, 8) // generate a size between 3 and 7
	var sb strings.Builder

	for i := 0; i < size; i++ {
		x := prn.Random() // random float 0-1
		idx := closestIndex(state.cumuFreq, float32(x))
		sb.WriteRune(state.options[idx]) // use WriteRune instead of WriteByte
	}

	return sb.String()
}
func (r *regexer) CacheRegex(regex string) {
	tokens := r.ExpandClass(regex)
	biasarr := make([]float32, 0, len(tokens))
	var sum float32
	for _, token := range tokens {
		bias := float32(r.Bias(token))
		biasarr = append(biasarr, bias)
		sum += bias
	}
	for i := range biasarr {
		biasarr[i] /= sum
	}

	cdf := make([]float32, len(biasarr))
	cum := float32(0)
	for i, w := range biasarr {
		cum += w
		cdf[i] = cum
	}

	r.cachedRex[regex] = cacheRexState{cumuFreq: cdf, options: tokens}
}

// closestIndex finds the first index in cdf where cdf[idx] >= x
func closestIndex(cdf []float32, x float32) int {
	for i, val := range cdf {
		if x <= val {
			return i
		}
	}
	return len(cdf) - 1 // fallback
}
func (r *regexer) ExpandClass(class string) []rune {
	var chars []rune
	runes := []rune(class)

	for i := 0; i < len(runes); i++ {
		if runes[i] == '\\' && i+1 < len(runes) {
			chars = append(chars, decodeClassEscape(runes[i+1]))
			i++
			continue
		}
		if i+2 < len(runes) && runes[i+1] == '-' && runes[i+2] != '\\' {
			// range a-z
			for c := runes[i]; c <= runes[i+2]; c++ {
				chars = append(chars, c)
			}
			i += 2
		} else {
			chars = append(chars, runes[i])
		}
	}
	return chars
}

func decodeClassEscape(r rune) rune {
	switch r {
	case 'n':
		return '\n'
	case 't':
		return '\t'
	case 'r':
		return '\r'
	default:
		// \], \-, \^ and friends mean the literal character.
		return r
	}
}

// Bias returns an integer weight for a rune based on usage frequency
func (r *regexer) Bias(c rune) int {
	// Uppercase letters: slightly less likely than lowercase. Check this
	// before the lowercase switch, otherwise 'E' matches the 'e' case.
	if unicode.IsUpper(c) {
		b := r.Bias(unicode.ToLower(c)) / 2
		if b < 1 {
			b = 1
		}
		return b
	}

	// Lowercase letters: frequency in English words (roughly)
	switch c {
	case 'e':
		return 12
	case 'a', 'i', 'o':
		return 9
	case 'n', 'r', 't', 's', 'l':
		return 6
	case 'c', 'd', 'm', 'u', 'p', 'b', 'g':
		return 4
	case 'f', 'h', 'v', 'k', 'w', 'y':
		return 3
	case 'j', 'x', 'q', 'z':
		return 1
	}

	// Digits: moderately common
	if unicode.IsDigit(c) {
		return 3
	}

	// Underscore: quite common in identifiers
	if c == '_' {
		return 5
	}

	// Everything else: low probability
	return 1
}
