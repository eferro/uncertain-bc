package domain

// Check is the outcome of an assertion. A comparison keeps its operands, so
// a failure can show what was compared, not only that it was false (ADR 0029).
type Check struct {
	Value       Value
	Left, Right Value
	Operator    string
}

// Holds is true only for a certain condition: a probability of 0.99 is a
// failure, because a check in a script must not pass by luck.
func (c Check) Holds() bool {
	x, _ := scalarOf(c.Value)
	return x == 1
}

// Assert evaluates one expression as a condition without binding anything,
// so assertions see the session but leave it as it was.
func Assert(source string, env Env) (Check, error) {
	expr, err := parseExpression(source)
	if err != nil {
		return Check{}, err
	}
	result, err := evaluateCheck(expr, env)
	if err != nil {
		return Check{}, err
	}
	if !isCondition(result.Value) {
		return Check{}, evalErrorAt(position{line: 1, col: 1}, "an assertion must be a condition or a probability")
	}
	return result, nil
}

func evaluateCheck(expr node, env Env) (Check, error) {
	c, isComparison := expr.(comparison)
	if !isComparison {
		value, err := expr.evaluate(env)
		return Check{Value: Shown(value)}, err
	}
	left, right, err := c.operands(env)
	if err != nil {
		return Check{}, err
	}
	return Check{Value: Probability(mean(c.apply(left, right))), Left: left, Right: right, Operator: c.symbol}, nil
}

func isCondition(v Value) bool {
	switch v := v.(type) {
	case Probability:
		return true
	case Scalar:
		return v == 0 || v == 1
	}
	return false
}
