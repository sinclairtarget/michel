package content_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sinclairtarget/michel/internal/content"
)

func TestLoadFromMarkdown(t *testing.T) {
	const fileContents = `---
title: My Blog Post
---
This is a blog post. Here is the first paragraph.

## Subheading
Here is the second paragraph.
`
	tmpdir := t.TempDir()
	filename := filepath.Join(tmpdir, "test-content.md")
	err := os.WriteFile(filename, []byte(fileContents), 0o644)
	if err != nil {
		t.Fatalf("failed to write content file to tmp dir: %v", err)
	}

	content, err := content.LoadFromMarkdown(tmpdir, filename)
	if err != nil {
		t.Fatalf("failed to load content: %v", err)
	}

	if content.Path != filename {
		t.Errorf(
			"content path incorrect; wanted %s, got %s",
			filename,
			content.Path,
		)
	}

	expectedName := "test-content"
	if content.Name != expectedName {
		t.Errorf(
			"content name incorrect; wanted %s, got %s",
			expectedName,
			content.Name,
		)
	}

	expectedTitle := "My Blog Post"
	if content.Frontmatter.Title != expectedTitle {
		t.Errorf(
			"title incorrect; wanted %s, got %s",
			expectedTitle,
			content.Frontmatter.Title,
		)
	}

	expectedHtml := `<p>This is a blog post. Here is the first paragraph.</p>
<h2>Subheading</h2>
<p>Here is the second paragraph.</p>
`
	if content.Html != expectedHtml {
		t.Errorf(
			"html incorrect; wanted:\n%s\ngot:\n%s\n",
			expectedHtml,
			content.Html,
		)
	}
}
