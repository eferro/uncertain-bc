package cli

import (
	"fmt"
	"strconv"
	"strings"
)

type command func(r replSession, arg string) (quit bool, err error)

var commands = map[string]command{
	"help":    showHelp,
	"vars":    showVariables,
	"seed":    reseed,
	"samples": resample,
	"quit":    func(replSession, string) (bool, error) { return true, nil },
}

const replHelp = `Commands:
  :help          show this help
  :vars          list variables with their summaries
  :seed N        restart the random stream with seed N
  :samples N     change the sample count (clears variables)
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
