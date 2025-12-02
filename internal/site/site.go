package site

import (
	"io/fs"
	"iter"
	"os"
	"path/filepath"
)

type Config struct {
	name string
}

type Site struct {
	baseDir string
	config  Config
}

func Load(dir string) Site {
	config := Config{
		name: "my site",
	}

	return Site{
		baseDir: dir,
		config:  config,
	}
}

func (s Site) Paths() (iter.Seq[string], func() error) {
	var iterErr error
	seq := func(yield func(string) bool) {
		fsys := os.DirFS(s.baseDir)
		matches, err := fs.Glob(fsys, "*")
		if err != nil {
			iterErr = err
			return
		}

		for _, path := range matches {
			finalPath := filepath.Join(s.baseDir, path)
			if !yield(finalPath) {
				return
			}
		}
	}

	finish := func() error {
		return iterErr
	}

	return seq, finish
}
