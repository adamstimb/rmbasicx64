package main

import "unicode"

// Scanner describes a type with methods for scanning a line of source
// code and generating tokens.  To tokenize a line of source code all
// we do is this:
// 		s := &Scanner{}
//		tokens, err := s.ScanTokens(source)
// To handle errors err will be zero is successful otherwise it will correspond
// to one of the error constants.  To get the exact position where the error was detected
// use s.CurrentPosition
type Scanner struct {
	Source          string
	Tokens          []Token
	CurrentPosition int
}

// isAtEnd returns true if the offset is at the end of the source code
func (s *Scanner) isAtEnd() bool {
	return s.CurrentPosition >= len(s.Source)
}

// advanced consumes the rune at the current position and moves the current position forward
func (s *Scanner) advance() rune {
	s.CurrentPosition++
	return rune(s.Source[s.CurrentPosition-1])
}

// peek returns the current rune but does not consume it.  If we're at the end of the source
// then newline \n is returned instead.
func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return rune('\n')
	} else {
		return rune(s.Source[s.CurrentPosition])
	}
}

// peekNext returns the next rune but does not consume it. If we're at the end of the source
// then newline \n is returned instead.
func (s *Scanner) peekNext() rune {
	if s.CurrentPosition + 1 >= len(s.Source) {
		return rune('\n')
	} else {
		return rune(s.Source[s.CurrentPosition + 1])
	}
}

// match checks if a particular rune is next in the source code, returning true and consuming
// the rune if so
func (s *Scanner) match(r rune) bool {
	if s.isAtEnd() {
		return false
	}
	if rune(s.Source[s.CurrentPosition]) == r {
		s.CurrentPosition++
		return true
	}
	return false
}

// addToken creates a new token and adds it the slice of tokens
func (s *Scanner) addToken(tokenType int, lexeme string, literal string) {
	s.Tokens = append(s.Tokens, Token{tokenType, lexeme, literal, s.CurrentPosition - 1})
}

// getString extracts a string literal from the source code
func (s *Scanner) getString() {
	stringVal := ""
	for s.peek() != '"' && !s.isAtEnd() {
		stringVal = append(stringVal, s.advance())
	}
	// handle unterminated string
	if s.isAtEnd() {
		logMsg("Unterminated string") // error handling tbc
	}
	// otherwise add the token
	s.addToken(StringLiteral, stringVal, "")
}

// getNumber extracts a numerical literal from the source code
func (s *Scanner) getNumber(firstRune rune) {
	// collect first rune
	stringVal := ""
	stringVal = append(stringVal, firstRune)
	// then collect the rest of the literal
	for {
		if unicode.IsDigit(s.peek()) {
			// consume this digit
			stringVal = append(stringVal, s.advance())	
			continue
		}
		// look for fractional part
		if s.peek() == '.' && unicode.IsDigit(s.peekNext()) {
			// consume the .
			stringVal = append(stringVal, s.advance())
			continue
		}
		// not a valid rune for numerical literal so stop collecting
		break
	}
	// test if parsable number then add token
	s.addToken(NumericalLiteral, stringVal, "")
}

// scanToken generates a token for the current rune
func (s *Scanner) scanToken() {
	switch r := s.advance(); r {
	// ignore whitespace
	case ' ':
		return
	// evaluate two-character tokens first
	case ':' && s.match('='):
		s.addToken(Assign, "", "")
		return
	case '/' && s.match('/'):
		s.addToken(IntegerDivision, "", "")
		return
	case '<' && s.match('>'):
		s.addToken(Inequality1)
		return
	case '>' && s.match('<'):
		s.addToken(Inequality2)
		return
	case '<' && s.match('='):
		s.addToken(LessThanEqualTo1)
		return
	case '=' && s.match('<'):
		s.addToken(LessThanEqualTo2)
		return
	case '>' && s.match('='):
		s.addToken(GreaterThanEqualTo1)
		return
	case '=' && s.match('>'):
		s.addToken(GreaterThanEqualTo2)
		return
	case '=' && s.match('='):
		s.addToken(InterestinglyEqual)
		return
	// then single-character
	case '(':
		s.addToken(LeftParen, "", "")
		return
	case ')':
		s.addToken(RightParen, "", "")
		return
	case ',':
		s.addToken(Comma, "", "")
		return
	case '.':
		s.addToken(Dot, "", "")
		return
	case '-':
		s.addToken(Minus, "", "")
		return
	case '+':
		s.addToken(Plus, "", "")
		return
	case ';':
		s.addToken(Semicolon, "", "")
		return
	case '/':
		s.addToken(ForwardSlash, "", "")
		return
	case '\':
		s.addToken(BackSlash, "", "")
		return
	case '*':
		s.addToken(Star, "", "")
		return
	case '^':
		s.addToken(Exponential, "", "")
		return
	case '<':
		s.addToken(LessThan, "", "")
		return
	case '>':
		s.addToken(GreaterThan, "", "")
		return
	// string literal
	case '"':
		s.getString()
		return
	default:
		// numerical literal
		if unicode.IsDigit(r) {
			// number(r)
			return
		}
		if unicode.IsLetter(r) {
			// identifier()
			return
		}
	}
}

// ScanTokens scans the source code and returns a slice of tokens
func (s *Scanner) ScanTokens(source string) []Token {
	s.Source = source
	s.Tokens = []Token{}
	s.CurrentPosition = 0
	for !s.isAtEnd() {
		s.scanToken()
	}
	// All done
	return s.Tokens
}
