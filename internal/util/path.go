package util

import (
	"fmt"
	"io/fs"
	"iter"
	"path/filepath"
	"strings"
)

var allowedCompoundExtensions = [...]string{".html.tmpl", ".go.html"}

// Returns the filename in a path, removing all leading directories and the
// extension.
//
// Compound extensions commonly used for Go template files are supported.
func BaseWithoutExt(path string) string {
	base := filepath.Base(path)

	for _, ext := range allowedCompoundExtensions {
		if strings.HasSuffix(base, ext) {
			return strings.TrimSuffix(base, ext)
		}
	}

	return strings.TrimSuffix(base, filepath.Ext(base))
}

// Returns a key for a file loaded from disk.
//
// These keys are used throughout Michel to identify content, templates, etc.
//
// The key is the path to the file, relative to the containing directory, with
// no extension.
func KeyFromPath(dir string, path string) string {
	relative, err := filepath.Rel(dir, path)
	if err != nil {
		panic("path could not be made relative to directory")
	}
	dirPart := filepath.Dir(relative)
	base := BaseWithoutExt(path)
	return filepath.Join(dirPart, base)
}

// Returns an iterator over all paths under the given directory.
func WalkPaths(dir string) (iter.Seq[string], func() error) {
	var iterErr error
	seq := func(yield func(string) bool) {
		walkFunc := func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				if !yield(path) {
					return fs.SkipAll
				}
			}

			return nil
		}

		iterErr = filepath.WalkDir(dir, walkFunc)
		if iterErr != nil {
			return
		}
	}

	finish := func() error {
		if iterErr != nil {
			return fmt.Errorf("failed to walk paths: %w", iterErr)
		}

		return nil
	}

	return seq, finish
}
