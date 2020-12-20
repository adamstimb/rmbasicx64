package main

import "github.com/adamstimb/nimgobus"

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

// numberValueToString receives a token with a numeric value and returns
// the value as a string
//func numberValueToString(token Token) string {
//	logMsg("NumberValueToString")
//	if token.Type == LiFloat || token.Type == MaVariableFloat ||
//		token.Type == LiInteger || token.Type == MaVariableInteger {
//		return strconv.FormatFloat(token.ValueFloat, 'f', -1, 64)
//	}
//	// need to handle unexpected better than this:
//	return ""
//}

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

// resolveVariable tries to get a variable's value from the store
func resolveVariable(g *Game, token Token) (string, int) {
	item, exists := g.Store[token.Symbol]
	if exists {
		// variable exists and therefore has a value
		return item[0].Value, 0
	} else {
		// variable does not exist so return error code
		return "", ErVariableWithoutAnyValue
	}
}

func parseStringExpression(g *Game, tokens []Token) (string, int) {
	logMsg("ParseStringExpression")
	// The only string expression supported by RM Basic is string concatenation so we look for this
	// pattern only: | [string variable / string literal] ... + |
	result := ""
	expectValue := true
	for _, token := range tokens {
		if expectValue {
			// If we expect a value then it can be a literal or a variable.  Anything else
			// will be invalid
			if token.Type != LiString && token.Type != LiNumber &&
				token.Type != MaVariableString && token.Type != MaVariableInteger &&
				token.Type != MaVariableFloat {
				return "", ErInvalidExpressionFound
			}
			// Otherwise handle the literal or variable:
			if token.Type == LiString || token.Type == LiNumber {
				// is a string literal so concat its symbol to result
				result = result + token.Symbol
			}
			// if it's a variable then resolve its value and concat
			if token.Type == MaVariableString || token.Type == MaVariableFloat ||
				token.Type == MaVariableInteger {
				// is a string variable so try to get value and concat it to result
				value, err := resolveVariable(g, token)
				if err != 0 {
					return "", err
				} else {
					result = result + value
				}
			}
			// set flag to expect addition symbol in next token
			expectValue = false
		} else {
			// if we don't expect a value then it must be an addition symbol
			if token.Type != MaAddition {
				return "", ErInvalidExpressionFound
			}
			// set flag to expect string valyue in next token
			expectValue = true
		}
	}
	logMsg("result=" + result)
	return result, 0
}

func parseVariableAssignment(g *Game, tokens []Token) int {
	logMsg("ParseVariableAssignment")
	// Token 0 is defines the variable that will store the result.  Token 1 is the assignment
	// symbol (this has already been checked by the time this function has been called).  Token 2
	// is either a literal or a value.  At least 3 tokens are therefore required.
	if len(tokens) < 3 {
		return ErInvalidExpressionFound
	}
	// if the assignment is to a string variable then parse string expression
	if tokens[0].Type == MaVariableString {
		value, err := parseStringExpression(g, tokens[2:])
		if err != 0 {
			return err
		} else {
			g.Store[tokens[0].Symbol] = []nimgobus.StoreItem{{0, value}}
		}
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
		if DEBUG {
			for _, thisToken := range statement {
				PrintToken(thisToken)
			}
		}
		// Get statement type - raise error and return if invalid
		statementType := getStatementType(statement)
		if statementType == 0 {
			return ErUnknownCommandProcedure
		}
		// Otherwise continue parsing the statement
		if statementType == StaInternalProcedureCall {
			return parseInternalProcedureCall(g, statement)
		}
		if statementType == StaVariableAssignment {
			return parseVariableAssignment(g, statement)
		}
	}
	return 0
}
