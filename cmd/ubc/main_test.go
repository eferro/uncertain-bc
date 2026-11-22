package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"ubc": main,
	})
}

func TestScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{
		Dir:   "testdata/script",
		Setup: copyExamples,
		Cmds:  map[string]func(*testscript.TestScript, bool, []string){"unescape": unescape},
	})
}

func copyExamples(env *testscript.Env) error {
	examples, err := filepath.Glob("../../examples/*.ubc")
	if err != nil {
		return err
	}
	dir := filepath.Join(env.WorkDir, "examples")
	if err := os.Mkdir(dir, 0o700); err != nil {
		return err
	}
	for _, path := range examples {
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(dir, filepath.Base(path)), content, 0o600); err != nil {
			return err
		}
	}
	return nil
}

// unescape rewrites a one-line file holding Go escapes (\n, \x01) into the
// raw bytes they stand for, so terminal keystrokes stay readable in scripts.
func unescape(ts *testscript.TestScript, neg bool, args []string) {
	if neg || len(args) != 1 {
		ts.Fatalf("usage: unescape file")
	}
	escaped := strings.TrimSuffix(ts.ReadFile(args[0]), "\n")
	raw, err := strconv.Unquote(`"` + escaped + `"`)
	ts.Check(err)
	ts.Check(os.WriteFile(ts.MkAbs(args[0]), []byte(raw), 0o600))
}
