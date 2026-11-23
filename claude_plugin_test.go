package install_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const skillDir = "claude-plugin/skills/ubc"

type marketplace struct {
	Name    string `json:"name"`
	Plugins []struct {
		Name   string `json:"name"`
		Source string `json:"source"`
	} `json:"plugins"`
}

func readJSON(t *testing.T, path string, into any) {
	t.Helper()
	content, err := os.ReadFile(filepath.FromSlash(path))
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(content, into))
}

func readText(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.FromSlash(path))
	require.NoError(t, err)
	return string(content)
}

func frontmatterField(t *testing.T, skill, field string) string {
	t.Helper()
	match := regexp.MustCompile(`(?m)^` + field + `: (.*)$`).FindStringSubmatch(skill)
	require.NotNil(t, match, "SKILL.md has no %s", field)
	return match[1]
}

func TestClaudePlugin_MarketplaceListsThePluginByTheNameOfItsManifest(t *testing.T) {
	var market marketplace
	readJSON(t, ".claude-plugin/marketplace.json", &market)
	require.Len(t, market.Plugins, 1)
	entry := market.Plugins[0]
	var manifest struct {
		Name string `json:"name"`
	}

	readJSON(t, filepath.Join(entry.Source, ".claude-plugin", "plugin.json"), &manifest)

	assert.Equal(t, "uncertain-bc", market.Name)
	assert.Equal(t, entry.Name, manifest.Name)
	assert.True(t, strings.HasPrefix(entry.Source, "./"), "a relative source must start with ./")
}

func TestClaudePlugin_SkillIsNamedAfterItsDirectory(t *testing.T) {
	skill := readText(t, skillDir+"/SKILL.md")

	assert.Equal(t, filepath.Base(skillDir), frontmatterField(t, skill, "name"))
}

func TestClaudePlugin_SkillDescriptionFitsTheAgentSkillsLimit(t *testing.T) {
	skill := readText(t, skillDir+"/SKILL.md")

	assert.LessOrEqual(t, len(frontmatterField(t, skill, "description")), 1024)
}

func TestClaudePlugin_SkillAndInstallScriptRequireTheSameVersion(t *testing.T) {
	script := readText(t, skillDir+"/scripts/ensure-ubc.sh")
	match := regexp.MustCompile(`(?m)^minimum=(\S+)$`).FindStringSubmatch(script)
	require.NotNil(t, match)

	assert.Contains(t, readText(t, skillDir+"/SKILL.md"), "older than "+match[1]+".")
	assert.Contains(t, readText(t, "README.md"), "ubc "+match[1]+" or later")
}

func TestClaudePlugin_SkillLinksToEveryReference(t *testing.T) {
	skill := readText(t, skillDir+"/SKILL.md")
	references, err := filepath.Glob(filepath.Join(filepath.FromSlash(skillDir), "references", "*.md"))
	require.NoError(t, err)
	require.NotEmpty(t, references)

	for _, reference := range references {
		assert.Contains(t, skill, "(references/"+filepath.Base(reference)+")")
	}
}
