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
	Content string   // Key naming content associated with this page
}

// Metadata for a Michel page available on disk.
type PageMetadata struct {
	key      string // unique id for the page
	Filepath string // source filepath for this file
	relURL   string
	absURL   string
	// From frontmatter
	Layouts    []string
	ContentKey string
}

func (m PageMetadata) Key() string { return m.key }

func (m PageMetadata) RelURL() string { return m.relURL }

func (m PageMetadata) AbsURL() string {
	if m.absURL == "" {
		slog.Warn(
			"no AbsURL for page; did you configure baseURL?",
			"key",
			m.Key(),
		)
		return m.relURL
	}

	return m.absURL
}

// A page fully loaded into memory.
type Page struct {
	PageMetadata
	TemplateText string
}

func LoadPageMetadata(
	dir string,
	path string,
	baseURL string,
) (PageMetadata, error) {
	slog.Debug("loading page from disk (metadata only)", "path", path)

	var (
		metadata PageMetadata
		err      error
	)

	if !isPagePath(path) {
		panic("called LoadPageMetadata() on non-page path")
	}

	metadata.key = util.KeyFromPath(dir, path)
	metadata.Filepath = path
	metadata.relURL = RelURL(metadata.key+".html", baseURL)
	if baseURL != "" {
		metadata.absURL = AbsURL(metadata.key+".html", baseURL)
	}

	f, err := os.Open(metadata.Filepath)
	if err != nil {
		return metadata, err
	}
	defer f.Close()

	result, err := load.ReadFile[frontmatter](
		metadata.Filepath,
		load.Opts{FrontmatterOnly: true},
	)
	if err != nil {
		return metadata, err
	}

	// Load frontmatter fields
	metadata.Layouts = result.Frontmatter.Layouts
	metadata.ContentKey = result.Frontmatter.Content

	return metadata, nil
}

// Load page fully.
func LoadPage(m PageMetadata) (Page, error) {
	slog.Debug("loading page from disk", "path", m.Filepath)

	page := Page{PageMetadata: m}

	result, err := load.ReadFile[frontmatter](m.Filepath, load.Opts{})
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
