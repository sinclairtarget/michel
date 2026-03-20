package template

import (
	"html/template"
	"io"
	"time"

	"github.com/sinclairtarget/michel/internal/config"
	"github.com/sinclairtarget/michel/internal/content"
	"github.com/sinclairtarget/michel/internal/content/myst"
	"github.com/sinclairtarget/michel/internal/page"
	"github.com/sinclairtarget/michel/internal/util"
)

// Defines the data structures available for access via '.' in Michel
// templates.
//
// This is basically the public API of Michel.
type Dot struct {
	Config  *config.Config
	Content *util.Collection[content.Content]
	Page    *page.Page
	Now     time.Time // Should be when the build started
}

// Defines the functions available in Michel templates.
func (d Dot) FuncMap(tmpl *template.Template, w io.Writer) template.FuncMap {
	return template.FuncMap{
		"html": myst.RenderHTML,
		"partial": func(key string, data any) error {
			return executePartial(tmpl, w, key, data)
		},
	}
}

func executePartial(
	tmpl *template.Template,
	w io.Writer,
	key string,
	data any,
) error {
	execName := TemplateName("partials", key)
	return tmpl.ExecuteTemplate(w, execName, data)
}
