package build

import (
	"html/template"
	"os"
	"path/filepath"
)

// Loads all partials templates. Namespaces them with a "partials/" prefix.
//
// Returns nil if there are no partials.
func loadPartials(dir string) (*template.Template, error) {
	pattern := filepath.Join(dir, "*.html.tmpl")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	var tmpl *template.Template
	for _, filename := range matches {
		partialName := partialNameFromFilename(filename)

		var partialTmpl *template.Template
		if tmpl == nil {
			tmpl = template.New(partialName)
			partialTmpl = tmpl
		} else {
			partialTmpl = tmpl.New(partialName)
		}

		b, err := os.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		s := string(b)
		_, err = partialTmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}

func partialNameFromFilename(filename string) string {
	return "partials/article"
}
