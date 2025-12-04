package content

import (
	"strings"
)

// Turns plain text into a sequence of HTML paragraph elements.
//
// We close a paragraph whenever we see an empty line in the source
// text.
func parsePlainText(text string) (html string, err error) {
	var (
		builder         strings.Builder
		isParagraphOpen bool
	)

	seq := strings.Lines(text)
	for line := range seq {
		if len(strings.TrimSpace(line)) > 0 {
			if !isParagraphOpen {
				builder.WriteString("<p>")
				isParagraphOpen = true
			}

			builder.WriteString(line)
		} else if isParagraphOpen {
			builder.WriteString("</p>\n")
			isParagraphOpen = false
		}
	}

	if isParagraphOpen {
		builder.WriteString("</p>\n")
		isParagraphOpen = false
	}

	return builder.String(), nil
}
