package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestTailEnds_KeepTwoSignificantFiguresOfTheirDistance(t *testing.T) {
	cases := []struct {
		lo, hi         float64
		wantLo, wantHi string
	}{
		{75.53, 378.3, "80", "380"},
		{8.642, 23.03, "9", "23"},
		{11458340, 11741129, "11.46M", "11.74M"},
		{997.9, 1012.05, "998", "1012"},
		{-2.105, 12.04, "-2", "12"},
		{0.01234, 0.0987, "0.012", "0.099"},
		{-1.5, 310, "0", "310"},
	}
	for _, c := range cases {
		lo, hi := tailEnds(c.lo, c.hi)

		assert.Equal(t, []string{c.wantLo, c.wantHi}, []string{lo, hi}, "%g %g", c.lo, c.hi)
	}
}

func TestTailEnds_KeepThreeSignificantFigures_WhenBothEndsAreEqual(t *testing.T) {
	lo, hi := tailEnds(1234, 1234)

	assert.Equal(t, []string{"1230", "1230"}, []string{lo, hi})
}
