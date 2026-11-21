package cli

import (
	"errors"
	"fmt"
	"image/png"
	"os"
	"path/filepath"
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
	"plot":    plot,
	"quit":    func(replSession, string) (bool, error) { return true, nil },
}

const replHelp = `Commands:
  :help          show this help
  :vars          list variables with their summaries
  :hist EXPR     draw the distribution of EXPR
  :plot EXPR F   write the distribution of EXPR to F (.svg or .png)
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
  triangular(lo, mode, hi)  beta(a, b)  pert(lo, mode, hi)
  exponential(rate)  gamma(shape, scale)
  mean(x) median(x) stdev(x) min(x) max(x) p(x, q)
  sqrt exp log log10 abs round floor ceil, min/max of several values
  x = expr   f(a, b) = expr   last   pi   e   # comments
  [3, 5, 8, 13]                       empirical data, resampled
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
	samples, err := r.samplesOf(expr, ":hist")
	if err != nil {
		return false, err
	}
	fmt.Fprint(r.stdout, histogramChart(samples, r.style))
	return false, nil
}

func plot(r replSession, arg string) (bool, error) {
	cut := strings.LastIndex(arg, " ")
	if cut < 0 {
		return false, errors.New(":plot expects an expression and a file name")
	}
	expr, path := strings.TrimSpace(arg[:cut]), arg[cut+1:]
	write, ok := plotWriters[filepath.Ext(path)]
	if !ok {
		return false, fmt.Errorf(":plot writes .svg or .png files, got '%s'", path)
	}
	samples, err := r.samplesOf(expr, ":plot")
	if err != nil {
		return false, err
	}
	return false, write(path, samples, expr)
}

var plotWriters = map[string]func(path string, s domain.SampleSet, title string) error{
	".svg": func(path string, s domain.SampleSet, title string) error {
		return os.WriteFile(path, []byte(svgHistogram(s, title)), 0o644)
	},
	".png": writePNG,
}

func writePNG(path string, s domain.SampleSet, title string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	if err := png.Encode(file, pngHistogram(s, title)); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}

func (r replSession) samplesOf(expr, command string) (domain.SampleSet, error) {
	results, err := r.session.Evaluate(expr)
	if err != nil {
		return nil, err
	}
	samples, ok := results[len(results)-1].Value.(domain.SampleSet)
	if !ok {
		return nil, fmt.Errorf("%s needs an uncertain value, got a plain number", command)
	}
	return samples, nil
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
	printResults(results, src.prefixLines, r.style, r.stdout)
	if err != nil {
		reportError(src, err, r.stderr)
	}
	return false, nil
}
