package myst

import (
	atrus "github.com/sinclairtarget/libatrus-go"
)

func Transform(node *Node) (*Node, error) {
	transformed, err := transform(&node.ASTNode)
	if err != nil {
		return node, err
	}

	return &Node{*transformed}, nil
}

func transform(node *atrus.ASTNode) (*atrus.ASTNode, error) {
	if node.Type() == "mystDirective" {
		value := "<div><p>Hi! It's-a me-ah!</p></div>"
		html, err := atrus.CreateHTMLNode(value)
		if err != nil {
			return node, err
		}

		return html, nil
	}

	for i, child := range node.Children() {
		transformed, err := transform(child)
		if err != nil {
			return node, err
		}

		if transformed != child {
			// We have to replace a node
			node.ReplaceChild(uint32(i), transformed)
		}
	}

	return node, nil
}
