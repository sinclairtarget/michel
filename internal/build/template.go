/*
* There are three kinds of templates in Michel:
*
*   1. Page templates (files under the site/ directory)
*   2. Layouts (used if referenced in the YAML frontmatter for a page)
*   3. Partials (pulled in if invoked in a page template)
*
* Layouts and partials must be referenced using their keys. The key for a
* layout or partial is the filepath to that layout or partial relative to the
* layouts/ or partials/ directory respectively, excluding the file extension.
*
* Layouts are rendered in the order they are listed in the YAML frontmatter.
*
* All Michel templates have access to certain Michel data structures exposed
* via the '.' (dot).
 */
package build

import (
	"fmt"
	"html/template"
	"os"

	"github.com/sinclairtarget/michel/internal/util"
)

// A layout or a partial (collectively referred to as "stencils" for lack of a
// better term).
type stencil struct {
	key          string // unique id for layout
	path         string // path it was loaded from
	templateText string
}

type Partial struct {
	stencil
}

func (p Partial) templateName() string {
	return templateName("partials", p.key)
}

type Layout struct {
	stencil
}

func (l Layout) templateName() string {
	return templateName("layouts", l.key)
}

// Returns the namespaced key that will be used to identify the template in the
// final template association / parse tree.
func templateName(namespace string, key string) string {
	return namespace + "/" + key
}

func loadPartials(dir string) ([]Partial, error) {
	partials := []Partial{}

	stencils, err := loadStencils(dir)
	if err != nil {
		return partials, err
	}

	for _, s := range stencils {
		partials = append(partials, Partial{s})
	}
	return partials, nil
}

func loadLayouts(dir string) ([]Layout, error) {
	layouts := []Layout{}

	stencils, err := loadStencils(dir)
	if err != nil {
		return layouts, err
	}

	for _, s := range stencils {
		layouts = append(layouts, Layout{s})
	}
	return layouts, nil
}

// Parse and add all partials to association.
func parsePartials(
	tmpl *template.Template,
	partials []Partial,
) (*template.Template, error) {
	for _, partial := range partials {
		tmpl = tmpl.New(partial.templateName())
		_, err := tmpl.Parse(partial.templateText)
		if err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}

// Parse and add named layouts to association.
func parseLayouts(
	tmpl *template.Template,
	layouts []Layout,
	keys []string,
) (*template.Template, error) {
	lookup := map[string]Layout{}
	for _, layout := range layouts {
		lookup[layout.key] = layout
	}

	for _, key := range keys {
		layout, ok := lookup[key]
		if !ok {
			return nil, fmt.Errorf("layout \"%s\" not found", key)
		}

		tmpl = tmpl.New(layout.templateName())
		_, err := tmpl.Parse(layout.templateText)
		if err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}

func loadStencils(dir string) ([]stencil, error) {
	stencils := []stencil{}

	seq, finish := util.WalkFiles(dir)
	for path := range seq {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		stencil := stencil{
			key:          util.KeyFromPath(dir, path),
			path:         path,
			templateText: string(b),
		}
		stencils = append(stencils, stencil)
	}

	err := finish()
	if err != nil {
		return stencils, err
	}

	return stencils, nil
}
