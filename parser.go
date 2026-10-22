package resrap

import (
	"fmt"
	"math"
	"slices"
	"strconv"
)

type parser struct {
	funcPtr      uint32
	printPtr     uint32
	nameMap      map[string]uint32  // Maps func names to their ids
	revNameMap   map[uint32]string  // Help in debugging
	defCheck     map[uint32]bool    // To check if a function exists
	charmap      map[uint32]string  // To store the print values corresponding to ids
	interRep     map[uint32][]token // Intermediate Representation
	words        map[string][]string
	tokens       []token
	errors       []error
	index        int
	graph        syntaxGraph
	regexhandler regexer
}

func newParser() parser {
	return parser{
		printPtr: math.MaxInt32, // grows downward
		funcPtr:  1000,          // IDs below 1000 reserved for core graph nodes

		// allocate maps to avoid nil panics
		nameMap:      make(map[string]uint32),
		defCheck:     make(map[uint32]bool),
		charmap:      make(map[uint32]string),
		interRep:     make(map[uint32][]token),
		revNameMap:   make(map[uint32]string),
		words:        make(map[string][]string),
		tokens:       []token{},
		errors:       []error{},
		graph:        newSyntaxGraph(),
		regexhandler: newRegexer(),
	}
}

func (i *parser) getPrintPtr() uint32 {
	i.printPtr--
	return i.printPtr
}
func (i *parser) getFuncPtr() uint32 {
	i.funcPtr++
	return i.funcPtr
}
func (i *parser) curr() token {
	if i.index < 0 || i.index >= len(i.tokens) {
		return token{typ: padding}
	}
	return i.tokens[i.index]
}

func (i *parser) errAt(t token, format string, args ...any) {
	msg := fmt.Sprintf(format, args...)
	i.errors = append(i.errors, fmt.Errorf("line %d:%d: %s", t.line, t.col, msg))
}

func (i *parser) errHere(format string, args ...any) {
	i.errAt(i.curr(), format, args...)
}

func (i *parser) match(word tokenType, expec []tokenType) bool {
	return slices.Contains(expec, word)
}
func (i *parser) expect(expected []tokenType, errmsg string) bool {
	if i.index >= len(i.tokens) {
		i.errors = append(i.errors, fmt.Errorf("%s (unexpected end of grammar)", errmsg))
		return true
	}
	if !i.match(i.curr().typ, expected) {
		i.errHere("%s", errmsg)
		i.index++
		return true
	}
	i.index++
	return false
}

func (i *parser) getIndex(name string) uint32 {
	value, ok := i.nameMap[name]
	if ok {
		return value
	}
	value = i.getFuncPtr()
	i.defCheck[value] = false // tb set back to true in the subject, else remain false
	i.nameMap[name] = value
	i.revNameMap[value] = name
	return value
}
func (i *parser) parseGrammar() {
	for i.index < len(i.tokens) {
		i.parseSubject()
		if len(i.errors) > 0 {
			return // Crash on errors for now
		}
	}
}
func (i *parser) parseSubject() {
	subject := i.curr()

	if i.expect([]tokenType{identifier}, "expected subject at start of statement") {
		return
	}
	if i.expect([]tokenType{colon}, "expected colon after subject") {
		return
	}
	id := i.getIndex(subject.text)
	if i.defCheck[id] { // If map is already set to true
		i.errAt(subject, "multiple definitions for %s", subject.text)
	}

	i.defCheck[id] = true
	startnode := i.graph.GetNode(uint32(start), start)
	startnode.AddEdgeNext(&i.graph, i.graph.GetNode(id, header), 1)

	// Every rule gets its own end node. Sharing one global end node made the
	// '^' operator jump between unrelated rules.
	endNode := i.graph.GetNode(i.getFuncPtr(), end)
	i.graph.ruleEnds[id] = endNode.id

	i.parseRules(id, false, endNode)
}

