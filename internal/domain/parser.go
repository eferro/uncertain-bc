package domain

import "strconv"

type parser struct {
	tokens []token
	pos    int
}

func parse(tokens []token) (node, error) {
	p := &parser{tokens: tokens}
	expr, err := p.additive()
	if err != nil {
		return nil, err
	}
	if tok := p.peek(); tok.kind != tokenEOF {
		return nil, unexpected(tok)
	}
	return expr, nil
}

func (p *parser) additive() (node, error) {
	left, err := p.primary()
	for err == nil && p.peek().kind == tokenPlus {
		p.pos++
		var right node
		right, err = p.primary()
		left = add{left: left, right: right}
	}
	return left, err
}

func (p *parser) primary() (node, error) {
	tok := p.peek()
	if tok.kind != tokenNumber {
		return nil, unexpected(tok)
	}
	p.pos++
	value, err := strconv.ParseFloat(tok.text, 64)
	if err != nil {
		return nil, syntaxErrorAt(tok.col, "invalid number '%s'", tok.text)
	}
	return number(value), nil
}

func (p *parser) peek() token {
	return p.tokens[p.pos]
}

func unexpected(tok token) *SyntaxError {
	if tok.kind == tokenEOF {
		return syntaxErrorAt(tok.col, "unexpected end of input")
	}
	return syntaxErrorAt(tok.col, "unexpected '%s'", tok.text)
}
