package build

import (
	"path/filepath"

	"github.com/sinclairtarget/michel/internal/content"
)

type ContentLibrary struct {
	m map[string]content.Content
}

func (lib ContentLibrary) Get(name string) content.Content {
	return lib.m[name]
}

func loadContent(dir string) (ContentLibrary, error) {
	var library ContentLibrary

	// Load plain text content
	contentMap := map[string]content.Content{}

	pattern := filepath.Join(dir, "*.txt")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return library, err
	}

	for _, match := range matches {
		c, err := content.LoadFromPlainText(match)
		if err != nil {
			return library, err
		}

		contentMap[c.Name] = c
	}

	library.m = contentMap
	return library, nil
}
