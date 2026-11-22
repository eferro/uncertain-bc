package cli_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInterval_SetsTheWidthOfThePrintedInterval(t *testing.T) {
	got := run(t, "", "--interval", "80", "-e", "10 to 20")

	assert.Regexp(t, `^≈ 14\.2   80% \[10\.8 … 18\.4\]   98% \[8\.6 … 24\]   [▁▂▃▄▅▆▇█]{16}\n$`, got.stdout)
}

func TestInterval_AcceptsFractionalWidths(t *testing.T) {
	assert.Regexp(t, `^≈ 14\.2   99\.5% \[8\.03 … 26\.1\]   `, run(t, "", "--interval", "99.5", "-e", "10 to 20").stdout)
}

func TestInterval_DefaultsToNinetyPercent(t *testing.T) {
	assert.Equal(t, run(t, "", "-e", "10 to 20"), run(t, "", "--interval", "90", "-e", "10 to 20"))
}

func TestInterval_AppliesToAssignmentsInTheRepl(t *testing.T) {
	assert.Regexp(t, `^ubc> x = ≈ 14\.2   50% \[`, interactive(t, "x = 10 to 20\n", "--interval", "50").stdout)
}

func TestInterval_AddsItsEndsToJSONQuantiles(t *testing.T) {
	object := decodeLines(t, run(t, "", "--json", "--interval", "80", "-e", "10 to 20").stdout)[0]
	same := decodeLines(t, run(t, "", "--json", "--quantiles", "10,90", "-e", "10 to 20").stdout)[0]

	require.Contains(t, same, "quantiles")
	assert.Equal(t, same, object)
}

func TestInterval_LeavesJSONAloneByDefault(t *testing.T) {
	object := decodeLines(t, run(t, "", "--json", "-e", "10 to 20").stdout)[0]

	assert.NotContains(t, object, "quantiles")
}

func TestInterval_CannotBeCombinedWithQuantiles(t *testing.T) {
	want := outcome{code: 2, stderr: "ubc: --interval and --quantiles cannot be combined\n"}

	assert.Equal(t, want, run(t, "", "--interval", "80", "--quantiles", "10,90", "-e", "1"))
}

func TestInterval_RejectsWhatIsNotAWidth(t *testing.T) {
	for _, width := range []string{"0", "100", "-5", "x", "NaN"} {
		want := outcome{code: 2, stderr: "ubc: --interval expects a percentage between 0 and 100, got '" + width + "'\n"}

		assert.Equal(t, want, run(t, "", "--interval", width, "-e", "1"), width)
	}
}

func TestInterval_AppliesToFailedAssertions(t *testing.T) {
	got := run(t, "", "--interval", "50", "--assert", "x < 15", "-e", "x = 10 to 20")

	assert.Regexp(t, `^ubc: --assert x < 15 failed: ≈ 14\.2   50% \[\S+ … \S+\]   98% \[8\.6 … 24\] < 15 \(P = `, got.stderr)
}
