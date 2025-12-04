package site

import (
	"io/fs"
	"iter"
	"os"
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
		fsys := os.DirFS(s.BaseDir)
		matches, err := fs.Glob(fsys, "*")
		if err != nil {
			iterErr = err
			return
		}

		for _, path := range matches {
			finalPath := filepath.Join(s.BaseDir, path)
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
