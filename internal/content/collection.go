package content

import (
	"fmt"
	"iter"
	"maps"

	"github.com/sinclairtarget/michel/internal/util"
)

// Collection of site content loaded from the content directory.
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
func LoadCollection(dir string) (Collection, error) {
	var collection Collection

	contentMap := map[string]Content{}
	seq, finish := util.WalkFiles(dir)
	for path := range seq {
		c, err := LoadFromMarkdown(dir, path)
		if err != nil {
			return collection, err
		}

		contentMap[c.Key] = c
	}

	err := finish()
	if err != nil {
		return collection, err
	}

	collection.loaded = contentMap
	return collection, nil
}
