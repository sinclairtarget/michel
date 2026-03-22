package site

import (
	"fmt"
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

// Makes calling Site.Pages.Get or Site.Assets.Get possible in templates.
type Shim[T any] struct {
	metadata   map[string]T
	collection string // for error messages
}

func (s Site) NumPages() int {
	return len(s.pageMetadata)
}

func (s Site) NumAssets() int {
	return len(s.assetMetadata)
}

func (s Site) Pages() Shim[PageMetadata] {
	return Shim[PageMetadata]{
		metadata:   s.pageMetadata,
		collection: "page",
	}
}

func (s Site) Assets() Shim[AssetMetadata] {
	return Shim[AssetMetadata]{
		metadata:   s.assetMetadata,
		collection: "asset",
	}
}

func (s Shim[T]) Get(key string) (T, error) {
	metadata, ok := s.metadata[key]
	if !ok {
		return metadata, fmt.Errorf(
			"%s with key \"%s\" not found",
			s.collection,
			key,
		)
	}

	return metadata, nil
}

func (s Shim[T]) All() iter.Seq[T] {
	return maps.Values(s.metadata)
}
