package cli_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcho_PrintsAssignmentsWithTheirNames(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "x = 2\n6\n"}, run(t, "", "--echo", "-e", "x = 2; x * 3"))
}

func TestEcho_KeepsLabelsForExpressionsInPrograms(t *testing.T) {
	path := writeFile(t, "model.ubc", "a = 2\na * 1\nb = a * 3\nb + 1\n")

	assert.Equal(t, outcome{code: 0, stdout: "a = 2\na * 1 = 2\nb = 6\nb + 1 = 7\n"}, run(t, "", "--echo", path))
}

func TestEcho_LabelsPlainValues(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "x = 2\n6\n"}, run(t, "", "--echo", "--plain", "-e", "x = 2; x * 3"))
}

func TestEcho_NamesAssignmentsInJSON(t *testing.T) {
	want := `{"name":"x","kind":"scalar","value":2}` + "\n" + `{"kind":"scalar","value":6}` + "\n"

	assert.Equal(t, outcome{code: 0, stdout: want}, run(t, "", "--echo", "--json", "-e", "x = 2; x * 3"))
}

func TestEcho_IsOffByDefault(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "6\n"}, run(t, "", "-e", "x = 2; x * 3"))
}
