package domain_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestEvaluate_ReportsWhereTheSyntaxIsWrong(t *testing.T) {
	cases := []struct {
		program string
		col     int
		message string
	}{
		{"4 $ 2", 3, "unexpected character '$'"},
		{"1 to", 5, "unexpected end of input"},
		{"1 +", 4, "unexpected end of input"},
		{"1 2", 3, "unexpected '2'"},
		{"+ 1", 1, "unexpected '+'"},
		{"(1 + 2", 7, "unexpected end of input"},
		{"(1 + 2 3", 8, "unexpected '3'"},
		{")", 1, "unexpected ')'"},
		{"(*)", 2, "unexpected '*'"},
		{"1 * * 2", 5, "unexpected '*'"},
		{"2 ^", 4, "unexpected end of input"},
		{"1 +\n2", 4, "unexpected end of line"},
		{"x =", 4, "unexpected end of input"},
		{"= 3", 1, "unexpected '='"},
		{"2 = 3", 3, "unexpected '='"},
		{"to = 3", 1, "unexpected 'to'"},
		{"1 + ; 2", 5, "unexpected ';'"},
		{"1.2.3", 1, "invalid number '1.2.3'"},
		{"1e", 2, "unexpected 'e'"},
		{"1kw", 2, "unexpected 'kw'"},
		{"1 k", 3, "unexpected 'k'"},
		{"1m", 2, "unexpected 'm'"},
		{"5%%", 3, "unexpected character '%'"},
		{"%", 1, "unexpected character '%'"},
		{"1e+", 2, "unexpected 'e'"},
	}
	for _, c := range cases {
		t.Run(c.program, func(t *testing.T) {
			_, err := domain.Evaluate(c.program, domain.NewEnv(1, samples))

			var syntaxErr *domain.SyntaxError
			require.ErrorAs(t, err, &syntaxErr)
			assert.Equal(t, domain.SyntaxError{Line: 1, Col: c.col, Msg: c.message}, *syntaxErr)
		})
	}
}

func TestSyntaxError_ReadsAsLineColumnAndMessage(t *testing.T) {
	err := &domain.SyntaxError{Line: 2, Col: 5, Msg: "unexpected '+'"}

	assert.EqualError(t, err, "line 2, col 5: unexpected '+'")
}

func TestEvaluate_RejectsTheCharactersJustOutsideDigitAndLetterRanges(t *testing.T) {
	for _, neighbour := range []string{":", "@", "[", "`", "{"} {
		t.Run(neighbour, func(t *testing.T) {
			_, err := domain.Evaluate(neighbour, domain.NewEnv(1, samples))

			var syntaxErr *domain.SyntaxError
			require.ErrorAs(t, err, &syntaxErr)
			assert.Equal(t, "unexpected character '"+neighbour+"'", syntaxErr.Msg)
		})
	}
}
