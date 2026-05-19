package export

import (
	"fmt"
	"io"

	"github.com/sinclairtarget/michel/internal/build"
	"github.com/sinclairtarget/michel/internal/content"
	"github.com/sinclairtarget/michel/internal/content/myst"
)

// Loads the given content file and exports it in the given format, writing to
// the given writer.
//
// To export MyST files, we parse them into the AST then pass them to a
// renderer. We should not be applying any Michel-specific transform to the
// AST.
func Export(filepath string, format string, w io.Writer) error {
	metadata, err := content.LoadMetadata(build.ContentDir, filepath)
	if err != nil {
		return err
	}

	c, err := content.LoadContent(metadata)
	if err != nil {
		return err
	}

	switch format {
	case "json":
		err = exportJSON(c, w)
	default:
		return fmt.Errorf("unsupported format \"%s\"", format)
	}
	if err != nil {
		return err
	}

	return nil
}

func exportJSON(c content.Content, w io.Writer) error {
	s, err := myst.RenderJSON(c.Root)
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
