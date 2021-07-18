package lexer

import (
	"regexp"
	"strings"
	"unicode"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
)

// Lexer describes a type with methods for lexing a line of source code and generating
// tokens.  To tokenize a line of source code all we do is this:
// 		s := &Scanner{}
//		tokens := s.ScanTokens(source)
type Lexer struct {
	Source               string        // source code string
	Tokens               []token.Token // buffer of tokens created by Scan()
	CurrentPosition      int           // position in the string
	currentTokenPosition int           // position of the buffer
}

// isAtEnd returns true if the offset is at the end of the source code
func (s *Lexer) isAtEnd() bool {
	return s.CurrentPosition >= len(s.Source)
}

// advanced consumes the rune at the current position and moves the current position forward
func (s *Lexer) advance() rune {
	s.CurrentPosition++
	return rune(s.Source[s.CurrentPosition-1])
}

// peek returns the current rune but does not consume it.  If we're at the end of the source
// then newline \n is returned instead.
func (s *Lexer) peek() rune {
	if s.isAtEnd() {
		return rune('\n')
	} else {
		return rune(s.Source[s.CurrentPosition])
	}
}

// peekNext returns the next rune but does not consume it. If we're at the end of the source
// then newline \n is returned instead.
func (s *Lexer) peekNext() rune {
	if s.CurrentPosition+1 >= len(s.Source) {
		return rune('\n')
	} else {
		return rune(s.Source[s.CurrentPosition+1])
	}
}

