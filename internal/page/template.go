package page

import (
	"fmt"
	"html/template"
	"os"

	"github.com/sinclairtarget/michel/internal/util"
)

// A layout or a partial (collectively referred to as "stencils" for lack of a
// better term).
type stencil struct {
	Key          string // unique id for layout
	Path         string // path it was loaded from
	TemplateText string
}

type Partial struct {
	stencil
}

func (p Partial) TemplateName() string {
	return TemplateName("partials", p.Key)
}

type Layout struct {
	stencil
}

func (l Layout) TemplateName() string {
	return TemplateName("layouts", l.Key)
}

// Returns the namespaced key that will be used to identify the template in the
// final template association / parse tree.
func TemplateName(namespace string, key string) string {
	return namespace + "/" + key
}

func LoadPartials(dir string) ([]Partial, error) {
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

func LoadLayouts(dir string) ([]Layout, error) {
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
func AddPartials(
	tmpl *template.Template,
	partials []Partial,
) (*template.Template, error) {
	for _, partial := range partials {
		tmpl = tmpl.New(partial.TemplateName())
		_, err := tmpl.Parse(partial.TemplateText)
		if err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}

// Parse and add named layouts to association.
func AddLayouts(
	tmpl *template.Template,
	layouts []Layout,
	keys []string,
) (*template.Template, error) {
	lookup := map[string]Layout{}
	for _, layout := range layouts {
		lookup[layout.Key] = layout
	}

	for _, key := range keys {
		layout, ok := lookup[key]
		if !ok {
			return nil, fmt.Errorf("layout \"%s\" not found", key)
		}

		tmpl = tmpl.New(layout.TemplateName())
		_, err := tmpl.Parse(layout.TemplateText)
		if err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}

func loadStencils(dir string) ([]stencil, error) {
	stencils := []stencil{}

	seq, finish := util.WalkPaths(dir)
	for path := range seq {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		stencil := stencil{
			Key:          util.KeyFromPath(dir, path),
			Path:         path,
			TemplateText: string(b),
		}
		stencils = append(stencils, stencil)
	}

	err := finish()
	if err != nil {
		return stencils, err
	}

	return stencils, nil
}
