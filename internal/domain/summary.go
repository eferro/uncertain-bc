package domain

import "math"

type summary func(v Value) float64

var summaries = map[string]summary{
	"mean":   mean,
	"median": func(v Value) float64 { return quantile(v, 0.5) },
	"stdev":  stdev,
}

func (s summary) call(c call, _ Env, args []Value) (Value, error) {
	if err := c.expectArgs(args, 1); err != nil {
		return nil, err
	}
	return Scalar(s(args[0])), nil
}

func quantile(v Value, q float64) float64 {
	if s, ok := Samples(v); ok {
		return s.Quantile(q)
	}
	return at(v, 0)
}

func stdev(v Value) float64 {
	m := mean(v)
	values := finiteValues(v)
	squares := 0.0
	for _, x := range values {
		squares += (x - m) * (x - m)
	}
	return math.Sqrt(squares / float64(len(values)))
}

func percentile(c call, _ Env, args []Value) (Value, error) {
	if err := c.expectArgs(args, 2); err != nil {
		return nil, err
	}
	q, ok := scalarOf(args[1])
	if !ok {
		return nil, evalErrorAt(c.at, "p: q must be a plain number")
	}
	if q < 0 || q > 1 {
		return nil, evalErrorAt(c.at, "p: q must be between 0 and 1, got %g", q)
	}
	return Scalar(quantile(args[0], q)), nil
}

// extreme is min or max: with one argument it summarises it, with several it
// picks the extreme sample by sample (SPEC §5).
func extreme(pick func(a, b float64) float64, end float64) builtin {
	return func(c call, _ Env, args []Value) (Value, error) {
		switch len(args) {
		case 0:
			return nil, evalErrorAt(c.at, "%s expects at least 1 argument, got 0", c.name)
		case 1:
			return Scalar(quantile(args[0], end)), nil
		}
		result := args[0]
		for _, arg := range args[1:] {
			result = combine(pick, result, arg)
		}
		return result, nil
	}
}

type elementWise func(float64) float64

var mathFunctions = map[string]elementWise{
	"sqrt":  math.Sqrt,
	"exp":   math.Exp,
	"log":   math.Log,
	"log10": math.Log10,
	"abs":   math.Abs,
	"round": math.Round,
	"floor": math.Floor,
	"ceil":  math.Ceil,
}

func (f elementWise) call(c call, _ Env, args []Value) (Value, error) {
	if err := c.expectArgs(args, 1); err != nil {
		return nil, err
	}
	return combine(func(a, _ float64) float64 { return f(a) }, args[0], Scalar(0)), nil
}
