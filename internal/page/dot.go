package page

import (
	"github.com/sinclairtarget/michel/internal/config"
	"github.com/sinclairtarget/michel/internal/content"
)

// Defines the data structures available for access via '.' in Michel
// templates.
//
// This is basically the public API of Michel.
type Dot struct {
	Config  *config.Config
	Content *content.Collection
}
