package glob

func Compile(pattern string) (*node, error) {
	tokens, err := lex(pattern)
	if err != nil {
		return nil, err
	}
	return Parse(tokens)
}
