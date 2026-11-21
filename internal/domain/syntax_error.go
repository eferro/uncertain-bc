package domain

import "fmt"

type SyntaxError struct {
	Line int
	Col  int
	Msg  string
}

func (e *SyntaxError) Error() string {
	return fmt.Sprintf("line %d, col %d: %s", e.Line, e.Col, e.Msg)
}

func syntaxErrorAt(col int, format string, args ...any) *SyntaxError {
	return &SyntaxError{Line: 1, Col: col, Msg: fmt.Sprintf(format, args...)}
}
