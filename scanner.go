package main

// Scanner describes a type with methods for scanning a line of source
// code and generating tokens.  To tokenize a line of source code all
// we do is this:
// 		s := &Scanner{}
//		s.New(source)
//		tokens, err := s.ScanTokens()
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

// advanced consumes the rune at the current position and moves the current position forward
func (s *Scanner) advance() rune {
	s.currentPosition++
	return rune(s.Source[s.currentPosition-1])
}

// peek returns the current rune but does not consume it.  If we're at the end of the source
// then newline \n is returned instead.
func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return rune('\n')
	} else {
		return rune(s.Source[s.currentPosition])
	}
}

// match checks if a particular rune is next in the source code, returning true and consuming
// the rune if so
func (s *Scanner) match(r rune) bool {
	if s.isAtEnd() {
		return false
	}
	if rune(s.Source[s.currentPosition]) == r {
		s.currentPosition++
		return true
	}
	return false
}

// addToken creates a new token and adds it the slice of tokens
func (s *Scanner) addToken(tokenType int, lexeme string, literal string) {
	s.Tokens = append(s.Tokens, Token{tokenType, lexeme, literal, s.currentPosition - 1})
}

// scanToken generates a token for the current rune
func (s *Scanner) scanToken() {
	switch r := s.advance(); r {
	// evaluate two-character tokens first
	case ':' && s.match('='):
		s.addToken(Assign, "", "")
	case '/' && s.match('/'):
		s.addToken(IntegerDivision, "", "")
	case '<' && s.match('>'):
		s.addToken(Inequality1)
	case '>' && s.match('<'):
		s.addToken(Inequality2)
	case '<' && s.match('='):
		s.addToken(LessThanEqualTo1)
	case '=' && s.match('<'):
		s.addToken(LessThanEqualTo2)
	case '>' && s.match('='):
		s.addToken(GreaterThanEqualTo1)
	case '=' && s.match('>'):
		s.addToken(GreaterThanEqualTo2)
	case '=' && s.match('='):
		s.addToken(InterestinglyEqual)
	// then single-character
	case '(':
		s.addToken(LeftParen, "", "")
	case ')':
		s.addToken(RightParen, "", "")
	case ',':
		s.addToken(Comma, "", "")
	case '.':
		s.addToken(Dot, "", "")
	case '-':
		s.addToken(Minus, "", "")
	case '+':
		s.addToken(Plus, "", "")
	case ';':
		s.addToken(Semicolon, "", "")
	case '/':
		s.addToken(ForwardSlash, "", "")
	case '\':
		s.addToken(BackSlash, "", "")
	case '*':
		s.addToken(Star, "", "")
	case '^':
		s.addToken(Exponential, "", "")
	case '<':
		s.addToken(LessThan, "", "")
	case '>':
		s.addToken(GreaterThan, "", "")
	}
}

// ScanTokens scans the source code and returns a slice of tokens
func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.scanToken()
	}
	// All done
	return s.Tokens
}
