package fileext

import (
	"path/filepath"
	"strings"
)

var allowedCompoundExtensions = [...]string{".html.tmpl", ".go.html"}

// Returns the filename in a path, removing all leading directories and the
// extension.
//
// Compound extensions commonly used for Go template files are supported.
func BaseWithoutExt(path string) string {
	base := filepath.Base(path)

	for _, ext := range allowedCompoundExtensions {
		if strings.HasSuffix(base, ext) {
			return strings.TrimSuffix(base, ext)
		}
	}

	return strings.TrimSuffix(base, filepath.Ext(base))
}
