// Represents content loaded from the content directory.
package content

import (
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"time"

	"github.com/sinclairtarget/michel/internal/content/myst"
	"github.com/sinclairtarget/michel/internal/load"
	"github.com/sinclairtarget/michel/internal/util"
)

// Frontmatter loaded from the beginning of a content file.
type frontmatter struct {
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
func (f frontmatter) ParsedDate() (time.Time, error) {
	if f.Date == "" {
		return fallbackDate, nil
	}

	t, err := time.ParseInLocation("2006-01-02", f.Date, time.Local)
	if err != nil {
		return fallbackDate, err
	}

	return t, nil
}

// Metadata describing a piece of Michel content available on disk.
type Metadata struct {
	key      string // unique id for the content
	Filepath string // source filepath for this file
	// From frontmatter
	Title       string
	Description string
	Date        time.Time
}

func (m Metadata) Key() string { return m.key }

// Slugifies the title and returns it.
//
// Returns an error if the title is an empty string.
//
// TODO: Improve this. Match Hugo behavior.
func (m Metadata) Slug() (string, error) {
	if m.Title == "" {
		return "", fmt.Errorf(
			"cannot slugify without title for content \"%s\"",
			m.Key(),
		)
	}

	lowered := strings.ToLower(m.Title)

	punct := regexp.MustCompile(`[!?'":]+`)
	unpunctuated := punct.ReplaceAllString(lowered, "")
	hyphenated := strings.ReplaceAll(unpunctuated, " ", "-")
	return hyphenated, nil
}

// Content fully loaded into memory and parsed.
type Content struct {
	Metadata
	Root *myst.Node
}

// Loads content partially into memory, reading only the YAML frontmatter.
func LoadMetadata(contentDir string, path string) (Metadata, error) {
	slog.Debug("loading content from disk (metadata only)", "path", path)

	var (
		metadata Metadata
		err      error
	)

	metadata.key = util.KeyFromPath(contentDir, path)
	metadata.Filepath = path

	result, err := load.ReadFile[frontmatter](
		path,
		load.Opts{FrontmatterOnly: true},
	)
	if err != nil {
		return metadata, err
	}

	// Load frontmatter fields
	metadata.Title = result.Frontmatter.Title
	metadata.Description = result.Frontmatter.Description

	metadata.Date, err = result.Frontmatter.ParsedDate()
	if err != nil {
		return metadata, fmt.Errorf(
			"failed to parse frontmatter date in content file \"%s\": %w",
			path,
			err,
		)
	}

	return metadata, nil
}

// Loads and parses content.
func LoadContent(m Metadata) (Content, error) {
	slog.Debug("loading content from disk", "path", m.Filepath)

	content := Content{Metadata: m}

	result, err := load.ReadFile[frontmatter](m.Filepath, load.Opts{})
	if err != nil {
		return content, err
	}

	// Parse MyST
	content.Root, err = myst.Parse(result.Text)
	if err != nil {
		return content, fmt.Errorf(
			"failed to parse content file \"%s\": %w",
			m.Filepath,
			err,
		)
	}

	return content, nil
}
