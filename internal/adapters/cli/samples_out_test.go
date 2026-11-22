package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/adapters/cli"
)

func readOutput(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(path)
	require.NoError(t, err)
	return string(content)
}

func TestSamplesOut_WritesOneCSVColumnPerVariable_AndRepeatsScalars(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.csv")

	got := run(t, "", "-n", "3", "--samples-out", path, "-e", "x = [2, 2.5]", "-e", "k = 1e21", "-e", "x * 2")

	assert.Equal(t, 0, got.code)
	assert.Equal(t, "x,k,last\n2,1e+21,4\n2.5,1e+21,5\n2,1e+21,4\n", readOutput(t, path))
}

func TestSamplesOut_WritesAJSONObjectOfSampleArrays_InVariableOrder(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.json")

	got := run(t, "", "-n", "2", "--samples-out", path, "-e", "z = 3", "-e", "a = P(1 < 2)")

	assert.Equal(t, 0, got.code)
	assert.Equal(t, "{\"z\":[3,3],\"a\":[1,1]}\n", readOutput(t, path))
}

func TestSamplesOut_WritesNonFiniteSamplesAsMissingValues(t *testing.T) {
	dir := t.TempDir()
	csv, json := filepath.Join(dir, "out.csv"), filepath.Join(dir, "out.json")
	program := []string{"-n", "1", "-e", "z = 0", "-e", "1 / z", "-e", "inf = last", "-e", "z / z"}

	require.Equal(t, 0, run(t, "", append([]string{"--samples-out", csv}, program...)...).code)
	require.Equal(t, 0, run(t, "", append([]string{"--samples-out", json}, program...)...).code)

	assert.Equal(t, "z,inf,last\n0,,\n", readOutput(t, csv))
	assert.Equal(t, "{\"z\":[0],\"inf\":[null],\"last\":[null]}\n", readOutput(t, json))
}

func TestSamplesOut_IncludesDataAndSetOverrides(t *testing.T) {
	data := writeFile(t, "times.txt", "4\n")
	path := filepath.Join(t.TempDir(), "out.csv")

	got := run(t, "", "-n", "1", "--samples-out", path, "--data", "lead="+data, "--set", "k=2", "-e", "lead * k")

	assert.Equal(t, 0, got.code)
	assert.Equal(t, "lead,k,last\n4,2,8\n", readOutput(t, path))
}

func TestSamplesOut_IsWrittenEvenWhenAnAssertionFails(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.csv")

	got := run(t, "", "-n", "1", "--samples-out", path, "--assert", "x > 5", "-e", "x = 1")

	assert.Equal(t, 4, got.code)
	assert.Equal(t, "x\n1\n", readOutput(t, path))
}

func TestSamplesOut_WritesNothing_WhenTheProgramFails(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.csv")

	got := run(t, "", "--samples-out", path, "-e", "x = 1", "-e", "y")

	assert.Equal(t, 1, got.code)
	assert.NoFileExists(t, path)
}

func TestSamplesOut_ExitsWithEvalError_WhenThereIsNothingToExport(t *testing.T) {
	path := filepath.Join(t.TempDir(), "out.json")

	got := run(t, "", "--samples-out", path, "-e", "f(x) = x")

	assert.Equal(t, outcome{code: 1, stderr: "ubc: --samples-out " + path + ": nothing to export\n"}, got)
	assert.NoFileExists(t, path)
}

func TestSamplesOut_ExitsWithEvalError_WhenTheFileCannotBeWritten(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing", "out.csv")

	got := run(t, "", "--samples-out", path, "-e", "1")

	assert.Equal(t, 1, got.code)
	assert.Equal(t, "ubc: --samples-out "+path+": "+missingFileError(t, path)+"\n", got.stderr)
}

func TestSamplesOut_ExitsWithUsageError_WhenTheExtensionIsNotCSVOrJSON(t *testing.T) {
	for _, path := range []string{"out.txt", "out", "out.CSV"} {
		t.Run(path, func(t *testing.T) {
			want := "ubc: --samples-out writes .csv or .json files, got '" + path + "'\n"

			assert.Equal(t, outcome{code: 2, stderr: want}, run(t, "", "--samples-out", path, "-e", "1"))
		})
	}
}

func TestSamplesOut_ExitsWithUsageError_InTheREPL(t *testing.T) {
	got := interactive(t, "", "--samples-out", "out.csv")

	assert.Equal(t, outcome{code: 2, stderr: "ubc: --samples-out needs a program: -e, files or stdin\n"}, got)
}

func TestSamplesOut_IsWrittenAgainOnEveryWatchRun(t *testing.T) {
	model := writeFile(t, "model.ubc", "x = 4\n")
	path := filepath.Join(filepath.Dir(model), "out.csv")
	rewrite(t, model, "x = 4\n", 0)
	edited := false
	app := cli.App{Wait: func() bool {
		if !edited {
			assert.Equal(t, "x\n4\n", readOutput(t, path))
			rewrite(t, model, "x = 9\n", 1)
		}
		edited = !edited
		return edited
	}}
	var stdout, stderr bytes.Buffer

	code := app.Run([]string{"-n", "1", "--samples-out", path, "--watch", model}, strings.NewReader(""), &stdout, &stderr)

	assert.Equal(t, 0, code)
	assert.Equal(t, "x\n9\n", readOutput(t, path))
}
