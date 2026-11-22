package domain

import (
	"maps"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReference_DocumentsEveryBuiltInFunction(t *testing.T) {
	reference, err := os.ReadFile(filepath.Join("..", "..", "docs", "REFERENCE.md"))
	require.NoError(t, err)
	names := slices.Concat(
		slices.Collect(maps.Keys(functions)),
		slices.Collect(maps.Keys(summaries)),
		slices.Collect(maps.Keys(mathFunctions)),
		slices.Collect(maps.Keys(distributions)),
	)

	for _, name := range names {
		assert.Contains(t, string(reference), "`"+name+"(", "docs/REFERENCE.md does not document %s", name)
	}
}
