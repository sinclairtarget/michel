package content

import (
	"fmt"
	"iter"

	atrus "github.com/sinclairtarget/libatrus-go"
)

type MySTNode struct {
	*atrus.ASTNode
}

// Returns an iterator over all nodes in the AST rooted at the given node that
// match the given type.
//
// A pre-order traversal of the AST is performed and matching nodes are
// returned in that order.
func (n *MySTNode) All(nodeType string) iter.Seq[*MySTNode] {
	seq := func(yield func(*MySTNode) bool) {
		if n.Type() == nodeType {
			if !yield(n) {
				return
			}
		}

		for child := range n.Children() {
			wrapped := MySTNode{child}
			for match := range wrapped.All(nodeType) {
				if !yield(match) {
					return
				}
			}
		}
	}

	return seq
}

// Returns the first node of the matching type found in AST during pre-order
// traversal, or nil if no matching node is found.
func (n *MySTNode) First(nodeType string) *MySTNode {
	for match := range n.All(nodeType) {
		return match
	}

	return nil
}

func RenderMyST(node *MySTNode) (string, error) {
	html, err := atrus.RenderHTML(node.ASTNode)
	if err != nil {
		return "", fmt.Errorf("libatrus render error: %w", err)
	}

	return html, nil
}

func parseMyST(text string) (*MySTNode, error) {
	opts := atrus.ParseOpts{
		ParseLevel: atrus.ParseLevelPost,
	}
	root, err := atrus.Parse(text, opts)
	if err != nil {
		return nil, fmt.Errorf("libatrus parse error: %w", err)
	}

	return &MySTNode{root}, nil
}
