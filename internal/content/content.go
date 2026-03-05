package content

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"

	atrus "github.com/sinclairtarget/libatrus-go"

	"github.com/sinclairtarget/michel/internal/frontmatter"
	"github.com/sinclairtarget/michel/internal/util/fileext"
)

type ContentFrontmatter struct {
	Title string
}

// A file with content for the site.
type Content struct {
	Key         string // unique id for the content
	Path        string // path content was loaded from
	Frontmatter ContentFrontmatter
	Html        string
}

func (c Content) Body() template.HTML {
	return template.HTML(c.Html)
}

func LoadFromMarkdown(contentDir string, path string) (Content, error) {
	var (
		content Content
		err     error
	)

	content.Path = path
	content.Key = contentKeyFromPath(contentDir, content.Path)

	f, err := os.Open(content.Path)
	if err != nil {
		return content, err
	}
	defer f.Close()

	result, err := frontmatter.ReadFile[ContentFrontmatter](f)
	if err != nil {
		return content, err
	}

	content.Frontmatter = result.Frontmatter
	content.Html, err = parseMystMarkdown(result.Text)
	if err != nil {
		return content, fmt.Errorf(
			"failed to parse content file \"%s\": %w",
			path,
			err,
		)
	}

	return content, nil
}

func contentKeyFromPath(contentDir string, path string) string {
	relPath, err := filepath.Rel(contentDir, path)
	if err != nil {
		panic(err)
	}

	return filepath.Join(filepath.Dir(relPath), fileext.BaseWithoutExt(path))
}

func parseMystMarkdown(text string) (string, error) {
	ast, err := atrus.ParseAST(text)
	if err != nil {
		return "", fmt.Errorf("libatrus parse error: %w", err)
	}

	html, err := atrus.RenderHTML(ast)
	if err != nil {
		return "", fmt.Errorf("libatrus render error: %w", err)
	}

	return html, nil
}
