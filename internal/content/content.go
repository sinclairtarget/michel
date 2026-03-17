// Represents content loaded from the content directory.
package content

import (
	"fmt"
	"os"

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

// Loads content file into memory, parsing the markdown.
func LoadFromMarkdown(contentDir string, path string) (Content, error) {
	var (
		content Content
		err     error
	)

	content.Path = path
	content.Key = util.KeyFromPath(contentDir, content.Path)

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
