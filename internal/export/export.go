// Package export handles exporting content files as other formats.
package export

import (
	"fmt"
	"io"

	"github.com/sinclairtarget/michel/internal/build"
	"github.com/sinclairtarget/michel/internal/content"
	"github.com/sinclairtarget/michel/internal/content/myst"
)

type Options struct {
	Format   string // Output format
	AddTitle bool   // Whether to add title metadata to content body as heading
}

// Loads the given content file and exports it in the given format, writing to
// the given writer.
//
// To export MyST files, we parse them into the AST then pass them to a
// renderer. We should not be applying any Michel-specific transform to the
// AST.
func Export(
	filepath string,
	options Options,
	w io.Writer,
) error {
	metadata, err := content.LoadMetadata(build.ContentDir, filepath)
	if err != nil {
		return err
	}

	c, err := content.LoadContent(metadata)
	if err != nil {
		return err
	}

	node := c.Root
	if options.AddTitle {
		node, err = myst.TransformAddTitle(node, c.Title)
		if err != nil {
			return fmt.Errorf("add title transform failed: %w", err)
		}
	}

	switch options.Format {
	case "json":
		err = exportJSON(node, w)
	case "typst":
		err = exportTypst(node, w)
	default:
		return fmt.Errorf("unsupported format \"%s\"", options.Format)
	}
	if err != nil {
		return err
	}

	return nil
}

func exportJSON(node *myst.Node, w io.Writer) error {
	s, err := myst.RenderJSON(node)
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, s)
	if err != nil {
		return err
	}

	// Make sure we follow the JSON string with a newline
	fmt.Fprintln(w)

	return nil
}

func exportTypst(node *myst.Node, w io.Writer) error {
	s, err := myst.RenderTypst(node)
	if err != nil {
		return err
	}

	_, err = io.WriteString(w, s)
	if err != nil {
		return err
	}

	return nil
}
