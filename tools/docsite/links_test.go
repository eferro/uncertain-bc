package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeRepo map[string]targetKind

func (r fakeRepo) kind(path string) (targetKind, bool) {
	kind, ok := r[path]
	return kind, ok
}

var repo = fakeRepo{
	"README.md":          file,
	"docs/GUIDE.md":      file,
	"docs/SPEC.md":       file,
	"LICENSE":            file,
	"examples":           directory,
	"../other/README.md": file,
	"..":                 directory,
}

var site = config{
	repoURL: "https://github.com/owner/repo",
	docsDir: "build/docs",
	pages:   []string{"README.md", "docs/GUIDE.md"},
}

func TestRewriteLinks_PointsLinksToUnpublishedFilesAtTheRepository(t *testing.T) {
	tests := []struct{ name, page, text, want string }{
		{"published page", "README.md", "[g](docs/GUIDE.md)", "[g](docs/GUIDE.md)"},
		{"published page upwards", "docs/GUIDE.md", "[h](../README.md#install)", "[h](../README.md#install)"},
		{"file", "README.md", "[l](LICENSE)", "[l](https://github.com/owner/repo/blob/main/LICENSE)"},
		{"file with anchor", "docs/GUIDE.md", "[s](SPEC.md#5-functions)", "[s](https://github.com/owner/repo/blob/main/docs/SPEC.md#5-functions)"},
		{"directory", "docs/GUIDE.md", "[e](../examples/)", "[e](https://github.com/owner/repo/tree/main/examples/)"},
		{"web", "README.md", "[w](https://example.com/x.md)", "[w](https://example.com/x.md)"},
		{"anchor", "README.md", "[a](#install)", "[a](#install)"},
		{"several", "README.md", "[l](LICENSE) and [g](docs/GUIDE.md)", "[l](https://github.com/owner/repo/blob/main/LICENSE) and [g](docs/GUIDE.md)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := rewriteLinks(tt.page, tt.text, site, repo)

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRewriteLinks_Fails_WhenATargetIsNotInTheRepository(t *testing.T) {
	tests := []struct{ name, page, text, want string }{
		{"missing", "docs/GUIDE.md", "[m](MISSING.md)", "broken link to docs/MISSING.md"},
		{"outside", "README.md", "[o](../other/README.md)", "broken link to ../other/README.md"},
		{"parent", "README.md", "[p](..)", "broken link to .."},
		{"first of several", "docs/GUIDE.md", "[a](A.md) and [b](B.md)", "broken link to docs/A.md"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := rewriteLinks(tt.page, tt.text, site, repo)

			assert.EqualError(t, err, tt.want)
		})
	}
}
