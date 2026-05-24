package myst

import (
	"bytes"
	"fmt"
	"html/template"

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

const codeBlockTmpl = `
<div class="code-block-scroll-container">
  <div class="code-block-container">
    <div class="code-block">
    {{ if .Filename }}
    <div class="code-block-filename">{{ .Filename }}</div>
    {{ end }}
    {{ .Content }}
    </div>
  </div>
</div>
`

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

		// syntax highlight
		highlightedCode, err := highlight.Highlight(
			code.Value,
			code.Lang,
			code.ShowLineNumbers,
			code.EmphasizeLines,
		)
		if err != nil {
			return node, fmt.Errorf("failed to syntax highlight: %w", err)
		}

		// generate final html
		t := template.Must(template.New("code").Parse(codeBlockTmpl))
		data := struct {
			Filename string
			Content  template.HTML
		}{
			Filename: code.Filename,
			Content:  template.HTML(highlightedCode),
		}

		var buf bytes.Buffer
		t.Execute(&buf, data)
		html := buf.String()

		htmlNode, err := atrus.CreateHTMLNode(html)
		if err != nil {
			return node, err
		}

		return &Node{htmlNode}, nil
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

// Given a title, adds the title as a new heading node with depth 1 at the very
// beginning of the document.
//
// Is a no-op if the given title is an empty string.
func TransformAddTitle(node *Node, title string) (*Node, error) {
	if title == "" {
		return node, nil
	}

	// If the AST has blocks, we should add it under the first block.
	// Otherwise we can add it to the root.
	parent := node.First("block")
	if parent == nil {
		parent = node
	}

	headingNode, err := atrus.CreateHeadingNode(1)
	if err != nil {
		return node, err
	}

	textNode, err := atrus.CreateTextNode(title)
	if err != nil {
		return node, err
	}

	headingNode.AppendChild(textNode)
	parent.PrependChild(headingNode)
	return node, nil
}
