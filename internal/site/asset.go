package site

import (
	"github.com/sinclairtarget/michel/internal/util"
)

type AssetMetadata struct {
	key  string
	Path string
}

func (m AssetMetadata) Key() string { return m.key }

func NewAsset(dir string, path string) AssetMetadata {
	return AssetMetadata{
		key:  util.KeyFromPath(dir, path),
		Path: path,
	}
}
