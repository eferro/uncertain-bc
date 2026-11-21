package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/rogpeppe/go-internal/testscript"
)

func TestMain(m *testing.M) {
	testscript.Main(m, map[string]func(){
		"ubc": main,
	})
}

func TestScripts(t *testing.T) {
	testscript.Run(t, testscript.Params{Dir: "testdata/script", Setup: copyExamples})
}

func copyExamples(env *testscript.Env) error {
	examples, err := filepath.Glob("../../examples/*.ubc")
	if err != nil {
		return err
	}
	for _, path := range examples {
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(env.WorkDir, filepath.Base(path)), content, 0o600); err != nil {
			return err
		}
	}
	return nil
}
