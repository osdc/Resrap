package resrap

import (
	"bufio"
	"os"
	"strings"
)

// ParseFile reads a file and returns statements split by lines ending with ';',
// skipping lines starting with "//".
func parseFile(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()

	var statements []string
	var current strings.Builder

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// skip comment or empty lines
		if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		// accumulate this line
		current.WriteString(line + " ")

		// if line ends with ';', finalize statement
		if strings.HasSuffix(line, ";") {
			stmt := strings.TrimSpace(current.String())
			statements = append(statements, stmt)
			current.Reset()
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// A final statement does not have to end with ';' at the file level; the
	// grammar parser will report a useful error if it is incomplete. Dropping
	// it silently made an entire trailing rule disappear.
	if strings.TrimSpace(current.String()) != "" {
		statements = append(statements, strings.TrimSpace(current.String()))
	}

	return statements, nil
}

type lang struct {
	graph *syntaxGraph
}

func newLang() lang {
	return lang{}
}
func (l *lang) GetGraph() *syntaxGraph {
	return l.graph
}
func (l *lang) ParserFile(filename string) error {
	// Read the file as a whole instead of reassembling statements line by
	// line. The scanner understands comments and @words directives, and this
	// avoids losing a trailing rule that has no final semicolon.
	data, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return l.ParserString(string(data))
}

func (l *lang) ParserString(data string) error {
	gb := newGraphBuilder()
	err := gb.startGeneration(data)
	l.graph = &gb.pars.graph
	return err
}
