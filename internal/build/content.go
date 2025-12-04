package build

import (
	"path/filepath"

	"github.com/sinclairtarget/michel/internal/content"
)

func loadContent(dir string) (map[string]content.Content, error) {
	// Load plain text content
	contentMap := map[string]content.Content{}

	pattern := filepath.Join(dir, "*.txt")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return contentMap, err
	}

	for _, match := range matches {
		c, err := content.LoadFromPlainText(match)
		if err != nil {
			return contentMap, err
		}

		contentMap[c.Name] = c
	}

	return contentMap, nil
}
