package page

import (
	"os"
	"strings"

	"github.com/sinclairtarget/michel/internal/load"
	"github.com/sinclairtarget/michel/internal/util"
)

type frontmatter struct {
	Layouts []string // Keys naming the layouts that should be used
}

// An HTML page in the site, possibly templated.
type Page struct {
	Key          string // unique id for the page
	Path         string // path page was loaded from
	TemplateText string
	// From frontmatter
	Layouts []string
}

func LoadPage(dir string, path string) (Page, error) {
	var (
		page Page
		err  error
	)

	if !IsPage(path) {
		panic("called LoadPage() on non-page path")
	}

	page.Path = path

	f, err := os.Open(page.Path)
	if err != nil {
		return page, err
	}
	defer f.Close()

	page.Key = util.KeyFromPath(dir, page.Path)

	result, err := load.ReadFile[frontmatter](page.Path, load.Opts{})
	if err != nil {
		return page, err
	}

	// Load frontmatter fields
	page.Layouts = result.Frontmatter.Layouts

	page.TemplateText = result.Text
	return page, nil
}

func IsPage(path string) bool {
	for _, ext := range []string{".html", ".tmpl", ".gohtml"} {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	return false
}
