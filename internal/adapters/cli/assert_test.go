package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eferro/uncertain-bc/internal/adapters/cli"
)

func TestAssert_ExitsWithSuccess_WhenEveryAssertionHolds(t *testing.T) {
	got := run(t, "", "--assert", "x < 3", "--assert", "x == 2", "-e", "x = 2")

	assert.Equal(t, outcome{code: 0}, got)
}

func TestAssert_ExitsWithAssertionFailure_AndShowsTheComparedValues(t *testing.T) {
	want := outcome{code: 4, stdout: "4\n", stderr: "ubc: --assert x > 3 failed: 2 > 3\n"}

	assert.Equal(t, want, run(t, "", "--assert", "x > 3", "-e", "x = 2; x * 2"))
}

func TestAssert_ShowsTheValue_WhenItIsNotAComparison(t *testing.T) {
	want := outcome{code: 4, stderr: "ubc: --assert P(normal(0, 1) > 0) failed: 0.50\n"}

	assert.Equal(t, want, run(t, "", "--assert", "P(normal(0, 1) > 0)", "-e", "x = 1"))
}

func TestAssert_ShowsTheSummaryAndTheProbability_WhenTheComparisonIsUncertain(t *testing.T) {
	want := "ubc: --assert normal(0, 1) < 1 failed: ~ 0.00436   90% [-1.65 ... 1.6] < 1 (P = 0.85)\n"

	assert.Equal(t, want, run(t, "", "--ascii", "--plain", "--assert", "normal(0, 1) < 1", "-e", "x = 1").stderr)
}

func TestAssert_ReportsEveryFailure(t *testing.T) {
	want := "ubc: --assert 1 > 2 failed: 1 > 2\nubc: --assert 0 failed: 0\n"

	assert.Equal(t, want, run(t, "", "--assert", "1 > 2", "--assert", "1 < 2", "--assert", "0", "-e", "x = 1").stderr)
}

func TestAssert_ExitsWithEvaluationError_WhenTheValueIsNotACondition(t *testing.T) {
	want := outcome{code: 1, stderr: "ubc: --assert x: line 1, col 1: an assertion must be a condition or a probability\n"}

	assert.Equal(t, want, run(t, "", "--assert", "x", "--assert", "x > 5", "-e", "x = 2"))
}

func TestAssert_ExitsWithParseError_WhenTheExpressionIsNotValid(t *testing.T) {
	want := outcome{code: 3, stderr: "ubc: --assert 1 <: line 1, col 4: unexpected end of input\n"}

	assert.Equal(t, want, run(t, "", "--assert", "1 <", "-e", "x = 2"))
}

func TestAssert_IsSkipped_WhenTheProgramFails(t *testing.T) {
	want := outcome{code: 1, stderr: "ubc: line 1, col 1: unknown variable 'nobody'\n"}

	assert.Equal(t, want, run(t, "", "--assert", "1 > 2", "-e", "nobody"))
}

func TestAssert_RunsAfterTheFiles(t *testing.T) {
	path := writeFile(t, "model.ubc", "y = x * 10\n")

	assert.Equal(t, 4, run(t, "", "--assert", "y < 20", "-e", "x = 2", path).code)
}

func TestAssert_RunsAfterStdin(t *testing.T) {
	assert.Equal(t, 4, run(t, "x = 2\ny = 30\n", "--assert", "y < 20").code)
}

func TestAssert_IsAUsageError_InTheRepl(t *testing.T) {
	want := outcome{code: 2, stderr: "ubc: --assert needs a program: -e, files or stdin\n"}

	assert.Equal(t, want, interactive(t, "1\n", "--assert", "1 < 2"))
}

func TestAssert_ReportsFailuresOnEveryWatchRun(t *testing.T) {
	path := writeFile(t, "model.ubc", "x = 1\n")
	rewrite(t, path, "x = 1\n", 0)
	polls := 0
	app := cli.App{Wait: func() bool {
		polls++
		if polls == 1 {
			rewrite(t, path, "x = 5\n", polls)
		}
		return polls == 1
	}}
	var stdout, stderr bytes.Buffer

	code := app.Run([]string{"--assert", "x < 3", "--watch", path}, strings.NewReader(""), &stdout, &stderr)

	assert.Equal(t, 4, code)
	assert.Equal(t, "ubc: --assert x < 3 failed: 5 < 3\n", stderr.String())
}

func TestAssert_IsListedInTheHelp(t *testing.T) {
	assert.Contains(t, run(t, "", "--help").stdout, "--assert EXPR")
}

func TestAssert_ShowsUnusualOperands_AsTheOutputDoes(t *testing.T) {
	cases := map[string]string{
		"x - x > 1":       "0 > 1",
		"(x - x) / 0 > 1": "no finite samples (100 dropped) > 1",
		"normal(0, 1) / (uniform(0, 1) > 0.5) < 5": "≈ 0.193   90% [-1.41 … 1.44]   (54 non-finite dropped) < 5 (P = 0.72)",
	}
	for expr, compared := range cases {
		t.Run(expr, func(t *testing.T) {
			got := run(t, "", "-n", "100", "--assert", expr, "-e", "x = 1 to 5")

			assert.Equal(t, "ubc: --assert "+expr+" failed: "+compared+"\n", got.stderr)
		})
	}
}
