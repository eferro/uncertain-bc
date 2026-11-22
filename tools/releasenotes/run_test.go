package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun_ExitsWithUsageCode2_WhenTheArgumentsAreWrong(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"CHANGELOG.md"}, &stdout, &stderr)

	assert.Equal(t, 2, code)
	assert.Empty(t, stdout.String())
	assert.Equal(t, "usage: releasenotes CHANGELOG VERSION\n", stderr.String())
}

func TestRun_ExitsWith1_WhenTheChangelogCannotBeRead(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := run([]string{"missing.md", "v1.0.0"}, &stdout, &stderr)

	assert.Equal(t, 1, code)
	assert.Empty(t, stdout.String())
}
