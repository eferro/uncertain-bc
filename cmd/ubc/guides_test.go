package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
	"github.com/rogpeppe/go-internal/txtar"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGuides_EveryConsoleBlockMatchesTheBinary(t *testing.T) {
	for _, guide := range []string{"docs/TUTORIAL.md", "docs/MODELLING.md", "docs/REFERENCE.md", "docs/HOWTO.md", "README.md"} {
		t.Run(guide, func(t *testing.T) {
			archive, err := guideArchive(repoDoc(t, guide))
			require.NoError(t, err)
			require.NotEmpty(t, archive.Files)
			dir := t.TempDir()
			require.NoError(t, os.WriteFile(filepath.Join(dir, "guide.txtar"), txtar.Format(archive), 0o600))

			testscript.Run(t, testscript.Params{Dir: dir, Setup: copyModels})
		})
	}
}

func TestGuides_AreLinkedFromTheDocsThatPointReadersToThem(t *testing.T) {
	tests := []struct{ doc, link string }{
		{"README.md", "(docs/MODELLING.md)"},
		{"docs/TUTORIAL.md", "(MODELLING.md)"},
		{"scenarios/README.md", "(../docs/MODELLING.md)"},
		{"README.md", "(docs/REFERENCE.md)"},
		{"docs/SPEC.md", "(REFERENCE.md)"},
		{"docs/TUTORIAL.md", "(REFERENCE.md)"},
		{"README.md", "(docs/HOWTO.md)"},
		{"docs/TUTORIAL.md", "(HOWTO.md)"},
		{"docs/REFERENCE.md", "(HOWTO.md)"},
		{"scenarios/README.md", "(../docs/HOWTO.md)"},
		{"README.md", "(docs/TUTORIAL.md)"},
		{"README.md", "(scenarios/README.md)"},
	}
	for _, tt := range tests {
		doc, link := tt.doc, tt.link
		t.Run(doc+" "+link, func(t *testing.T) {
			assert.True(t, strings.Contains(repoDoc(t, doc), link), "%s does not link to %s", doc, link)
		})
	}
}

func repoDoc(t *testing.T, path string) string {
	t.Helper()
	content, err := os.ReadFile(filepath.Join("..", "..", filepath.FromSlash(path)))
	require.NoError(t, err)
	return string(content)
}

func TestGuideArchive_ChecksTheOutputOfEveryCommand(t *testing.T) {
	archive, err := guideArchive("Intro\n\n```console\n$ ubc -e '1 + 1'\n2\n$ ubc examples/a.ubc\nL1: 1\nL2: 2\n```\n")

	require.NoError(t, err)
	assert.Equal(t, "exec ubc -e '1 + 1'\ncmp stdout output1\nexec ubc examples/a.ubc\ncmp stdout output2\n", string(archive.Comment))
	assert.Equal(t, []txtar.File{
		{Name: "output1", Data: []byte("2\n")},
		{Name: "output2", Data: []byte("L1: 1\nL2: 2\n")},
	}, archive.Files)
}

func TestGuideArchive_ExpectsNoOutput_WhenTheCommandShowsNone(t *testing.T) {
	archive, err := guideArchive("```console\n$ ubc -e 'x = 1'\n```\n")

	require.NoError(t, err)
	assert.Equal(t, "exec ubc -e 'x = 1'\n! stdout .\n", string(archive.Comment))
	assert.Empty(t, archive.Files)
}

func TestGuideArchive_ExpectsTheCommandToFail_WhenItShowsAnError(t *testing.T) {
	archive, err := guideArchive("```console\n$ ubc -e '10 to 5'\nubc: line 1, col 4: bad range\n```\n")

	require.NoError(t, err)
	assert.Equal(t, "! exec ubc -e '10 to 5'\n! stdout .\ncmp stderr error1\n", string(archive.Comment))
	assert.Equal(t, []txtar.File{{Name: "error1", Data: []byte("ubc: line 1, col 4: bad range\n")}}, archive.Files)
}

func TestGuideArchive_ChecksStdoutAndStderrApart_WhenTheCommandShowsBoth(t *testing.T) {
	archive, err := guideArchive("```console\n$ ubc --assert 'x > 3' -e 'x = 1'\n1\nubc: --assert x > 3 failed: 1 > 3\n$ ubc -e '2'\n2\n```\n")

	require.NoError(t, err)
	assert.Equal(t, "! exec ubc --assert 'x > 3' -e 'x = 1'\ncmp stdout output1\ncmp stderr error1\nexec ubc -e '2'\ncmp stdout output2\n", string(archive.Comment))
	assert.Equal(t, []txtar.File{
		{Name: "output1", Data: []byte("1\n")},
		{Name: "error1", Data: []byte("ubc: --assert x > 3 failed: 1 > 3\n")},
		{Name: "output2", Data: []byte("2\n")},
	}, archive.Files)
}

