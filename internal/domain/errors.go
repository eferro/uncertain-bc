package domain

import "fmt"

type SyntaxError struct {
	Line int
	Col  int
	Msg  string
}

func (e *SyntaxError) Error() string {
	return position(e.Line, e.Col, e.Msg)
}

type EvalError struct {
	Line int
	Col  int
	Msg  string
}

func (e *EvalError) Error() string {
	return position(e.Line, e.Col, e.Msg)
}

func position(line, col int, msg string) string {
	return fmt.Sprintf("line %d, col %d: %s", line, col, msg)
}

func syntaxErrorAt(col int, format string, args ...any) *SyntaxError {
	return &SyntaxError{Line: 1, Col: col, Msg: fmt.Sprintf(format, args...)}
}

func evalErrorAt(col int, format string, args ...any) *EvalError {
	return &EvalError{Line: 1, Col: col, Msg: fmt.Sprintf(format, args...)}
}
