package cli

import (
	"math"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eferro/uncertain-bc/internal/domain"
)

func TestSignificant_RoundsToThreeSignificantFigures(t *testing.T) {
	cases := map[float64]string{
		14.142:   "14.1",
		9.996:    "10",
		299.6:    "300",
		1234:     "1230",
		0.012345: "0.0123",
		-14.142:  "-14.1",
		0:        "0",
		20:       "20",
	}
	for value, want := range cases {
		assert.Equal(t, want, significant(value), "significant(%v)", value)
	}
}

func TestSignificant_UsesASuffixFromTenThousandUp(t *testing.T) {
	cases := map[float64]string{
		9999:    "10k",
		10000:   "10k",
		22412:   "22.4k",
		-15000:  "-15k",
		999500:  "1M",
		1234567: "1.23M",
		2.5e9:   "2.5B",
		3e12:    "3T",
		4.56e15: "4560T",
	}
	for value, want := range cases {
		assert.Equal(t, want, significant(value), "significant(%v)", value)
	}
}

func TestSignificant_ShowsNoFloatNoise_WhenRoundingToTensOrMore(t *testing.T) {
	cases := map[float64]string{
		11560770:  "11.6M",
		11499000:  "11.5M",
		11699000:  "11.7M",
		1.03e17:   "103000T",
		-11560770: "-11.6M",
	}
	for value, want := range cases {
		assert.Equal(t, want, significant(value), "significant(%v)", value)
	}
}

func TestRoundForSummary_KeepsTwoSignificantFiguresOfTheWidth_WhenFinerThanThreeOfItsOwn(t *testing.T) {
	cases := []struct{ x, width, want float64 }{
		{1001234, 2012, 1001200},
		{-1001234, 2012, -1001200},
		{0.999871, 0.00046, 0.99987},
		{1000.4, 10.4, 1000},
		{1004.6, 99.4, 1005},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, roundForSummary(c.x, c.width), "roundForSummary(%v, %v)", c.x, c.width)
	}
}

func TestRoundForSummary_KeepsThreeSignificantFiguresOfItsOwn_WhenTheWidthIsWider(t *testing.T) {
	cases := []struct{ x, width, want float64 }{
		{14.142, 9.8, 14.1},
		{170.3, 200, 170},
		{1234, 0, 1230},
		{1234, 100, 1230},
		{0, 5, 0},
	}
	for _, c := range cases {
		assert.Equal(t, c.want, roundForSummary(c.x, c.width), "roundForSummary(%v, %v)", c.x, c.width)
	}
}

func TestSummaryNumbers_TellTheEndsOfANarrowIntervalApart(t *testing.T) {
	got := summaryNumbers(ramp(1_000_000, 1_002_000), unicodeStyle)

	assert.Equal(t, "≈ 1.001M   90% [1.0001M … 1.0019M]   98% [1M … 1.002M]", got)
}

func TestSummaryNumbers_TellTheEndsOfANarrowInnerIntervalApart_WithAnotherWidth(t *testing.T) {
	st := unicodeStyle
	st.interval = 5

	got := summaryNumbers(ramp(1_000_000, 1_002_000), st)

	assert.Equal(t, "≈ 1.001M   5% [1.00095M … 1.00105M]   98% [1M … 1.002M]", got)
}

func TestSummaryNumbers_TellNarrowQuantilesApart_UsingTheWidthOfTheInnerInterval(t *testing.T) {
	st := unicodeStyle
	st.quantiles = []float64{5, 95}

	got := summaryNumbers(ramp(1_000_000, 1_002_000), st)

	assert.Equal(t, "≈ 1.001M   p5 1.0001M   p95 1.0019M", got)
}

func TestProbability_ShowsTwoDecimalsAndFlagsTheExtremes(t *testing.T) {
	cases := map[float64]string{
		0.5:   "0.50",
		0.16:  "0.16",
		0.004: "<0.01",
		0.996: ">0.99",
		0:     "0",
		1:     "1",
	}
	for value, want := range cases {
		assert.Equal(t, want, probability(value), "probability(%v)", value)
	}
}

func TestPlain_ShowsAtMostTenSignificantFigures(t *testing.T) {
	cases := map[float64]string{
		4:                 "4",
		0.1:               "0.1",
		3.141592653589793: "3.141592654",
		3000:              "3000",
		1e12:              "1000000000000",
		-2.5:              "-2.5",
		0.000123456789012: "0.000123456789",
	}
	for value, want := range cases {
		assert.Equal(t, want, plain(value), "plain(%v)", value)
	}
}

