// Package domain holds the pure core of ubc: the language and its evaluation.
package domain

type expression interface {
	evaluate(env Env) (Value, error)
}

type node interface {
	eval() float64
}

type number float64

func (n number) eval() float64 { return float64(n) }

type add struct {
	left, right node
}

func (a add) eval() float64 { return a.left.eval() + a.right.eval() }

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
