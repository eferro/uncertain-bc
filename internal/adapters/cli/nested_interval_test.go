package cli_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNestedInterval_FollowsTheNinetyPercentOneByDefault(t *testing.T) {
	got := run(t, "", "-e", "(8 to 25) * 12")

	assert.Regexp(t, `^≈ 170   90% \[96 … 296\]   98% \[75 … 380\]   [▁▂▃▄▅▆▇█]{16}\n$`, got.stdout)
}

func TestNestedInterval_StaysAtNinetyEightPercent_WhenTheInnerWidthChanges(t *testing.T) {
	got := run(t, "", "--interval", "50", "-e", "(8 to 25) * 12")

	assert.Regexp(t, `^≈ 170   50% \[\S+ … \S+\]   98% \[75 … 380\]   `, got.stdout)
}

func TestNestedInterval_IsLeftOut_WhenTheIntervalIsAtLeastNinetyEightPercent(t *testing.T) {
	for _, width := range []string{"98", "99.5"} {
		got := run(t, "", "--interval", width, "-e", "(8 to 25) * 12")

		assert.Regexp(t, `^≈ 170   `+width+`% \[\S+ … \S+\]   [▁▂▃▄▅▆▇█]{16}\n$`, got.stdout, width)
	}
}

func TestNestedInterval_IsShown_WhenTheIntervalIsJustBelowNinetyEightPercent(t *testing.T) {
	got := run(t, "", "--interval", "97.9", "-e", "(8 to 25) * 12")

	assert.Regexp(t, `^≈ 170   97\.9% \[\S+ … \S+\]   98% \[75 … 380\]   `, got.stdout)
}

func TestNestedInterval_IsReplacedByQuantiles(t *testing.T) {
	got := run(t, "", "--quantiles", "10,90", "-e", "(8 to 25) * 12")

	assert.Regexp(t, `^≈ 170   p10 109   p90 262   [▁▂▃▄▅▆▇█]{16}\n$`, got.stdout)
}

func TestNestedInterval_AddsP1AndP99ToJSON(t *testing.T) {
	object := decodeLines(t, run(t, "", "--json", "-e", "(8 to 25) * 12").stdout)[0]

	assert.InDelta(t, 75.5, object["p1"], 0.1)
	assert.InDelta(t, 378.3, object["p99"], 0.1)
}

func TestNestedInterval_AppliesToFailedAssertions(t *testing.T) {
	got := run(t, "", "--assert", "x < 15", "-e", "x = 10 to 20")

	assert.Regexp(t, `^ubc: --assert x < 15 failed: ≈ 14\.2   90% \[10 … 19\.8\]   98% \[\S+ … \S+\] < 15 \(P = `, got.stderr)
}

func TestNestedInterval_RoundsItsEndsOutwardsToTheSpread_WhenTheDistributionIsNarrow(t *testing.T) {
	got := run(t, "", "-e", "1000 to 1010")

	assert.Regexp(t, `^≈ 1005   90% \[1000 … 1009\.9\]   98% \[997 … 1013\]   `, got.stdout)
}
