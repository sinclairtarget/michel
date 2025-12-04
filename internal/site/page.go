package site

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/sinclairtarget/michel/internal/frontmatter"
)

type PageFrontmatter struct {
	Layouts []string
}

type Page struct {
	Path         string
	Frontmatter  PageFrontmatter
	TemplateText string
}

func IsPage(path string) bool {
	return strings.HasSuffix(path, ".html") || strings.HasSuffix(path, ".tmpl")
}

func LoadPage(path string) (Page, error) {
	var (
		page Page
		err  error
	)

	if !IsPage(path) {
		panic("called LoadPage() on non-page path")
	}

	page.Path, err = filepath.Abs(path)
	if err != nil {
		return page, err
	}

	f, err := os.Open(page.Path)
	if err != nil {
		return page, err
	}
	defer f.Close()

	result, err := frontmatter.ReadFile[PageFrontmatter](f)
	if err != nil {
		return page, err
	}

	page.Frontmatter = result.Frontmatter
	page.TemplateText = result.Text
	return page, nil
}
