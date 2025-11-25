package glob

type scanner struct {
	pattern      string
	position     int
	readPosition int
	char         byte
}

func newScanner(pattern string) *scanner {
	s := &scanner{pattern: pattern}
	s.readChar()
	return s
}

func (s *scanner) readChar() (byte, int, int) {
	if s.readPosition >= len(s.pattern) {
		s.char = 0
	} else {
		s.char = s.pattern[s.readPosition]
	}
	s.position = s.readPosition
	s.readPosition++
	return s.char, s.position, s.readPosition
}

func (s *scanner) peekChar() byte {
	if s.readPosition >= len(s.pattern) {
		return 0
	}
	return s.pattern[s.readPosition]
}

func (s *scanner) eof() bool {
	return s.position >= len(s.pattern)
}
