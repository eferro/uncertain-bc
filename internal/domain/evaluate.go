// Package domain holds the pure core of ubc: the language and its evaluation.
package domain

type node interface {
	eval() float64
}

type number float64

func (n number) eval() float64 { return float64(n) }

type add struct {
	left, right node
}

func (a add) eval() float64 { return a.left.eval() + a.right.eval() }

func Evaluate(program string) (float64, error) {
	tokens, err := lex(program)
	if err != nil {
		return 0, err
	}
	tree, err := parse(tokens)
	if err != nil {
		return 0, err
	}
	return tree.eval(), nil
}
