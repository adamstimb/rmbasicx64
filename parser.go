package main

// findToken finds matching token types in a list of tokens and returns their indexes
func findToken(tokens []Token, tokenType int) []int {
	indexes := []int{}
	for i, thisToken := range tokens {
		if thisToken.Type == tokenType {
			indexes = append(indexes, i)
		}
	}
	return indexes
}

// getStatements splits a list of tokens into statements
func getStatements(tokens []Token) [][]Token {
	logMsg("GetStatements")
	statements := [][]Token{}
	separators := findToken(tokens, PnStatementSeparator)
	// handle no separators, therefore token list is just one statement
	if len(separators) == 0 {
		statements = append(statements, tokens)
		return statements
	}
	// others split the token list into statements
	lastSepIndex := 0
	for i, sepIndex := range separators {
		statements = append(statements, tokens[lastSepIndex:sepIndex])
		lastSepIndex = sepIndex
		if i < len(separators) {
			statements = append(statements, tokens[lastSepIndex:])
		}
	}
	return statements
}

// parseTokens receives a list of tokens and produces one or more syntax trees
func parseTokens(tokens []Token) {
	// Get statements from code
	logMsg("ParseTokens")
	statements := getStatements(tokens)
	for _, statement := range statements {
		// Create syntax tree for each statement
		logMsg("New statement:")
		if DEBUG {
			for _, thisToken := range statement {
				PrintToken(thisToken)
			}
		}

	}
}
