package site

import (
	"log/slog"
	"os"
	"strings"

	"github.com/sinclairtarget/michel/internal/load"
	"github.com/sinclairtarget/michel/internal/util"
)

type frontmatter struct {
	Layouts []string // Keys naming the layouts that should be used
}

// Metadata for a Michel page available on disk.
type PageMetadata struct {
	Key  string // unique id for the page
	Path string // filepath for this file
	// From frontmatter
	Layouts []string
}

// A page fully loaded into memory.
type Page struct {
	PageMetadata
	TemplateText string
}

func LoadPageMetadata(dir string, path string) (PageMetadata, error) {
	slog.Debug("loading page from disk (metadata only)", "path", path)

	var (
		metadata PageMetadata
		err      error
	)

	if !isPagePath(path) {
		panic("called LoadPageMetadata() on non-page path")
	}

	metadata.Path = path

	f, err := os.Open(metadata.Path)
	if err != nil {
		return metadata, err
	}
	defer f.Close()

	metadata.Key = util.KeyFromPath(dir, metadata.Path)

	result, err := load.ReadFile[frontmatter](
		metadata.Path,
		load.Opts{FrontmatterOnly: true},
	)
	if err != nil {
		return metadata, err
	}

	// Load frontmatter fields
	metadata.Layouts = result.Frontmatter.Layouts

	return metadata, nil
}

// Load page fully.
func LoadPage(m PageMetadata) (Page, error) {
	slog.Debug("loading page from disk", "path", m.Path)

	page := Page{PageMetadata: m}

	result, err := load.ReadFile[frontmatter](m.Path, load.Opts{})
	if err != nil {
		return page, err
	}

	page.TemplateText = result.Text
	return page, nil
}

func isPagePath(path string) bool {
	for _, ext := range []string{".html", ".tmpl", ".gohtml"} {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	return false
}
