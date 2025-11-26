package glob

func lex(pattern string) ([]token, error) {
	scanner := newScanner(pattern)
	tokens := make([]token, 0, len(pattern)/2+1)

	for !scanner.eof() {
		switch scanner.char {
		case '*':
			if scanner.peekChar() == '*' {
				scanner.readChar()
				tokens = append(tokens, token{Type: tokenDoubleStar, Literal: "**"})
			} else {
				tokens = append(tokens, token{Type: tokenStar, Literal: "*"})
			}
		case '/':
			tokens = append(tokens, token{Type: tokenSlash, Literal: "/"})
		case '.':
			tokens = append(tokens, token{Type: tokenDot, Literal: "."})
		case '!':
			tokens = append(tokens, token{Type: tokenNegate, Literal: "!"})
		default:
			if !isSpecialChar(scanner.char) {
				startPos := scanner.position
				for !scanner.eof() && !isSpecialChar(scanner.char) {
					scanner.readChar()
				}
				literal := pattern[startPos:scanner.position]
				tokens = append(tokens, token{Type: tokenLiteral, Literal: literal})
				continue
			}
		}
		scanner.readChar()
	}

	return tokens, nil
}
