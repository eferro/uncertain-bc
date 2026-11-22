package cli

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"path"
	"strings"
	"unicode"
	"unicode/utf8"
)

// formulas is the shared file every scenario may call into; it is loaded
// before any other scenario (ADR 0043, ADR 0054).
const formulas = "formulas"

const scenarioExt = ".ubc"

const questionPrefix = "# Question:"

type unknownScenarioError struct{ name string }

func (e unknownScenarioError) Error() string {
	return fmt.Sprintf("unknown scenario '%s'; ubc --scenarios lists them", e.name)
}

// catalogue falls back to an empty one, so an App built without scenarios
// (tests, other embedders) reports every name as unknown instead of panicking.
func (a App) catalogue() fs.FS {
	if a.Scenarios == nil {
		return embed.FS{}
	}
	return a.Scenarios
}

func scenarioNames(catalogue fs.FS) []string {
	var names []string
	_ = fs.WalkDir(catalogue, ".", func(file string, entry fs.DirEntry, err error) error {
		if err == nil && !entry.IsDir() && path.Ext(file) == scenarioExt && path.Dir(file) != "." {
			names = append(names, strings.TrimSuffix(file, scenarioExt))
		}
		return nil
	})
	return names
}

func readScenario(catalogue fs.FS, name string) (string, error) {
	content, err := fs.ReadFile(catalogue, name+scenarioExt)
	if err != nil {
		return "", unknownScenarioError{name}
	}
	return string(content), nil
}

// scenarioSources is the program --scenario runs: the shared formulas, then
// the scenario itself.
func scenarioSources(catalogue fs.FS, name string) ([]source, error) {
	if name == "" {
		return nil, nil
	}
	var sources []source
	for _, part := range scenarioParts(name) {
		program, err := readScenario(catalogue, part)
		if err != nil {
			return nil, unknownScenarioError{name}
		}
		sources = append(sources, source{name: part, program: program})
	}
	return sources, nil
}

func scenarioParts(name string) []string {
	if name == formulas {
		return []string{formulas}
	}
	return []string{formulas, name}
}

func listScenarios(catalogue fs.FS, stdout io.Writer) int {
	names := scenarioNames(catalogue)
	width := 0
	for _, name := range names {
		width = max(width, len(name))
	}
	for _, name := range names {
		program, _ := readScenario(catalogue, name)
		line := fmt.Sprintf("%-*s   %s", width, name, question(program))
		fmt.Fprintln(stdout, strings.TrimRight(line, " "))
	}
	return exitOK
}

func printScenarioSource(catalogue fs.FS, name string, stdout, stderr io.Writer) int {
	program, err := readScenario(catalogue, name)
	if err != nil {
		fmt.Fprintf(stderr, "ubc: %v\n", err)
		return exitUsageError
	}
	fmt.Fprint(stdout, program)
	return exitOK
}

// question reads the "# Question:" paragraph of the scenario template
// (ADR 0042): the comment lines from that prefix up to the first empty one.
func question(program string) string {
	_, rest, found := strings.Cut(program, questionPrefix)
	if !found {
		return ""
	}
	first, rest, _ := strings.Cut(rest, "\n")
	words := strings.Fields(first)
	for _, line := range strings.Split(rest, "\n") {
		text, comment := strings.CutPrefix(line, "#")
		if !comment || strings.TrimSpace(text) == "" {
			break
		}
		words = append(words, strings.Fields(text)...)
	}
	return capitalise(strings.Join(words, " "))
}

func capitalise(text string) string {
	first, size := utf8.DecodeRuneInString(text)
	if size == 0 {
		return text
	}
	return string(unicode.ToUpper(first)) + text[size:]
}
