package domain

import "errors"

// Override binds name to the value of expr and keeps it there: later
// assignments to name are skipped without being evaluated, so a value given
// from outside a model wins over the model's own (ADR 0028).
func Override(name, expr string, env Env) error {
	value, err := evaluateExpression(expr, env)
	if err != nil {
		return err
	}
	env.bindOverride(name, value)
	return nil
}

// OverrideData binds name to observed values, resampled like a list literal,
// and keeps it there as Override does (ADR 0031).
func OverrideData(name string, values []float64, env Env) error {
	if len(values) == 0 {
		return errors.New("no values to resample")
	}
	env.bindOverride(name, env.resample(values))
	return nil
}

func (e Env) bindOverride(name string, value Value) {
	e.variables[name] = value
	e.overridden[name] = true
}

// IsName reports whether name can be assigned to in a program.
func IsName(name string) bool {
	tokens, err := lex(name)
	return err == nil && tokens[0].kind == tokenIdent && tokens[0].text == name
}

func evaluateExpression(source string, env Env) (Value, error) {
	expr, err := parseExpression(source)
	if err != nil {
		return nil, err
	}
	return statement{expr: expr}.evaluate(env)
}

// parseExpression reads a source that must hold exactly one expression, as
// the command-line flags do.
func parseExpression(source string) (node, error) {
	tokens, err := lex(source)
	if err != nil {
		return nil, err
	}
	p := &parser{tokens: tokens}
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if tok := p.peek(); tok.kind != tokenEOF {
		return nil, unexpected(tok)
	}
	return expr, nil
}
