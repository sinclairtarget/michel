package build

import (
	"path/filepath"

	"github.com/sinclairtarget/michel/internal/util/fileext"
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
	pageName := fileext.BaseWithoutExt(path)
	filename := pageName + ".html"
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
