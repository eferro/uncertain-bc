package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eferro/uncertain-bc/internal/adapters/cli"
)

func TestSet_OverridesTheModelsAssignment(t *testing.T) {
	path := writeFile(t, "model.ubc", "team = 3\nteam * 10\n")

	assert.Equal(t, outcome{code: 0, stdout: "50\n"}, run(t, "", "--set", "team=5", path))
}

func TestSet_AppliesToStdin(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "50\n"}, run(t, "team = 3\nteam * 10\n", "--set", "team=5"))
}

func TestSet_TrimsSpacesAroundTheName(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "10\n"}, run(t, "", "--set", " team = 5", "-e", "team * 2"))
}

func TestSet_SplitsAtTheFirstEqualsSign(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "1\n"}, run(t, "", "--set", "ok=1==1", "-e", "ok"))
}

func TestSet_PrintsNothingByItself(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "x = 2\n"}, run(t, "", "--set", "team=5", "--echo", "-e", "x = 2"))
}

func TestSet_ExitsWithUsageError_WhenMalformed(t *testing.T) {
	cases := map[string]string{
		"team":    "ubc: --set expects name=expr, got 'team'\n",
		"=5":      "ubc: --set: '' is not a variable name\n",
		"2x=5":    "ubc: --set: '2x' is not a variable name\n",
		"a b=5":   "ubc: --set: 'a b' is not a variable name\n",
		"to=5":    "ubc: --set: 'to' is not a variable name\n",
		"team;=5": "ubc: --set: 'team;' is not a variable name\n",
	}
	for value, want := range cases {
		t.Run(value, func(t *testing.T) {
			assert.Equal(t, outcome{code: 2, stderr: want}, run(t, "", "--set", value, "-e", "1"))
		})
	}
}

func TestSet_ExitsWithParseError_WhenTheExpressionIsNotValid(t *testing.T) {
	want := outcome{code: 3, stderr: "ubc: --set team: line 1, col 4: unexpected end of input\n"}

	assert.Equal(t, want, run(t, "", "--set", "team=1 +", "-e", "1"))
}

func TestSet_ExitsWithEvaluationError_WhenTheExpressionFails(t *testing.T) {
	want := outcome{code: 1, stderr: "ubc: --set team: line 1, col 4: range start 10 must be below its end 5\n"}

	assert.Equal(t, want, run(t, "", "--set", "team=10 to 5", "-e", "1"))
}

func TestSet_StopsAtTheFirstFailingOverride(t *testing.T) {
	got := run(t, "", "--set", "a=nobody", "--set", "b=nobody", "-e", "1")

	assert.Equal(t, outcome{code: 1, stderr: "ubc: --set a: line 1, col 1: unknown variable 'nobody'\n"}, got)
}

func TestSet_AppliesToTheRepl(t *testing.T) {
	want := outcome{code: 0, stdout: "ubc> team = 5\nubc> 10\nubc> \n"}

	assert.Equal(t, want, interactive(t, "team = 3\nteam * 2\n", "--set", "team=5"))
}

func TestSet_ShowsInTheReplVariables(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "ubc> team = 5\nubc> \n"}, interactive(t, ":vars\n", "--set", "team=5"))
}

func TestSet_ReportsAFailingOverride_BeforeTheRepl(t *testing.T) {
	want := outcome{code: 1, stderr: "ubc: --set team: line 1, col 1: unknown variable 'nobody'\n"}

	assert.Equal(t, want, interactive(t, "1\n", "--set", "team=nobody"))
}

func TestSet_AppliesToEveryWatchRun(t *testing.T) {
	path := writeFile(t, "model.ubc", "team = 3\nteam * 10\n")
	rewrite(t, path, "team = 3\nteam * 10\n", 0)
	polls := 0
	app := cli.App{Wait: func() bool {
		polls++
		if polls == 1 {
			rewrite(t, path, "team = 3\nteam * 100\n", polls)
		}
		return polls == 1
	}}
	var stdout, stderr bytes.Buffer

	code := app.Run([]string{"--set", "team=5", "--watch", path}, strings.NewReader(""), &stdout, &stderr)

	assert.Equal(t, 0, code)
	assert.Equal(t, "50\n--- "+path+" changed ---\n500\n", stdout.String())
}

func TestSet_IsListedInTheHelp(t *testing.T) {
	assert.Contains(t, run(t, "", "--help").stdout, "--set NAME=EXPR")
}
