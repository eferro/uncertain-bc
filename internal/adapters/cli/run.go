// Package cli is the command-line adapter: it turns arguments, files and
// stdin into evaluations, and results into text.
package cli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"golang.org/x/term"

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
	Version     string
	RandomSeed  func() uint64
	Interactive bool
	HistoryPath string
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
	plain       bool
	json        bool
	echo        bool
	hist        bool
	quantiles   string
}

type multiString []string

func (m *multiString) String() string { return fmt.Sprint(*m) }

func (m *multiString) Set(value string) error {
	*m = append(*m, value)
	return nil
}

const usage = `Usage:
  ubc [flags] [FILE...]

Evaluates expressions with uncertain numbers. Without -e or files, it reads
the program from stdin. Flags go before files.

Flags:
  -e, --eval EXPR     evaluate EXPR (repeatable)
  -n, --samples N     Monte Carlo sample count (default 5000)
  -s, --seed N        random seed (default 1)
      --random        seed from the clock
      --ascii         print ASCII instead of Unicode
      --plain         print only the central value (median), for scripts
      --json          print one JSON object per result
      --echo          also print assignments, with their names
      --hist          draw a histogram under every uncertain result
      --quantiles L   show these percentiles instead of the 90% interval (10,50,90)
  -h, --help          show this help
      --version       print the version
`

func (a App) Run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	opts, code, proceed := readFlags(args, stdout, stderr)
	if !proceed {
		return code
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
	return a.execute(opts, settings, stdin, stdout, stderr)
}

func (a App) execute(opts options, settings application.Settings, stdin io.Reader, stdout, stderr io.Writer) int {
	if a.Interactive && len(opts.expressions) == 0 && len(opts.files) == 0 {
		return a.interactive(application.NewSession(settings), styleFor(opts), stdin, stdout, stderr)
	}
	sources, err := collectSources(opts, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "ubc: %v\n", err)
		return exitEvalError
	}
	return evaluateAll(sources, settings, styleFor(opts), stdout, stderr)
}

func validateFlags(opts options) error {
	switch {
	case opts.samples <= 0:
		return fmt.Errorf("--samples must be positive, got %d", opts.samples)
	case opts.json && opts.plain:
		return errors.New("--json and --plain cannot be combined")
	case opts.hist && (opts.json || opts.plain):
		return errors.New("--hist cannot be combined with --json or --plain")
	case opts.random && opts.seedGiven:
		return errors.New("--seed and --random cannot be combined")
	}
	_, err := parseQuantiles(opts.quantiles)
	return err
}

func readFlags(args []string, stdout, stderr io.Writer) (options, int, bool) {
	opts, err := parseFlags(args, stderr)
	switch {
	case errors.Is(err, flag.ErrHelp):
		fmt.Fprint(stdout, usage)
		return opts, exitOK, false
	case err != nil:
		fmt.Fprint(stderr, usage)
		return opts, exitUsageError, false
	}
	return opts, exitOK, true
}

func (a App) interactive(session *application.Session, st style, stdin io.Reader, stdout, stderr io.Writer) int {
	run := func(lines lineReader, out, errOut io.Writer) int { return repl(session, st, lines, out, errOut) }
	if f, ok := stdin.(*os.File); ok && term.IsTerminal(int(f.Fd())) {
		return a.terminalRepl(f, stdout, stderr, run)
	}
	return run(promptedScanner{lines: bufio.NewScanner(stdin), stdout: stdout}, stdout, stderr)
}

func parseFlags(args []string, stderr io.Writer) (options, error) {
	opts := options{samples: application.DefaultSettings().Samples, seed: application.DefaultSettings().Seed}
	fs := flag.NewFlagSet("ubc", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
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
	fs.BoolVar(&opts.plain, "plain", false, "print only the central value")
	fs.BoolVar(&opts.json, "json", false, "print one JSON object per result")
	fs.BoolVar(&opts.echo, "echo", false, "also print assignments, with their names")
	fs.BoolVar(&opts.hist, "hist", false, "draw a histogram under every uncertain result")
	fs.StringVar(&opts.quantiles, "quantiles", "", "comma-separated percentiles to show instead of the 90% interval")
	fs.BoolVar(&opts.version, "version", false, "print the version and exit")
}

func (a App) settings(opts options) (application.Settings, error) {
	if err := validateFlags(opts); err != nil {
		return application.Settings{}, err
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
		printResults(results, src.prefixLines, st, stdout)
		if err != nil {
			reportError(src, err, stderr)
			return exitCode(err)
		}
	}
	return exitOK
}

// printResults prints expressions, and assignments too with --echo. In a
// program that prints several expressions each one says its line.
func printResults(results []domain.Result, prefixLines bool, st style, stdout io.Writer) {
	prefix := prefixLines && !st.plain && !st.json && countExpressions(results) > 1
	for _, result := range results {
		if line, shown := resultLine(result, prefix, st); shown {
			fmt.Fprintln(stdout, line)
		}
	}
}

func resultLine(result domain.Result, prefix bool, st style) (string, bool) {
	switch {
	case result.Name != "":
		return display(result, st), st.echo
	case prefix:
		return fmt.Sprintf("L%d: %s", result.Line, display(result, st)), true
	}
	return display(result, st), true
}

func countExpressions(results []domain.Result) int {
	expressions := 0
	for _, result := range results {
		if result.Name == "" {
			expressions++
		}
	}
	return expressions
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
