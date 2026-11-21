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
