package domain

type tokenKind int

const (
	tokenEOF tokenKind = iota
	tokenNewline
	tokenSemicolon
	tokenNumber
	tokenIdent
	tokenTo
	tokenAssign
	tokenPlus
	tokenMinus
	tokenStar
	tokenSlash
	tokenCaret
	tokenLParen
	tokenRParen
	tokenComma
	tokenLess
	tokenGreater
	tokenLessEqual
	tokenGreaterEqual
	tokenEqual
	tokenNotEqual
)

var operators = map[byte]tokenKind{
	';': tokenSemicolon,
	'=': tokenAssign,
	'+': tokenPlus,
	'-': tokenMinus,
	'*': tokenStar,
	'/': tokenSlash,
	'^': tokenCaret,
	'(': tokenLParen,
	')': tokenRParen,
	',': tokenComma,
	'<': tokenLess,
	'>': tokenGreater,
}

var twoCharOperators = map[string]tokenKind{
	"<=": tokenLessEqual,
	">=": tokenGreaterEqual,
	"==": tokenEqual,
	"!=": tokenNotEqual,
}

type position struct {
	line, col int
}

type token struct {
	kind tokenKind
	text string
	position
}

type lexer struct {
	src       string
	i         int
	line      int
	lineStart int
	tokens    []token
}

func lex(src string) ([]token, error) {
	l := &lexer{src: src, line: 1}
	for l.i < len(l.src) {
		if err := l.next(); err != nil {
			return nil, err
		}
	}
	l.emit(tokenEOF, "", l.i)
	return l.tokens, nil
}

func (l *lexer) next() error {
	switch c := l.src[l.i]; c {
	case '\n':
		l.newline()
	case ' ', '\t', '\r':
		l.i++
	case '#':
		l.i = scanWhile(l.src, l.i, func(c byte) bool { return c != '\n' })
	default:
		return l.token(c)
	}
	return nil
}

func (l *lexer) token(c byte) error {
	if l.operator() {
		return nil
	}
	switch {
	case isDigit(c) || c == '.':
		l.emitAndSkip(tokenNumber, scanNumber(l.src, l.i))
	case isIdentStart(c):
		end := scanWhile(l.src, l.i, isIdentPart)
		l.emitAndSkip(identKind(l.src[l.i:end]), end)
	default:
		return syntaxErrorAt(l.position(l.i), "unexpected character '%c'", c)
	}
	return nil
}

func (l *lexer) operator() bool {
	if l.i+1 < len(l.src) {
		if kind, ok := twoCharOperators[l.src[l.i:l.i+2]]; ok {
			l.emitAndSkip(kind, l.i+2)
			return true
		}
	}
	if kind, ok := operators[l.src[l.i]]; ok {
		l.emitAndSkip(kind, l.i+1)
		return true
	}
	return false
}

func (l *lexer) newline() {
	l.emit(tokenNewline, "", l.i)
	l.i++
	l.line++
	l.lineStart = l.i
}

func (l *lexer) emitAndSkip(kind tokenKind, end int) {
	l.emit(kind, l.src[l.i:end], l.i)
	l.i = end
}

func (l *lexer) emit(kind tokenKind, text string, at int) {
	l.tokens = append(l.tokens, token{kind: kind, text: text, position: l.position(at)})
}

func (l *lexer) position(at int) position {
	return position{line: l.line, col: at - l.lineStart + 1}
}

func identKind(text string) tokenKind {
	if text == "to" {
		return tokenTo
	}
	return tokenIdent
}

func scanNumber(src string, start int) int {
	end := scanWhile(src, start, func(c byte) bool { return isDigit(c) || c == '.' })
	return scanSuffix(src, scanExponent(src, end))
}

var multipliers = map[byte]float64{'k': 1e3, 'M': 1e6, 'B': 1e9, 'T': 1e12}

func scanSuffix(src string, i int) int {
	if i >= len(src) {
		return i
	}
	_, isMultiplier := multipliers[src[i]]
	endsHere := i+1 >= len(src) || !isIdentPart(src[i+1])
	if src[i] == '%' || (isMultiplier && endsHere) {
		return i + 1
	}
	return i
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

func isIdentStart(c byte) bool {
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isIdentPart(c byte) bool {
	return isIdentStart(c) || isDigit(c)
}
