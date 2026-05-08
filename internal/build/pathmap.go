package build

import (
	"path/filepath"

	"github.com/sinclairtarget/michel/internal/site"
)

func mapPage(page site.PageMetadata, outDir string) string {
	return filepath.Join(outDir, page.Target())
}

func mapAsset(asset site.AssetMetadata, outDir string) string {
	return filepath.Join(outDir, asset.Target())
}
