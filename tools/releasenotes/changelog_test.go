package main

import (
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	releasedVersion = regexp.MustCompile(`^\d+\.\d+\.\d+$`)
	releaseDate     = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	githubLink      = regexp.MustCompile(`(?m)^\[([^\]]+)\]: https://github\.com/eferro/uncertain-bc/\S+$`)
)

func TestChangelog_FollowsKeepAChangelog(t *testing.T) {
	changelog := repoChangelog(t)
	all := sections(changelog)

	assert.True(t, strings.HasPrefix(changelog, "# Changelog\n"))
	require.NotEmpty(t, all)
	assert.Equal(t, "Unreleased", all[0].version)
	assert.Empty(t, all[0].date)
	assert.Equal(t, "0.1.0", all[len(all)-1].version)
}

func TestChangelog_EveryReleaseHasADateAndNotes(t *testing.T) {
	for _, s := range sections(repoChangelog(t))[1:] {
		assert.Regexp(t, releasedVersion, s.version)
		assert.Regexp(t, releaseDate, s.date, s.version)
		assert.NotEmpty(t, s.body, s.version)
	}
}

func TestChangelog_ListsReleasesFromNewestToOldest(t *testing.T) {
	var versions []string
	for _, s := range sections(repoChangelog(t))[1:] {
		versions = append(versions, s.version)
	}

	assert.True(t, slices.IsSortedFunc(versions, compareVersionsDescending), versions)
	assert.Equal(t, versions, slices.Compact(slices.Clone(versions)), "a version is repeated")
}

func TestChangelog_LinksEveryVersionToItsChanges(t *testing.T) {
	changelog := repoChangelog(t)
	linked := linkedVersions(changelog)
	for _, s := range sections(changelog) {
		assert.Contains(t, linked, s.version)
	}
}

func linkedVersions(changelog string) []string {
	var versions []string
	for _, match := range githubLink.FindAllStringSubmatch(changelog, -1) {
		versions = append(versions, match[1])
	}
	return versions
}

func repoChangelog(t *testing.T) string {
	t.Helper()
	content, err := os.ReadFile("../../CHANGELOG.md")
	require.NoError(t, err)
	return string(content)
}

func compareVersionsDescending(a, b string) int {
	return slices.Compare(versionNumbers(b), versionNumbers(a))
}

func versionNumbers(version string) []int {
	var major, minor, patch int
	_, _ = fmt.Sscanf(version, "%d.%d.%d", &major, &minor, &patch)
	return []int{major, minor, patch}
}

func TestCompareVersionsDescending_OrdersByNumberNotByText(t *testing.T) {
	assert.Negative(t, compareVersionsDescending("0.10.0", "0.9.1"))
	assert.Positive(t, compareVersionsDescending("0.9.1", "1.0.0"))
	assert.Zero(t, compareVersionsDescending("0.7.2", "0.7.2"))
}
