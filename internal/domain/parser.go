package domain

import (
	"slices"
	"strconv"
)

type statement struct {
	line int
	name string
	expr node
}

type parser struct {
	tokens []token
	pos    int
}

func parse(tokens []token) ([]statement, error) {
	p := &parser{tokens: tokens}
	var statements []statement
	for p.skipSeparators(); p.peek().kind != tokenEOF; p.skipSeparators() {
		stmt, err := p.statement()
		if err != nil {
			return nil, err
		}
		if tok := p.peek(); !isSeparator(tok) && tok.kind != tokenEOF {
			return nil, unexpected(tok)
		}
		statements = append(statements, stmt)
	}
	return statements, nil
}

func (p *parser) statement() (statement, error) {
	line := p.peek().line
	name := ""
	if p.peek().kind == tokenIdent && p.tokens[p.pos+1].kind == tokenAssign {
		name = p.advance().text
		p.advance()
	}
	expr, err := p.rangeExpression()
	return statement{line: line, name: name, expr: expr}, err
}

func (p *parser) skipSeparators() {
	for isSeparator(p.peek()) {
		p.advance()
	}
}

func isSeparator(tok token) bool {
	return tok.kind == tokenNewline || tok.kind == tokenSemicolon
}

func (p *parser) rangeExpression() (node, error) {
	from, err := p.additive()
	if err != nil || p.peek().kind != tokenTo {
		return from, err
	}
	at := p.advance().position
	to, err := p.additive()
	return rangeNode{from: from, to: to, at: at}, err
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
	case tokenIdent:
		return variable{name: tok.text, at: tok.position}, nil
	case tokenLParen:
		return p.parenthesised()
	default:
		return nil, unexpected(tok)
	}
}

func (p *parser) parenthesised() (node, error) {
	inner, err := p.rangeExpression()
	if err != nil {
		return nil, err
	}
	if tok := p.advance(); tok.kind != tokenRParen {
		return nil, unexpected(tok)
	}
	return inner, nil
}

func parseNumber(tok token) (node, error) {
	digits, scale := splitSuffix(tok.text)
	value, err := strconv.ParseFloat(digits, 64)
	if err != nil {
		return nil, syntaxErrorAt(tok.position, "invalid number '%s'", tok.text)
	}
	return number(scale(value)), nil
}

func splitSuffix(text string) (string, func(float64) float64) {
	last := text[len(text)-1]
	if last == '%' {
		return text[:len(text)-1], func(v float64) float64 { return v / 100 }
	}
	if multiplier, ok := multipliers[last]; ok {
		return text[:len(text)-1], func(v float64) float64 { return v * multiplier }
	}
	return text, func(v float64) float64 { return v }
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
	switch tok.kind {
	case tokenEOF:
		return syntaxErrorAt(tok.position, "unexpected end of input")
	case tokenNewline:
		return syntaxErrorAt(tok.position, "unexpected end of line")
	default:
		return syntaxErrorAt(tok.position, "unexpected '%s'", tok.text)
	}
}
