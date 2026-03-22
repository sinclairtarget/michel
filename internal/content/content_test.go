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

	m, err := content.LoadMetadata(tmpdir, filename)
	if err != nil {
		t.Fatalf("failed to load content: %v", err)
	}

	if m.Path != filename {
		t.Errorf(
			"content path incorrect; wanted %s, got %s",
			filename,
			m.Path,
		)
	}

	expectedKey := "test-content"
	if m.Key() != expectedKey {
		t.Errorf(
			"content name incorrect; wanted %s, got %s",
			expectedKey,
			m.Key(),
		)
	}

	expectedTitle := "My Blog Post"
	if m.Title != expectedTitle {
		t.Errorf(
			"title incorrect; wanted %s, got %s",
			expectedTitle,
			m.Title,
		)
	}

	expectedDescription := "Foobar"
	if m.Description != expectedDescription {
		t.Errorf(
			"description incorrect; wanted %s, got %s",
			expectedDescription,
			m.Description,
		)
	}

	expectedDate := time.Date(2025, 12, 04, 0, 0, 0, 0, time.Local)
	if m.Date != expectedDate {
		t.Errorf(
			"date incorrect; wanted %s, got %s",
			expectedDate,
			m.Date,
		)
	}

	c, err := content.LoadContent(m)
	if err != nil {
		t.Fatal(err)
	}

	expectedHtml := `<p>This is a blog post. Here is the first paragraph.</p>
<h2>Subheading</h2>
<p>Here is the second paragraph.</p>
`
	output, err := myst.RenderHTML(c.Root)
	if err != nil {
		t.Errorf("failed to render to HTML: %v", err)
	}
	if string(output) != expectedHtml {
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

	_, err = content.LoadMetadata(tmpdir, filename)
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

	m, err := content.LoadMetadata(tmpdir, filename)
	if err != nil {
		t.Fatalf("failed to load content: %v", err)
	}

	if m.Date.Year() != 1 {
		t.Errorf("default year should have been 1; got %d", m.Date.Year())
	}
}
