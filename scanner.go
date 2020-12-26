package main

// Scanner describes a type with methods for scanning a line of source
// code and generating tokens
type Scanner struct {
	Source string
	Tokens []Token
}
