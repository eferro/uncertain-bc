package domain

type comparison struct {
	op          tokenKind
	symbol      string
	left, right node
}

func (c comparison) evaluate(env Env) (Value, error) {
	left, right, err := c.operands(env)
	if err != nil {
		return nil, err
	}
	return c.apply(left, right), nil
}

func (c comparison) operands(env Env) (left, right Value, err error) {
	if left, err = c.left.evaluate(env); err != nil {
		return nil, nil, err
	}
	right, err = c.right.evaluate(env)
	return left, right, err
}

func (c comparison) apply(left, right Value) Value {
	return combine(comparators[c.op], left, right)
}

var comparators = map[tokenKind]func(a, b float64) float64{
	tokenLess:         func(a, b float64) float64 { return indicator(a < b) },
	tokenGreater:      func(a, b float64) float64 { return indicator(a > b) },
	tokenLessEqual:    func(a, b float64) float64 { return indicator(a <= b) },
	tokenGreaterEqual: func(a, b float64) float64 { return indicator(a >= b) },
	tokenEqual:        func(a, b float64) float64 { return indicator(a == b) },
	tokenNotEqual:     func(a, b float64) float64 { return indicator(a != b) },
}

func indicator(holds bool) float64 {
	if holds {
		return 1
	}
	return 0
}

func isIndicator(v Value) bool {
	for i := range size(v) {
		if x := at(v, i); x != 0 && x != 1 {
			return false
		}
	}
	return true
}
