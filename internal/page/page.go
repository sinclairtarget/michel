/*
* Package page implements all functionality for handling HTML pages and
* templating.
*
* In Michel, all files under the site/ directory get copied into the output
* directory during a build. Files that have an extension of .html, .tmpl, or
* .gohtml are considered pages. (Files that do not are considered assets.) All
* pages are run through the templating system.
*
* There are three kinds of templates in Michel:
*
*   1. Page templates (files under the site/ directory)
*   2. Layouts (used if referenced in the YAML frontmatter for a page)
*   3. Partials (pulled in if invoked in a page template)
*
* Layouts and partial must be referenced using their keys. The key for a
* layout or partial is the filepath to that layout or partial relative to the
* layouts/ or partials/ directory respectively, excluding the file extension.
*
* Layouts are rendered in the order they are listed in the YAML frontmatter.
*
* All Michel templates have access to certain Michel data structures exposed
* via the '.' (dot).
 */
package page

import (
	"os"
	"strings"

	"github.com/sinclairtarget/michel/internal/frontmatter"
	"github.com/sinclairtarget/michel/internal/util"
)

type PageFrontmatter struct {
	Layouts []string // Keys naming the layouts that should be used
}

// An HTML page in the site, possibly templated.
type Page struct {
	Key          string // unique id for the page
	Path         string // path page was loaded from
	Frontmatter  PageFrontmatter
	TemplateText string
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

	result, err := frontmatter.ReadFile[PageFrontmatter](f)
	if err != nil {
		return page, err
	}

	page.Frontmatter = result.Frontmatter
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
