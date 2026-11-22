package cli_test

import (
	"bytes"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"

	"github.com/eferro/uncertain-bc/internal/adapters/cli"
)

var catalogue = fstest.MapFS{
	"formulas.ubc":          {Data: []byte("# Reusable formulas.\ndouble(x) = 2 * x\nrate = 3\n")},
	"team/replace.ubc":      {Data: []byte("# Question: how much does it cost\n# to replace a person?\n#\n# Run it:\ncost = 10\ncost\n")},
	"team/hire.ubc":         {Data: []byte("# Question: how long does hiring take?\nweeks = double(3)\nweeks\n")},
	"delivery/date.ubc":     {Data: []byte("# Just a model, no question.\n1\n")},
	"delivery/blank.ubc":    {Data: []byte("# Question:\n#\n1\n")},
	"delivery/broken.ubc":   {Data: []byte("1 +\n")},
	"delivery/notes.txt":    {Data: []byte("not a scenario\n")},
	"risk/loss/nested.ubc":  {Data: []byte("# Question: deeper?\n2\n")},
	"risk/README.md":        {Data: []byte("# Question: not a scenario\n")},
	"product/uses-last.ubc": {Data: []byte("# Question:   spaced   out?\nlast\n")},
}

func runWithCatalogue(t *testing.T, args ...string) outcome {
	t.Helper()
	return runApp(t, cli.App{Scenarios: catalogue}, "", args...)
}

func runApp(t *testing.T, app cli.App, stdin string, args ...string) outcome {
	t.Helper()
	var stdout, stderr bytes.Buffer
	code := app.Run(args, strings.NewReader(stdin), &stdout, &stderr)
	return outcome{code: code, stdout: stdout.String(), stderr: stderr.String()}
}

func TestScenarios_ListsEveryScenarioWithItsQuestion(t *testing.T) {
	want := "" +
		"delivery/blank\n" +
		"delivery/broken\n" +
		"delivery/date\n" +
		"product/uses-last   Spaced out?\n" +
		"risk/loss/nested    Deeper?\n" +
		"team/hire           How long does hiring take?\n" +
		"team/replace        How much does it cost to replace a person?\n"

	assert.Equal(t, outcome{code: 0, stdout: want}, runWithCatalogue(t, "--scenarios"))
}

func TestScenarios_ListsNothing_WhenTheBinaryHasNoCatalogue(t *testing.T) {
	assert.Equal(t, outcome{code: 0}, run(t, "", "--scenarios"))
}

func TestScenario_RunsTheScenario(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "10\n"}, runWithCatalogue(t, "--scenario", "team/replace"))
}

func TestScenario_LoadsTheSharedFormulasFirst(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "6\n"}, runWithCatalogue(t, "--scenario", "team/hire"))
}

func TestScenario_RunsBeforeExpressionsAndFiles(t *testing.T) {
	path := writeFile(t, "model.ubc", "cost * 3\n")

	got := runWithCatalogue(t, "--echo", "--scenario", "team/replace", "-e", "cost * 2", path)

	assert.Equal(t, outcome{code: 0, stdout: "rate = 3\ncost = 10\n10\n20\n30\n"}, got)
}

func TestScenario_CombinesWithSet(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "4\n"}, runWithCatalogue(t, "--set", "cost=4", "--scenario", "team/replace"))
}

func TestScenario_Formulas_LoadsOnlyTheFormulasForYourOwnModel(t *testing.T) {
	path := writeFile(t, "model.ubc", "double(21)\n")

	assert.Equal(t, outcome{code: 0, stdout: "rate = 3\n42\n"}, runWithCatalogue(t, "--echo", "--scenario", "formulas", path))
}

func TestScenario_DoesNotReadStdin(t *testing.T) {
	got := runApp(t, cli.App{Scenarios: catalogue}, "99\n", "--scenario", "team/replace")

	assert.Equal(t, outcome{code: 0, stdout: "10\n"}, got)
}

func TestScenario_RunsInsteadOfTheRepl_WhenStdinIsATerminal(t *testing.T) {
	got := runApp(t, cli.App{Scenarios: catalogue, Interactive: true}, "99\n", "--scenario", "team/replace")

	assert.Equal(t, outcome{code: 0, stdout: "10\n"}, got)
}

func TestScenario_NamesTheScenario_WhenItFails(t *testing.T) {
	want := outcome{code: 3, stderr: "ubc: delivery/broken: line 1, col 4: unexpected end of line\n"}

	assert.Equal(t, want, runWithCatalogue(t, "--scenario", "delivery/broken"))
}

func TestScenario_ExitsWithUsageError_WhenTheNameIsUnknown(t *testing.T) {
	for _, name := range []string{"team/nobody", "team/replace.ubc", "team", "../team/replace", "delivery/notes", "risk/README"} {
		t.Run(name, func(t *testing.T) {
			want := outcome{code: 2, stderr: "ubc: unknown scenario '" + name + "'; ubc --scenarios lists them\n"}

			assert.Equal(t, want, runWithCatalogue(t, "--scenario", name, "-e", "1"))
		})
	}
}

func TestScenario_ExitsWithUsageError_WhenTheBinaryHasNoCatalogue(t *testing.T) {
	want := outcome{code: 2, stderr: "ubc: unknown scenario 'team/replace'; ubc --scenarios lists them\n"}

	assert.Equal(t, want, run(t, "", "--scenario", "team/replace"))
}

func TestScenarioSource_PrintsTheScenarioAsItIs(t *testing.T) {
	want := outcome{code: 0, stdout: string(catalogue["team/replace.ubc"].Data)}

	assert.Equal(t, want, runWithCatalogue(t, "--scenario-source", "team/replace"))
}

func TestScenarioSource_PrintsTheFormulas(t *testing.T) {
	want := outcome{code: 0, stdout: string(catalogue["formulas.ubc"].Data)}

	assert.Equal(t, want, runWithCatalogue(t, "--scenario-source", "formulas"))
}

func TestScenarioSource_ExitsWithUsageError_WhenTheNameIsUnknown(t *testing.T) {
	want := outcome{code: 2, stderr: "ubc: unknown scenario 'nobody'; ubc --scenarios lists them\n"}

	assert.Equal(t, want, runWithCatalogue(t, "--scenario-source", "nobody"))
}
