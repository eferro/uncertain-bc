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
		{"forty", 1, "unexpected 'forty'"},
		{"4 $ 2", 3, "unexpected character '$'"},
		{"a", 1, "unexpected 'a'"},
		{"z", 1, "unexpected 'z'"},
		{"A", 1, "unexpected 'A'"},
		{"Z", 1, "unexpected 'Z'"},
		{"_x9", 1, "unexpected '_x9'"},
		{"1 to", 5, "unexpected end of input"},
		{"1 +", 4, "unexpected end of input"},
		{"", 1, "unexpected end of input"},
		{"1 2", 3, "unexpected '2'"},
		{"+ 1", 1, "unexpected '+'"},
		{"1.2.3", 1, "invalid number '1.2.3'"},
		{"1e", 2, "unexpected 'e'"},
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
	for _, neighbour := range []string{"/", ":", "@", "[", "`", "{"} {
		t.Run(neighbour, func(t *testing.T) {
			_, err := domain.Evaluate(neighbour, domain.NewEnv(1, samples))

			var syntaxErr *domain.SyntaxError
			require.ErrorAs(t, err, &syntaxErr)
			assert.Equal(t, "unexpected character '"+neighbour+"'", syntaxErr.Msg)
		})
	}
}
