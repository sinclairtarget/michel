package content

import (
	"path/filepath"
)

type ContentLibrary struct {
	m map[string]Content
}

func (lib ContentLibrary) Get(name string) Content {
	return lib.m[name]
}

func LoadContent(dir string) (ContentLibrary, error) {
	var library ContentLibrary

	// Load plain text content
	contentMap := map[string]Content{}

	pattern := filepath.Join(dir, "*.txt")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return library, err
	}

	for _, match := range matches {
		c, err := LoadFromPlainText(match)
		if err != nil {
			return library, err
		}

		contentMap[c.Name] = c
	}

	library.m = contentMap
	return library, nil
}
