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
	Version    string
	RandomSeed func() uint64
}

type options struct {
	expressions multiString
	files       []string
	version     bool
	samples     int
	seed        uint64
	seedGiven   bool
	random      bool
	ascii       bool
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
	settings, err := a.settings(opts)
	if err != nil {
		fmt.Fprintf(stderr, "ubc: %v\n", err)
		return exitUsageError
	}
	sources, err := collectSources(opts, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "ubc: %v\n", err)
		return exitEvalError
	}
	return evaluateAll(sources, settings, styleFor(opts), stdout, stderr)
}

func parseFlags(args []string, stderr io.Writer) (options, error) {
	opts := options{samples: application.DefaultSettings().Samples, seed: application.DefaultSettings().Seed}
	fs := flag.NewFlagSet("ubc", flag.ContinueOnError)
	fs.SetOutput(stderr)
	registerFlags(fs, &opts)
	err := fs.Parse(args)
	opts.files = fs.Args()
	fs.Visit(func(f *flag.Flag) { opts.seedGiven = opts.seedGiven || f.Name == "s" || f.Name == "seed" })
	return opts, err
}

func registerFlags(fs *flag.FlagSet, opts *options) {
	for _, name := range []string{"e", "eval"} {
		fs.Var(&opts.expressions, name, "evaluate EXPR (repeatable)")
	}
	for _, name := range []string{"n", "samples"} {
		fs.IntVar(&opts.samples, name, opts.samples, "Monte Carlo sample count")
	}
	for _, name := range []string{"s", "seed"} {
		fs.Uint64Var(&opts.seed, name, opts.seed, "random seed")
	}
	fs.BoolVar(&opts.random, "random", false, "seed from the clock")
	fs.BoolVar(&opts.ascii, "ascii", false, "print ASCII instead of Unicode")
	fs.BoolVar(&opts.version, "version", false, "print the version and exit")
}

func (a App) settings(opts options) (application.Settings, error) {
	if opts.samples <= 0 {
		return application.Settings{}, fmt.Errorf("--samples must be positive, got %d", opts.samples)
	}
	if opts.random && opts.seedGiven {
		return application.Settings{}, errors.New("--seed and --random cannot be combined")
	}
	seed := opts.seed
	if opts.random {
		seed = a.RandomSeed()
	}
	return application.Settings{Seed: seed, Samples: opts.samples}, nil
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

func evaluateAll(sources []source, settings application.Settings, st style, stdout, stderr io.Writer) int {
	session := application.NewSession(settings)
	for _, src := range sources {
		results, err := session.Evaluate(src.program)
		printExpressions(results, src.prefixLines, st, stdout)
		if err != nil {
			reportError(src, err, stderr)
			return exitCode(err)
		}
	}
	return exitOK
}

func printExpressions(results []domain.Result, prefixLines bool, st style, stdout io.Writer) {
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
		fmt.Fprintln(stdout, format(result.Value, st))
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
