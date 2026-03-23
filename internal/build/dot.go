package build

import (
	"fmt"
	"html/template"
	"io"
	"iter"
	"slices"
	"time"

	"github.com/sinclairtarget/michel/internal/config"
	"github.com/sinclairtarget/michel/internal/content"
	"github.com/sinclairtarget/michel/internal/content/myst"
	"github.com/sinclairtarget/michel/internal/site"
	"github.com/sinclairtarget/michel/internal/util"
)

// Defines the data structures available for access via '.' in Michel
// templates.
//
// This is basically the public API of Michel.
type Dot struct {
	Config  config.Config
	Content content.Corpus
	Site    site.Site
	Page    site.PageMetadata // Currently rendering page
	Now     time.Time         // Should be when the build started
}

// Defines the functions available in Michel templates.
func (d Dot) funcMap(tmpl *template.Template, w io.Writer) template.FuncMap {
	return template.FuncMap{
		"renderHTML": myst.RenderHTML,
		"partial": func(key string, data any) error {
			return executePartial(tmpl, w, key, data)
		},
		"select":  selectAny,
		"reject":  rejectAny,
		"collect": collectAny,
		"reverse": reverseAny,
		"relURL": func(suffix string) string {
			return site.RelURL(suffix, d.Config.BaseURL)
		},
		"absURL": func(suffix string) string {
			return site.AbsURL(suffix, d.Config.BaseURL)
		},
	}
}

func executePartial(
	tmpl *template.Template,
	w io.Writer,
	key string,
	data any,
) error {
	execName := templateName("partials", key)
	return tmpl.ExecuteTemplate(w, execName, data)
}

func selectAny(field string, pattern string, seq any) iter.Seq[util.Keyed] {
	switch v := seq.(type) {
	case iter.Seq[util.Keyed]:
		return util.Select(v, field, pattern)
	case iter.Seq[content.Entry]:
		return util.CoerceSeq[content.Entry, util.Keyed](
			util.Select[content.Entry](v, field, pattern),
		)
	case iter.Seq[site.PageMetadata]:
		return util.CoerceSeq[site.PageMetadata, util.Keyed](
			util.Select[site.PageMetadata](v, field, pattern),
		)
	case iter.Seq[site.AssetMetadata]:
		return util.CoerceSeq[site.AssetMetadata, util.Keyed](
			util.Select[site.AssetMetadata](v, field, pattern),
		)
	default:
		msg := fmt.Sprintf("select used with unknown type %T", v)
		panic(msg)
	}
}

func rejectAny(field string, pattern string, seq any) iter.Seq[util.Keyed] {
	switch v := seq.(type) {
	case iter.Seq[util.Keyed]:
		return util.Reject(v, field, pattern)
	case iter.Seq[content.Entry]:
		return util.CoerceSeq[content.Entry, util.Keyed](
			util.Reject[content.Entry](v, field, pattern),
		)
	case iter.Seq[site.PageMetadata]:
		return util.CoerceSeq[site.PageMetadata, util.Keyed](
			util.Reject[site.PageMetadata](v, field, pattern),
		)
	case iter.Seq[site.AssetMetadata]:
		return util.CoerceSeq[site.AssetMetadata, util.Keyed](
			util.Reject[site.AssetMetadata](v, field, pattern),
		)
	default:
		msg := fmt.Sprintf("reject used with unknown type %T", v)
		panic(msg)
	}
}

func collectAny(seq any) []util.Keyed {
	switch v := seq.(type) {
	case iter.Seq[util.Keyed]:
		return slices.Collect(v)
	case iter.Seq[content.Entry]:
		return slices.Collect(
			util.CoerceSeq[content.Entry, util.Keyed](v),
		)
	case iter.Seq[site.PageMetadata]:
		return slices.Collect(
			util.CoerceSeq[site.PageMetadata, util.Keyed](v),
		)
	case iter.Seq[site.AssetMetadata]:
		return slices.Collect(
			util.CoerceSeq[site.AssetMetadata, util.Keyed](v),
		)
	default:
		msg := fmt.Sprintf("collect used with unknown type %T", v)
		panic(msg)
	}
}

func reverseAny(seq any) iter.Seq[util.Keyed] {
	return func(yield func(util.Keyed) bool) {
		for _, elem := range slices.Backward(collectAny(seq)) {
			if !yield(elem) {
				return
			}
		}
	}
}
