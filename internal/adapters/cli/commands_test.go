package cli_test

import (
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	assert.Equal(t, withoutPrompts(fresh.stdout), withoutPrompts(reseeded.stdout))
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

func TestSave_WritesTheSuccessfulStatementsToAFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "model.ubc")

	got := interactive(t, "x = 2\ny\nx * 3\n:save "+path+"\n")

	saved, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "x = 2\nx * 3\n", string(saved))
	assert.Equal(t, 0, got.code)
}

func TestLoad_EvaluatesAFileIntoTheSession(t *testing.T) {
	path := writeFile(t, "model.ubc", "a = 4\na + 1\n")

	assert.Equal(t, outcome{code: 0, stdout: "ubc> 5\nubc> 4\nubc> \n"}, interactive(t, ":load "+path+"\na\n"))
}

func TestLoad_NamesTheFile_WhenItHasAnError(t *testing.T) {
	path := writeFile(t, "model.ubc", "1\ny\n")

	got := interactive(t, ":load "+path+"\n")

	assert.Equal(t, "ubc> 1\nubc> \n", got.stdout)
	assert.Equal(t, "ubc: "+path+": line 2, col 1: unknown variable 'y'\n", got.stderr)
}

func TestSaveAndLoad_ReportFileProblems(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "missing.ubc")
	cases := map[string]string{
		":save":                     "ubc: :save expects a file name\n",
		":load":                     "ubc: :load expects a file name\n",
		":load " + missing:          "ubc: open " + missing + ": no such file or directory\n",
		":save " + missing + "/x/y": "ubc: open " + missing + "/x/y: no such file or directory\n",
	}
	for input, message := range cases {
		t.Run(input, func(t *testing.T) {
			assert.Equal(t, message, interactive(t, input+"\n").stderr)
		})
	}
}

func withoutPrompts(screen string) string {
	return strings.ReplaceAll(screen, "ubc> ", "")
}

func TestPlot_WritesAnSVGHistogramOfAnExpression(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cost.svg")

	got := interactive(t, "x = 10 to 20\n:plot x * 2 "+path+"\n")

	require.Empty(t, got.stderr)
	svg, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Contains(t, string(svg), "<svg")
	assert.Contains(t, string(svg), ">x * 2</text>")
}

func TestPlot_WritesAPNGHistogramOfAnExpression(t *testing.T) {
	path := filepath.Join(t.TempDir(), "cost.png")

	got := interactive(t, "x = 10 to 20\n:plot x "+path+"\n")

	require.Empty(t, got.stderr)
	file, err := os.Open(path)
	require.NoError(t, err)
	defer func() { _ = file.Close() }()
	img, err := png.Decode(file)
	require.NoError(t, err)
	assert.Equal(t, 640, img.Bounds().Dx())
}

func TestPlot_ReportsWhatItCannotDo(t *testing.T) {
	dir := t.TempDir()
	cases := map[string]string{
		":plot":                             "ubc: :plot expects an expression and a file name\n",
		":plot x.svg":                       "ubc: :plot expects an expression and a file name\n",
		":plot 1 to 2 " + dir + "/c.gif":    "ubc: :plot writes .svg or .png files, got '" + dir + "/c.gif'\n",
		":plot 5 " + dir + "/c.svg":         "ubc: :plot needs an uncertain value, got a plain number\n",
		":plot y " + dir + "/c.svg":         "ubc: line 1, col 1: unknown variable 'y'\n",
		":plot 1 to 2 " + dir + "/no/c.svg": "ubc: open " + dir + "/no/c.svg: no such file or directory\n",
	}
	for input, message := range cases {
		t.Run(input, func(t *testing.T) {
			assert.Equal(t, message, interactive(t, input+"\n").stderr)
		})
	}
}
