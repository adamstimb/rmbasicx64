package main

import (
	"strings"
	"unicode"
)

// Scanner describes a type with methods for scanning a line of source code and generating
// tokens.  To tokenize a line of source code all we do is this:
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
	if s.CurrentPosition+1 >= len(s.Source) {
		return rune('\n')
	} else {
		return rune(s.Source[s.CurrentPosition+1])
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
func (s *Scanner) addToken(TokenType int, lexeme string, literal string) {
	s.Tokens = append(s.Tokens, Token{TokenType, lexeme, literal, s.CurrentPosition - 1})
}

// getString extracts a string literal from the source code
func (s *Scanner) getString() {
	stringVal := []rune{}
	for s.peek() != '"' && !s.isAtEnd() {
		stringVal = append(stringVal, s.advance())
	}
	// handle string termination then add the token
	if s.peek() == '"' {
		s.advance()
	}
	s.addToken(StringLiteral, string(stringVal), "")
}

// getNumber extracts a numerical literal from the source code
func (s *Scanner) getNumber(firstRune rune) {
	// collect first rune
	stringVal := []rune{}
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
	s.addToken(NumericalLiteral, string(stringVal), "")
}

// getIdentifier extracts an identifier (keyword, variable, etc) from the source code
func (s *Scanner) getIdentifier(firstRune rune) {
	// collect first rune
	stringVal := []rune{}
	stringVal = append(stringVal, firstRune)
	// then collect rest of the identifier
	for unicode.IsDigit(s.peek()) || unicode.IsLetter(s.peek()) {
		stringVal = append(stringVal, s.advance())
	}
	// Got the full identifier so now see if it matches a keyword.  If it matches a
	// keyword add token with the token type identifying the keyword, otherwise add
	// a variable (which in RM Basic then requires checking for a trailing $ or % to
	// get the type, if any)
	keywords := keywordMap()
	if t, found := keywords[strings.ToUpper(string(stringVal))]; found {
		// is a keyword
		s.addToken(t, strings.ToUpper(string(stringVal)), "")
	} else {
		// Is another kind of identifier.  Check for trailing $ and % before adding token:
		if s.peek() == '$' || s.peek() == '%' {
			// consume this char and add token
			stringVal = append(stringVal, s.advance())
		}
		s.addToken(Identifier, strings.Title(string(stringVal)), "")
	}
}

// scanToken generates a token for the current rune
func (s *Scanner) scanToken() {
	switch r := s.advance(); r {
	// ignore whitespace
	case ' ':
		return
	// one- and two-character tokens
	case ':':
		if s.match('=') {
			s.addToken(Assign, ":=", "")
			return
		}
		s.addToken(Colon, ":", "")
		return
	case '/':
		if s.match('/') {
			s.addToken(IntegerDivision, "//", "")
			return
		}
		s.addToken(ForwardSlash, "/", "")
		return
	case '<':
		if s.match('>') {
			s.addToken(Inequality1, "<>", "")
			return
		}
		if s.match('=') {
			s.addToken(LessThanEqualTo1, "<=", "")
			return
		}
		s.addToken(LessThan, "<", "")
		return
	case '>':
		if s.match('<') {
			s.addToken(Inequality2, "><", "")
			return
		}
		if s.match('=') {
			s.addToken(GreaterThanEqualTo2, ">=", "")
		}
		s.addToken(GreaterThan, ">", "")
		return
	case '=':
		if s.match('<') {
			s.addToken(LessThanEqualTo2, "=<", "")
			return
		}
		if s.match('>') {
			s.addToken(GreaterThanEqualTo2, "=>", "")
			return
		}
		if s.match('=') {
			s.addToken(InterestinglyEqual, "==", "")
			return
		}
		s.addToken(Equal, "=", "")
		return
	// then single-character
	case '(':
		s.addToken(LeftParen, "(", "")
		return
	case ')':
		s.addToken(RightParen, ")", "")
		return
	case ',':
		s.addToken(Comma, ",", "")
		return
	case '.':
		s.addToken(Dot, ".", "")
		return
	case '-':
		s.addToken(Minus, "-", "")
		return
	case '+':
		s.addToken(Plus, "+", "")
		return
	case ';':
		s.addToken(Semicolon, ";", "")
		return
	case '\\':
		s.addToken(BackSlash, "\\", "")
		return
	case '*':
		s.addToken(Star, "*", "")
		return
	case '^':
		s.addToken(Exponential, "^", "")
		return
	// string literal
	case '"':
		s.getString()
		return
	default:
		// numerical literal
		if unicode.IsDigit(r) {
			s.getNumber(r)
			return
		}
		// identifier
		if unicode.IsLetter(r) {
			s.getIdentifier(r)
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
