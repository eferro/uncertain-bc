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
		Setup: copyModels,
		Cmds:  map[string]func(*testscript.TestScript, bool, []string){"unescape": unescape},
	})
}

// copyModels gives every scenario the models a reader runs from the
// repository root, so commands read exactly as in the docs.
func copyModels(env *testscript.Env) error {
	for _, dir := range []string{"examples", "scenarios"} {
		if err := copyTree(filepath.Join("..", "..", dir), filepath.Join(env.WorkDir, dir)); err != nil {
			return err
		}
	}
	return nil
}

func copyTree(from, to string) error {
	files, err := modelFiles(from)
	if err != nil {
		return err
	}
	for _, file := range files {
		content, err := os.ReadFile(filepath.Join(from, file))
		if err != nil {
			return err
		}
		target := filepath.Join(to, filepath.FromSlash(file))
		if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
			return err
		}
		if err := os.WriteFile(target, content, 0o600); err != nil {
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
