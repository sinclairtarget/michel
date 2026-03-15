package page

import (
	"html/template"

	"github.com/sinclairtarget/michel/internal/config"
	"github.com/sinclairtarget/michel/internal/content"
	"github.com/sinclairtarget/michel/internal/content/myst"
)

// Defines the data structures available for access via '.' in Michel
// templates.
//
// This is basically the public API of Michel.
type Dot struct {
	Config  *config.Config
	Content *content.Collection
}

func (d Dot) FuncMap() template.FuncMap {
	return template.FuncMap{
		"html": myst.RenderHTML,
	}
}
