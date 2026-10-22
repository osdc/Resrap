// Package resrap is a seedable, grammar-based code snippet generator.
// Instead of parsing code, it generates code from formal grammars (ABNF-style),
// producing endless, realistic-looking snippets.
package resrap

import (
	"fmt"
	"os"
)

// Resrap is the main accesspoint for singlethreaded uses
// Pretty Much Collection of graphs which can be generated using parsing grammar
type Resrap struct {
	languageGraph map[string]lang
}

// NewResrap creates and returns a new Resrap instance.
// The returned instance starts with no loaded grammars.
func NewResrap() *Resrap {
	return &Resrap{
		languageGraph: make(map[string]lang),
	}
}

// ParseGrammar parses a grammar string and stores it under the given name.
// name: a unique identifier for this grammar (e.g., "C"), should be in ABNF format(Check osdc/resrap for more info on that).
// Returns error generated while parsing
func (r *Resrap) ParseGrammar(name, grammar string) error {
	l := newLang()
	err := l.ParserString(grammar)
	if err != nil {
		return err
	}
	r.languageGraph[name] = l
	r.languageGraph[name].graph.Normalize()
	return nil
}

// ParseGrammarFile parses a grammar from a file and stores it under the given name.
// name: a unique identifier for this grammar (e.g., "C"), should be in ABNF format(Check osdc/resrap for more info on that).
// location: path to the grammar file.
// Returns error generated while parsing
func (r *Resrap) ParseGrammarFile(name, location string) error {
	l := newLang()
	err := l.ParserFile(location)
	if err != nil {
		return err
	}
	r.languageGraph[name] = l
	r.languageGraph[name].graph.Normalize()

	return nil
}

// AddWords attaches a list of words to a rule in a parsed grammar. Any
// reference to that rule will emit one of the words instead of expanding the
// rule. This is useful for identifiers, keywords and other "word-like"
// terminals that look better than random regex output.
func (r *Resrap) AddWords(name, rule string, words []string) error {
	l, ok := r.languageGraph[name]
	if !ok || l.graph == nil {
		return fmt.Errorf("resrap: grammar %q is not loaded", name)
	}
	return l.graph.SetWords(rule, words)
}

// AddKeywords is an alias for AddWords.
func (r *Resrap) AddKeywords(name, rule string, words []string) error {
	return r.AddWords(name, rule, words)
}

// GenerateComplete generates content but only stops at the boundary of the
// starting rule once the token budget is reached. Unlike GenerateRandom, this
// does not cut the output in the middle of a statement or expression, so it is
// suitable for generating source code that should be parsed or compiled.
func (r *Resrap) GenerateComplete(name, startingNode string, tokens int) string {
	prng := newPRNG(0)
	return r.generateComplete(name, startingNode, &prng, tokens)
}

// GenerateCompleteWithSeeded is the deterministic variant of GenerateComplete.
func (r *Resrap) GenerateCompleteWithSeeded(name, startingNode string, seed uint64, tokens int) string {
	prng := newPRNG(seed)
	return r.generateComplete(name, startingNode, &prng, tokens)
}

func (r *Resrap) generateComplete(name, startingNode string, prng *prng, tokens int) string {
	l, ok := r.languageGraph[name]
	if !ok || l.graph == nil {
		return ""
	}
	return l.graph.GraphWalkComplete(prng, startingNode, tokens)
}

// GenerateRandom generates content from the grammar identified by 'name'.
// startingNode: the starting heading in the grammar for generation.
// Returns a string containing the generated content.
// The generation is non-deterministic (random).
func (r *Resrap) GenerateRandom(name, startingNode string, tokens int) string {
	prng := newPRNG(0)
	return r.generate(name, startingNode, &prng, tokens)
}

// GenerateWithSeeded generates content from the grammar identified by 'name'.
// startingNode: the starting symbol in the grammar for generation.
// seed: a numeric seed to make generation deterministic.
// Returns a string containing the generated content.
func (r *Resrap) GenerateWithSeeded(name, startingNode string, seed uint64, tokens int) string {
	prng := newPRNG(seed)
	return r.generate(name, startingNode, &prng, tokens)
}

func (r *Resrap) generate(name, startingNode string, prng *prng, tokens int) string {
	l, ok := r.languageGraph[name]
	if !ok || l.graph == nil {
		return ""
	}
	return l.graph.GraphWalk(prng, startingNode, tokens)
}

// GenerateCodebase takes a config like one below  and generates a complete codebase
// src/ # src is the root folder, with code and more inside of it
//
//	code/
//	  c[10x20 code_*.c] #Generate 10 c files named code_[sth_unique].c each 20 tokens long in src/code
//	  core/
//	    c[5x1000 *.c]
//	more/
//	  sql[10x100 *.sql]
func (r *Resrap) GenerateCodebase(configLoc, target string) error {
	config, err := os.ReadFile(configLoc)
	if err != nil {
		return err
	}
	parent, err := ParseDSL(string(config))
	if err != nil {
		return err
	}
	return parent.generateNode(r, target)
}
