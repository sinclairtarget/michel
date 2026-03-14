package build

import (
	"path/filepath"

	"github.com/sinclairtarget/michel/internal/util"
)

func mapPagePath(path string, pagesDir string, targetDir string) string {
	targetFilepath := util.KeyFromPath(pagesDir, path) + ".html"
	return filepath.Join(targetDir, targetFilepath)
}

func mapAssetPath(path string, pagesDir string, targetDir string) string {
	relative, err := filepath.Rel(pagesDir, path)
	if err != nil {
		panic("asset path could not be made relative to given directory")
	}

	return filepath.Join(targetDir, relative)
}
