package cli_test

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/adapters/cli"
)

type outcome struct {
	code   int
	stdout string
	stderr string
}

func run(t *testing.T, stdin string, args ...string) outcome {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := cli.App{Version: "v1.2.3"}.Run(args, strings.NewReader(stdin), &stdout, &stderr)
	return outcome{code: code, stdout: stdout.String(), stderr: stderr.String()}
}

func writeFile(t *testing.T, name, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	return path
}

func TestRun_PrintsTheResult_WhenGivenAnExpression(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "42\n"}, run(t, "", "-e", "42"))
}

func TestRun_ExitsWithUsageError_WhenAFlagIsUnknown(t *testing.T) {
	got := run(t, "", "--nope")

	assert.Equal(t, 2, got.code)
	assert.Empty(t, got.stdout)
	assert.Contains(t, got.stderr, "-nope")
}

func TestRun_ExitsWithParseError_WhenTheExpressionIsNotValid(t *testing.T) {
	want := outcome{code: 3, stderr: "ubc: line 1, col 4: unexpected end of input\n"}

	assert.Equal(t, want, run(t, "", "-e", "1 +"))
}

func TestRun_ExitsWithEvaluationError_WhenTheExpressionCannotBeEvaluated(t *testing.T) {
	want := outcome{code: 1, stderr: "ubc: line 1, col 4: range start 10 must be below its end 5\n"}

	assert.Equal(t, want, run(t, "", "-e", "10 to 5"))
}

func TestRun_SummarisesASampleSetAsMedianAndNinetyPercentInterval(t *testing.T) {
	got := run(t, "", "-e", "10 to 20")

	assert.Equal(t, 0, got.code)
	assert.Regexp(t, `^≈ \S+   90% \[\S+ … \S+\]\n$`, got.stdout)
}

func TestRun_PrintsTheVersion_WhenAskedFor(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "ubc v1.2.3\n"}, run(t, "", "--version"))
}

func TestRun_PrintsEveryExpressionButNoAssignment(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "6\n3\n"}, run(t, "", "-e", "x = 2; x * 3; y = x; y + 1"))
}

func TestRun_PrintsAPlainNumber_WhenEverySampleIsTheSame(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "0\n"}, run(t, "", "-e", "x = 1 to 5; x - x"))
}

func TestRun_PrintsTheResultsBeforeAnError(t *testing.T) {
	want := outcome{code: 1, stdout: "1\n", stderr: "ubc: line 1, col 7: range start 10 must be below its end 5\n"}

	assert.Equal(t, want, run(t, "", "-e", "1; 10 to 5; 3"))
}

func TestRun_EvaluatesRepeatedExpressionsInOneSession(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "2\n6\n"}, run(t, "", "-e", "x = 2; x", "-e", "x * 3"))
}

func TestRun_StopsAtTheFirstFailingExpression(t *testing.T) {
	got := run(t, "", "-e", "1", "-e", "y", "-e", "3")

	assert.Equal(t, outcome{code: 1, stdout: "1\n", stderr: "ubc: line 1, col 1: unknown variable 'y'\n"}, got)
}

func TestRun_ReadsTheProgramFromStdin_WhenThereAreNoExpressionsOrFiles(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "3\n"}, run(t, "1 + 2\n"))
}

func TestRun_PrefixesTheLine_WhenAProgramPrintsMoreThanOneResult(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "L1: 1\nL3: 3\n"}, run(t, "1\nx = 2\nx + 1\n"))
}

func TestRun_EvaluatesFilesInOrderInOneSession(t *testing.T) {
	first := writeFile(t, "first.ubc", "x = 2\n")
	second := writeFile(t, "second.ubc", "x + 1\n")

	assert.Equal(t, outcome{code: 0, stdout: "3\n"}, run(t, "ignored", first, second))
}

func TestRun_NamesTheFile_WhenAFileHasAnError(t *testing.T) {
	path := writeFile(t, "model.ubc", "1\ny\n")

	want := outcome{code: 1, stdout: "1\n", stderr: "ubc: " + path + ": line 2, col 1: unknown variable 'y'\n"}
	assert.Equal(t, want, run(t, "", path))
}

func TestRun_ExitsWithEvaluationError_WhenAFileCannotBeRead(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.ubc")

	got := run(t, "", missing)

	assert.Equal(t, 1, got.code)
	assert.Regexp(t, `^ubc: .*missing\.ubc.*\n$`, got.stderr)
}

func TestRun_ExitsWithEvaluationError_WhenStdinCannotBeRead(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := cli.App{}.Run(nil, failingReader{}, &stdout, &stderr)

	assert.Equal(t, 1, code)
	assert.Equal(t, "ubc: stdin: broken pipe\n", stderr.String())
}

type failingReader struct{}

func (failingReader) Read([]byte) (int, error) { return 0, errors.New("broken pipe") }

func TestRun_AcceptsTheLongEvalFlag(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "3\n"}, run(t, "", "--eval", "1 + 2"))
}

func TestRun_EvaluatesExpressionsBeforeFiles(t *testing.T) {
	path := writeFile(t, "model.ubc", "x * 10\n")

	assert.Equal(t, outcome{code: 0, stdout: "20\n"}, run(t, "", "-e", "x = 2", path))
}
