package domain

type comparison struct {
	op          tokenKind
	left, right node
}

func (c comparison) evaluate(env Env) (Value, error) {
	left, err := c.left.evaluate(env)
	if err != nil {
		return nil, err
	}
	right, err := c.right.evaluate(env)
	if err != nil {
		return nil, err
	}
	return combine(comparators[c.op], left, right), nil
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
