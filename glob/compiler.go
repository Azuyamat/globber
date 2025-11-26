package glob

func compile(pattern string) (*node, error) {
	tokens, err := lex(pattern)
	if err != nil {
		return nil, err
	}
	return parse(tokens)
}