func (i *parser) parseRules(root uint32, isDeep bool, endNode *syntaxNode) (*syntaxNode, *syntaxNode) {
	rootnode := i.graph.GetNode(root, idk)
	bufferNode := rootnode
	if endNode == nil {
		typ := end
		if isDeep {
			// Parenthesised groups are inline constructs, not rule calls.
			// Using the rule end node type here made the walk return to the
			// caller and silently skip tokens that followed the group.
			typ = groupEnd
		}
		endNode = i.graph.GetNode(i.getFuncPtr(), typ)
	}
	var startBuffer *syntaxNode
	for i.index < len(i.tokens) {
		switch i.curr().typ {
		case identifier:
			// Means its a reference to a different Subject(presumably)
			pointerid := i.getIndex(i.tokens[i.index].text)
			node := i.graph.GetNode(i.getFuncPtr(), pointer)
			node.pointer = pointerid
			bufferNode.AddEdgeNext(&i.graph, node, i.getProbability())
			jumpNode := i.graph.GetNode(i.getFuncPtr(), jump)
			node.AddEdgeNext(&i.graph, jumpNode, 1)
			startBuffer = bufferNode
			bufferNode = jumpNode
			// Basically just add the word and next to it its jump node
			// So when generating, the control will pass to the node at the location and save the exit in a stack
			// Then when it reached its local collapse node, then the control will automatically come back to default
		case character, regex:
			index := i.getPrintPtr()
			i.charmap[index] = i.tokens[i.index].text
			var leafnode *syntaxNode
			if i.tokens[i.index].typ == character {
				leafnode = i.graph.GetNode(index, ch)
			} else {
				leafnode = i.graph.GetNode(index, rx)
				i.regexhandler.CacheRegex(i.curr().text)
			}
			bufferNode.AddEdgeNext(&i.graph, leafnode, i.getProbability())
			jumpNode := i.graph.GetNode(i.getFuncPtr(), jump)
			leafnode.AddEdgeNext(&i.graph, jumpNode, 1)
			startBuffer = bufferNode
			bufferNode = jumpNode
		case colon:
			// Colon is not allowed here
			i.errHere("missing semicolon before ':'")
			return nil, nil
		case maybe:
			if startBuffer == nil {
				i.errHere("'?' has nothing to repeat")
				return nil, nil
			}
			startBuffer.AddEdgeNext(&i.graph, bufferNode, 1-i.getProbability()) // An option to skip to the end
		case oneormore:
			if startBuffer == nil {
				i.errHere("'+' has nothing to repeat")
				return nil, nil
			}
			bufferNode.AddEdgeNext(&i.graph, startBuffer, i.getProbability()) // An option to go to the start
		case anyno:
			if startBuffer == nil {
				i.errHere("'*' has nothing to repeat")
				return nil, nil
			}
			startBuffer.AddEdgeNext(&i.graph, bufferNode, 1-i.getProbability()) // Well both of them combined
			bufferNode.AddEdgeNext(&i.graph, startBuffer, i.getProbability())
		case option:
			// In case of an option, no need to really do anything, just set the buffer settings back to the parent
			bufferNode.AddEdgeNext(&i.graph, endNode, i.getProbability())
			bufferNode = rootnode
			startBuffer = nil
		case padding:
			bufferNode.AddEdgeNext(&i.graph, endNode, 1)
			if isDeep {
				i.errHere("unexpected ';' inside group")
				i.index++
				return nil, nil
			}
			i.index++
			return nil, nil // End of this statement
		case bracopen:
			i.index++
			startBuffer, bufferNode = i.parseRules(bufferNode.id, true, nil)
			if startBuffer == nil || bufferNode == nil {
				return nil, nil
			}
		case bracclose:
			if isDeep {
				bufferNode.AddEdgeNext(&i.graph, endNode, 1)
				return rootnode, endNode
			}
			i.errHere("stray ')' found")
			i.index++
			return nil, nil
		case infinite:
			if startBuffer == nil {
				i.errHere("'^' has nothing to loop")
				return nil, nil
			}
			// '^' repeats the preceding construct in place. The explicit loop
			// node lets complete generation stop at the budget and then emit
			// whatever follows the loop exactly once.
			loopNode := i.graph.GetNode(i.getFuncPtr(), loop)
			loopNode.pointer = startBuffer.id
			bufferNode.AddEdgeNext(&i.graph, loopNode, 1)
			bufferNode = loopNode
			startBuffer = nil
		case probability:
			i.errHere("probability '<...>' must follow a token")
			return nil, nil
		default:
			i.errHere("unexpected token %q", i.curr().text)
			i.index++
			return nil, nil
		}
		i.index++
	}
	if isDeep {
		i.errors = append(i.errors, fmt.Errorf("unterminated '(' in grammar"))
	} else {
		i.errors = append(i.errors, fmt.Errorf("unexpected end of grammar: missing ';'"))
	}
	return nil, nil
}
func (i *parser) getProbability() float32 {
	i.index++
	if i.index >= len(i.tokens) {
		i.index-- // Reverting
		return 0.5
	}
	if i.tokens[i.index].typ == probability {
		num := i.tokens[i.index].text
		numf, err := strconv.ParseFloat(num, 32)
		if err != nil {
			i.errAt(i.tokens[i.index], "invalid probability %q: %v", num, err)
			i.index--
			return 0
		}
		if numf < 0 {
			i.errAt(i.tokens[i.index], "negative probability found")
			i.index--
			return 0
		}
		if math.IsNaN(numf) || math.IsInf(numf, 0) {
			i.errAt(i.tokens[i.index], "invalid probability %q", num)
			i.index--
			return 0
		}
		return float32(numf)
	}
	i.index-- // Reverting
	return 0.5
}

// attachWords resolves @words rule names to graph ids.
func (i *parser) attachWords() {
	for name, words := range i.words {
		id, ok := i.nameMap[name]
		if !ok {
			// A word-only rule does not need a grammar definition. Give it an
			// id so it can be generated like any other rule.
			id = i.getIndex(name)
			i.defCheck[id] = true
		}
		i.graph.words[id] = append([]string(nil), words...)
	}
}

func (i *parser) ValidateGraph() []error {
	if len(i.errors) > 0 {
		return nil
	}
	var errors []error
	for key, val := range i.defCheck {
		if !val && len(i.graph.words[key]) == 0 {
			errors = append(errors, fmt.Errorf("definition of '%s' not found", i.revNameMap[key]))
		}
	}
	return errors
}
