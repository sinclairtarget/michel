package site

import (
	"os"
	"strings"

	"github.com/sinclairtarget/michel/internal/frontmatter"
	"github.com/sinclairtarget/michel/internal/util/fileext"
)

type PageFrontmatter struct {
	Layouts []string
}

// Let user write just the simple layout name, but adjust here because the proper
// name includes the `layouts` prefix.
func (pm PageFrontmatter) LayoutsFullName() []string {
	var adjustedNames []string
	for _, name := range pm.Layouts {
		adjustedNames = append(adjustedNames, "layouts/"+name)
	}

	return adjustedNames
}

type Page struct {
	Path         string
	Name         string
	Frontmatter  PageFrontmatter
	TemplateText string
}

func IsPage(path string) bool {
	for _, ext := range []string{".html", ".tmpl", ".gohtml"} {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	return false
}

func LoadPage(path string) (Page, error) {
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

	page.Name = PageNameFromPath(page.Path)

	result, err := frontmatter.ReadFile[PageFrontmatter](f)
	if err != nil {
		return page, err
	}

	page.Frontmatter = result.Frontmatter
	page.TemplateText = result.Text
	return page, nil
}

func PageNameFromPath(path string) string {
	return fileext.BaseWithoutExt(path)
}
