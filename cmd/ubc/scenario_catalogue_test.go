package main

import (
	"io/fs"
	"path/filepath"
	"strings"
)

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

func missingFromIndex(index string, files []string) []string {
	var missing []string
	for _, file := range files {
		if !strings.Contains(index, "]("+file+")") {
			missing = append(missing, file)
		}
	}
	return missing
}
