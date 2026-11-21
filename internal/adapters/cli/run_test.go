package cli_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/eferro/uncertain-bc/internal/adapters/cli"
)

func TestRun_PrintsTheResult_WhenGivenAnExpression(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := cli.Run([]string{"-e", "42"}, &stdout, &stderr)

	assert.Equal(t, 0, code)
	assert.Equal(t, "42\n", stdout.String())
	assert.Empty(t, stderr.String())
}

func TestRun_ExitsWithUsageError_WhenAFlagIsUnknown(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := cli.Run([]string{"--nope"}, &stdout, &stderr)

	assert.Equal(t, 2, code)
	assert.Empty(t, stdout.String())
	assert.Contains(t, stderr.String(), "-nope")
}

func TestRun_ExitsWithParseError_WhenTheExpressionIsNotValid(t *testing.T) {
	var stdout, stderr bytes.Buffer

	code := cli.Run([]string{"-e", "forty-two"}, &stdout, &stderr)

	assert.Equal(t, 3, code)
	assert.Empty(t, stdout.String())
	assert.Regexp(t, `^ubc: .+\n$`, stderr.String())
}
