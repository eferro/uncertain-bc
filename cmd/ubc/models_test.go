package main

import (
	"io/fs"
	"os"
	"path/filepath"

	"github.com/rogpeppe/go-internal/testscript"
)

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

// modelFiles relies on WalkDir visiting entries in lexical order, so the
// result is sorted.
func modelFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() || filepath.Ext(path) != ".ubc" {
			return err
		}
		rel, err := filepath.Rel(root, path)
		files = append(files, filepath.ToSlash(rel))
		return err
	})
	return files, err
}
