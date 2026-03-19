// Represents content loaded from the content directory.
package content

import (
	"fmt"
	"os"
	"time"

	"github.com/sinclairtarget/michel/internal/content/myst"
	"github.com/sinclairtarget/michel/internal/frontmatter"
	"github.com/sinclairtarget/michel/internal/util"
)

type contentFrontmatter struct {
	Title       string
	Description string
	Date        string
}

// If Date is missing, fallback to Jan 1, year 1.
//
// Hugo does this and it seems like a reasonable choice.
var fallbackDate time.Time = time.Date(1, 1, 1, 0, 0, 0, 0, time.Local)

// Parses the Date field into a time.Time.
//
// We use this to represent the calendar date associated with the content. But
// this is just a calendar date and not an instant in time.
func (c contentFrontmatter) ParsedDate() (time.Time, error) {
	if c.Date == "" {
		return fallbackDate, nil
	}

	t, err := time.Parse("2006-01-02", c.Date)
	if err != nil {
		return fallbackDate, err
	}

	return t, nil
}

// A file with content for the site.
type Content struct {
	Key  string // unique id for the content
	Path string // path content was loaded from
	Root *myst.Node
	// From frontmatter
	Title       string
	Description string
	Date        time.Time
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

	result, err := frontmatter.ReadFile[contentFrontmatter](f)
	if err != nil {
		return content, err
	}

	// Load frontmatter fields
	content.Title = result.Frontmatter.Title
	content.Description = result.Frontmatter.Description

	content.Date, err = result.Frontmatter.ParsedDate()
	if err != nil {
		return content, fmt.Errorf(
			"failed to parse frontmatter date in content file \"%s\": %w",
			err,
		)
	}

	// Parse MyST
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
