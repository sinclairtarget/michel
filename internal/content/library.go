package content

import (
	"fmt"
	"io/fs"
	"path/filepath"
)

type ContentLibrary struct {
	m map[string]Content
}

func (lib ContentLibrary) Get(name string) (Content, error) {
	content, ok := lib.m[name]
	if !ok {
		return content, fmt.Errorf("content \"%s\" not found", name)
	}

	return content, nil
}

// Load all content into memory.
func LoadContent(dir string) (ContentLibrary, error) {
	var library ContentLibrary

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

			contentMap[c.Name] = c
		}

		return nil
	}

	err := filepath.WalkDir(dir, walkFunc)
	if err != nil {
		return library, err
	}

	library.m = contentMap
	return library, nil
}
