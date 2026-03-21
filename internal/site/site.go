package site

import (
	"iter"
	"maps"

	"github.com/sinclairtarget/michel/internal/util"
)

type Site struct {
	pageMetadata  map[string]PageMetadata
	assetMetadata map[string]AssetMetadata
}

func LoadSite(dir string) (Site, error) {
	site := Site{
		pageMetadata:  map[string]PageMetadata{},
		assetMetadata: map[string]AssetMetadata{},
	}

	seq, finish := util.WalkFiles(dir)
	for path := range seq {
		if isPagePath(path) {
			m, err := LoadPageMetadata(dir, path)
			if err != nil {
				return site, err
			}

			site.pageMetadata[m.Key()] = m
		} else {
			m := NewAsset(dir, path)
			site.assetMetadata[m.Key()] = m
		}
	}

	err := finish()
	if err != nil {
		return site, err
	}

	return site, nil
}

func (s Site) Pages() iter.Seq[PageMetadata] {
	return maps.Values(s.pageMetadata)
}

func (s Site) Assets() iter.Seq[AssetMetadata] {
	return maps.Values(s.assetMetadata)
}

func (s Site) NumPages() int {
	return len(s.pageMetadata)
}

func (s Site) NumAssets() int {
	return len(s.assetMetadata)
}
