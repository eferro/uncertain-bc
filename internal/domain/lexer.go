package domain

type tokenKind int

const (
	tokenEOF tokenKind = iota
	tokenNumber
	tokenPlus
	tokenMinus
	tokenStar
	tokenSlash
	tokenCaret
	tokenLParen
	tokenRParen
	tokenTo
	tokenIdent
)

var operators = map[byte]tokenKind{
	'+': tokenPlus,
	'-': tokenMinus,
	'*': tokenStar,
	'/': tokenSlash,
	'^': tokenCaret,
	'(': tokenLParen,
	')': tokenRParen,
}

type token struct {
	kind tokenKind
	text string
	col  int
}

func lex(src string) ([]token, error) {
	var tokens []token
	for i := 0; i < len(src); {
		next, tok, err := lexOne(src, i)
		if err != nil {
			return nil, err
		}
		if tok != nil {
			tokens = append(tokens, *tok)
		}
		i = next
	}
	return append(tokens, token{kind: tokenEOF, col: len(src) + 1}), nil
}

func lexOne(src string, i int) (int, *token, error) {
	c := src[i]
	if kind, ok := operators[c]; ok {
		return i + 1, &token{kind: kind, text: string(c), col: i + 1}, nil
	}
	switch {
	case c == ' ' || c == '\t':
		return i + 1, nil, nil
	case isDigit(c) || c == '.':
		end := scanNumber(src, i)
		return end, &token{kind: tokenNumber, text: src[i:end], col: i + 1}, nil
	case isIdentStart(c):
		end := scanWhile(src, i, isIdentPart)
		return end, identOrKeyword(src[i:end], i+1), nil
	}
	return 0, nil, syntaxErrorAt(i+1, "unexpected character '%c'", c)
}

func scanNumber(src string, start int) int {
	end := scanWhile(src, start, func(c byte) bool { return isDigit(c) || c == '.' })
	return scanExponent(src, end)
}

func scanExponent(src string, i int) int {
	if i >= len(src) || (src[i] != 'e' && src[i] != 'E') {
		return i
	}
	j := i + 1
	if j < len(src) && (src[j] == '+' || src[j] == '-') {
		j++
	}
	if j >= len(src) || !isDigit(src[j]) {
		return i
	}
	return scanWhile(src, j, isDigit)
}

func scanWhile(src string, i int, accept func(byte) bool) int {
	for i < len(src) && accept(src[i]) {
		i++
	}
	return i
}

func isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func identOrKeyword(text string, col int) *token {
	if text == "to" {
		return &token{kind: tokenTo, text: text, col: col}
	}
	return &token{kind: tokenIdent, text: text, col: col}
}

func isIdentStart(c byte) bool {
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isIdentPart(c byte) bool {
	return isIdentStart(c) || isDigit(c)
}
