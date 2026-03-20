package site

import (
	"github.com/sinclairtarget/michel/internal/util"
)

type AssetMetadata struct {
	Key  string
	Path string
}

func NewAsset(dir string, path string) AssetMetadata {
	return AssetMetadata{
		Key:  util.KeyFromPath(dir, path),
		Path: path,
	}
}
