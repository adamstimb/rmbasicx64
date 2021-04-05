package main

import (
	"regexp"
	"strings"
	"unicode"
)

// Scanner describes a type with methods for scanning a line of source code and generating
// tokens.  To tokenize a line of source code all we do is this:
// 		s := &Scanner{}
//		tokens := s.ScanTokens(source)
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
func (s *Scanner) addToken(TokenType int, literal string) {
	s.Tokens = append(s.Tokens, Token{TokenType, literal})
}

// getString extracts a string literal from the source code
func (s *Scanner) getString() {
	stringVal := []rune{}
	for !s.isAtEnd() {
		// handle literal double-double quote "" or terminating double quote "
		if s.peek() == '"' && s.peekNext() == '"' {
			// is double-double quote so consume both and continue collecting
			stringVal = append(stringVal, s.advance())
			stringVal = append(stringVal, s.advance())
			continue
		}
		if s.peek() == '"' && s.peekNext() != '"' {
			// is double quote so stop collecting
			break
		} else {
			// is part of the string literal so collect it
			stringVal = append(stringVal, s.advance())
		}
	}
	// handle string termination then add the token
	if s.peek() == '"' {
		s.advance()
	}
	s.addToken(StringLiteral, string(stringVal))
}

// getHexLiteral extracts a hex literal from the source code
func (s *Scanner) getHexLiteral(firstRune rune) {
	hexVal := []rune{}
	hexVal = append(hexVal, firstRune)
	for {
		r := s.peek()
		isHex, _ := regexp.Match("[0-9a-fA-F]", []byte(string(r)))
		if isHex {
			// is hex so consume and add to value
			hexVal = append(hexVal, s.advance())
		} else {
			// not hex so stop collection
			break
		}
	}
	s.addToken(HexLiteral, strings.ToUpper(string(hexVal)))
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
	s.addToken(NumericalLiteral, string(stringVal))
}

// getIdentifier extracts an identifier (keyword, variable, etc) from the source code
func (s *Scanner) getIdentifier(firstRune rune) {
	// collect first rune
	stringVal := []rune{}
	stringVal = append(stringVal, firstRune)
	// then collect rest of the identifier
	for unicode.IsDigit(s.peek()) || unicode.IsLetter(s.peek()) || s.peek() == '_' {
		stringVal = append(stringVal, s.advance())
	}
	// Got the full identifier so now see if it matches a keyword.  If it matches a
	// keyword add token with the token type identifying the keyword, otherwise add
	// a variable (which in RM Basic then requires checking for a trailing $ or % to
	// get the type, if any)
	keywords := keywordMap()
	if t, found := keywords[strings.ToUpper(string(stringVal))]; found {
		// is a keyword
		s.addToken(t, strings.ToUpper(string(stringVal)))
		// Handle special case of REM (comment)
		if t == REM {
			s.getComment()
		}
	} else {
		// Is another kind of identifier.  Check for trailing $ and % before adding token:
		if s.peek() == '$' || s.peek() == '%' {
			// consume this char and add token
			stringVal = append(stringVal, s.advance())
		}
		// Enforce the Rm_Basic_Camel_Case_Thing by splitting around _, titling the words
		// and recombining
		newStringVal := ""
		subwords := strings.Split(string(stringVal), "_")
		if len(subwords) == 0 {
			s.addToken(IdentifierLiteral, strings.Title(strings.ToLower(string(stringVal))))
		} else {
			for _, subword := range subwords {
				if newStringVal == "" {
					newStringVal = newStringVal + strings.Title(strings.ToLower(subword))
				} else {
					newStringVal = newStringVal + "_" + strings.Title(strings.ToLower(subword))
				}
			}
			s.addToken(IdentifierLiteral, newStringVal)
		}
	}
}

// getComment assumes all remaining code is a comment, and puts it into a final token
func (s *Scanner) getComment() {
	stringVal := s.Source[s.CurrentPosition+1:]
	s.advance()
	s.addToken(Comment, stringVal)
	s.CurrentPosition = len(s.Source)
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
			s.addToken(Assign, ":=")
			return
		}
		s.addToken(Colon, ":")
		return
	case '/':
		if s.match('/') {
			s.addToken(IntegerDivision, "//")
			return
		}
		s.addToken(ForwardSlash, "/")
		return
	case '<':
		if s.match('>') {
			s.addToken(Inequality1, "<>")
			return
		}
		if s.match('=') {
			s.addToken(LessThanEqualTo1, "<=")
			return
		}
		s.addToken(LessThan, "<")
		return
	case '>':
		if s.match('<') {
			s.addToken(Inequality2, "><")
			return
		}
		if s.match('=') {
			s.addToken(GreaterThanEqualTo1, ">=")
			return
		}
		s.addToken(GreaterThan, ">")
		return
	case '=':
		if s.match('<') {
			s.addToken(LessThanEqualTo2, "=<")
			return
		}
		if s.match('>') {
			s.addToken(GreaterThanEqualTo2, "=>")
			return
		}
		if s.match('=') {
			s.addToken(InterestinglyEqual, "==")
			return
		}
		s.addToken(Equal, "=")
		return
	// then single-character
	case '(':
		s.addToken(LeftParen, "(")
		return
	case ')':
		s.addToken(RightParen, ")")
		return
	case ',':
		s.addToken(Comma, ",")
		return
	case '.':
		s.addToken(Dot, ".")
		return
	case '-':
		s.addToken(Minus, "-")
		return
	case '+':
		s.addToken(Plus, "+")
		return
	case ';':
		s.addToken(Semicolon, ";")
		return
	case '\\':
		s.addToken(BackSlash, "\\")
		return
	case '*':
		s.addToken(Star, "*")
		return
	case '^':
		s.addToken(Exponential, "^")
		return
	// string literal
	case '"':
		s.getString()
		return
	// hex literal
	case '&':
		// if next char is hex-ish then assume is hex literal and get it
		nextChar := s.peekNext()
		isHex, _ := regexp.Match("[0-9a-fA-F]", []byte(string(nextChar)))
		if isHex {
			s.getHexLiteral(r)
		}
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
		// unexpected chars are tokenized as Illegal
		s.addToken(Illegal, string(r))
	}
}

// Scan scans the source code and returns a slice of tokens
func (s *Scanner) Scan(source string) []Token {
	s.Source = source
	s.Tokens = []Token{}
	s.CurrentPosition = 0
	// Handle special case of only whitespace as input
	if strings.TrimSpace(s.Source) == "" {
		// is just whitespace so don't scan
	} else {
		for !s.isAtEnd() {
			s.scanToken()
		}
	}
	// All done - add end of line token and return
	s.addToken(EndOfLine, "")
	return s.Tokens
}
