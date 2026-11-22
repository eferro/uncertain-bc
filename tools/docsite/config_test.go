package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const validConfig = `site_name: ubc
repo_url: https://github.com/owner/repo
docs_dir: build/docs
nav:
  - Home: README.md
  - Guides:
    - Tutorial: docs/TUTORIAL.md
`

func TestParseConfig_ReadsTheRepositoryTheDocsDirAndEveryPageOfTheNav(t *testing.T) {
	got, err := parseConfig(validConfig)

	require.NoError(t, err)
	assert.Equal(t, config{
		repoURL: "https://github.com/owner/repo",
		docsDir: "build/docs",
		pages:   []string{"README.md", "docs/TUTORIAL.md"},
	}, got)
}

func TestParseConfig_Fails_WhenAKeyTheSiteNeedsIsMissing(t *testing.T) {
	tests := []struct{ name, text, want string }{
		{"repo_url", "docs_dir: d\nnav:\n  - Home: README.md\n", "no repo_url"},
		{"docs_dir", "repo_url: u\nnav:\n  - Home: README.md\n", "no docs_dir"},
		{"nav", "repo_url: u\ndocs_dir: d\n", "no pages in nav"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseConfig(tt.text)

			assert.EqualError(t, err, tt.want)
		})
	}
}
