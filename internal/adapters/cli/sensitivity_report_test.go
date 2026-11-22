package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eferro/uncertain-bc/internal/application"
)

func TestSensitivityReport_DrawsABarAsLongAsEachCorrelationIsStrong(t *testing.T) {
	drivers := []application.Driver{{Name: "units", Correlation: 0.96}, {Name: "price", Correlation: -0.52}, {Name: "c", Correlation: 0.024}, {Name: "d", Correlation: 0.026}}

	got := sensitivityReport("profit", drivers, unicodeStyle)

	assert.Equal(t, "sensitivity of profit (rank correlation):\n"+
		"  units  +0.96  ███████████████████\n"+
		"  price  -0.52  ██████████\n"+
		"  c      +0.02\n"+
		"  d      +0.03  █\n", got)
}

func TestSensitivityReport_DrawsASCIIBars_WithASCII(t *testing.T) {
	got := sensitivityReport("x", []application.Driver{{Name: "y", Correlation: 1}}, asciiStyle)

	assert.Equal(t, "sensitivity of x (rank correlation):\n  y  +1.00  ####################\n", got)
}
