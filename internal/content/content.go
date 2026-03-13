// Represents content loaded from the content directory.
package content

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sinclairtarget/michel/internal/content/myst"
	"github.com/sinclairtarget/michel/internal/frontmatter"
	"github.com/sinclairtarget/michel/internal/util"
)

type ContentFrontmatter struct {
	Title string
}

// A file with content for the site.
type Content struct {
	Key         string // unique id for the content
	Path        string // path content was loaded from
	Frontmatter ContentFrontmatter
	Root        *myst.Node
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
	content.Root, err = myst.Parse(result.Text)
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

	return filepath.Join(filepath.Dir(relPath), util.BaseWithoutExt(path))
}
