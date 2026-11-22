// Package cli is the command-line adapter: it turns arguments, files and
// stdin into evaluations, and results into text.
package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"slices"

	"golang.org/x/term"

	"github.com/eferro/uncertain-bc/internal/application"
	"github.com/eferro/uncertain-bc/internal/domain"
)

const (
	exitOK         = 0
	exitEvalError  = 1
	exitUsageError = 2
	exitParseError = 3
	exitAssertion  = 4
)

type App struct {
	Version     string
	RandomSeed  func() uint64
	Interactive bool
	HistoryPath string
	Wait        func() bool
}

type options struct {
	expressions multiString
	sets        multiString
	data        multiString
	asserts     multiString
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
	percentiles []float64
	interval    string
	width       float64
	presets     []preset
	watch       bool
	completion  string
	samplesPath string
	samplesOut  *samplesOut
	sens        string
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
      --set NAME=EXPR set NAME to EXPR, overriding the model's value (repeatable)
      --data NAME=FILE
                      set NAME to the numbers in FILE, one per line (repeatable)
      --assert EXPR   after the program, exit 4 unless EXPR certainly holds (repeatable)
      --samples-out FILE
                      after the program, write every variable's samples to FILE (.csv or .json)
      --sens EXPR     after the program, rank the variables that drive EXPR's uncertainty
  -n, --samples N     Monte Carlo sample count (default 5000)
  -s, --seed N        random seed (default 1)
      --random        seed from the clock
      --ascii         print ASCII instead of Unicode
      --plain         print only the central value (median), for scripts
      --json          print one JSON object per result
      --echo          also print assignments, with their names
      --hist          draw a histogram under every uncertain result
      --watch         evaluate the files again whenever they change
      --interval W    print a W% interval instead of the 90% one (default 90)
      --quantiles L   show these percentiles instead of the intervals (10,50,90)
  -h, --help          show this help
      --version       print the version
      --completion S  print a completion script for bash, zsh or fish
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
	if opts.completion != "" {
		return printCompletion(opts.completion, stdout, stderr)
	}
	opts, err := resolveFlags(opts)
	if err != nil {
		fmt.Fprintf(stderr, "ubc: %v\n", err)
		return exitUsageError
	}
	return a.execute(opts, a.settings(opts), stdin, stdout, stderr)
}

func (a App) execute(opts options, settings application.Settings, stdin io.Reader, stdout, stderr io.Writer) int {
	if a.Interactive && len(opts.expressions) == 0 && len(opts.files) == 0 {
		if name := programOnlyFlag(opts); name != "" {
			fmt.Fprintf(stderr, "ubc: %s needs a program: -e, files or stdin\n", name)
			return exitUsageError
		}
		session, code := startSession(opts, settings, stderr)
		if session == nil {
			return code
		}
		return a.interactive(session, styleFor(opts), stdin, stdout, stderr)
	}
	if opts.watch {
		return a.watch(opts, settings, styleFor(opts), stdout, stderr)
	}
	return evaluateSources(opts, settings, styleFor(opts), stdin, stdout, stderr)
}

func evaluateSources(opts options, settings application.Settings, st style, stdin io.Reader, stdout, stderr io.Writer) int {
	sources, err := collectSources(opts, stdin)
	if err != nil {
		fmt.Fprintf(stderr, "ubc: %v\n", err)
		return exitEvalError
	}
	session, code := startSession(opts, settings, stderr)
	if session == nil {
		return code
	}
	if code := evaluateAll(session, sources, st, stdout, stderr); code != exitOK {
		return code
	}
	if code := exportSamples(session, opts.samplesOut, opts.samples, stderr); code != exitOK {
		return code
	}
	if code := printSensitivity(session, opts.sens, st, stdout, stderr); code != exitOK {
		return code
	}
	return checkAssertions(session, opts.asserts, st, stderr)
}

var flagRules = []struct {
	broken  func(o options) bool
	message string
}{
	{func(o options) bool { return o.json && o.plain }, "--json and --plain cannot be combined"},
	{func(o options) bool { return o.hist && (o.json || o.plain) }, "--hist cannot be combined with --json or --plain"},
	{func(o options) bool { return o.sens != "" && (o.json || o.plain) }, "--sens cannot be combined with --json or --plain"},
	{func(o options) bool { return o.interval != "" && o.quantiles != "" }, "--interval and --quantiles cannot be combined"},
	{func(o options) bool { return o.random && o.seedGiven }, "--seed and --random cannot be combined"},
	{func(o options) bool { return o.watch && len(o.files) == 0 }, "--watch needs at least one file"},
}

// resolveFlags checks the flags and parses the values that need it once, so
// the rest of the adapter works with checked values.
func resolveFlags(opts options) (options, error) {
	if opts.samples <= 0 {
		return opts, fmt.Errorf("--samples must be positive, got %d", opts.samples)
	}
	for _, rule := range flagRules {
		if rule.broken(opts) {
			return opts, errors.New(rule.message)
		}
	}
	return parseFlagValues(opts)
}

