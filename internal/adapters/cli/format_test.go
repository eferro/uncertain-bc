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
