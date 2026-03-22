package build

import (
	"path/filepath"

	"github.com/sinclairtarget/michel/internal/site"
)

func mapPage(page site.PageMetadata, targetDir string) string {
	targetFilepath := page.Key() + ".html"
	return filepath.Join(targetDir, targetFilepath)
}

func mapAsset(asset site.AssetMetadata, targetDir string) string {
	targetFilepath := asset.Key()
	return filepath.Join(targetDir, targetFilepath)
}
