package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReference_DocumentsEveryReplCommand(t *testing.T) {
	reference, err := os.ReadFile(filepath.Join("..", "..", "..", "docs", "REFERENCE.md"))
	require.NoError(t, err)

	for name := range commands {
		assert.Contains(t, string(reference), "`:"+name, "docs/REFERENCE.md does not document :%s", name)
	}
}
