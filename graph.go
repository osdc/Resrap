package resrap

import (
	"sort"
	"strings"

	"github.com/golang-collections/collections/stack"
)

type nextoption struct {
	node        *syntaxNode
	probability float32
}
type nodeType int8

const (
	start nodeType = iota
	header
	jump
	end
	ch
	rx
	pointer
	idk
	loop
	groupEnd
)

type syntaxNode struct {
	next    []nextoption // All options from here
	cf      []float32    // Cumulative frequency of all the options
	id      uint32       // The id of the node
	typ     nodeType
	pointer uint32 // In case its a pointer type
}

func (s *syntaxNode) AddEdgeNext(g *syntaxGraph, node *syntaxNode, probability float32) {
	newNode := nextoption{node: node, probability: probability}
	s.next = append(s.next, newNode)
	g.nodeRef[node.id] = node
}

type syntaxGraph struct {
	nodeRef      map[uint32]*syntaxNode
	namemap      map[string]uint32
	charmap      map[uint32]string
	regexhandler regexer
	words        map[uint32][]string // rule id -> generated words/keywords
	ruleEnds     map[uint32]uint32   // rule id -> its own end node id
}

func (s *syntaxGraph) GetNode(id uint32, typ nodeType) *syntaxNode {
	if s.nodeRef[id] != nil {
		return s.nodeRef[id]
	}
	newNode := &syntaxNode{nil, nil, id, typ, 0}
	s.nodeRef[id] = newNode
	return newNode
}

func newSyntaxGraph() syntaxGraph {
	return syntaxGraph{
		nodeRef:  make(map[uint32]*syntaxNode),
		words:    make(map[uint32][]string),
		ruleEnds: make(map[uint32]uint32),
	}
}

// SetWords attaches a list of words to a rule. References to that rule will
// emit a randomly selected word instead of expanding its grammar.
func (s *syntaxGraph) SetWords(rule string, words []string) error {
	id, ok := s.namemap[rule]
	if !ok {
		return errUnknownRule(rule)
	}
	if len(words) == 0 {
		delete(s.words, id)
		return nil
	}
	s.words[id] = append([]string(nil), words...)
	return nil
}

type ruleNotFoundError string

func (e ruleNotFoundError) Error() string { return "rule '" + string(e) + "' not found" }
func errUnknownRule(rule string) error    { return ruleNotFoundError(rule) }

func (s *syntaxGraph) Normalize() {
	// We will even out all the children going through the whole graph
	// And also create a cumulative frequency graph that will help in traversing
	for _, node := range s.nodeRef {
		// Extract the elements in a diff array
		CF := make([]float32, 0, len(node.next))

		var sum float32
		for _, n := range node.next {
			p := n.probability
			if p < 0 {
				p = 0
			}
			CF = append(CF, p)
			sum += p
		}
		if !(sum > 0) {
			// A node with zero/NaN weights should still be traversable. Treat
			// it as a uniform choice instead of producing NaN cumulative
			// frequencies (which used to panic in GraphWalk).
			for i := range CF {
				CF[i] = 1
			}
			sum = float32(len(CF))
		}
		var cf float32
		for i := range CF {
			CF[i] = cf + CF[i]/sum
			cf = CF[i]
		}
		node.cf = CF
		// Now we have a cool little CF Array Normalized to 1
		// When picking random values, we pick one between 0 and 1
		// And then choose its closest value from the array
		// For probability based selections
	}
}

func (s *syntaxGraph) choice(current *syntaxNode, prng *prng) *syntaxNode {
	if len(current.next) == 0 {
		return nil
	}
	if len(current.cf) != len(current.next) {
		// Graph was not normalized (or was built by hand). Fall back to a
		// uniform choice rather than indexing an empty array.
		return current.next[prng.RandomInt(0, len(current.next))].node
	}
	value := float32(prng.Random())
	index := sort.Search(len(current.cf), func(i int) bool {
		return current.cf[i] >= value
	})
	if index >= len(current.next) {
		index = len(current.next) - 1
	}
	if index < 0 {
		index = 0
	}
	return current.next[index].node
}

func (s *syntaxGraph) GraphWalk(prng *prng, start string, tokens int) string {
	return s.walk(prng, start, tokens, false)
}

