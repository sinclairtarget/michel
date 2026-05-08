package site

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/sinclairtarget/michel/internal/load"
	"github.com/sinclairtarget/michel/internal/util"
)

type frontmatter struct {
	Key     string
	Layouts []string // Keys naming the layouts that should be used
}

// Metadata for a Michel page available on disk.
//
// The key that uniquely identifies a page is determined thusly:
//  1. If the frontmatter for the page includes an explicit key, use that.
//  2. If the page is an HTML file (according to the extension), then use the
//     filepath of the template relative to the site/ directory with the
//     extension stripped off.
//  3. Otherwise, the key is the same as in 2 except the extension is
//     retained. A trailing `.tmpl`, as in `foo.xml.tmpl`, is still removed.
//     (This allows index.html and index.xml to exist without collision).
type PageMetadata struct {
	key      string // unique id for the page
	Filepath string // source filepath for this file
	target   string // output filepath
	relURL   string
	absURL   string
	// From frontmatter
	Layouts []string
}

func (m PageMetadata) Key() string { return m.key }

func (m PageMetadata) RelURL() string { return m.relURL }

func (m PageMetadata) AbsURL() string {
	if m.absURL == "" {
		slog.Warn(
			"page did not have an AbsURL; did you configure baseURL?",
			"key",
			m.Key(),
		)
		return m.relURL
	}

	return m.absURL
}

// Output filepath.
func (m PageMetadata) Target() string { return m.target }

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
		metadata   PageMetadata
		defaultKey string
		err        error
	)

	if !isPagePath(path) {
		panic("called LoadPageMetadata() on non-page path")
	}

	metadata.Filepath = path

	if isHTMLPagePath(path) {
		defaultKey = util.RelpathWithoutExt(dir, path)
		metadata.target = defaultKey + ".html"
		metadata.relURL = RelURL(metadata.target, baseURL)
		if baseURL != "" {
			metadata.absURL = AbsURL(metadata.target, baseURL)
		}
	} else {
		defaultKey, err = filepath.Rel(dir, path) // keep extension
		if err != nil {
			panic("page path could not be made relative to site directory")
		}

		// except .tmpl, we want to get rid of that if it's there
		if strings.HasSuffix(defaultKey, ".tmpl") {
			defaultKey = strings.TrimSuffix(defaultKey, ".tmpl")
		}

		metadata.target = defaultKey
		metadata.relURL = RelURL(metadata.target, baseURL)
		if baseURL != "" {
			metadata.absURL = AbsURL(metadata.target, baseURL)
		}
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
	metadata.key = defaultKey
	if result.Frontmatter.Key != "" {
		metadata.key = result.Frontmatter.Key
	}

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
	for _, ext := range []string{".html", ".html.tmpl", ".xml", ".xml.tmpl"} {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	return false
}

func isHTMLPagePath(path string) bool {
	for _, ext := range []string{".html", ".html.tmpl"} {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	return false
}
