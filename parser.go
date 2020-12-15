package main

// Statement types defined here
const (
	StaVariableAssignment       = 3000
	StaInternalProcedureCall    = 3001
	StaUserDefinedProcedureCall = 3002
	StaInvalid                  = 3003
)

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
		logMsg("No separators found")
		return statements
	}
	// others split the token list into statements
	lastSepIndex := 0
	for i, sepIndex := range separators {
		statements = append(statements, tokens[lastSepIndex:sepIndex])
		lastSepIndex = sepIndex
		if i < len(separators) {
			statements = append(statements, tokens[lastSepIndex+1:])
		}
	}
	return statements
}

// getStatementType receives a statement and decides based on the first one or two tokens
// which type of statement is to be parsed
func getStatementType(tokens []Token) int {
	logMsg("GetStatementType:")
	// RM Basic has three types of statements:
	// Variable assignment			Foo% := Bar% + 1 or a
	// Internal procedure call		PLOT "Foobar", 0, 0 SIZE 2 BRUSH 4
	// User-defined procedure call  Draw_Shapes
	// So the token type of the first two token determines the statement type
	firstToken := tokens[0]
	secondTokenExists := false
	if len(tokens) > 1 {
		secondTokenExists = true
	}
	// Variable assignment or internal procedure call?
	if firstToken.Type == MaVariableFloat ||
		firstToken.Type == MaVariableInteger ||
		firstToken.Type == MaVariableString {
		// First token is a variable.  Since procedure or function names are tokenized
		// as variable names we have to check the second token type to figure out which
		// it is.  But we can only try to evaluate the second token type if it exists
		// so we have to break that down a little.
		if secondTokenExists {
			secondToken := tokens[1]
			if secondToken.Type == MaAssign {
				// it's variable assignment
				logMsg("StaVariableAssignment")
				return StaVariableAssignment
			} else {
				// it's user-defined procedure call
				logMsg("StaUserDefinedProcedureCall")
				return StaUserDefinedProcedureCall
			}
		} else {
			// also user-define procedure call, e.g. with no arguments
			logMsg("StaUserDefinedProcedureCall")
			return StaUserDefinedProcedureCall
		}
	}
	if firstToken.Type >= KwABS && firstToken.Type <= KwSET {
		// it's internal procedure call
		logMsg("StaInternalProcedureCall")
		return StaInternalProcedureCall
	}
	// Otherwise invalid
	logMsg("StaInvalid")
	return StaInvalid
}

func parseBye(g *Game, tokens []Token) int {
	logMsg("ParseBye")
	if len(tokens) > 1 {
		return ErEndOfInstructionExpected
	}
	return internalBye(g)
}

func parsePrint(g *Game, tokens []Token) int {
	logMsg("ParsePrint")
	if len(tokens) == 1 {
		return ErNotEnoughParameters
	}
	nextToken := tokens[1]
	// validate next token (must be "number of string")
	if nextToken.Type == MaVariableFloat ||
		nextToken.Type == MaVariableInteger ||
		nextToken.Type == MaVariableString ||
		nextToken.Type == LiFloat ||
		nextToken.Type == LiInteger ||
		nextToken.Type == LiString {
	} else {
		return ErNumberOrStringNeeded
	}
	// now convert to string if required then call print
	var text string
	if nextToken.Type == LiString {
		text = nextToken.ValueString
		return internalPrint(g, text)
	}
	return 0
}

func parseInternalProcedureCall(g *Game, tokens []Token) int {
	logMsg("ParseInternalProcedureCall")
	firstToken := tokens[0]
	if firstToken.Type == KwBYE {
		return parseBye(g, tokens)
	}
	if firstToken.Type == KwPRINT {
		return parsePrint(g, tokens)
	}
	return 0
}

// parseTokens receives a list of tokens and produces one or more syntax trees
func parseTokens(g *Game, tokens []Token) int {
	// Get statements from code.  RM BASIC supports multiple statements per line
	// seperated by the : symbol.  These will be parsed as if they were seperate lines
	logMsg("ParseTokens")
	statements := getStatements(tokens)
	for _, statement := range statements {
		// Parse this statement
		logMsg("New statement:")
		var err int
		if DEBUG {
			for _, thisToken := range statement {
				PrintToken(thisToken)
			}
		}
		// Get statement type - raise error and return if invalid
		statementType := getStatementType(statement)
		if statementType == 0 {
			err = ErUnknownCommandProcedure
		}
		// Continue parsing the statement
		if statementType == StaInternalProcedureCall {
			err = parseInternalProcedureCall(g, statement)
		}
		if err != 0 {
			return err
		}
	}
	return 0
}
