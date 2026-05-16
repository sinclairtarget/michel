package myst

import (
	"fmt"

	atrus "github.com/sinclairtarget/libatrus-go"

	"github.com/sinclairtarget/michel/internal/highlight"
)

// Run all Michel custom MyST transforms that we want to run before rendering
// to HTML.
func htmlRenderTransform(node *Node) (*Node, error) {
	transformed, err := transformCode(node)
	if err != nil {
		return node, err
	}

	return transformed, nil
}

// A MyST transform to handle syntax highlighting and a few other code block
// presentation niceties.
//
// The libatrus HTML renderer does not handle syntax highlighting for code
// nodes. We implement our own HTML renderer (just for code nodes) with a MyST
// transform that turns code nodes into HTML nodes.
//
// This applies only to code blocks, not inline code.
func transformCode(node *Node) (*Node, error) {
	t := node.Type()
	if t == "code" {
		code := node.Code()

		highlightedCode, err := highlight.Highlight(
			code.Value,
			code.Lang,
			code.ShowLineNumbers,
		)
		if err != nil {
			return node, fmt.Errorf("failed to syntax highlight: %w", err)
		}

		html, err := atrus.CreateHTMLNode(highlightedCode)
		if err != nil {
			return node, err
		}

		return &Node{html}, nil
	}

	for i, child := range node.Children() {
		wrapped := &Node{child}
		transformed, err := transformCode(wrapped)
		if err != nil {
			return node, err
		}

		if transformed != wrapped {
			// We have to replace a node
			node.ReplaceChild(uint32(i), transformed.ASTNode)
		}
	}

	return node, nil
}
