package glob

type matcher struct {
	ast *node
}

func Matcher(pattern string) *matcher {
	ast, err := Compile(pattern)
	if err != nil {
		panic(err)
	}
	matcher := &matcher{ast: ast}
	return matcher
}

func (m *matcher) AST() *node {
	return m.ast
}

func (m *matcher) Matches(path string) (bool, error) {
	tokens, err := lex(path)
	if err != nil {
		return false, err
	}
	log("Tokens: %+v\n", tokens)
	return m.ast.Match(tokens)
}
