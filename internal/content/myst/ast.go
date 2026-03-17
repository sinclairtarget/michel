package myst

import (
	"iter"

	atrus "github.com/sinclairtarget/libatrus-go"
)

// Represents a node in the parsed MyST AST.
// See https://mystmd.org/spec/myst-schema
//
// We wrap the basic node with helper methods for traversing the AST.
type Node struct {
	atrus.ASTNode
}

// Returns an iterator over all nodes in the AST rooted at the given node that
// match the given type.
//
// A pre-order traversal of the AST is performed and matching nodes are
// returned in that order.
func (n *Node) All(nodeType string) iter.Seq[*Node] {
	seq := func(yield func(*Node) bool) {
		if n.Type() == nodeType {
			if !yield(n) {
				return
			}
		}

		for _, child := range n.Children() {
			wrapped := Node{*child}
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
func (n *Node) First(nodeType string) *Node {
	for match := range n.All(nodeType) {
		return match
	}

	return nil
}
