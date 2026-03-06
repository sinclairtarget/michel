package util

import (
	"fmt"
	"io/fs"
	"iter"
	"path/filepath"
)

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
