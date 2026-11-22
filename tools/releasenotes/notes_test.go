package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const sample = "# Changelog\n\nIntro.\n\n## [Unreleased]\n\n## [1.1.0] - 2026-02-01\n\n### Fixed\n\n- B.\n\n## [1.0.0] - 2026-01-01\n\n- A.\n\n[Unreleased]: u\n[1.1.0]: v\n"

func TestSections_ListsEveryVersionInOrder_WithItsDateAndBody(t *testing.T) {
	assert.Equal(t, []section{
		{version: "Unreleased", body: ""},
		{version: "1.1.0", date: "2026-02-01", body: "### Fixed\n\n- B."},
		{version: "1.0.0", date: "2026-01-01", body: "- A."},
	}, sections(sample))
}

func TestSections_IsEmpty_WhenThereAreNoVersionHeadings(t *testing.T) {
	assert.Empty(t, sections("# Changelog\n\nNothing yet.\n"))
}

func TestNotesFor_ReturnsTheBodyOfTheVersion_WithAndWithoutTheVPrefix(t *testing.T) {
	for _, version := range []string{"v1.1.0", "1.1.0"} {
		notes, err := notesFor(sample, version)

		require.NoError(t, err)
		assert.Equal(t, "### Fixed\n\n- B.\n", notes)
	}
}

func TestNotesFor_KeepsLinkLikeLinesInsideTheBody(t *testing.T) {
	notes, err := notesFor("## [1.0.0] - 2026-01-01\n\n- See [ADR 1]: it explains.\n- [x]: y\n\n## [0.1.0] - 2025-01-01\n\n- C.\n", "1.0.0")

	require.NoError(t, err)
	assert.Equal(t, "- See [ADR 1]: it explains.\n- [x]: y\n", notes)
}

func TestNotesFor_Fails_WhenTheVersionHasNoSection(t *testing.T) {
	_, err := notesFor(sample, "v2.0.0")

	assert.EqualError(t, err, "no section for 2.0.0")
}

func TestNotesFor_Fails_WhenTheSectionIsEmpty(t *testing.T) {
	_, err := notesFor(sample, "Unreleased")

	assert.EqualError(t, err, "the Unreleased section is empty")
}
