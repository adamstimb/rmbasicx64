package main

// Scanner describes a type with methods for scanning a line of source
// code and generating tokens
type Scanner struct {
	Source          string
	Tokens          []Token
	currentPosition int
}

// New initializes a new scanner for a line of source code
func (s *Scanner) New(source string) {
	s.Source = source
	s.Tokens = []Token{}
	s.currentPosition = 0
}

// isAtEnd returns true if the offset is at the end of the source code
func (s *Scanner) isAtEnd() bool {
	return s.currentPosition >= len(s.Source)
}

// advanced consumes the byte at the current position and moves the current position forward
func (s *Scanner) advance() byte {
	s.currentPosition++
	return s.Source[s.currentPosition-1]
}

// addToken creates a new token and adds it the slice of tokens
func (s *Scanner) addToken(tokenType int, lexeme string, literal string, position int) {
	s.Tokens = append(s.Tokens, Token{tokenType, lexeme, literal, position})
}
