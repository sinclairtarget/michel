package util

import (
	"errors"
	"fmt"
	"io/fs"
	"iter"
	"os"
	"path/filepath"
	"strings"
)

var compoundExtensions = [...]string{".html.tmpl"}

// Returns the path relative to given base directory with the file extension
// stripped.
//
// Compound extensions commonly used for Go template files are supported.
func RelpathWithoutExt(dir string, path string) string {
	relative, err := filepath.Rel(dir, path)
	if err != nil {
		panic("path could not be made relative to directory")
	}
	dirPart := filepath.Dir(relative)
	base := baseWithoutExt(path)
	return filepath.Join(dirPart, base)
}

func baseWithoutExt(path string) string {
	base := filepath.Base(path)

	for _, ext := range compoundExtensions {
		if strings.HasSuffix(base, ext) {
			return strings.TrimSuffix(base, ext)
		}
	}

	return strings.TrimSuffix(base, filepath.Ext(base))
}

// Returns an iterator over all files under the given directory (including
// under subdirectories).
//
// If the given directory doesn't exist, returns an empty sequence.
func WalkFiles(dir string) (iter.Seq[string], func() error) {
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

		err := filepath.WalkDir(dir, walkFunc)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			iterErr = err
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

// Returns an iterator over all directories under the given directory,
// recursively.
//
// Will also yield the directory itself.
//
// If the directory doesn't exist, returns an empty sequence.
func WalkDirs(dir string) (iter.Seq[string], func() error) {
	var iterErr error
	seq := func(yield func(string) bool) {
		walkFunc := func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				if !yield(path) {
					return fs.SkipAll
				}
			}

			return nil
		}

		err := filepath.WalkDir(dir, walkFunc)
		if err != nil && !errors.Is(err, fs.ErrNotExist) {
			iterErr = err
		}
	}

	finish := func() error {
		if iterErr != nil {
			return fmt.Errorf("failed to walk directories: %w", iterErr)
		}

		return nil
	}

	return seq, finish
}

func IsDir(path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return info.IsDir(), nil
}

// Returns true if the given filepath is with dir or a subdirectory of dir,
// false otherwise.
func IsContained(dir string, path string) bool {
	rel, err := filepath.Rel(dir, path)
	if err != nil {
		panic("path could not be made relative to directory")
	}
	return !strings.HasPrefix(rel, "..")
}
