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

	c, err := content.LoadFromMarkdown(tmpdir, filename)
	if err != nil {
		t.Fatalf("failed to load content: %v", err)
	}

	if c.Path != filename {
		t.Errorf(
			"content path incorrect; wanted %s, got %s",
			filename,
			c.Path,
		)
	}

	expectedKey := "test-content"
	if c.Key != expectedKey {
		t.Errorf(
			"content name incorrect; wanted %s, got %s",
			expectedKey,
			c.Key,
		)
	}

	expectedTitle := "My Blog Post"
	if c.Frontmatter.Title != expectedTitle {
		t.Errorf(
			"title incorrect; wanted %s, got %s",
			expectedTitle,
			c.Frontmatter.Title,
		)
	}

	expectedHtml := `<p>This is a blog post. Here is the first paragraph.</p>
<h2>Subheading</h2>
<p>Here is the second paragraph.</p>
`
	output, err := content.RenderMyST(c.Root)
	if err != nil {
		t.Errorf("failed to render to HTML: %v", err)
	}
	if output != expectedHtml {
		t.Errorf(
			"html incorrect; wanted:\n%s\ngot:\n%s\n",
			expectedHtml,
			output,
		)
	}
}
