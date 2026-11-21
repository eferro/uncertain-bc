// Package cli is the command-line adapter: it turns arguments, files and
// stdin into evaluations, and results into text.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/eferro/uncertain-bc/internal/application"
	"github.com/eferro/uncertain-bc/internal/domain"
)

const (
	exitOK         = 0
	exitEvalError  = 1
	exitUsageError = 2
	exitParseError = 3
)

type App struct {
	Version string
}

type options struct {
	expressions multiString
	files       []string
	version     bool
}

type multiString []string

func (m *multiString) String() string { return fmt.Sprint(*m) }

func (m *multiString) Set(value string) error {
	*m = append(*m, value)
	return nil
}

func (a App) Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	opts, err := parseFlags(args, stderr)
	if err != nil {
		return exitUsageError
	}
	if opts.version {
		fmt.Fprintf(stdout, "ubc %s\n", a.Version)
		return exitOK
	}
	sources, err := collectSources(opts, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "ubc: %v\n", err)
		return exitEvalError
	}
	return evaluateAll(sources, stdout, stderr)
}

func parseFlags(args []string, stderr io.Writer) (options, error) {
	var opts options
	fs := flag.NewFlagSet("ubc", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Var(&opts.expressions, "e", "evaluate EXPR (repeatable)")
	fs.Var(&opts.expressions, "eval", "evaluate EXPR (repeatable)")
	fs.BoolVar(&opts.version, "version", false, "print the version and exit")
	err := fs.Parse(args)
	opts.files = fs.Args()
	return opts, err
}

type source struct {
	name        string
	program     string
	prefixLines bool
}

func collectSources(opts options, stdin io.Reader) ([]source, error) {
	var sources []source
	for _, expr := range opts.expressions {
		sources = append(sources, source{program: expr})
	}
	for _, path := range opts.files {
		content, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source{name: path, program: string(content), prefixLines: true})
	}
	if len(sources) > 0 {
		return sources, nil
	}
	content, err := io.ReadAll(stdin)
	if err != nil {
		return nil, fmt.Errorf("stdin: %w", err)
	}
	return []source{{program: string(content), prefixLines: true}}, nil
}

func evaluateAll(sources []source, stdout, stderr io.Writer) int {
	session := application.NewSession(application.DefaultSettings())
	for _, src := range sources {
		results, err := session.Evaluate(src.program)
		printExpressions(results, src.prefixLines, stdout)
		if err != nil {
			reportError(src, err, stderr)
			return exitCode(err)
		}
	}
	return exitOK
}

func printExpressions(results []domain.Result, prefixLines bool, stdout io.Writer) {
	var expressions []domain.Result
	for _, result := range results {
		if result.Name == "" {
			expressions = append(expressions, result)
		}
	}
	for _, result := range expressions {
		if prefixLines && len(expressions) > 1 {
			fmt.Fprintf(stdout, "L%d: ", result.Line)
		}
		fmt.Fprintln(stdout, format(result.Value))
	}
}

func reportError(src source, err error, stderr io.Writer) {
	if src.name != "" {
		fmt.Fprintf(stderr, "ubc: %s: %v\n", src.name, err)
		return
	}
	fmt.Fprintf(stderr, "ubc: %v\n", err)
}

func exitCode(err error) int {
	var syntaxErr *domain.SyntaxError
	if errors.As(err, &syntaxErr) {
		return exitParseError
	}
	return exitEvalError
}
