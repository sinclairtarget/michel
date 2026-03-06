package content

import (
	"fmt"
	"io/fs"
	"iter"
	"maps"
	"path/filepath"
)

// Collection of site content loaded from the content directory
type Collection struct {
	loaded map[string]Content
}

func (c Collection) Get(key string) (Content, error) {
	content, ok := c.loaded[key]
	if !ok {
		return content, fmt.Errorf("content with key \"%s\" not found", key)
	}

	return content, nil
}

func (c Collection) All() iter.Seq[Content] {
	return maps.Values(c.loaded)
}

// Load all content into memory.
func LoadAllContent(dir string) (Collection, error) {
	var collection Collection

	contentMap := map[string]Content{}

	walkFunc := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			c, err := LoadFromMarkdown(dir, path)
			if err != nil {
				return err
			}

			contentMap[c.Key] = c
		}

		return nil
	}

	err := filepath.WalkDir(dir, walkFunc)
	if err != nil {
		return collection, err
	}

	collection.loaded = contentMap
	return collection, nil
}
