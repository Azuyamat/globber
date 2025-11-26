package glob

func parse(tokens []token) (*node, error) {
	var firstNode *node
	var lastNode *node

	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		negate := false
		if token.Type == tokenNegate {
			negate = true
			i++
			if i >= len(tokens) {
				break
			}
			token = tokens[i]
		}

		node := &node{
			Type:   token.Type,
			Value:  token.Literal,
			Negate: negate,
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
