package cli_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQuit_EndsTheSessionAtOnce(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "ubc> "}, interactive(t, ":quit\n1\n"))
}

func TestVars_SaysSo_WhenThereAreNone(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "ubc> (no variables)\nubc> \n"}, interactive(t, ":vars\n"))
}

func TestVars_ListsEachVariableWithItsSummary(t *testing.T) {
	got := interactive(t, "x = 2\ny = 1 to 5\nx = 3\n:vars\n")

	lines := strings.Split(got.stdout, "\n")
	assert.Equal(t, "ubc> x = 3", lines[3])
	assert.Regexp(t, `^y = ≈ \S+   90% `, lines[4])
}

func TestSeed_RestartsTheRandomStream(t *testing.T) {
	reseeded := interactive(t, ":seed 2\n1 to 2\n")
	fresh := interactive(t, "1 to 2\n", "-s", "2")

	assert.Equal(t, strings.TrimPrefix(fresh.stdout, "ubc> "), strings.TrimPrefix(reseeded.stdout, "ubc> ubc> "))
}

func TestSamples_ChangesTheSampleCountQuietly_WhenThereAreNoVariables(t *testing.T) {
	for _, input := range []string{":samples 1\n1 to 2\n", ":samples   1  \n1 to 2\n"} {
		got := interactive(t, input)

		assert.Regexp(t, `^ubc> ubc> 1\.\d+\nubc> \n$`, got.stdout, input)
		assert.Empty(t, got.stderr, input)
	}
}

func TestSamples_WarnsThatItClearsVariables(t *testing.T) {
	got := interactive(t, "x = 1\n:samples 10\nx\n")

	assert.Equal(t, "ubc: variables cleared (sample count changed)\nubc: line 1, col 1: unknown variable 'x'\n", got.stderr)
}

func TestHelp_ListsCommandsAndFunctions(t *testing.T) {
	got := interactive(t, ":help\n")

	assert.Contains(t, got.stdout, ":vars")
	assert.Contains(t, got.stdout, ":samples N")
	assert.Contains(t, got.stdout, "normal(mean, sd)")
}

func TestCommands_ReportMistakesAndCarryOn(t *testing.T) {
	cases := map[string]string{
		":nope":        "ubc: unknown command ':nope' (try :help)\n",
		":seed":        "ubc: :seed expects a whole number, got ''\n",
		":seed x":      "ubc: :seed expects a whole number, got 'x'\n",
		":samples 0":   "ubc: :samples expects a positive whole number, got '0'\n",
		":samples 1.5": "ubc: :samples expects a positive whole number, got '1.5'\n",
	}
	for input, message := range cases {
		t.Run(input, func(t *testing.T) {
			got := interactive(t, input+"\n1\n")

			assert.Equal(t, outcome{code: 0, stdout: "ubc> ubc> 1\nubc> \n", stderr: message}, got)
		})
	}
}

func TestHist_DrawsTheDistributionOfAnExpression(t *testing.T) {
	got := interactive(t, "x = uniform(0, 10)\n:hist x * 2\n")

	lines := strings.Split(got.stdout, "\n")
	assert.Empty(t, got.stderr)
	assert.Regexp(t, `^ubc> \s+\S+ ┤█+$`, lines[1])
	assert.Regexp(t, `^\(P1 to P99 shown, \d+ samples outside\)$`, lines[21])
}

func TestHist_ReportsWhatItCannotDraw(t *testing.T) {
	cases := map[string]string{
		":hist":     "ubc: :hist expects an expression\n",
		":hist 5":   "ubc: :hist needs an uncertain value, got a plain number\n",
		":hist 1>0": "ubc: :hist needs an uncertain value, got a plain number\n",
		":hist y":   "ubc: line 1, col 1: unknown variable 'y'\n",
	}
	for input, message := range cases {
		t.Run(input, func(t *testing.T) {
			assert.Equal(t, message, interactive(t, input+"\n").stderr)
		})
	}
}
