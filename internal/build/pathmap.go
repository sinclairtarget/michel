package build

import (
	"path/filepath"

	"github.com/sinclairtarget/michel/internal/page"
)

func mapPagePath(
	path string,
	pagesDir string,
	targetDir string,
) (string, error) {
	relative, err := filepath.Rel(pagesDir, path)
	if err != nil {
		return "", err
	}

	dirPart := filepath.Dir(relative)
	filename := page.PageKeyFromPath(path) + ".html"
	return filepath.Join(targetDir, dirPart, filename), nil
}

func mapAssetPath(
	path string,
	pagesDir string,
	targetDir string,
) (string, error) {
	relative, err := filepath.Rel(pagesDir, path)
	if err != nil {
		return "", err
	}

	return filepath.Join(targetDir, relative), nil
}
