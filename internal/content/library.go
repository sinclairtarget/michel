package content

import (
	"fmt"
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
