package content

import (
	"fmt"

	atrus "github.com/sinclairtarget/libatrus-go"
)

func parseMystMarkdown(text string) (string, error) {
	ast, err := atrus.ParseAST(text)
	if err != nil {
		return "", fmt.Errorf("libatrus parse error: %w", err)
	}

	html, err := atrus.RenderHTML(ast)
	if err != nil {
		return "", fmt.Errorf("libatrus render error: %w", err)
	}

	return html, nil
}
