// Package highlight implements syntax highlighting (and some other
// presentation goodies) for code blocks.
package highlight

import (
	"fmt"
	"log/slog"
	"strings"

	chroma "github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
)

func Highlight(
	text string,
	lang string,
	showLineNumbers bool,
) (string, error) {
	lexer := lexers.Get(lang)
	if lexer == nil {
		slog.Warn("no chroma lexer found", "lang", lang)
		return "", nil
	}
	lexer = chroma.Coalesce(lexer)

	style := styles.Fallback
	formatter := html.New(
		html.WithClasses(true),
		html.WithLineNumbers(showLineNumbers),
	)

	it, err := lexer.Tokenise(nil, text)
	if err != nil {
		return "", fmt.Errorf("failed to tokenize: %w", err)
	}

	var builder strings.Builder
	err = formatter.Format(&builder, style, it)
	if err != nil {
		return "", fmt.Errorf("failed to format: %w", err)
	}

	return builder.String(), nil
}
