package domain

import "fmt"

type SyntaxError struct {
	Line int
	Col  int
	Msg  string
}

func (e *SyntaxError) Error() string {
	return located(e.Line, e.Col, e.Msg)
}

type EvalError struct {
	Line int
	Col  int
	Msg  string
}

func (e *EvalError) Error() string {
	return located(e.Line, e.Col, e.Msg)
}

func located(line, col int, msg string) string {
	return fmt.Sprintf("line %d, col %d: %s", line, col, msg)
}

func syntaxErrorAt(at position, format string, args ...any) *SyntaxError {
	return &SyntaxError{Line: at.line, Col: at.col, Msg: fmt.Sprintf(format, args...)}
}

func evalErrorAt(at position, format string, args ...any) *EvalError {
	return &EvalError{Line: at.line, Col: at.col, Msg: fmt.Sprintf(format, args...)}
}
