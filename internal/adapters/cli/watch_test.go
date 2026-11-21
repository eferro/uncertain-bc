package cli_test

import (
	"bytes"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/adapters/cli"
)

func rewrite(t *testing.T, path, content string, age int) {
	t.Helper()
	require.NoError(t, os.WriteFile(path, []byte(content), 0o600))
	stamp := time.Date(2026, 1, 1, 0, 0, age, 0, time.UTC)
	require.NoError(t, os.Chtimes(path, stamp, stamp))
}

func watching(t *testing.T, path string, edits ...string) outcome {
	t.Helper()
	polls := 0
	app := cli.App{Wait: func() bool {
		if polls < len(edits) && edits[polls] != "" {
			rewrite(t, path, edits[polls], polls+1)
		}
		polls++
		return polls <= len(edits)
	}}
	var stdout, stderr bytes.Buffer
	code := app.Run([]string{"--watch", path}, strings.NewReader(""), &stdout, &stderr)
	return outcome{code: code, stdout: stdout.String(), stderr: stderr.String()}
}

func TestWatch_ReEvaluatesTheFileWhenItChanges(t *testing.T) {
	path := writeFile(t, "model.ubc", "1 + 1\n")
	rewrite(t, path, "1 + 1\n", 0)

	got := watching(t, path, "2 + 2\n")

	assert.Equal(t, "2\n--- "+path+" changed ---\n4\n", got.stdout)
}

func TestWatch_StaysQuiet_WhileNothingChanges(t *testing.T) {
	path := writeFile(t, "model.ubc", "1 + 1\n")

	assert.Equal(t, "2\n", watching(t, path, "", "").stdout)
}

func TestWatch_KeepsWatchingAfterAnError(t *testing.T) {
	path := writeFile(t, "model.ubc", "y\n")
	rewrite(t, path, "y\n", 0)

	got := watching(t, path, "3\n")

	assert.Equal(t, "--- "+path+" changed ---\n3\n", got.stdout)
	assert.Equal(t, "ubc: "+path+": line 1, col 1: unknown variable 'y'\n", got.stderr)
	assert.Equal(t, 0, got.code)
}

func TestWatch_NeedsAFile(t *testing.T) {
	assert.Equal(t, outcome{code: 2, stderr: "ubc: --watch needs at least one file\n"}, run(t, "", "--watch", "-e", "1"))
}

func TestWatch_NoticesAChangeOfSize_WhenTheTimestampIsTheSame(t *testing.T) {
	path := writeFile(t, "model.ubc", "1 + 1\n")
	rewrite(t, path, "1 + 1\n", 0)
	polls := 0
	app := cli.App{Wait: func() bool {
		polls++
		if polls == 1 {
			rewrite(t, path, "20 + 20\n", 0)
		}
		return polls <= 1
	}}
	var stdout, stderr bytes.Buffer

	app.Run([]string{"--watch", path}, strings.NewReader(""), &stdout, &stderr)

	assert.Equal(t, "2\n--- "+path+" changed ---\n40\n", stdout.String())
}

func TestWatch_EvaluatesEachChangeOnce(t *testing.T) {
	path := writeFile(t, "model.ubc", "1 + 1\n")
	rewrite(t, path, "1 + 1\n", 0)

	got := watching(t, path, "2 + 2\n", "", "")

	assert.Equal(t, "2\n--- "+path+" changed ---\n4\n", got.stdout)
}