// match checks if a particular rune is next in the source code, returning true and consuming
// the rune if so
func (s *Lexer) match(r rune) bool {
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
func (s *Lexer) addToken(TokenType string, literal string) {
	s.Tokens = append(s.Tokens, token.Token{TokenType: TokenType, Literal: literal})
}

// getString extracts a string literal from the source code
func (s *Lexer) getString() {
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
	s.addToken(token.StringLiteral, string(stringVal))
}

// getHexLiteral extracts a hex literal from the source code
func (s *Lexer) getHexLiteral(firstRune rune) {
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
	s.addToken(token.HexLiteral, strings.ToUpper(string(hexVal)))
}

// getNumber extracts a numeric literal from the source code.  It also
// recognises scientific notation, e.g. 4e+7
func (s *Lexer) getNumber(firstRune rune) {
	// collect first rune
	stringVal := []rune{}
	stringVal = append(stringVal, firstRune)
	// then collect the rest of the literal
	gotExp := false     // ... if we've collected an(exponent symbol (e/E)
	gotExpSign := false // ... if we've collected the sign of the exponent (+/-)
	for {
		// accept only ONE e/E followed by ONE optional -/+
		if !gotExp && (s.peek() == 'e' || s.peek() == 'E') {
			// consume e/E
			gotExp = true
			stringVal = append(stringVal, s.advance())
			continue
		}
		// accept only ONE +/- immediately after e/E
		lastChar := stringVal[len(stringVal)-1]
		if !gotExpSign && gotExp && (lastChar == 'e' || lastChar == 'E') && (s.peek() == '+' || s.peek() == '-') {
			// cosume +/-
			gotExpSign = true
			stringVal = append(stringVal, s.advance())
			continue
		}
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
		// not a valid rune for numeric literal so stop collecting
		break
	}
	s.addToken(token.NumericLiteral, string(stringVal))
}

// This is a bit hacky:
var builtins = map[string]string{
	"LEN": "LEN",
	"ABS": "ABS",
	"ATN": "ATN",
	"COS": "COS",
	"EXP": "EXP",
	"INT": "INT",
	"LN":  "LN",
}

// getIdentifier extracts an identifier (keyword, variable, etc) from the source code
func (s *Lexer) getIdentifier(firstRune rune) {
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
	if token.IsKeyword(strings.ToUpper(string(stringVal))) {
		// is a keyword but if it corresponds to a built-in function we have to
		// bump it to identifier literal
		_, ok := builtins[strings.ToUpper(string(stringVal))]
		if ok {
			// is built-in
			s.addToken(token.IdentifierLiteral, strings.ToUpper(string(stringVal)))
		} else {
			// is keyword --- maybe we'll have to accept them all as identifiers...?
			s.addToken(strings.ToUpper(string(stringVal)), strings.ToUpper(string(stringVal)))
		}
		// Handle special case of REM (comment)
		if strings.ToUpper(string(stringVal)) == token.REM {
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
			s.addToken(token.IdentifierLiteral, strings.Title(strings.ToLower(string(stringVal))))
		} else {
			for _, subword := range subwords {
				if newStringVal == "" {
					newStringVal = newStringVal + strings.Title(strings.ToLower(subword))
				} else {
					newStringVal = newStringVal + "_" + strings.Title(strings.ToLower(subword))
				}
			}
			s.addToken(token.IdentifierLiteral, newStringVal)
		}
	}
}

// getComment assumes all remaining code is a comment, and puts it into a final token
func (s *Lexer) getComment() {
	stringVal := s.Source[s.CurrentPosition+1:]
	s.advance()
	s.addToken(token.Comment, stringVal)
	s.CurrentPosition = len(s.Source)
}

// scanToken generates a token for the current rune
func (s *Lexer) scanToken() {
	switch r := s.advance(); r {
	case '\n':
		s.addToken(token.NewLine, "\n")
		// ignore whitespace
	case ' ':
		return
	// one- and two-character tokens
	case ':':
		if s.match('=') {
			s.addToken(token.Assign, ":=")
			return
		}
		s.addToken(token.Colon, ":")
		return
	case '/':
		s.addToken(token.ForwardSlash, "/")
		return
	case '<':
		if s.match('>') {
			s.addToken(token.Inequality1, "<>")
			return
		}
		if s.match('=') {
			s.addToken(token.LessThanEqualTo1, "<=")
			return
		}
		s.addToken(token.LessThan, "<")
		return
	case '>':
		if s.match('<') {
			s.addToken(token.Inequality2, "><")
			return
		}
		if s.match('=') {
			s.addToken(token.GreaterThanEqualTo1, ">=")
			return
		}
		s.addToken(token.GreaterThan, ">")
		return
	case '=':
		if s.match('<') {
			s.addToken(token.LessThanEqualTo2, "=<")
			return
		}
		if s.match('>') {
			s.addToken(token.GreaterThanEqualTo2, "=>")
			return
		}
		if s.match('=') {
			s.addToken(token.InterestinglyEqual, "==")
			return
		}
		s.addToken(token.Equal, "=")
		return
	// then single-character
	case '(':
		s.addToken(token.LeftParen, "(")
		return
	case ')':
		s.addToken(token.RightParen, ")")
		return
	case ',':
		s.addToken(token.Comma, ",")
		return
	case '.':
		s.addToken(token.Dot, ".")
		return
	case '-':
		// Check if it represents an operator or negative number
		// If there was no previous token, or the previous token was an operator,
		// and the next char is a number then it's a negative number that needs
		// to be collected, otherwise it's an operator to collect.
		if len(s.Tokens) == 0 && unicode.IsDigit(s.peek()) {
			s.getNumber(r)
			return
		}
		if len(s.Tokens) >= 1 && (token.IsOperator(s.Tokens[len(s.Tokens)-1]) || s.Tokens[len(s.Tokens)-1].TokenType == token.LeftParen) {
			s.getNumber(r)
			return
		} else {
			// Is operator so just collect it
			s.addToken(token.Minus, "-")
			return
		}
	case '+':
		s.addToken(token.Plus, "+")
		return
	case ';':
		s.addToken(token.Semicolon, ";")
		return
	case '\\':
		s.addToken(token.BackSlash, "\\")
		return
	case '*':
		s.addToken(token.Star, "*")
		return
	case '^':
		s.addToken(token.Exponential, "^")
		return
	case '!':
		s.addToken(token.Exclamation, "!")
	case '#':
		s.addToken(token.Hash, "#")
	case '~':
		s.addToken(token.Tilde, "~")
	case '[':
		s.addToken(token.LeftSquareBrace, "[")
	case ']':
		s.addToken(token.RightSquareBrace, "]")
	case '{':
		s.addToken(token.LeftCurlyBrace, "{")
	case '}':
		s.addToken(token.RightCurlyBrace, "}")
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
		// numeric literal
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
		s.addToken(token.Illegal, string(r))
	}
}

// Scan scans the source code and returns a slice of tokens
func (s *Lexer) Scan(source string) []token.Token {
	s.Source = source
	s.Tokens = []token.Token{}
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
	//s.addToken(token.EndOfInstruction, token.EndOfInstruction)
	s.addToken(token.EOF, token.EOF)
	s.currentTokenPosition = 0
	return s.Tokens
}

// PeekToken returns the token at the current tokenPosition in the buffer
// but does not consume it.
func (s *Lexer) PeekToken() token.Token {
	return s.Tokens[s.currentTokenPosition]
}

// NextToken returns the token at the current tokenPosition in the buffer
// and consumes it, moving to the next token
func (s *Lexer) NextToken() token.Token {
	if s.currentTokenPosition >= len(s.Tokens) {
		s.currentTokenPosition = len(s.Tokens) - 1
		return s.Tokens[s.currentTokenPosition]
	}
	s.currentTokenPosition++
	return s.Tokens[s.currentTokenPosition-1]
}
