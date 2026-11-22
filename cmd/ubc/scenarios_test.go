package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/adapters/cli"
)

const scenariosDir = "../../scenarios"

func TestScenarios_EveryFileEvaluates(t *testing.T) {
	files, err := modelFiles(scenariosDir)
	require.NoError(t, err)
	require.NotEmpty(t, files)
	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			var app cli.App

			code := app.Run([]string{filepath.Join(scenariosDir, file)}, strings.NewReader(""), &stdout, &stderr)

			assert.Equal(t, 0, code, stderr.String())
			assert.Empty(t, stderr.String())
		})
	}
}

func TestScenarios_TheIndexListsEveryFile(t *testing.T) {
	files, err := modelFiles(scenariosDir)
	require.NoError(t, err)
	index, err := os.ReadFile(filepath.Join(scenariosDir, "README.md"))
	require.NoError(t, err)

	assert.Empty(t, missingFromIndex(string(index), files))
}

func TestScenarioFiles_FindsUbcFilesInEveryArea(t *testing.T) {
	root := t.TempDir()
	writeFile(t, root, "team/b.ubc")
	writeFile(t, root, "team/a.ubc")
	writeFile(t, root, "product/c.ubc")
	writeFile(t, root, "README.md")
	writeFile(t, root, "team/notes.txt")
	writeFile(t, root, "old.ubc/d.ubc")

	files, err := modelFiles(root)

	require.NoError(t, err)
	assert.Equal(t, []string{"old.ubc/d.ubc", "product/c.ubc", "team/a.ubc", "team/b.ubc"}, files)
}

func TestScenarioFiles_Fails_WhenTheDirectoryDoesNotExist(t *testing.T) {
	_, err := modelFiles(filepath.Join(t.TempDir(), "missing"))

	assert.Error(t, err)
}

func TestMissingFromIndex_ListsFilesWithoutALink(t *testing.T) {
	index := "| team | [a](team/a.ubc) | Q? |\nmentions team/b.ubc without a link\n"

	missing := missingFromIndex(index, []string{"team/a.ubc", "team/b.ubc", "team/c.ubc"})

	assert.Equal(t, []string{"team/b.ubc", "team/c.ubc"}, missing)
}

func writeFile(t *testing.T, root, name string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(name))
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o700))
	require.NoError(t, os.WriteFile(path, []byte("1\n"), 0o600))
}
