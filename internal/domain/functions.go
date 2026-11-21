package domain

import (
	"fmt"
	"math"
)

type call struct {
	name string
	args []node
	at   position
}

type builtin func(c call, env Env, args []Value) (Value, error)

var functions = map[string]builtin{
	"P":   probabilityOf,
	"p":   percentile,
	"min": extreme(math.Min, 0),
	"max": extreme(math.Max, 1),
}

func lookup(name string) (builtin, bool) {
	if fn, ok := functions[name]; ok {
		return fn, true
	}
	if s, ok := summaries[name]; ok {
		return s.call, true
	}
	if f, ok := mathFunctions[name]; ok {
		return f.call, true
	}
	if d, ok := distributions[name]; ok {
		return d.call, true
	}
	return nil, false
}

func (c call) evaluate(env Env) (Value, error) {
	fn, ok := lookup(c.name)
	if !ok {
		return nil, evalErrorAt(c.at, "unknown function '%s'", c.name)
	}
	args := make([]Value, len(c.args))
	for i, arg := range c.args {
		value, err := arg.evaluate(env)
		if err != nil {
			return nil, err
		}
		args[i] = value
	}
	return fn(c, env, args)
}

func (c call) expectArgs(args []Value, want int) error {
	if len(args) != want {
		return evalErrorAt(c.at, "%s expects %s, got %d", c.name, plural(want, "argument"), len(args))
	}
	return nil
}

func plural(n int, noun string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", noun)
	}
	return fmt.Sprintf("%d %ss", n, noun)
}

func (c call) plainNumbers(args []Value, want int) ([]float64, error) {
	if err := c.expectArgs(args, want); err != nil {
		return nil, err
	}
	numbers := make([]float64, len(args))
	for i, arg := range args {
		x, ok := scalarOf(arg)
		if !ok {
			return nil, evalErrorAt(c.at, "%s expects plain numbers", c.name)
		}
		numbers[i] = x
	}
	return numbers, nil
}

func probabilityOf(c call, _ Env, args []Value) (Value, error) {
	if err := c.expectArgs(args, 1); err != nil {
		return nil, err
	}
	if !isIndicator(args[0]) {
		return nil, evalErrorAt(c.at, "P expects a comparison")
	}
	return Probability(mean(args[0])), nil
}
