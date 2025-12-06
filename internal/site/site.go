package site

import (
	"fmt"
	"io/fs"
	"iter"
	"path/filepath"
)

type Config struct {
	Name string
}

type Site struct {
	BaseDir string
	Config  Config
}

func Load(dir string) Site {
	config := Config{
		Name: "my site",
	}

	return Site{
		BaseDir: dir,
		Config:  config,
	}
}

func (s Site) Paths() (iter.Seq[string], func() error) {
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

		iterErr = filepath.WalkDir(s.BaseDir, walkFunc)
		if iterErr != nil {
			return
		}
	}

	finish := func() error {
		if iterErr != nil {
			return fmt.Errorf("failed to iterate site paths: %w", iterErr)
		}

		return nil
	}

	return seq, finish
}
