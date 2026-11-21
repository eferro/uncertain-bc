package domain

import (
	"slices"
	"strconv"
)

type parser struct {
	tokens []token
	pos    int
}

func parse(tokens []token) (expression, error) {
	p := &parser{tokens: tokens}
	expr, err := p.rangeExpression()
	if err != nil {
		return nil, err
	}
	if tok := p.peek(); tok.kind != tokenEOF {
		return nil, unexpected(tok)
	}
	return expr, nil
}

func (p *parser) rangeExpression() (expression, error) {
	from, err := p.additive()
	if err != nil || p.peek().kind != tokenTo {
		return scalarExpr{node: from}, err
	}
	col := p.advance().col
	to, err := p.additive()
	return rangeExpr{from: from, to: to, col: col}, err
}

func (p *parser) additive() (node, error) {
	return p.leftAssociative(p.multiplicative, tokenPlus, tokenMinus)
}

func (p *parser) multiplicative() (node, error) {
	return p.leftAssociative(p.unary, tokenStar, tokenSlash)
}

func (p *parser) leftAssociative(operand func() (node, error), operators ...tokenKind) (node, error) {
	left, err := operand()
	for err == nil && slices.Contains(operators, p.peek().kind) {
		op := p.advance().kind
		var right node
		right, err = operand()
		left = binary{op: op, left: left, right: right}
	}
	return left, err
}

func (p *parser) unary() (node, error) {
	if p.peek().kind != tokenMinus {
		return p.power()
	}
	p.advance()
	operand, err := p.unary()
	return negate{operand: operand}, err
}

func (p *parser) power() (node, error) {
	base, err := p.primary()
	if err != nil || p.peek().kind != tokenCaret {
		return base, err
	}
	p.advance()
	exponent, err := p.unary()
	return binary{op: tokenCaret, left: base, right: exponent}, err
}

func (p *parser) primary() (node, error) {
	switch tok := p.advance(); tok.kind {
	case tokenNumber:
		return parseNumber(tok)
	case tokenLParen:
		return p.parenthesised()
	default:
		return nil, unexpected(tok)
	}
}

func (p *parser) parenthesised() (node, error) {
	inner, err := p.additive()
	if err != nil {
		return nil, err
	}
	if tok := p.advance(); tok.kind != tokenRParen {
		return nil, unexpected(tok)
	}
	return inner, nil
}

func parseNumber(tok token) (node, error) {
	value, err := strconv.ParseFloat(tok.text, 64)
	if err != nil {
		return nil, syntaxErrorAt(tok.col, "invalid number '%s'", tok.text)
	}
	return number(value), nil
}

func (p *parser) peek() token {
	return p.tokens[p.pos]
}

func (p *parser) advance() token {
	tok := p.tokens[p.pos]
	if tok.kind != tokenEOF {
		p.pos++
	}
	return tok
}

func unexpected(tok token) *SyntaxError {
	if tok.kind == tokenEOF {
		return syntaxErrorAt(tok.col, "unexpected end of input")
	}
	return syntaxErrorAt(tok.col, "unexpected '%s'", tok.text)
}