func parseFlagValues(opts options) (options, error) {
	var err error
	if opts.presets, err = parsePresets(opts); err != nil {
		return opts, err
	}
	if opts.samplesOut, err = parseSamplesOut(opts.samplesPath); err != nil {
		return opts, err
	}
	if opts.percentiles, err = parseQuantiles(opts.quantiles); err != nil {
		return opts, err
	}
	opts.width, err = parseInterval(opts.interval)
	return opts, err
}

func parsePresets(opts options) ([]preset, error) {
	overrides, err := parseOverrides(opts.sets)
	if err != nil {
		return nil, err
	}
	datasets, err := parseDatasets(opts.data)
	if err != nil {
		return nil, err
	}
	return slices.Concat(datasets, overrides), nil
}

// programOnlyFlag names a flag that acts after a program, which the REPL
// does not have.
func programOnlyFlag(opts options) string {
	switch {
	case len(opts.asserts) > 0:
		return "--assert"
	case opts.samplesOut != nil:
		return "--samples-out"
	case opts.sens != "":
		return "--sens"
	}
	return ""
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
	return run(newPromptedScanner(stdin, stdout), stdout, stderr)
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
	registerInputFlags(fs, opts)
	registerOutputFlags(fs, opts)
}

func registerInputFlags(fs *flag.FlagSet, opts *options) {
	for _, name := range []string{"e", "eval"} {
		fs.Var(&opts.expressions, name, "evaluate EXPR (repeatable)")
	}
	fs.Var(&opts.sets, "set", "set NAME=EXPR, overriding the model's value (repeatable)")
	fs.Var(&opts.data, "data", "set NAME to the numbers in FILE, one per line (repeatable)")
	fs.Var(&opts.asserts, "assert", "exit 4 unless EXPR holds after the program (repeatable)")
	fs.StringVar(&opts.samplesPath, "samples-out", "", "write every variable's samples to FILE (.csv or .json)")
	fs.StringVar(&opts.sens, "sens", "", "rank the variables that drive the uncertainty of EXPR")
	for _, name := range []string{"n", "samples"} {
		fs.IntVar(&opts.samples, name, opts.samples, "Monte Carlo sample count")
	}
	for _, name := range []string{"s", "seed"} {
		fs.Uint64Var(&opts.seed, name, opts.seed, "random seed")
	}
	fs.BoolVar(&opts.random, "random", false, "seed from the clock")
	fs.BoolVar(&opts.watch, "watch", false, "evaluate the files again whenever they change")
	fs.BoolVar(&opts.version, "version", false, "print the version and exit")
	fs.StringVar(&opts.completion, "completion", "", "print a completion script for bash, zsh or fish")
}

func registerOutputFlags(fs *flag.FlagSet, opts *options) {
	fs.BoolVar(&opts.ascii, "ascii", false, "print ASCII instead of Unicode")
	fs.BoolVar(&opts.plain, "plain", false, "print only the central value")
	fs.BoolVar(&opts.json, "json", false, "print one JSON object per result")
	fs.BoolVar(&opts.echo, "echo", false, "also print assignments, with their names")
	fs.BoolVar(&opts.hist, "hist", false, "draw a histogram under every uncertain result")
	fs.StringVar(&opts.interval, "interval", "", "width of the printed interval, a percentage (default 90)")
	fs.StringVar(&opts.quantiles, "quantiles", "", "comma-separated percentiles to show instead of the 90% and 98% intervals")
}

func (a App) settings(opts options) application.Settings {
	seed := opts.seed
	if opts.random {
		seed = a.RandomSeed()
	}
	return application.Settings{Seed: seed, Samples: opts.samples}
}

type source struct {
	name    string
	program string
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
		sources = append(sources, source{name: path, program: string(content)})
	}
	if len(sources) > 0 {
		return sources, nil
	}
	content, err := io.ReadAll(stdin)
	if err != nil {
		return nil, fmt.Errorf("stdin: %w", err)
	}
	return []source{{program: string(content)}}, nil
}

func evaluateAll(session *application.Session, sources []source, st style, stdout, stderr io.Writer) int {
	for _, src := range sources {
		results, err := session.Evaluate(src.program)
		printResults(results, st, stdout)
		if err != nil {
			reportError(src, err, stderr)
			return exitCode(err)
		}
	}
	return exitOK
}

// printResults prints expressions, and assignments too with --echo. In a
// program that prints several expressions each one says what it is (ADR 0053).
func printResults(results []domain.Result, st style, stdout io.Writer) {
	label := !st.plain && countExpressions(results) > 1
	for _, result := range results {
		if line, shown := resultLine(result, label, st); shown {
			fmt.Fprintln(stdout, line)
		}
	}
}

func resultLine(result domain.Result, label bool, st style) (string, bool) {
	switch {
	case result.Name != "":
		return display(result, st), st.echo
	case !label:
		return display(result, st), true
	}
	if named, ok := labelled(result); ok {
		return display(named, st), true
	}
	if st.json {
		return display(result, st), true
	}
	return fmt.Sprintf("L%d: %s", result.Line, display(result, st)), true
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

func printCompletion(shell string, stdout, stderr io.Writer) int {
	script, ok := completionScript(shell)
	if !ok {
		fmt.Fprintf(stderr, "ubc: --completion expects bash, zsh or fish, got '%s'\n", shell)
		return exitUsageError
	}
	fmt.Fprint(stdout, script)
	return exitOK
}
