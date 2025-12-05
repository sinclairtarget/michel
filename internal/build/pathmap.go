package build

import (
	"path/filepath"

	"github.com/sinclairtarget/michel/internal/site"
)

func mapPagePath(
	path string,
	siteDir string,
	targetDir string,
) (string, error) {
	relative, err := filepath.Rel(siteDir, path)
	if err != nil {
		return "", err
	}

	dirPart := filepath.Dir(relative)
	filename := site.PageNameFromPath(path) + ".html"
	return filepath.Join(targetDir, dirPart, filename), nil
}

func mapAssetPath(
	path string,
	siteDir string,
	targetDir string,
) (string, error) {
	relative, err := filepath.Rel(siteDir, path)
	if err != nil {
		return "", err
	}

	return filepath.Join(targetDir, relative), nil
}
