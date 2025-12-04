package site_test

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/sinclairtarget/michel/internal/site"
)

// We should be able to load a templated site page from disk.
//
// Site pages support YAML frontmatter, separated from the following content by
// `---`.
func TestLoadPage(t *testing.T) {
	const templateText = `{{ define "title" }}
<title>{{ .Content.Title }}</title>
{{ end }}
{{ define "main" }}
    {{ template "partials/article" . }}
{{ end }}
`
	const fileContents = `---
layouts:
  - base
  - article
---
` + templateText

	tmpdir := t.TempDir()
	filename := filepath.Join(tmpdir, "page.html.tmpl")
	err := os.WriteFile(filename, []byte(fileContents), 0o644)
	if err != nil {
		t.Fatalf("failed to write template to tmp dir: %v", err)
	}

	page, err := site.LoadPage(filename)
	if err != nil {
		t.Fatalf("failed to load template: %v", err)
	}

	if page.Path != filename {
		t.Errorf(
			"page path incorrect; wanted \"%s\" but got \"%s\"",
			filename,
			page.Path,
		)
	}

	expected := []string{"base", "article"}
	if !slices.Equal(page.Frontmatter.Layouts, expected) {
		t.Errorf(
			"frontmatter layouts incorrect; wanted %v but got %v",
			expected,
			page.Frontmatter.Layouts,
		)
	}

	if page.TemplateText != templateText {
		t.Errorf(
			"page template text incorrect; wanted:\n%s\ngot:\n%s\n",
			templateText,
			page.TemplateText,
		)
	}
}

// Frontmatter is optional.
func TestLoadPageNoFrontmatter(t *testing.T) {
	const fileContents = `{{ define "title" }}
<title>{{ .Content.Title }}</title>
{{ end }}
{{ define "main" }}
    {{ template "partials/article" . }}
{{ end }}
`

	tmpdir := t.TempDir()
	filename := filepath.Join(tmpdir, "page.html.tmpl")
	err := os.WriteFile(filename, []byte(fileContents), 0o644)
	if err != nil {
		t.Fatalf("failed to write template to tmp dir: %v", err)
	}

	page, err := site.LoadPage(filename)
	if err != nil {
		t.Fatalf("failed to load template: %v", err)
	}

	if page.Path != filename {
		t.Errorf(
			"page path incorrect; wanted \"%s\" but got \"%s\"",
			filename,
			page.Path,
		)
	}

	if len(page.Frontmatter.Layouts) > 0 {
		t.Error("page frontmatter layouts non-empty; wanted empty slice")
	}

	if page.TemplateText != fileContents {
		t.Errorf(
			"page template text incorrect; wanted:\n%s\ngot:\n%s\n",
			fileContents,
			page.TemplateText,
		)
	}
}
