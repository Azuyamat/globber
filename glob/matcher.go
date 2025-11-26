package glob

type matcher struct {
	ast *node

	isExactMatch   bool
	isSimpleSuffix bool
}

func Matcher(pattern string) *matcher {
	ast, err := compile(pattern)
	if err != nil {
		panic(err)
	}
	matcher := &matcher{ast: ast}
	if ast.Type == tokenLiteral && ast.Next == nil {
		matcher.isExactMatch = true
	}
	if ast.Type == tokenStar && ast.Next != nil && ast.Next.Type == tokenLiteral && ast.Next.Next == nil {
		matcher.isSimpleSuffix = true
	}
	return matcher
}

func (m *matcher) AST() *node {
	return m.ast
}

func (m *matcher) Matches(path string) (bool, error) {
	if m.isExactMatch {
		return path == m.ast.Value, nil
	}
	if m.isSimpleSuffix {
		suffix := m.ast.Next.Value
		return len(path) >= len(suffix) && path[len(path)-len(suffix):] == suffix, nil
	}
	return m.ast.MatchString(path, 0)
}
