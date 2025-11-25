package glob

type tokenType int

const (
	tokenEOF tokenType = iota
	tokenLiteral
	tokenStar
	tokenDoubleStar
	tokenSlash
	tokenLBracket
	tokenRBracket
	tokenQuestion
	tokenDot
	tokenNegate
	tokenPipe
	tokenComma
	tokenLParen
	tokenRParen
)

type token struct {
	Type    tokenType
	Literal string
}

func isSpecialChar(ch byte) bool {
	switch ch {
	case '*', '?', '[', ']', '/', '!', '|', ',', '(', ')', '.':
		return true
	default:
		return false
	}
}
