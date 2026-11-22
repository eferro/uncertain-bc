package cli_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLabels_NameAVariableResultAfterTheVariable(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "x = 2\ny = 3\n"}, run(t, "x = 2\ny = 3\nx\ny\n"))
}

func TestLabels_NameAShortExpressionAfterItsSourceText(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "x + 1 = 3\nP(x  >  1) = 1\n"}, run(t, "x = 2\nx + 1   # note\nP(x  >  1)\n"))
}

func TestLabels_NameAnExpressionOfThirtyCharacters(t *testing.T) {
	expression := "10" + strings.Repeat(" + 0", 7)
	require.Len(t, expression, 30)

	assert.Equal(t, outcome{code: 0, stdout: "1 = 1\n" + expression + " = 10\n"}, run(t, "1\n"+expression+"\n"))
}

func TestLabels_KeepTheLinePrefix_WhenTheExpressionIsLongerThanThirtyCharacters(t *testing.T) {
	expression := "100" + strings.Repeat(" + 0", 7)
	require.Len(t, expression, 31)

	assert.Equal(t, outcome{code: 0, stdout: "1 = 1\nL2: 100\n"}, run(t, "1\n"+expression+"\n"))
}

func TestLabels_NameAVariableResult_WhenItsNameIsLongerThanThirtyCharacters(t *testing.T) {
	name := "a" + strings.Repeat("_long_name", 4)

	assert.Equal(t, outcome{code: 0, stdout: "1 = 1\n" + name + " = 2\n"}, run(t, "1\n"+name+" = 2\n"+name+"\n"))
}

func TestLabels_NameTheResultsOfAnExpressionFlag_WhenItPrintsSeveral(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "x = 2\nx * 3 = 6\n"}, run(t, "", "-e", "x = 2; x; x * 3"))
}

func TestLabels_LeaveASingleResultBare(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "6\n"}, run(t, "x = 2\nx * 3\n"))
}

func TestLabels_AreTheNameFieldInJSON(t *testing.T) {
	want := `{"name":"x","kind":"scalar","value":2}` + "\n" + `{"kind":"scalar","value":100}` + "\n"

	assert.Equal(t, outcome{code: 0, stdout: want}, run(t, "x = 2\nx\n100"+strings.Repeat(" + 0", 7)+"\n", "--json"))
}

func TestLabels_AreLeftOutOfPlainOutput(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "2\n6\n"}, run(t, "x = 2\nx\nx * 3\n", "--plain"))
}

func TestLabels_KeepComparisonsReadableInJSON(t *testing.T) {
	want := `{"name":"1 < 2","kind":"probability","value":1}` + "\n" + `{"name":"P(2 > 1)","kind":"probability","value":1}` + "\n"

	assert.Equal(t, outcome{code: 0, stdout: want}, run(t, "1 < 2\nP(2 > 1)\n", "--json"))
}

func TestLabels_KeepTheLinePrefix_WhenALongExpressionIsANumber(t *testing.T) {
	number := strings.Repeat("0", 30) + "7"
	require.Len(t, number, 31)

	assert.Equal(t, outcome{code: 0, stdout: "1 = 1\nL2: 7\n"}, run(t, "1\n"+number+"\n"))
}
