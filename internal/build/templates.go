package build

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/sinclairtarget/michel/internal/util"
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

	tmpl := template.New("root")
	for _, filename := range matches {
		partialKey := partialKeyFromPath(filename)
		partialTmpl := tmpl.New(partialKey)

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

func partialKeyFromPath(path string) string {
	base := util.BaseWithoutExt(path)
	return "partials/" + base
}

func loadLayouts(
	dir string,
	paths []string,
	existingTmpl *template.Template,
) (*template.Template, error) {
	tmpl := existingTmpl
	for _, path := range paths {
		key := layoutKeyFromPath(path)
		tmpl = tmpl.New(key)

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

func layoutKeyFromPath(path string) string {
	base := util.BaseWithoutExt(path)
	return "layouts/" + base
}

func layoutPathFromKey(key string, layoutsDir string) (string, error) {
	pattern := filepath.Join(
		layoutsDir,
		strings.TrimPrefix(key, "layouts/"),
	) + "*"
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}

	if len(matches) == 0 {
		// Trust that this will fail later when we try to load this nonexistent
		// file.
		return key + defaultLayoutExt, nil
	}

	if len(matches) > 1 {
		return "", fmt.Errorf("more than one match for layout \"%s\"", key)
	}

	return matches[0], nil
}
