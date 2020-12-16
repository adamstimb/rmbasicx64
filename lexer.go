package main

import (
	"sort"
)

// maskSymbols masks tokenized symbols in a line of code
func maskSymbols(code string, tokens []Token) string {
	bytes := []byte(code)
	for _, thisToken := range tokens {
		for i := range bytes {
			if i >= thisToken.Location[0] && i <= thisToken.Location[1] {
				bytes[i] = ' '
			}
		}
	}
	return string(bytes)
}

// sortTokens returns the token list in the order they appear in the code
func sortTokens(tokens []Token) []Token {
	// first make a map of token locations to tokens and list the locations
	locations := []int{}
	mappedTokens := make(map[int]Token)
	for _, thisToken := range tokens {
		mappedTokens[thisToken.Location[0]] = thisToken
		locations = append(locations, thisToken.Location[0])
	}
	// now make a new list of tokens in the correct order
	sort.Ints(locations)
	sortedTokens := []Token{}
	for _, k := range locations {
		sortedTokens = append(sortedTokens, mappedTokens[k])
	}
	return sortedTokens
}

// Tokenize receives a line of code and returns a list of tokens
func tokenize(code string) []Token {
	logMsg("Tokenize")
	// tokens are collected in this slice
	tokens := []Token{}
	// String literals
	// tokenizing is much simpler if we pad the code and then search for keywords
	// that are enclosed by white space
	code = PadString(code)
	logMsg("String literals:")
	tokens = TokenizeStringLiterals(code)
	// Punctuation
	logMsg("Punctuation:")
	for _, thisToken := range TokenizePunctuation(code) {
		PrintToken(thisToken)
		tokens = append(tokens, thisToken)
	}
	// Mathematical
	logMsg("Mathematical:")
	for _, thisToken := range TokenizeMathematical(code) {
		PrintToken(thisToken)
		tokens = append(tokens, thisToken)
	}
	logMsg("Numerical literals:")
	for _, thisToken := range TokenizeNumericalLiterals(code) {
		PrintToken(thisToken)
		tokens = append(tokens, thisToken)
	}
	// Keywords
	logMsg("Keywords:")
	tokens = TokenizeKeywords(tokens, code)
	// Variables
	logMsg("Variables:")
	for _, thisToken := range TokenizeVariables(tokens, code) {
		PrintToken(thisToken)
		tokens = append(tokens, thisToken)
	}
	// done - sort tokens and return
	sortedTokens := sortTokens(tokens)
	if DEBUG {
		logMsg("Final list of tokens:")
		for _, thisToken := range sortedTokens {
			PrintToken(thisToken)
		}
	}
	return sortTokens(tokens)
}

// Format receives a line of code and tokens returns properly formatted code
func format(code string, tokens []Token) string {
	logMsg("Format")
	// Tokenize the code then reconstruct the code from the token's symbols
	// Symbols won't be in the correct order so first load the symbols into a map
	// where key is the start position of the symbol and value is in the symbol
	symbols := make(map[int]string)
	keys := make([]int, len(tokens))
	commentIndex := 0
	for i, thisToken := range tokens {
		symbols[thisToken.Location[0]] = thisToken.Symbol
		keys[i] = thisToken.Location[0]
		// If a REM was detected remember the start index of the subsequent comment
		if thisToken.Type == KwREM {
			commentIndex = thisToken.Location[1] + 1
		}
	}
	//tokens = sortTokens(tokens)
	var formatted string
	for i, thisToken := range tokens {
		formatted = formatted + thisToken.Symbol
		// If REM append the comment that we collected earlier and break
		if thisToken.Symbol == "REM" {
			formatted = formatted + " " + code[commentIndex:]
			break
		}
		// if this isn't the last token then decide if we need to add a space
		// before the next symbol
		spaceRequired := true
		if i < len(tokens)-1 {
			// parentheses
			if thisToken.Type == PnLeftParenthesis {
				spaceRequired = false
			}
			if (thisToken.Type == MaVariableFloat ||
				thisToken.Type == MaVariableInteger ||
				thisToken.Type == MaVariableString) &&
				tokens[i+1].Type == PnRightParenthesis {
				spaceRequired = false
			}
			// value lists
			if (thisToken.Type == MaVariableFloat ||
				thisToken.Type == MaVariableInteger ||
				thisToken.Type == MaVariableString) &&
				(tokens[i+1].Type == PnValueSeparator || tokens[i+1].Type == PnCoordinateSeperator) {
				spaceRequired = false
			}
		}
		if spaceRequired {
			formatted = formatted + " "
		}
	}
	logMsg("Formatted=" + formatted)
	return formatted
}
