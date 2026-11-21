package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eferro/uncertain-bc/internal/adapters/cli"
)

func interactive(t *testing.T, input string, args ...string) outcome {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := cli.App{Interactive: true}.Run(args, strings.NewReader(input), &stdout, &stderr)
	return outcome{code: code, stdout: stdout.String(), stderr: stderr.String()}
}

func TestRepl_EvaluatesEachLineAfterAPrompt(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "ubc> 3\nubc> \n"}, interactive(t, "1 + 2\n"))
}

func TestRepl_PrintsAssignmentsAndKeepsVariables(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "ubc> 2\nubc> 6\nubc> \n"}, interactive(t, "x = 2\nx * 3\n"))
}

func TestRepl_ReportsAnErrorAndCarriesOn(t *testing.T) {
	want := outcome{code: 0, stdout: "ubc> ubc> 1\nubc> \n", stderr: "ubc: line 1, col 1: unknown variable 'y'\n"}

	assert.Equal(t, want, interactive(t, "y\n1\n"))
}

func TestRepl_IgnoresBlankLines(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "ubc> ubc> 1\nubc> \n"}, interactive(t, "\n1\n"))
}

func TestRepl_EvaluatesTheLastLine_WhenItHasNoNewline(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "ubc> 7\nubc> \n"}, interactive(t, "7"))
}

func TestRepl_IsNotStarted_WhenThereAreExpressionsOrFiles(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "4\n"}, interactive(t, "ignored\n", "-e", "4"))
}

func TestRepl_HonoursTheOutputFlags(t *testing.T) {
	got := interactive(t, "10 to 20\n", "--ascii", "-n", "100")

	assert.Regexp(t, `^ubc> ~ `, got.stdout)
}

func TestRepl_ExitsWithEvaluationError_WhenStdinCannotBeRead(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := cli.App{Interactive: true}.Run(nil, failingReader{}, &stdout, &stderr)

	assert.Equal(t, 1, code)
	assert.Equal(t, "ubc: stdin: broken pipe\n", stderr.String())
}
