package cli_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eferro/uncertain-bc/internal/adapters/cli"
)

var app = cli.App{Version: "v1.2.3"}

func TestRun_PrintsTheResult_WhenGivenAnExpression(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := app.Run([]string{"-e", "42"}, &stdout, &stderr)

	assert.Equal(t, 0, code)
	assert.Equal(t, "42\n", stdout.String())
	assert.Empty(t, stderr.String())
}

func TestRun_ExitsWithUsageError_WhenAFlagIsUnknown(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := app.Run([]string{"--nope"}, &stdout, &stderr)

	assert.Equal(t, 2, code)
	assert.Empty(t, stdout.String())
	assert.Contains(t, stderr.String(), "-nope")
}

func TestRun_ExitsWithParseError_WhenTheExpressionIsNotValid(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := app.Run([]string{"-e", "1 +"}, &stdout, &stderr)

	assert.Equal(t, 3, code)
	assert.Empty(t, stdout.String())
	assert.Regexp(t, `^ubc: .+\n$`, stderr.String())
}

func TestRun_ExitsWithEvaluationError_WhenTheExpressionCannotBeEvaluated(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := app.Run([]string{"-e", "10 to 5"}, &stdout, &stderr)

	assert.Equal(t, 1, code)
	assert.Empty(t, stdout.String())
	assert.Regexp(t, `^ubc: .+\n$`, stderr.String())
}

func TestRun_SummarisesASampleSetAsMedianAndNinetyPercentInterval(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := app.Run([]string{"-e", "10 to 20"}, &stdout, &stderr)

	assert.Equal(t, 0, code)
	assert.Regexp(t, `^≈ \S+   90% \[\S+ … \S+\]\n$`, stdout.String())
}

func TestRun_PrintsTheVersion_WhenAskedFor(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := app.Run([]string{"--version"}, &stdout, &stderr)

	assert.Equal(t, 0, code)
	assert.Equal(t, "ubc v1.2.3\n", stdout.String())
	assert.Empty(t, stderr.String())
}

func TestRun_PrintsEveryExpressionButNoAssignment(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := app.Run([]string{"-e", "x = 2; x * 3; y = x; y + 1"}, &stdout, &stderr)

	assert.Equal(t, 0, code)
	assert.Equal(t, "6\n3\n", stdout.String())
}

func TestRun_PrintsAPlainNumber_WhenEverySampleIsTheSame(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := app.Run([]string{"-e", "x = 1 to 5; x - x"}, &stdout, &stderr)

	assert.Equal(t, 0, code)
	assert.Equal(t, "0\n", stdout.String())
}

func TestRun_PrintsTheResultsBeforeAnError(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := app.Run([]string{"-e", "1; 10 to 5; 3"}, &stdout, &stderr)

	assert.Equal(t, 1, code)
	assert.Equal(t, "1\n", stdout.String())
	assert.Equal(t, "ubc: line 1, col 7: range start 10 must be below its end 5\n", stderr.String())
}
