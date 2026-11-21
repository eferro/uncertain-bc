package cli

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/eferro/uncertain-bc/internal/domain"
)

type command func(r replSession, arg string) (quit bool, err error)

var commands = map[string]command{
	"help":    showHelp,
	"vars":    showVariables,
	"seed":    reseed,
	"samples": resample,
	"hist":    showHistogram,
	"save":    save,
	"load":    load,
	"quit":    func(replSession, string) (bool, error) { return true, nil },
}

const replHelp = `Commands:
  :help          show this help
  :vars          list variables with their summaries
  :hist EXPR     draw the distribution of EXPR
  :seed N        restart the random stream with seed N
  :samples N     change the sample count (clears variables)
  :save FILE     write the statements that worked to FILE
  :load FILE     evaluate FILE in this session
  :quit          leave the session (or Ctrl-D)

Language:
  1.5e3  1k 2.5M 3B 1T  5%          numbers and suffixes
  + - * / ^  ( )                      arithmetic
  < > <= >= == !=  P(cond)            comparisons and probabilities
  a to b                              90% interval, lognormal when both ends > 0
  normal(mean, sd)  lognormal(mu, sigma)  uniform(lo, hi)
  triangular(lo, mode, hi)  beta(a, b)
  mean(x) median(x) stdev(x) min(x) max(x) p(x, q)
  sqrt exp log log10 abs round floor ceil, min/max of several values
  x = expr   last   pi   e   # comments
`

func (r replSession) command(line string) (bool, error) {
	name, arg, _ := strings.Cut(strings.TrimSpace(line), " ")
	run, ok := commands[name]
	if !ok {
		return false, fmt.Errorf("unknown command ':%s' (try :help)", name)
	}
	return run(r, strings.TrimSpace(arg))
}

func showHelp(r replSession, _ string) (bool, error) {
	fmt.Fprint(r.stdout, replHelp)
	return false, nil
}

func showVariables(r replSession, _ string) (bool, error) {
	variables := r.session.Variables()
	if len(variables) == 0 {
		fmt.Fprintln(r.stdout, "(no variables)")
	}
	for _, v := range variables {
		r.printVariable(v)
	}
	return false, nil
}

func reseed(r replSession, arg string) (bool, error) {
	seed, err := strconv.ParseUint(arg, 10, 64)
	if err != nil {
		return false, fmt.Errorf(":seed expects a whole number, got '%s'", arg)
	}
	r.session.Reseed(seed)
	return false, nil
}

func resample(r replSession, arg string) (bool, error) {
	samples, err := strconv.Atoi(arg)
	if err != nil || samples <= 0 {
		return false, fmt.Errorf(":samples expects a positive whole number, got '%s'", arg)
	}
	if len(r.session.Variables()) > 0 {
		fmt.Fprintln(r.stderr, "ubc: variables cleared (sample count changed)")
	}
	r.session.SetSamples(samples)
	return false, nil
}

func showHistogram(r replSession, expr string) (bool, error) {
	if expr == "" {
		return false, errors.New(":hist expects an expression")
	}
	results, err := r.session.Evaluate(expr)
	if err != nil {
		return false, err
	}
	samples, ok := results[len(results)-1].Value.(domain.SampleSet)
	if !ok {
		return false, errors.New(":hist needs an uncertain value, got a plain number")
	}
	fmt.Fprint(r.stdout, histogramChart(samples, r.style))
	return false, nil
}

func save(r replSession, path string) (bool, error) {
	if path == "" {
		return false, errors.New(":save expects a file name")
	}
	var content strings.Builder
	for _, program := range r.session.Transcript() {
		content.WriteString(program + "\n")
	}
	return false, os.WriteFile(path, []byte(content.String()), 0o644)
}

func load(r replSession, path string) (bool, error) {
	if path == "" {
		return false, errors.New(":load expects a file name")
	}
	content, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	src := source{name: path, program: string(content), prefixLines: true}
	results, err := r.session.Evaluate(src.program)
	printExpressions(results, src.prefixLines, r.style, r.stdout)
	if err != nil {
		reportError(src, err, r.stderr)
	}
	return false, nil
}