// GraphWalkComplete is like GraphWalk but finishes the current top-level rule
// after the token budget is reached. This avoids cutting generated source code
// in the middle of a statement or expression. The budget is therefore a soft
// limit: output may contain slightly more than `tokens` terminals.
func (s *syntaxGraph) GraphWalkComplete(prng *prng, start string, tokens int) string {
	return s.walk(prng, start, tokens, true)
}

func (s *syntaxGraph) walk(prng *prng, start string, tokens int, complete bool) string {
	if tokens <= 0 {
		return ""
	}
	var result strings.Builder
	jumpStack := stack.New()
	startID, ok := s.namemap[start]
	if !ok {
		return ""
	}
	// A rule can itself be a word list, e.g. @words identifier = ...
	if words := s.words[startID]; len(words) > 0 {
		return words[prng.RandomInt(0, len(words))]
	}
	startingNode := s.nodeRef[startID]
	if startingNode == nil {
		return ""
	}

	printedTokens := 0
	startEnd := s.ruleEnds[startID]
	maxIterations := tokens*1000 + 10000
	if maxIterations < 10000 {
		maxIterations = 10000
	}
	iterations := 0
	current := startingNode
	for current != nil {
		iterations++
		if iterations > maxIterations {
			return result.String()
		}
		if tokens > 0 && printedTokens >= tokens {
			if !complete {
				return result.String()
			}
			// In complete mode only stop at the starting rule's boundary.
			if current.id == startEnd && jumpStack.Len() == 0 {
				return result.String()
			}
		}
		// Process logic only if name starts with ' or [

		if current.typ == ch {
			// Extract content between quotes and handle escape sequences
			content := s.charmap[current.id]
			result.WriteString(unescapeString(content))
			printedTokens++
		} else if current.typ == rx {
			result.WriteString(s.regexhandler.GenerateString(s.charmap[current.id], prng))
			printedTokens++
		} else if current.typ == loop {
			// '^' repeats the construct immediately before it until the token
			// budget is reached. In complete mode the walk then falls through
			// to whatever follows the loop, so constructs such as a closing
			// brace are emitted exactly once.
			if !complete || printedTokens < tokens {
				target := s.nodeRef[current.pointer]
				if target == nil {
					return result.String()
				}
				current = target
				continue
			}
		} else if current.typ == pointer {
			// If the target rule has a word list, emit a word directly and
			// return to the caller instead of traversing that rule's grammar.
			if words := s.words[current.pointer]; len(words) > 0 {
				result.WriteString(words[prng.RandomInt(0, len(words))])
				printedTokens++
				if len(current.next) > 0 {
					current = current.next[0].node
				} else {
					current = nil
				}
				continue
			}
			if len(current.next) == 0 {
				return result.String()
			}
			jumpStack.Push(current.next[0].node.id)
			current = s.GetNode(current.pointer, header)
			continue // Skip the normal next node selection
		} else if current.typ == end {
			if jumpStack.Len() != 0 {
				nameInt := jumpStack.Pop()
				id, ok := nameInt.(uint32)
				if !ok {
					break
				}
				current = s.GetNode(id, idk)
				continue // Skip the normal next node selection
			}
		}

		// move to next (randomly selected if multiple options)
		current = s.choice(current, prng)
	}
	return result.String()
}

// Helper function to handle escape sequences
func unescapeString(s string) string {
	result := make([]byte, 0, len(s))
	for i := 0; i < len(s); i++ {
		if s[i] == '\\' && i+1 < len(s) {
			switch s[i+1] {
			case 'n':
				result = append(result, '\n')
				i++ // skip the next character
			case 't':
				result = append(result, '\t')
				i++
			case 'r':
				result = append(result, '\r')
				i++
			case '\\':
				result = append(result, '\\')
				i++
			case '\'':
				result = append(result, '\'')
				i++
			case '"':
				result = append(result, '"')
				i++
			case '0':
				result = append(result, 0)
				i++
			default:
				// If it's not a recognized escape sequence, keep both characters
				result = append(result, s[i])
			}
		} else {
			result = append(result, s[i])
		}
	}
	return string(result)
}
