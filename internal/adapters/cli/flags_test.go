package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eferro/uncertain-bc/internal/adapters/cli"
)

func TestRun_DrawsTheRequestedNumberOfSamples(t *testing.T) {
	got := run(t, "", "-n", "1", "-e", "1 to 2")

	assert.Equal(t, 0, got.code)
	assert.Regexp(t, `^1\.\d+\n$`, got.stdout)
	assert.Equal(t, got, run(t, "", "--samples", "1", "-e", "1 to 2"))
}

func TestRun_ExitsWithUsageError_WhenTheSampleCountIsNotPositive(t *testing.T) {
	for _, n := range []string{"0", "-5"} {
		want := outcome{code: 2, stderr: "ubc: --samples must be positive, got " + n + "\n"}

		assert.Equal(t, want, run(t, "", "-n", n, "-e", "1"))
	}
}

func TestRun_UsesSeedOne_ByDefault(t *testing.T) {
	assert.Equal(t, run(t, "", "-s", "1", "-e", "1 to 2"), run(t, "", "-e", "1 to 2"))
}

func TestRun_ChangesTheSamples_WhenTheSeedChanges(t *testing.T) {
	assert.NotEqual(t, run(t, "", "-e", "1 to 2").stdout, run(t, "", "--seed", "2", "-e", "1 to 2").stdout)
}

func TestRun_SeedsFromTheClock_WhenRandomIsRequested(t *testing.T) {
	var stdout, stderr bytes.Buffer
	app := cli.App{RandomSeed: func() uint64 { return 2 }}

	code := app.Run([]string{"--random", "-e", "1 to 2"}, strings.NewReader(""), &stdout, &stderr)

	assert.Equal(t, 0, code)
	assert.Equal(t, run(t, "", "-s", "2", "-e", "1 to 2").stdout, stdout.String())
}

func TestRun_ExitsWithUsageError_WhenSeedAndRandomAreCombined(t *testing.T) {
	want := outcome{code: 2, stderr: "ubc: --seed and --random cannot be combined\n"}

	for _, seedFlag := range []string{"-s", "--seed"} {
		assert.Equal(t, want, run(t, "", seedFlag, "3", "--random", "-e", "1"), seedFlag)
	}
}

func TestRun_PrintsOnlyAscii_WhenAsked(t *testing.T) {
	got := run(t, "", "--ascii", "-e", "10 to 20")

	assert.Equal(t, 0, got.code)
	assert.Regexp(t, `^~ \S+   90% \[\S+ \.\.\. \S+\]   [ .:\-=+*#]{16}\n$`, got.stdout)
}

func TestRun_PrintsUsage_WhenHelpIsRequested(t *testing.T) {
	for _, flag := range []string{"-h", "--help"} {
		got := run(t, "", flag)

		assert.Equal(t, 0, got.code, flag)
		assert.Contains(t, got.stdout, "Usage:", flag)
		assert.Contains(t, got.stdout, "-n, --samples N", flag)
		assert.Empty(t, got.stderr, flag)
	}
}

func TestRun_PrintsUsageToStderr_WhenAFlagIsWrong(t *testing.T) {
	got := run(t, "", "--nope")

	assert.Equal(t, 2, got.code)
	assert.Contains(t, got.stderr, "flag provided but not defined: -nope")
	assert.Contains(t, got.stderr, "Usage:")
}

func TestRun_PrintsOnlyTheCentralValueWithFullPrecision_WhenPlain(t *testing.T) {
	got := run(t, "1 + 2\nx = 10 to 20\nx\nP(x > 15)\n", "--plain")

	assert.Equal(t, 0, got.code)
	assert.Regexp(t, `^3\n14\.\d{6,}\n0\.\d+\n$`, got.stdout)
}

func TestRun_PrintsACompletionScript(t *testing.T) {
	got := run(t, "", "--completion", "bash")

	assert.Equal(t, 0, got.code)
	assert.Contains(t, got.stdout, "complete -o filenames -F _ubc ubc")
}

func TestRun_ExitsWithUsageError_WhenTheShellIsUnknown(t *testing.T) {
	want := outcome{code: 2, stderr: "ubc: --completion expects bash, zsh or fish, got 'tcsh'\n"}

	assert.Equal(t, want, run(t, "", "--completion", "tcsh"))
}
