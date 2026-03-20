package site_test

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/sinclairtarget/michel/internal/site"
)

// We should be able to load a templated page from disk.
//
// Pages support YAML frontmatter, separated from the following content by
// `---`.
func TestLoadPage(t *testing.T) {
	const templateText = `{{ define "title" }}
<title>{{ .Content.Title }}</title>
{{ end }}
{{ define "main" }}
    {{ partial "article" . }}
{{ end }}
`
	const fileContents = `---
layouts:
  - base
  - blog
---
` + templateText

	tmpdir := t.TempDir()
	filename := filepath.Join(tmpdir, "page.html.tmpl")
	err := os.WriteFile(filename, []byte(fileContents), 0o644)
	if err != nil {
		t.Fatalf("failed to write template to tmp dir: %v", err)
	}

	metadata, err := site.LoadPageMetadata(tmpdir, filename)
	if err != nil {
		t.Fatalf("failed to load template: %v", err)
	}

	if metadata.Path != filename {
		t.Errorf(
			"page path incorrect; wanted \"%s\" but got \"%s\"",
			filename,
			metadata.Path,
		)
	}

	expected := []string{"base", "blog"}
	if !slices.Equal(metadata.Layouts, expected) {
		t.Errorf(
			"frontmatter layouts incorrect; wanted %v but got %v",
			expected,
			metadata.Layouts,
		)
	}

	page, err := site.LoadPage(metadata)
	if err != nil {
		t.Fatal(err)
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
    {{ partial "article" . }}
{{ end }}
`

	tmpdir := t.TempDir()
	filename := filepath.Join(tmpdir, "page.html.tmpl")
	err := os.WriteFile(filename, []byte(fileContents), 0o644)
	if err != nil {
		t.Fatalf("failed to write template to tmp dir: %v", err)
	}

	metadata, err := site.LoadPageMetadata(tmpdir, filename)
	if err != nil {
		t.Fatalf("failed to load template: %v", err)
	}

	if metadata.Path != filename {
		t.Errorf(
			"page path incorrect; wanted \"%s\" but got \"%s\"",
			filename,
			metadata.Path,
		)
	}

	if len(metadata.Layouts) > 0 {
		t.Error("page frontmatter layouts non-empty; wanted empty slice")
	}

	page, err := site.LoadPage(metadata)
	if err != nil {
		t.Fatal(err)
	}

	if page.TemplateText != fileContents {
		t.Errorf(
			"page template text incorrect; wanted:\n%s\ngot:\n%s\n",
			fileContents,
			page.TemplateText,
		)
	}
}
