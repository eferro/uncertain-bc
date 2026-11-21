package cli_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHistFlag_DrawsAHistogramUnderEveryUncertainResult(t *testing.T) {
	got := run(t, "", "--hist", "-e", "10 to 20")

	lines := strings.Split(strings.TrimSuffix(got.stdout, "\n"), "\n")
	require.Len(t, lines, 22)
	assert.Regexp(t, `^≈ `, lines[0])
	assert.Regexp(t, `^\s+\S+ ┤█*$`, lines[1])
	assert.Regexp(t, `^\(P1 to P99 shown`, lines[21])
}

func TestHistFlag_LeavesPlainNumbersAlone(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "4\n"}, run(t, "", "--hist", "-e", "2 + 2"))
}

func TestHistFlag_CannotBeCombinedWithMachineOutput(t *testing.T) {
	for _, other := range []string{"--json", "--plain"} {
		want := outcome{code: 2, stderr: "ubc: --hist cannot be combined with --json or --plain\n"}

		assert.Equal(t, want, run(t, "", "--hist", other, "-e", "1"), other)
	}
}

func TestQuantiles_ReplaceTheNinetyPercentInterval(t *testing.T) {
	got := run(t, "", "--quantiles", "10,50,90", "-e", "10 to 20")

	assert.Regexp(t, `^≈ \S+   p10 \S+   p50 \S+   p90 \S+   [▁▂▃▄▅▆▇█]{16}\n$`, got.stdout)
}

func TestQuantiles_AcceptFractionalPercentages(t *testing.T) {
	assert.Regexp(t, `   p2\.5 \S+   p97\.5 \S+   `, run(t, "", "--quantiles", "2.5,97.5", "-e", "10 to 20").stdout)
}

func TestQuantiles_AreAddedToJSON(t *testing.T) {
	object := decodeLines(t, run(t, "", "--json", "--quantiles", "10,90", "-e", "10 to 20").stdout)[0]

	quantiles, ok := object["quantiles"].(map[string]any)
	require.True(t, ok)
	assert.InDelta(t, 10.8, quantiles["p10"], 0.4)
	assert.InDelta(t, 18.5, quantiles["p90"], 0.6)
	assert.Contains(t, object, "p5")
}

func TestQuantiles_RejectWhatIsNotAPercentage(t *testing.T) {
	for _, list := range []string{"101", "-1", "x", "10,,90"} {
		want := outcome{code: 2, stderr: "ubc: --quantiles expects percentages between 0 and 100, got '" + list + "'\n"}

		assert.Equal(t, want, run(t, "", "--quantiles", list, "-e", "1"), list)
	}
}

func TestHistFlag_LeavesCertainSampleSetsAlone(t *testing.T) {
	assert.Equal(t, outcome{code: 0, stdout: "0\n"}, run(t, "", "--hist", "-e", "x = 1 to 5; x - x"))
}

func TestHistFlag_WorksInTheRepl(t *testing.T) {
	assert.Regexp(t, `^ubc> ≈ .+\n\s+\S+ ┤█`, interactive(t, "10 to 20\n", "--hist").stdout)
}
