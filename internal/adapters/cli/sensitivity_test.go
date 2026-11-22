package cli_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const perfectDrivers = "sensitivity of x (rank correlation):\n" +
	"  double    +1.00  ████████████████████\n" +
	"  opposite  -1.00  ████████████████████\n"

func TestSens_ListsTheVariablesThatDriveAnExpression(t *testing.T) {
	got := interactive(t, "x = 1 to 10\ndouble = x * 2\nopposite = 0 - x\n:sens x\n")

	assert.Equal(t, "", got.stderr)
	assert.Contains(t, withoutPrompts(got.stdout), perfectDrivers)
}

func TestSens_LeavesLastAsItWas(t *testing.T) {
	got := interactive(t, "x = 1 to 10\ny = x * 2\n7\n:sens x + 1\nlast\n")

	assert.True(t, strings.HasSuffix(got.stdout, "ubc> 7\nubc> \n"), got.stdout)
}

func TestSens_ReportsWhatItCannotRank(t *testing.T) {
	cases := map[string]string{
		":sens":     "ubc: :sens expects an expression\n",
		":sens 5":   "ubc: :sens needs an uncertain value, got a plain number\n",
		":sens x":   "ubc: :sens needs other uncertain variables to compare with\n",
		":sens y":   "ubc: line 1, col 1: unknown variable 'y'\n",
		":sens x +": "ubc: line 1, col 4: unexpected end of input\n",
	}
	for input, want := range cases {
		got := interactive(t, "x = 1 to 2\n"+input+"\n")

		assert.Equal(t, want, got.stderr, input)
	}
}

func TestSensFlag_ListsTheDriversAfterTheProgram(t *testing.T) {
	got := run(t, "", "--sens", "x", "-e", "x = 1 to 10", "-e", "double = x * 2", "-e", "opposite = 0 - x", "-e", "3")

	assert.Equal(t, 0, got.code)
	assert.Equal(t, "3\n"+perfectDrivers, got.stdout)
}

func TestSensFlag_RunsBeforeTheAssertions(t *testing.T) {
	got := run(t, "", "--sens", "x", "--assert", "x < 0", "-e", "x = 1 to 10", "-e", "y = x")

	assert.Equal(t, 4, got.code)
	assert.Contains(t, got.stdout, "  y  +1.00")
}

func TestSensFlag_IsAnEvaluationError_WhenItCannotRank(t *testing.T) {
	got := run(t, "", "--sens", "x", "--assert", "x < 0", "-e", "x = 1 to 10")

	assert.Equal(t, outcome{code: 1, stderr: "ubc: --sens x: needs other uncertain variables to compare with\n"}, got)
}

func TestSensFlag_IsAParseError_WhenTheExpressionIsBroken(t *testing.T) {
	got := run(t, "", "--sens", "x +", "-e", "x = 1 to 10")

	assert.Equal(t, 3, got.code)
}

func TestSensFlag_NeedsAProgramAndATextOutput(t *testing.T) {
	cases := map[string][]string{
		"ubc: --sens needs a program: -e, files or stdin\n":       {"--sens", "x"},
		"ubc: --sens cannot be combined with --json or --plain\n": {"--sens", "x", "--plain", "-e", "1"},
	}
	for want, args := range cases {
		got := interactive(t, "", args...)

		assert.Equal(t, outcome{code: 2, stderr: want}, got, args)
	}
}
