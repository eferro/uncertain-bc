package domain

import (
	"slices"
	"strconv"
)

type statement struct {
	line     int
	name     string
	expr     node
	source   span
	function *definition
}

// span is where a statement's expression sits in the program, so a result
// can be labelled with the text the user wrote (ADR 0053).
type span struct {
	from, to int
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
	if p.isDefinition() {
		def, err := p.definition()
		return statement{line: line, function: def}, err
	}
	name := ""
	if p.peek().kind == tokenIdent && p.tokens[p.pos+1].kind == tokenAssign {
		name = p.advance().text
		p.advance()
	}
	from := p.peek().offset
	expr, err := p.expression()
	return statement{line: line, name: name, expr: expr, source: span{from: from, to: p.tokens[p.pos-1].end()}}, err
}

// isDefinition looks ahead for name(...) followed by '=', which tells a
// function definition from a call.
func (p *parser) isDefinition() bool {
	if p.peek().kind != tokenIdent || p.tokens[p.pos+1].kind != tokenLParen {
		return false
	}
	end, found := p.closingParen(p.pos + 1)
	return found && p.tokens[end+1].kind == tokenAssign
}

func (p *parser) closingParen(open int) (int, bool) {
	depth := 0
	for i := open; p.tokens[i].kind != tokenEOF && !isSeparator(p.tokens[i]); i++ {
		depth += parenDepth[p.tokens[i].kind]
		if depth == 0 {
			return i, true
		}
	}
	return 0, false
}

var parenDepth = map[tokenKind]int{tokenLParen: 1, tokenRParen: -1}

func (p *parser) definition() (*definition, error) {
	name := p.advance()
	p.advance()
	params, err := p.parameters()
	if err != nil {
		return nil, err
	}
	p.advance()
	body, err := p.expression()
	return &definition{name: name.text, params: params, body: body, at: name.position}, err
}

func (p *parser) parameters() ([]token, error) {
	var params []token
	if p.peek().kind == tokenRParen {
		p.advance()
		return params, nil
	}
	for {
		param := p.advance()
		if param.kind != tokenIdent {
			return nil, unexpected(param)
		}
		params = append(params, param)
		if next := p.advance(); next.kind == tokenRParen {
			return params, nil
		} else if next.kind != tokenComma {
			return nil, unexpected(next)
		}
	}
}

func (p *parser) skipSeparators() {
	for isSeparator(p.peek()) {
		p.advance()
	}
}

func isSeparator(tok token) bool {
	return tok.kind == tokenNewline || tok.kind == tokenSemicolon
}

var comparisons = []tokenKind{tokenLess, tokenGreater, tokenLessEqual, tokenGreaterEqual, tokenEqual, tokenNotEqual}

func (p *parser) expression() (node, error) {
	left, err := p.rangeExpression()
	if err != nil || !slices.Contains(comparisons, p.peek().kind) {
		return left, err
	}
	op := p.advance()
	right, err := p.rangeExpression()
	return comparison{op: op.kind, symbol: op.text, left: left, right: right}, err
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
		if p.peek().kind == tokenLParen {
			return p.call(tok)
		}
		return variable{name: tok.text, at: tok.position}, nil
	case tokenLParen:
		return p.parenthesised()
	case tokenLBracket:
		elements, err := p.arguments(tokenRBracket)
		return list{elements: elements, at: tok.position}, err
	default:
		return nil, unexpected(tok)
	}
}

func (p *parser) call(name token) (node, error) {
	p.advance()
	args, err := p.arguments(tokenRParen)
	if err != nil {
		return nil, err
	}
	return call{name: name.text, args: args, at: name.position}, nil
}

func (p *parser) arguments(closing tokenKind) ([]node, error) {
	var args []node
	if p.peek().kind == closing {
		p.advance()
		return args, nil
	}
	for {
		arg, err := p.expression()
		if err != nil {
			return nil, err
		}
		args = append(args, arg)
		if tok := p.advance(); tok.kind == closing {
			return args, nil
		} else if tok.kind != tokenComma {
			return nil, unexpected(tok)
		}
	}
}

func (p *parser) parenthesised() (node, error) {
	inner, err := p.expression()
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
