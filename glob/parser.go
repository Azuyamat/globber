package glob

func Parse(tokens []token) (*node, error) {
	var firstNode *node
	var lastNode *node

	for _, token := range tokens {
		node := &node{
			Type:  token.Type,
			Value: token.Literal,
		}
		if firstNode == nil {
			firstNode = node
		}

		if lastNode != nil {
			lastNode.Next = node
		}

		lastNode = node
	}

	return firstNode, nil
}
