package main

import (
	"sort"
)

// Tokenize receives a line of code and returns a list of tokens
func Tokenize(code string) []Token {
	// tokens are collected in this slice
	tokens := []Token{}
	// String literals
	print("TOKENIZE: ")
	println(code)
	// finding tokenizing is much simpler if we pad the code and then search for keywords
	// that are enclosed by white space
	code = PadString(code)
	println("string literals:")
	for _, thisToken := range TokenizeStringLiterals(code) {
		PrintToken(thisToken)
		tokens = append(tokens, thisToken)
		// Mask this string so we don't tokenize inside them
		bytes := []byte(code)
		for i := range bytes {
			if i >= thisToken.Location[0] && i <= thisToken.Location[1] {
				bytes[i] = ' '
			}
			code = string(bytes)
		}
	}
	// Punctuation
	println("punctuation:")
	for _, thisToken := range TokenizePunctuation(code) {
		PrintToken(thisToken)
		tokens = append(tokens, thisToken)
	}
	// Mathematical
	println("mathematical:")
	for _, thisToken := range TokenizeMathematical(code) {
		PrintToken(thisToken)
		tokens = append(tokens, thisToken)
	}
	println("numerical literals:")
	for _, thisToken := range TokenizeNumericalLiterals(code) {
		PrintToken(thisToken)
		tokens = append(tokens, thisToken)
	}
	// Keywords
	println("keywords:")
	for _, thisToken := range TokenizeKeywords(code) {
		PrintToken(thisToken)
		tokens = append(tokens, thisToken)
	}
	// Variables
	println("variables:")
	for _, thisToken := range TokenizeVariables(tokens, code) {
		PrintToken(thisToken)
		tokens = append(tokens, thisToken)
	}
	// done
	return tokens
}

// Format receives a line of code and returns properly formatted code
func Format(code string) string {
	// Tokenize the code then reconstruct the code from the token's symbols
	// Symbols won't be in the correct order so first load the symbols into a map
	// where key is the start position of the symbol and value is in the symbol
	symbols := make(map[int]string)
	tokens := Tokenize(code)
	keys := make([]int, len(tokens))
	commentIndex := 0
	print("FORMAT: ")
	println(code)
	for i, thisToken := range tokens {
		PrintToken(thisToken)
		symbols[thisToken.Location[0]] = thisToken.Symbol
		keys[i] = thisToken.Location[0]
		// If a REM was detected remember the start index of the subsequent comment
		if thisToken.Type == KwREM {
			commentIndex = thisToken.Location[1] + 1
		}
	}
	// sort keys
	sort.Ints(keys)
	// create the formatted string
	var formatted string
	for i, k := range keys {
		formatted = formatted + symbols[k]
		// If REM append the comment that we collected earlier and break
		if symbols[k] == "REM" {
			formatted = formatted + " " + code[commentIndex:]
			break
		}
		// add a space if not at end of string
		if i < len(keys) {
			formatted = formatted + " "
		}
	}
	return formatted
}
