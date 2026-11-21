package domain

// Override binds name to the value of expr and keeps it there: later
// assignments to name are skipped without being evaluated, so a value given
// from outside a model wins over the model's own (ADR 0028).
func Override(name, expr string, env Env) error {
	value, err := evaluateExpression(expr, env)
	if err != nil {
		return err
	}
	env.variables[name] = value
	env.overridden[name] = true
	return nil
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
