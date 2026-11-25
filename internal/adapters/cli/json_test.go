package cli_test

import (
	"encoding/json"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func decodeLines(t *testing.T, stdout string) []map[string]any {
	t.Helper()
	var objects []map[string]any
	for _, line := range strings.Split(strings.TrimSuffix(stdout, "\n"), "\n") {
		var object map[string]any
		require.NoError(t, json.Unmarshal([]byte(line), &object), line)
		objects = append(objects, object)
	}
	return objects
}

func TestJSON_DescribesAPlainNumber(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: `{"kind":"scalar","value":4}` + "\n"}, run(t, "", "--json", "-e", "2 + 2"))
}

func TestJSON_DescribesAProbability(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: `{"kind":"probability","value":1}` + "\n"}, run(t, "", "--json", "-e", "P(2 > 1)"))
}

func TestJSON_DescribesAnAssignedComparison_AsAProbability(t *testing.T) {
	got := run(t, "", "--json", "--echo", "-e", "hit = 0 to 10 > 5")

	require.Equal(t, 0, got.code)
	object := decodeLines(t, got.stdout)[0]
	assert.Equal(t, "hit", object["name"])
	assert.Equal(t, "probability", object["kind"])
	assert.InDelta(t, 0.5, object["value"], 0.03)
}

func TestJSON_SummarisesASampleSet(t *testing.T) {
	got := run(t, "", "--json", "-e", "10 to 20")

	require.Equal(t, 0, got.code)
	assert.Regexp(t, `^\{"kind":"samples","median":[\d.]+,"mean":[\d.]+,"stdev":[\d.]+,"p1":[\d.]+,"p5":[\d.]+,"p95":[\d.]+,"p99":[\d.]+,"n":5000,"dropped":0\}\n$`, got.stdout)
	object := decodeLines(t, got.stdout)[0]
	assert.InDelta(t, 14.14, object["median"], 0.3)
	assert.InDelta(t, 10, object["p5"], 0.3)
	assert.InDelta(t, 20, object["p95"], 0.6)
}

func TestJSON_RoundsEveryStatisticToTenSignificantFigures(t *testing.T) {
	got := run(t, "", "--json", "-e", "10 to 20")

	require.Equal(t, 0, got.code)
	assert.Contains(t, got.stdout, `"median":14.15512332,`)
	for _, number := range regexp.MustCompile(`\d+\.\d+`).FindAllString(got.stdout, -1) {
		assert.LessOrEqual(t, len(strings.TrimLeft(strings.Replace(number, ".", "", 1), "0")), 10, number)
	}
}

func TestJSON_PrintsOneLinePerExpression(t *testing.T) {
	objects := decodeLines(t, run(t, "", "--json", "-e", "x = 2; x; x * 3").stdout)

	require.Len(t, objects, 2)
	assert.Equal(t, 6.0, objects[1]["value"])
}

func TestJSON_CountsDroppedSamplesAndKeepsSummariesFinite(t *testing.T) {
	object := decodeLines(t, run(t, "", "--json", "-e", "x = 0 to 1; x / (x > 0.5)").stdout)[0]

	assert.Greater(t, object["dropped"], 1000.0)
	assert.IsType(t, 0.0, object["median"])
}

func TestJSON_UsesNullForValuesThatAreNotFinite(t *testing.T) {
	assert.Equal(t, `{"kind":"scalar","value":null}`+"\n", run(t, "", "--json", "-e", "1 / 0").stdout)
	assert.Equal(t, `{"kind":"samples","median":null,"mean":null,"stdev":null,"p1":null,"p5":null,"p95":null,"p99":null,"n":5000,"dropped":5000}`+"\n",
		run(t, "", "--json", "-e", "(1 to 2) * 0 / 0").stdout)
}

func TestJSON_CannotBeCombinedWithPlain(t *testing.T) {
	want := outcome{code: 2, stderr: "ubc: --json and --plain cannot be combined\n"}

	assert.Equal(t, want, run(t, "", "--json", "--plain", "-e", "1"))
}

func TestJSON_StaysParseable_WhenAProgramPrintsSeveralResults(t *testing.T) {
	got := run(t, "1\n2\n", "--json")

	assert.Equal(t, `{"name":"1","kind":"scalar","value":1}`+"\n"+`{"name":"2","kind":"scalar","value":2}`+"\n", got.stdout)
}
