package cli

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHistory_StartsEmpty_WhenTheFileDoesNotExist(t *testing.T) {
	assert.Equal(t, 0, loadHistory(filepath.Join(t.TempDir(), "history")).Len())
}

func TestHistory_KeepsTheMostRecentEntryFirst(t *testing.T) {
	h := loadHistory("")
	h.Add("1 + 2")
	h.Add("x = 3")

	assert.Equal(t, 2, h.Len())
	assert.Equal(t, "x = 3", h.At(0))
	assert.Equal(t, "1 + 2", h.At(1))
}

func TestHistory_IgnoresBlankEntries(t *testing.T) {
	h := loadHistory("")
	h.Add("   ")

	assert.Equal(t, 0, h.Len())
}

func TestHistory_PersistsEntriesForTheNextSession(t *testing.T) {
	path := filepath.Join(t.TempDir(), "state", "ubc", "history")
	first := loadHistory(path)
	first.Add("1 + 2")
	first.Add("x = 3")

	second := loadHistory(path)

	assert.Equal(t, 2, second.Len())
	assert.Equal(t, "x = 3", second.At(0))
}

func TestHistory_KeepsWorkingInMemory_WhenTheFileCannotBeWritten(t *testing.T) {
	blocker := filepath.Join(t.TempDir(), "not-a-directory")
	require.NoError(t, os.WriteFile(blocker, nil, 0o600))
	h := loadHistory(filepath.Join(blocker, "history"))

	h.Add("1 + 2")

	assert.Equal(t, "1 + 2", h.At(0))
}

func TestHistory_KeepsWorkingInMemory_WhenTheFileIsADirectory(t *testing.T) {
	path := t.TempDir()
	h := loadHistory(path)

	h.Add("1 + 2")

	assert.Equal(t, "1 + 2", h.At(0))
}
