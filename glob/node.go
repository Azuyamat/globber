package glob

import "fmt"

type node struct {
	Type  tokenType
	Value string
	Next  *node

	Children []*node
}

func (n *node) String() string {
	return fmt.Sprintf("Node{Type: %d, Value: %s}\n", n.Type, n.Value)
}

func (n *node) Tree() string {
	sb := ""
	sb += fmt.Sprintf("- %s\n", n.String())
	for _, child := range n.Children {
		sb += "  " + child.Tree()
	}
	if n.Next != nil {
		sb += n.Next.Tree()
	}
	return sb
}

func (n *node) Match(tokens []token) (bool, error) {
	if len(tokens) == 0 {
		log("No more tokens to match.")
		return false, nil
	}
	currentToken := tokens[0]
	log("Matching Node: %s with Token: %+v\n", n.String(), currentToken)
	switch n.Type {
	case tokenLiteral:
		if n.Value == string(currentToken.Literal) {
			if n.Next == nil {
				return true, nil
			}
			return n.Next.Match(tokens[1:])
		}
	case tokenDoubleStar:
		if n.Next == nil {
			return true, nil
		}
		for i := 0; i <= len(tokens); i++ {
			match, err := n.Next.Match(tokens[i:])
			if err != nil {
				return false, err
			}
			if match {
				return true, nil
			}
		}
		return n.Next.Match(tokens[1:])
	case tokenStar:
		if n.Next == nil {
			return true, nil
		}
		return n.Next.Match(tokens[1:])
	case tokenDot:
		if n.Value == string(currentToken.Literal) {
			if n.Next == nil {
				return true, nil
			}
			return n.Next.Match(tokens[1:])
		}
	case tokenSlash:
		if n.Value == string(currentToken.Literal) {
			if n.Next == nil {
				return true, nil
			}
			return n.Next.Match(tokens[1:])
		}
	default:
		log("Unhandled node type:", n.Type)
		return false, nil
	}
	log("No match for Node:", n.String())
	return false, nil
}
