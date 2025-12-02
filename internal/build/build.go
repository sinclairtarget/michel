package build

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	sitePkg "github.com/sinclairtarget/michel/internal/site"
)

type Options struct {
	SiteDir     string
	TargetDir   string
	ShouldClean bool
}

func Build(logger *slog.Logger, options Options) error {
	if options.ShouldClean {
		logger.Debug("cleaning target directory")
		err := clean(options.TargetDir)
		if err != nil {
			return fmt.Errorf("failed to clean target directory: %v", err)
		}
	}

	logger.Debug("loading site")
	site := sitePkg.Load(options.SiteDir)

	logger.Debug("processing site pages")
	seq, finish := site.Paths()
	for sitePath := range seq {
		targetPath, err := target(options.SiteDir, options.TargetDir, sitePath)
		if err != nil {
			return fmt.Errorf("could not map path: %v", err)
		}

		err = process(sitePath, targetPath)
		if err != nil {
			return fmt.Errorf("failed to process \"%s\": %v", sitePath, err)
		}
	}

	err := finish()
	if err != nil {
		return err
	}

	return nil
}

func clean(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}

	err = os.Mkdir(dir, 0o755)
	if err != nil {
		return err
	}

	return nil
}

// Returns output path under target dir given path in source dir.
func target(siteDir string, targetDir string, path string) (string, error) {
	relative, err := filepath.Rel(siteDir, path)
	if err != nil {
		return "", err
	}

	return filepath.Join(targetDir, relative), nil
}

func process(sourcePath string, targetPath string) error {
	source, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer source.Close()

	target, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer target.Close()

	_, err = io.Copy(target, source)
	return err
}
