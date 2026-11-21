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
