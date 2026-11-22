package main

import (
	"bytes"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/adapters/cli"
)

const (
	scenariosDir = "../../scenarios"
	formulasFile = "formulas.ubc"
)

func TestScenarios_EveryFileEvaluates(t *testing.T) {
	files, err := modelFiles(scenariosDir)
	require.NoError(t, err)
	require.NotEmpty(t, files)
	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			var stdout, stderr bytes.Buffer
			var app cli.App

			code := app.Run(scenarioArgs(scenariosDir, file), strings.NewReader(""), &stdout, &stderr)

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

func TestScenarios_EveryIndexLinkPointsToAFile(t *testing.T) {
	index, err := os.ReadFile(filepath.Join(scenariosDir, "README.md"))
	require.NoError(t, err)

	assert.Empty(t, brokenLinks(string(index), scenariosDir))
}

func TestScenarioArgs_RunsTheSharedFormulasFirst(t *testing.T) {
	args := scenarioArgs("root", "team/a.ubc")

	assert.Equal(t, []string{filepath.Join("root", "formulas.ubc"), filepath.Join("root", "team", "a.ubc")}, args)
}

func TestScenarioArgs_RunsTheSharedFormulasAlone(t *testing.T) {
	args := scenarioArgs("root", "formulas.ubc")

	assert.Equal(t, []string{filepath.Join("root", "formulas.ubc")}, args)
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

func TestBrokenLinks_ListsUbcLinksToMissingFiles(t *testing.T) {
	root := filepath.Join(t.TempDir(), "scenarios")
	writeFile(t, root, "team/a.ubc")
	writeFile(t, root, "../examples/b.ubc")
	index := "[a](team/a.ubc) [b](../examples/b.ubc) [c](team/c.ubc) [d](../examples/d.ubc) [tutorial](../docs/TUTORIAL.md)"

	broken := brokenLinks(index, root)

	assert.Equal(t, []string{"team/c.ubc", "../examples/d.ubc"}, broken)
}

func scenarioArgs(root, file string) []string {
	formulas := filepath.Join(root, formulasFile)
	if file == formulasFile {
		return []string{formulas}
	}
	return []string{formulas, filepath.Join(root, file)}
}

func missingFromIndex(index string, files []string) []string {
	links := indexLinks(index)
	var missing []string
	for _, file := range files {
		if !slices.Contains(links, file) {
			missing = append(missing, file)
		}
	}
	return missing
}

var ubcLink = regexp.MustCompile(`\]\(([^)]+\.ubc)\)`)

func indexLinks(index string) []string {
	var links []string
	for _, match := range ubcLink.FindAllStringSubmatch(index, -1) {
		links = append(links, match[1])
	}
	return links
}

func brokenLinks(index, root string) []string {
	var broken []string
	for _, link := range indexLinks(index) {
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(link))); err != nil {
			broken = append(broken, link)
		}
	}
	return broken
}

func writeFile(t *testing.T, root, name string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(name))
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o700))
	require.NoError(t, os.WriteFile(path, []byte("1\n"), 0o600))
}
