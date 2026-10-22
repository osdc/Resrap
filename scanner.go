package resrap

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// scanner will take all the components and convert them to a series of tokens

type scanner struct {
	Input  string
	pos    int // current byte offset
	width  int // width of last rune
	currR  rune
	line   int // line of the *next* rune
	col    int // column of the *next* rune
	currLn int // line of the current rune
	currCl int // column of the current rune
	tokens []token
}

func extracttokens(inp string) ([]token, []ScanError) {
	sc := scanner{Input: inp, line: 1, col: 1}
	return sc.scan()
}

// ScanError describes a scanning failure: the location where it occurred and a
// human-readable message.
type ScanError struct {
	Pos  int
	Line int
	Col  int
	Msg  string
}

func (e ScanError) String() string {
	return fmt.Sprintf("line %d:%d: %s", e.Line, e.Col, e.Msg)
}

// next advances to and returns the next rune.
func (s *scanner) next() rune {
	if s.pos >= len(s.Input) {
		s.width = 0
		s.currR = -1 // EOF
		return -1
	}
	r, w := utf8.DecodeRuneInString(s.Input[s.pos:])
	s.width = w
	s.pos += w
	s.currLn = s.line
	s.currCl = s.col
	if r == '\n' {
		s.line++
		s.col = 1
	} else {
		s.col++
	}
	s.currR = r
	return r
}

// peek looks at the next rune without consuming it.
func (s *scanner) peek() rune {
	if s.pos >= len(s.Input) {
		return -1
	}
	r, _ := utf8.DecodeRuneInString(s.Input[s.pos:])
	return r
}

func (s *scanner) add(typ tokenType, text string) {
	s.tokens = append(s.tokens, token{
		typ:  typ,
		text: text,
		line: s.currLn,
		col:  s.currCl,
	})
}

// scanDelimited reads until closeRune, honouring backslash escapes so literals
// such as 'it\'s' and regex classes such as [a\\-z] keep their payload intact.
// The escape sequences themselves are decoded later by unescapeString / ExpandClass.
func (s *scanner) scanDelimited(open, closeRune rune, _ bool) (string, *ScanError) {
	startPos := s.pos - s.width
	startLine, startCol := s.currLn, s.currCl

	var buf strings.Builder
	for {
		r := s.next()
		if r == -1 {
			return "", &ScanError{
				Pos:  startPos,
				Line: startLine,
				Col:  startCol,
				Msg:  fmt.Sprintf("unterminated %q", open),
			}
		}
		if r == closeRune {
			return buf.String(), nil
		}
		if r == '\\' {
			nr := s.next()
			if nr == -1 {
				return "", &ScanError{
					Pos:  startPos,
					Line: startLine,
					Col:  startCol,
					Msg:  fmt.Sprintf("unterminated %q", open),
				}
			}
			buf.WriteRune('\\')
			buf.WriteRune(nr)
			continue
		}
		buf.WriteRune(r)
	}
}

// skipLine consumes the rest of the current line (used for # and // comments).
func (s *scanner) skipLine() {
	for {
		r := s.peek()
		if r == -1 || r == '\n' {
			return
		}
		s.next()
	}
}

func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z')
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func isIdentStart(r rune) bool {
	return isAlpha(r) || r == '_'
}

func isIdentPart(r rune) bool {
	return isAlpha(r) || isDigit(r) || r == '_'
}

func (s *scanner) scanIdentifier() string {
	var buf strings.Builder
	buf.WriteRune(s.currR)

	// Keep reading while next runes are valid identifier parts.
	for {
		r := s.peek()
		if !isIdentPart(r) {
			break
		}
		s.next()
		buf.WriteRune(s.currR)
	}
	return buf.String()
}

func (s *scanner) scan() ([]token, []ScanError) {
	var errs []ScanError
	for s.next() != -1 {
		switch s.currR {
		case ' ', '\t', '\r', '\n':
			// whitespace
		case '#':
			s.skipLine()
		case '/':
			if s.peek() == '/' {
				s.skipLine()
			} else {
				errs = append(errs, ScanError{s.pos - s.width, s.currLn, s.currCl, fmt.Sprintf("unexpected character %q", s.currR)})
			}
		case '+':
			s.add(oneormore, "")
		case '*':
			s.add(anyno, "")
		case '^':
			s.add(infinite, "")
		case '?':
			s.add(maybe, "")
		case '|':
			s.add(option, "")
		case ';':
			s.add(padding, "")
		case '(':
			s.add(bracopen, "")
		case ')':
			s.add(bracclose, "")
		case ':':
			s.add(colon, "")
		case '\'':
			val, err := s.scanDelimited('\'', '\'', false)
			if err != nil {
				errs = append(errs, *err)
			} else {
				s.add(character, val)
			}

		case '<':
			val, err := s.scanDelimited('<', '>', false)
			if err != nil {
				errs = append(errs, *err)
			} else {
				s.add(probability, val)
			}

		case '[':
			val, err := s.scanDelimited('[', ']', false)
			if err != nil {
				errs = append(errs, *err)
			} else {
				s.add(regex, val)
			}

		default:
			if isIdentStart(s.currR) {
				buff := s.scanIdentifier()
				if buff != "" {
					s.add(identifier, buff)
				}
			} else {
				errs = append(errs, ScanError{s.pos - s.width, s.currLn, s.currCl, fmt.Sprintf("unexpected character %q", s.currR)})
			}
		}
	}
	return s.tokens, errs
}

func (s *scanner) PrintTokens() {
	for _, token := range s.tokens {
		fmt.Println(token.text)
	}
}
