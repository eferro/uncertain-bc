package cli_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eferro/uncertain-bc/internal/adapters/cli"
)

func TestData_BindsTheNumbersInTheFile_AsAnEmpiricalSampleSet(t *testing.T) {
	path := writeFile(t, "times.txt", "3\n5\n8\n13\n")

	assert.Equal(t, run(t, "", "-e", "[3, 5, 8, 13]"), run(t, "", "--data", "lead="+path, "-e", "lead"))
}

func TestData_OverridesTheModelsAssignment(t *testing.T) {
	data := writeFile(t, "times.txt", "4\n")
	model := writeFile(t, "model.ubc", "lead = 1 to 2\nmax(lead)\n")

	assert.Equal(t, outcome{code: 0, stdout: "4\n"}, run(t, "", "--data", "lead="+data, model))
}

func TestData_AppliesToStdin(t *testing.T) {
	path := writeFile(t, "times.txt", "4\n")

	assert.Equal(t, outcome{code: 0, stdout: "8\n"}, run(t, "max(lead) * 2\n", "--data", "lead="+path))
}

func TestData_IgnoresBlankLinesAndSpaces_AndAcceptsWindowsLineEnds(t *testing.T) {
	path := writeFile(t, "times.txt", "\r\n  -2.5 \r\n\t\r\n1e3\r\n")

	got := run(t, "", "--data", "lead="+path, "-e", "min(lead)", "-e", "max(lead)")

	assert.Equal(t, outcome{code: 0, stdout: "-2.5\n1000\n"}, got)
}

func TestData_IsLoadedBeforeSet_SoOverridesCanUseIt(t *testing.T) {
	path := writeFile(t, "times.txt", "4\n")

	got := run(t, "", "--set", "worst=max(lead) * 2", "--data", "lead="+path, "-e", "worst")

	assert.Equal(t, outcome{code: 0, stdout: "8\n"}, got)
}

func TestData_LosesToSetForTheSameName_AndTheLastDataWins(t *testing.T) {
	first := writeFile(t, "first.txt", "4\n")
	second := writeFile(t, "second.txt", "6\n")

	assert.Equal(t, outcome{code: 0, stdout: "6\n"}, run(t, "", "--data", "lead="+first, "--data", "lead="+second, "-e", "max(lead)"))
	assert.Equal(t, outcome{code: 0, stdout: "1\n"}, run(t, "", "--set", "lead=1", "--data", "lead="+first, "-e", "lead"))
}

func TestData_TrimsSpacesAroundTheName(t *testing.T) {
	path := writeFile(t, "times.txt", "4\n")

	assert.Equal(t, outcome{code: 0, stdout: "4\n"}, run(t, "", "--data", " lead ="+path, "-e", "max(lead)"))
}

func TestData_ExitsWithUsageError_WhenMalformed(t *testing.T) {
	cases := map[string]string{
		"lead":       "ubc: --data expects name=file, got 'lead'\n",
		"=times.txt": "ubc: --data: '' is not a variable name\n",
		"2x=a.txt":   "ubc: --data: '2x' is not a variable name\n",
		"to=a.txt":   "ubc: --data: 'to' is not a variable name\n",
		"lead=":      "ubc: --data lead: no file name\n",
	}
	for value, want := range cases {
		t.Run(value, func(t *testing.T) {
			assert.Equal(t, outcome{code: 2, stderr: want}, run(t, "", "--data", value, "-e", "1"))
		})
	}
}

func TestData_ExitsWithEvalError_WhenALineIsNotANumber(t *testing.T) {
	cases := map[string]string{
		"3\n3 days\n": "line 2: '3 days' is not a number",
		"1k\n":        "line 1: '1k' is not a number",
		"\nNaN\n":     "line 2: 'NaN' is not a number",
		"inf\n":       "line 1: 'inf' is not a number",
		"1e400\n":     "line 1: '1e400' is not a number",
		"3,4\n":       "line 1: '3,4' is not a number",
	}
	for content, want := range cases {
		t.Run(content, func(t *testing.T) {
			path := writeFile(t, "bad.txt", content)

			got := run(t, "", "--data", "lead="+path, "-e", "1")

			assert.Equal(t, outcome{code: 1, stderr: "ubc: --data lead: " + path + ": " + want + "\n"}, got)
		})
	}
}

func TestData_ExitsWithEvalError_WhenTheFileHasNoNumbers(t *testing.T) {
	path := writeFile(t, "empty.txt", "\n  \n")

	got := run(t, "", "--data", "lead="+path, "-e", "1")

	assert.Equal(t, outcome{code: 1, stderr: "ubc: --data lead: " + path + ": no values to resample\n"}, got)
}

func TestData_ExitsWithEvalError_WhenTheFileCannotBeRead(t *testing.T) {
	got := run(t, "", "--data", "lead=missing.txt", "-e", "1")

	assert.Equal(t, outcome{code: 1, stderr: "ubc: --data lead: " + missingFileError(t, "missing.txt") + "\n"}, got)
}

func TestData_StopsAtTheFirstFailingFile_BeforeAnyProgramRuns(t *testing.T) {
	good := writeFile(t, "good.txt", "4\n")

	got := run(t, "", "--data", "lead=missing.txt", "--data", "size="+good, "-e", "1")

	assert.Equal(t, 1, got.code)
	assert.Empty(t, got.stdout)
	assert.Equal(t, "ubc: --data lead: "+missingFileError(t, "missing.txt")+"\n", got.stderr)
}

func TestData_IsReadAgainOnEveryWatchRun(t *testing.T) {
	data := writeFile(t, "times.txt", "4\n")
	model := writeFile(t, "model.ubc", "max(lead)\n")
	rewrite(t, model, "max(lead)\n", 0)
	edited := false
	app := cli.App{Wait: func() bool {
		if !edited {
			rewrite(t, data, "9\n", 1)
			rewrite(t, model, "max(lead) \n", 1)
		}
		edited = !edited
		return edited
	}}
	var stdout, stderr bytes.Buffer

	code := app.Run([]string{"--data", "lead=" + data, "--watch", model}, strings.NewReader(""), &stdout, &stderr)

	assert.Equal(t, 0, code)
	assert.Equal(t, "4\n--- "+model+" changed ---\n9\n", stdout.String())
}
