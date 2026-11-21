// Package domain holds the pure core of ubc: the language and its evaluation.
package domain

import "math"

type node interface {
	evaluate(env Env) (Value, error)
}

type number float64

func (n number) evaluate(Env) (Value, error) { return Scalar(n), nil }

type negate struct {
	operand node
}

func (n negate) evaluate(env Env) (Value, error) {
	value, err := n.operand.evaluate(env)
	if err != nil {
		return nil, err
	}
	return combine(func(a, _ float64) float64 { return -a }, value, Scalar(0)), nil
}

type binary struct {
	op          tokenKind
	left, right node
}

func (b binary) evaluate(env Env) (Value, error) {
	left, err := b.left.evaluate(env)
	if err != nil {
		return nil, err
	}
	right, err := b.right.evaluate(env)
	if err != nil {
		return nil, err
	}
	return combine(arithmetic[b.op], left, right), nil
}

var arithmetic = map[tokenKind]func(a, b float64) float64{
	tokenPlus:  func(a, b float64) float64 { return a + b },
	tokenMinus: func(a, b float64) float64 { return a - b },
	tokenStar:  func(a, b float64) float64 { return a * b },
	tokenSlash: func(a, b float64) float64 { return a / b },
	tokenCaret: math.Pow,
}

type variable struct {
	name string
	at   position
}

func (v variable) evaluate(env Env) (Value, error) {
	value, ok := env.variables[v.name]
	if !ok {
		return nil, evalErrorAt(v.at, "unknown variable '%s'", v.name)
	}
	return value, nil
}

type rangeNode struct {
	from, to node
	at       position
}

func (r rangeNode) evaluate(env Env) (Value, error) {
	from, to, err := r.bounds(env)
	if err != nil {
		return nil, err
	}
	if from >= to {
		return nil, evalErrorAt(r.at, "range start %g must be below its end %g", from, to)
	}
	return ninetyPercentRange(env, from, to), nil
}

func (r rangeNode) bounds(env Env) (from, to float64, err error) {
	fromValue, err := r.from.evaluate(env)
	if err != nil {
		return 0, 0, err
	}
	toValue, err := r.to.evaluate(env)
	if err != nil {
		return 0, 0, err
	}
	from, fromOK := scalarOf(fromValue)
	to, toOK := scalarOf(toValue)
	if !fromOK || !toOK {
		return 0, 0, evalErrorAt(r.at, "range bounds must be plain numbers")
	}
	return from, to, nil
}

type Result struct {
	Line  int
	Name  string
	Value Value
}

func Evaluate(program string, env Env) ([]Result, error) {
	tokens, err := lex(program)
	if err != nil {
		return nil, err
	}
	statements, err := parse(tokens)
	if err != nil {
		return nil, err
	}
	return run(statements, env)
}

func run(statements []statement, env Env) ([]Result, error) {
	results := make([]Result, 0, len(statements))
	for _, stmt := range statements {
		if stmt.function != nil {
			if err := env.define(*stmt.function); err != nil {
				return results, err
			}
			continue
		}
		value, err := stmt.evaluate(env)
		if err != nil {
			return results, err
		}
		env.variables[stmt.binding()] = value
		results = append(results, Result{Line: stmt.line, Name: stmt.name, Value: value})
	}
	return results, nil
}

// binding is the variable a statement sets: its own name, or last for an
// expression, as in bc.
func (s statement) binding() string {
	if s.name == "" {
		return "last"
	}
	return s.name
}

// evaluate reports a bare top-level comparison as a probability (SPEC §4 rule 5).
func (s statement) evaluate(env Env) (Value, error) {
	value, err := s.expr.evaluate(env)
	if _, isComparison := s.expr.(comparison); isComparison && err == nil {
		return Probability(mean(value)), nil
	}
	return value, err
}
