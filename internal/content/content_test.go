package content_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/sinclairtarget/michel/internal/content"
	"github.com/sinclairtarget/michel/internal/content/myst"
)

func TestLoadFromMarkdown(t *testing.T) {
	const fileContents = `---
title: My Blog Post
description: Foobar
date: 2025-12-04
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
	if c.Title != expectedTitle {
		t.Errorf(
			"title incorrect; wanted %s, got %s",
			expectedTitle,
			c.Title,
		)
	}

	expectedDescription := "Foobar"
	if c.Description != expectedDescription {
		t.Errorf(
			"description incorrect; wanted %s, got %s",
			expectedDescription,
			c.Description,
		)
	}

	expectedDate := time.Date(2025, 12, 04, 0, 0, 0, 0, time.Local)
	if c.Date != expectedDate {
		t.Errorf(
			"date incorrect; wanted %s, got %s",
			expectedDate,
			c.Date,
		)
	}

	expectedHtml := `<p>This is a blog post. Here is the first paragraph.</p>
<h2>Subheading</h2>
<p>Here is the second paragraph.</p>
`
	output, err := myst.RenderHTML(c.Root)
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

func TestLoadFromMarkdownBadDate(t *testing.T) {
	const fileContents = `---
title: My Blog Post
date: foo
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

	_, err = content.LoadFromMarkdown(tmpdir, filename)
	if err == nil {
		t.Error("expected error but got nil")
	}
}

func TestLoadFromMarkdownNoDate(t *testing.T) {
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

	if c.Date.Year() != 1 {
		t.Errorf("default year should have been 1; got %d", c.Date.Year())
	}
}
