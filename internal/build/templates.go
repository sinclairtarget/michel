package build

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/sinclairtarget/michel/internal/util/fileext"
)

var defaultLayoutExt = ".html.tmpl"

// Loads all partials templates. Namespaces them with a "partials/" prefix.
//
// Returns nil if there are no partials.
func loadPartials(dir string) (*template.Template, error) {
	pattern := filepath.Join(dir, "*")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	var tmpl *template.Template
	for _, filename := range matches {
		partialName := partialNameFromPath(filename)

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

func partialNameFromPath(path string) string {
	base := fileext.BaseWithoutExt(path)
	return "partials/" + base
}

func loadLayouts(
	dir string,
	paths []string,
	existingTmpl *template.Template,
) (*template.Template, error) {
	tmpl := existingTmpl
	for _, path := range paths {
		name := layoutNameFromPath(path)
		if tmpl != nil {
			tmpl = tmpl.New(name)
		} else {
			tmpl = template.New(name)
		}

		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		s := string(b)
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}

func layoutNameFromPath(path string) string {
	base := fileext.BaseWithoutExt(path)
	return "layouts/" + base
}

func layoutPathFromName(name string, layoutsDir string) (string, error) {
	pattern := filepath.Join(
		layoutsDir,
		strings.TrimPrefix(name, "layouts/"),
	) + "*"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		// Trust that this will fail later when we try to load this nonexistent
		// file.
		return name + defaultLayoutExt, nil
	}

	if len(matches) > 1 {
		return "", fmt.Errorf("more than one match for layout \"%s\"", name)
	}

	return matches[0], nil
}