func TestGuideArchive_ReadsALineAsOutput_WhenItStartsWithUbcButNotTheErrorPrefix(t *testing.T) {
	archive, err := guideArchive("```console\n$ ubc --scenario-source x\nubc:3\n```\n")

	require.NoError(t, err)
	assert.Equal(t, "exec ubc --scenario-source x\ncmp stdout output1\n", string(archive.Comment))
}

func TestGuideArchive_ComparesTheFileWithTheOutput_WhenTheCommandIsCat(t *testing.T) {
	archive, err := guideArchive("```console\n$ ubc -n 2 --samples-out s.csv -e 'x = 1'\n1\n$ cat s.csv\nx\n1\n1\n```\n")

	require.NoError(t, err)
	assert.Equal(t, "exec ubc -n 2 --samples-out s.csv -e 'x = 1'\ncmp stdout output1\ncmp s.csv output2\n", string(archive.Comment))
	assert.Equal(t, []txtar.File{
		{Name: "output1", Data: []byte("1\n")},
		{Name: "output2", Data: []byte("x\n1\n1\n")},
	}, archive.Files)
}

func TestGuideArchive_ChecksTheShownSourceAgainstTheFileTheNextBlockRuns(t *testing.T) {
	archive, err := guideArchive("```ubc\nx = 1\nx\n```\n\nRun it:\n\n```console\n$ ubc --echo examples/x.ubc\nx = 1\n```\n")

	require.NoError(t, err)
	assert.Equal(t, "cmp examples/x.ubc source1.ubc\nexec ubc --echo examples/x.ubc\ncmp stdout output1\n", string(archive.Comment))
	assert.Equal(t, []txtar.File{
		{Name: "source1.ubc", Data: []byte("x = 1\nx\n")},
		{Name: "output1", Data: []byte("x = 1\n")},
	}, archive.Files)
}

func TestGuideArchive_IgnoresOtherBlocks(t *testing.T) {
	archive, err := guideArchive("```\nplain\n```\n\n```sh\n$ ls\n```\n")

	require.NoError(t, err)
	assert.Empty(t, string(archive.Comment))
	assert.Empty(t, archive.Files)
}

func TestGuideArchive_ReadsFenceLinesInsideABlockAsContent(t *testing.T) {
	archive, err := guideArchive("```text\n```ubc\n```\n```console\n$ ubc -e '1'\n1\n```\n")

	require.NoError(t, err)
	assert.Equal(t, "exec ubc -e '1'\ncmp stdout output1\n", string(archive.Comment))
}

func TestGuideArchive_ChecksTheLastUbcFileOfTheCommand(t *testing.T) {
	archive, err := guideArchive("```ubc\nx\n```\n```console\n$ ubc examples/a.ubc examples/b.ubc\n```\n")

	require.NoError(t, err)
	assert.Contains(t, string(archive.Comment), "cmp examples/b.ubc source1.ubc\n")
}

func TestGuideArchive_Fails(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		message  string
	}{
		{"WhenACommandIsNotUbc", "```console\n$ ls\n```\n", `console command "ls" runs neither ubc nor cat FILE`},
		{"WhenCatShowsMoreThanOneFile", "```console\n$ cat a.csv b.csv\nx\n```\n", `console command "cat a.csv b.csv" runs neither ubc nor cat FILE`},
		{"WhenCatNamesNoFile", "```console\n$ cat\nx\n```\n", `console command "cat" runs neither ubc nor cat FILE`},
		{"WhenCatShowsNothing", "```console\n$ cat a.csv\n```\n", `console command "cat a.csv" shows no content`},
		{"WhenAConsoleBlockIsEmpty", "```console\n```\n", "console block is empty"},
		{"WhenAConsoleBlockStartsWithOutput", "```console\n2\n$ ubc -e '2'\n```\n", `console block starts with "2" instead of a $ command`},
		{"WhenTheSourceIsNeverRun", "```ubc\nx\n```\n", "ubc block is not followed by a console block that runs it"},
		{"WhenTheSourceIsFollowedByAnotherBlock", "```ubc\nx\n```\n```\ntext\n```\n```console\n$ ubc examples/x.ubc\n```\n", "ubc block is not followed by a console block that runs it"},
		{"WhenTheCommandAfterTheSourceRunsNoFile", "```ubc\nx\n```\n```console\n$ ubc -e 'x'\n```\n", `console command "ubc -e 'x'" does not run the ubc file shown before it`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := guideArchive(tt.markdown)

			assert.EqualError(t, err, tt.message)
		})
	}
}
