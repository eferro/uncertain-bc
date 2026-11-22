package main

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/eferro/uncertain-bc/internal/adapters/cli"
	"github.com/eferro/uncertain-bc/scenarios"
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

func TestScenarios_TheBinaryEmbedsEveryFile(t *testing.T) {
	files, err := modelFiles(scenariosDir)
	require.NoError(t, err)
	var embedded []string
	require.NoError(t, fs.WalkDir(scenarios.Files, ".", func(path string, entry fs.DirEntry, err error) error {
		if err == nil && !entry.IsDir() {
			embedded = append(embedded, path)
		}
		return err
	}))

	assert.Equal(t, files, embedded)
}

func TestScenarios_TheIndexListsEveryFile(t *testing.T) {
	files, err := modelFiles(scenariosDir)
	require.NoError(t, err)
	index := catalogueIndex(t)

	assert.Empty(t, missingFromIndex(index, files))
}

func TestScenarios_EveryIndexLinkPointsToAFile(t *testing.T) {
	index := catalogueIndex(t)

	assert.Empty(t, brokenLinks(index, scenariosDir))
}

func TestScenarios_TheIndexStartsWithFiveFeaturedScenarios(t *testing.T) {
	index := catalogueIndex(t)

	assert.Len(t, indexLinks(section(index, "### Start here")), 5)
}

func TestScenarios_EveryAreaHasItsOwnSection(t *testing.T) {
	files, err := modelFiles(scenariosDir)
	require.NoError(t, err)
	index := catalogueIndex(t)

	for area, inArea := range filesByArea(files) {
		t.Run(area, func(t *testing.T) {
			assert.Empty(t, missingFromIndex(section(index, areaHeading(area)), inArea))
		})
	}
}

func TestScenarios_EveryScenarioNamesItsSources(t *testing.T) {
	for file, content := range scenarioContents(t) {
		t.Run(file, func(t *testing.T) {
			assert.Regexp(t, `(?m)^# Sources: \S`, content)
		})
	}
}

func TestScenarios_EveryAssumptionHasAComment(t *testing.T) {
	for file, content := range scenarioContents(t) {
		t.Run(file, func(t *testing.T) {
			assert.Empty(t, uncommentedAssumptions(content))
		})
	}
}

func TestSection_ReturnsTheLinesUpToTheNextHeading(t *testing.T) {
	index := "# Title\n### Team\na\nb\n### Cloud\nc\n## Formulas\nd\n"

	assert.Equal(t, "a\nb\n", section(index, "### Team"))
	assert.Equal(t, "c\n", section(index, "### Cloud"))
	assert.Empty(t, section(index, "### Risk"))
}

func TestFilesByArea_GroupsScenariosAndSkipsTheFormulas(t *testing.T) {
	groups := filesByArea([]string{"formulas.ubc", "team/a.ubc", "cloud/b.ubc", "team/c.ubc"})

	assert.Equal(t, map[string][]string{"team": {"team/a.ubc", "team/c.ubc"}, "cloud": {"cloud/b.ubc"}}, groups)
}

func TestAreaHeading_CapitalisesTheArea(t *testing.T) {
	assert.Equal(t, "### Reliability", areaHeading("reliability"))
}

func TestUncommentedAssumptions_ListsAssignmentsWithoutAComment(t *testing.T) {
	model := "# Question: q?\nx = 1\n\n# Assumptions\na = 1 to 2   # why\nb = 3 to 4\n# a note\nc = 5\n# Calculation\nd = a + b\n"

	assert.Equal(t, []string{"b = 3 to 4", "c = 5"}, uncommentedAssumptions(model))
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

func catalogueIndex(t *testing.T) string {
	t.Helper()
	index, err := os.ReadFile(filepath.Join(scenariosDir, "README.md"))
	require.NoError(t, err)
	return string(index)
}

func scenarioContents(t *testing.T) map[string]string {
	t.Helper()
	files, err := modelFiles(scenariosDir)
	require.NoError(t, err)
	contents := map[string]string{}
	for _, file := range scenarioFiles(files) {
		content, err := os.ReadFile(filepath.Join(scenariosDir, file))
		require.NoError(t, err)
		contents[file] = string(content)
	}
	require.NotEmpty(t, contents)
	return contents
}

func section(index, heading string) string {
	var lines strings.Builder
	inside := false
	for line := range strings.Lines(index) {
		if isHeading(line) {
			inside = strings.TrimSpace(line) == heading
			continue
		}
		if inside {
			lines.WriteString(line)
		}
	}
	return lines.String()
}

func isHeading(line string) bool {
	return strings.HasPrefix(line, "## ") || strings.HasPrefix(line, "### ")
}

func scenarioFiles(files []string) []string {
	var scenarios []string
	for _, file := range files {
		if strings.Contains(file, "/") {
			scenarios = append(scenarios, file)
		}
	}
	return scenarios
}

func filesByArea(files []string) map[string][]string {
	groups := map[string][]string{}
	for _, file := range scenarioFiles(files) {
		area, _, _ := strings.Cut(file, "/")
		groups[area] = append(groups[area], file)
	}
	return groups
}

func areaHeading(area string) string {
	return "### " + strings.ToUpper(area[:1]) + area[1:]
}

func uncommentedAssumptions(model string) []string {
	var missing []string
	inside := false
	for line := range strings.Lines(model) {
		line = strings.TrimSpace(line)
		switch line {
		case "# Assumptions":
			inside = true
		case "# Calculation":
			inside = false
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		if inside && strings.Contains(line, "=") && !strings.Contains(line, "#") {
			missing = append(missing, line)
		}
	}
	return missing
}

func writeFile(t *testing.T, root, name string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(name))
	require.NoError(t, os.MkdirAll(filepath.Dir(path), 0o700))
	require.NoError(t, os.WriteFile(path, []byte("1\n"), 0o600))
}

func TestCopyTree_CopiesModelsAndDataFiles(t *testing.T) {
	from, to := t.TempDir(), t.TempDir()
	writeFile(t, from, "a.ubc")
	writeFile(t, from, "data/weeks.txt")
	writeFile(t, from, "README.md")

	require.NoError(t, copyTree(from, to))

	assert.FileExists(t, filepath.Join(to, "a.ubc"))
	data, err := os.ReadFile(filepath.Join(to, "data", "weeks.txt"))
	require.NoError(t, err)
	assert.Equal(t, "1\n", string(data))
	assert.NoFileExists(t, filepath.Join(to, "README.md"))
}
