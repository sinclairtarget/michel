package content

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/sinclairtarget/michel/internal/frontmatter"
	"github.com/sinclairtarget/michel/internal/util/fileext"
)

type Frontmatter struct {
	Title string
}

type Content struct {
	Path        string
	Name        string
	Frontmatter Frontmatter
	Html        string
}

func (c Content) Body() template.HTML {
	return template.HTML(c.Html)
}

func IsPlaintext(path string) bool {
	return strings.HasSuffix(path, ".txt")
}

func LoadFromPlainText(path string) (Content, error) {
	var (
		content Content
		err     error
	)

	if !IsPlaintext(path) {
		panic("called LoadFromPlainText() on non-plain text file")
	}

	content.Path, err = filepath.Abs(path)
	if err != nil {
		return content, err
	}

	content.Name = contentNameFromPath(path)

	f, err := os.Open(content.Path)
	if err != nil {
		return content, err
	}
	defer f.Close()

	result, err := frontmatter.ReadFile[Frontmatter](f)
	if err != nil {
		return content, err
	}

	content.Frontmatter = result.Frontmatter
	content.Html, err = parsePlainText(result.Text)
	if err != nil {
		return content, fmt.Errorf(
			"failed to parse content file \"%s\": %w",
			path,
			err,
		)
	}

	return content, nil
}

func contentNameFromPath(path string) string {
	return fileext.BaseWithoutExt(path)
}
