// Package domain holds the pure core of ubc: the language and its evaluation.
package domain

import "math"

type expression interface {
	evaluate(env Env) (Value, error)
}

type node interface {
	eval() float64
}

type number float64

func (n number) eval() float64 { return float64(n) }

type negate struct {
	operand node
}

func (n negate) eval() float64 { return -n.operand.eval() }

type binary struct {
	op          tokenKind
	left, right node
}

func (b binary) eval() float64 {
	return arithmetic[b.op](b.left.eval(), b.right.eval())
}

var arithmetic = map[tokenKind]func(a, b float64) float64{
	tokenPlus:  func(a, b float64) float64 { return a + b },
	tokenMinus: func(a, b float64) float64 { return a - b },
	tokenStar:  func(a, b float64) float64 { return a * b },
	tokenSlash: func(a, b float64) float64 { return a / b },
	tokenCaret: math.Pow,
}

type scalarExpr struct {
	node node
}

func (s scalarExpr) evaluate(Env) (Value, error) {
	return Scalar(s.node.eval()), nil
}

type rangeExpr struct {
	from, to node
	col      int
}

func (r rangeExpr) evaluate(env Env) (Value, error) {
	from, to := r.from.eval(), r.to.eval()
	if from >= to {
		return nil, evalErrorAt(r.col, "range start %g must be below its end %g", from, to)
	}
	return ninetyPercentRange(env, from, to), nil
}

func Evaluate(program string, env Env) (Value, error) {
	tokens, err := lex(program)
	if err != nil {
		return nil, err
	}
	expr, err := parse(tokens)
	if err != nil {
		return nil, err
	}
	return expr.evaluate(env)
}
