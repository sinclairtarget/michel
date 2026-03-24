// MyST Markdown handling.
package myst

import (
	"fmt"
	"html/template"

	atrus "github.com/sinclairtarget/libatrus-go"
)

// Parse MyST markdown into a MyST AST.
func Parse(text string) (*Node, error) {
	opts := atrus.ParseOpts{
		ParseLevel: atrus.ParseLevelPost,
	}
	root, err := atrus.Parse(text, opts)
	if err != nil {
		return nil, fmt.Errorf("libatrus parse error: %w", err)
	}

	return &Node{*root}, nil
}

// Render MyST AST to HTML.
func RenderHTML(node *Node) (template.HTML, error) {
	html, err := atrus.RenderHTML(&node.ASTNode)
	if err != nil {
		return "", fmt.Errorf("libatrus render error: %w", err)
	}

	return template.HTML(html), nil
}

// Render MyST AST to JSON.
func RenderJSON(node *Node) (string, error) {
	opts := atrus.JSONOpts{
		Whitespace: atrus.JSONIndent2,
	}
	json, err := atrus.RenderJSON(&node.ASTNode, opts)
	if err != nil {
		return "", fmt.Errorf("libatrus render error: %w", err)
	}

	return json, nil
}