func TestPlain_ShowsNoFloatNoise_WhenRoundingToTensOrMore(t *testing.T) {
	cases := map[float64]string{
		100000010000000:  "100000010000000",
		1.03000001e18:    "1030000010000000000",
		-100000010000000: "-100000010000000",
	}
	for value, want := range cases {
		assert.Equal(t, want, plain(value), "plain(%v)", value)
	}
}

func TestTailEnds_RoundOutwardsToTwoSignificantFiguresOfTheirDistance(t *testing.T) {
	cases := []struct {
		lo, hi         float64
		wantLo, wantHi float64
	}{
		{11458340, 11741129, 11450000, 11750000},
		{997.9, 1012.05, 997, 1013},
		{0.01234, 0.0987, 0.012, 0.099},
		{0.29, 0.57, 0.29, 0.57},
		{610, 15000, 610, 15000},
	}
	for _, c := range cases {
		lo, hi := tailEnds(c.lo, c.hi, c.lo, c.hi)

		assert.Equal(t, []float64{c.wantLo, c.wantHi}, []float64{lo, hi}, "%g %g", c.lo, c.hi)
	}
}

func TestTailEnds_KeepTwoSignificantFiguresOfTheirOwn_WhenTheDistanceIsCoarser(t *testing.T) {
	cases := []struct {
		lo, hi         float64
		wantLo, wantHi float64
	}{
		{615.9, 15978.1, 610, 16000},
		{75.53, 378.3, 75, 380},
		{8.642, 23.03, 8.6, 24},
		{-2.105, 12.04, -2.2, 13},
		{-1.5, 310, -1.5, 310},
		{-310, -1.5, -310, -1.5},
	}
	for _, c := range cases {
		lo, hi := tailEnds(c.lo, c.hi, c.lo, c.hi)

		assert.Equal(t, []float64{c.wantLo, c.wantHi}, []float64{lo, hi}, "%g %g", c.lo, c.hi)
	}
}

func TestSummaryNumbers_KeepThe98PercentAroundThe90Percent_WhenTheInnerEndsRoundFurtherOut(t *testing.T) {
	cases := map[string]domain.SampleSet{
		"≈ 1000   90% [1000 … 1000]   98% [1000 … 1007]": repeated(map[float64]int{1001: 3, 1004: 94, 1007: 3}),
		"≈ 1050   90% [1050 … 1050]   98% [1043 … 1050]": repeated(map[float64]int{1043: 3, 1046: 94, 1047: 3}),
	}
	for want, samples := range cases {
		assert.Equal(t, want, summaryNumbers(samples, unicodeStyle))
	}
}

func TestTailEnds_KeepThreeSignificantFigures_WhenBothEndsAreEqual(t *testing.T) {
	lo, hi := tailEnds(1234, 1234, 1230, 1230)

	assert.Equal(t, []float64{1230, 1230}, []float64{lo, hi})
}

func TestTailEnds_ContainTheTailsAndTheInnerInterval_ForSkewedSamples(t *testing.T) {
	for seed := range uint64(40) {
		samples := skewedSamples(seed)
		p1, p99 := samples.Quantile(0.01), samples.Quantile(0.99)
		innerLo, innerHi := roundSignificant(samples.Quantile(0.05), 3), roundSignificant(samples.Quantile(0.95), 3)

		lo, hi := tailEnds(p1, p99, innerLo, innerHi)

		assert.LessOrEqual(t, lo, math.Min(p1, innerLo), "seed %d", seed)
		assert.GreaterOrEqual(t, hi, math.Max(p99, innerHi), "seed %d", seed)
	}
}

func repeated(counts map[float64]int) domain.SampleSet {
	var samples domain.SampleSet
	for value, count := range counts {
		for range count {
			samples = append(samples, value)
		}
	}
	return samples
}

func skewedSamples(seed uint64) domain.SampleSet {
	r := rand.New(rand.NewPCG(seed, seed))
	samples := make(domain.SampleSet, 5000)
	for i := range samples {
		samples[i] = 1000 * math.Exp(r.NormFloat64())
	}
	return samples
}

func ramp(lo, hi float64) domain.SampleSet {
	samples := make(domain.SampleSet, 1001)
	for i := range samples {
		samples[i] = lo + (hi-lo)*float64(i)/1000
	}
	return samples
}
