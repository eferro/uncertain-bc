package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRun_ExitsWithUsageCode2_WhenTheArgumentsAreWrong(t *testing.T) {
	var stderr bytes.Buffer

	code := run([]string{"a.yml", "b.yml"}, &stderr)

	assert.Equal(t, 2, code)
	assert.Equal(t, "usage: docsite MKDOCS_CONFIG\n", stderr.String())
}

func TestRun_ExitsWith1_WhenTheConfigLacksWhatTheSiteNeeds(t *testing.T) {
	config := filepath.Join(t.TempDir(), "mkdocs.yml")
	require.NoError(t, os.WriteFile(config, []byte("site_name: x\n"), 0o600))
	var stderr bytes.Buffer

	code := run([]string{config}, &stderr)

	assert.Equal(t, 1, code)
	assert.Equal(t, "docsite: "+config+": no repo_url\n", stderr.String())
}

func TestAssemble_PublishesEveryPageOfTheRepositorySite_WithoutBrokenLinks(t *testing.T) {
	text, err := os.ReadFile("../../mkdocs.yml")
	require.NoError(t, err)
	site, err := parseConfig(string(text))
	require.NoError(t, err)
	out := t.TempDir()

	require.NoError(t, assemble("../..", site, out))

	for _, page := range site.pages {
		assert.FileExists(t, filepath.Join(out, filepath.FromSlash(page)))
	}
}
