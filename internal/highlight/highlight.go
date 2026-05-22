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

// Returns an HTML `pre` element as a string with syntax highlighting for the
// given text.
func Highlight(
	text string,
	lang string,
	showLineNumbers bool,
	emphasizeLines []uint,
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
		html.HighlightLines(toRanges(emphasizeLines)),
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

// Turns a list of integers into a list of ranges covering all the integers.
//
// The list of integers must be sorted.
func toRanges(lineNums []uint) [][2]int {
	if len(lineNums) == 0 {
		return [][2]int{}
	}

	lineRanges := [][2]int{}
	var openNum int
	for i, lineNum := range lineNums {
		if i == 0 {
			openNum = int(lineNum)
			continue
		}

		closeNum := int(lineNums[i-1])
		if int(lineNum) > closeNum+1 {
			lineRanges = append(lineRanges, [2]int{openNum, closeNum})
		}
	}

	closeNum := int(lineNums[len(lineNums)-1])
	lineRanges = append(lineRanges, [2]int{openNum, closeNum})
	return lineRanges
}
