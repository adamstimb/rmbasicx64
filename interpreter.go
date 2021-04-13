package main

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"
)

// Interpreter is the BASIC interpreter itself and behaves as a state machine that
// can receive, store and interpret BASIC code and execute the code to update its
// own state.
type Interpreter struct {
	store          map[string]interface{} // A map for storing variables and array (the key is the variable name)
	program        map[int]string         // A map for storing a program (the key is the line number)
	currentTokens  []Token                // A line of tokens for immediate execution
	errorCode      int                    // The current errorCode
	lineNumber     int                    // The current line number being executed (-1 indicates immediate-mode, therefore no line number)
	badTokenIndex  int                    // If there was an error, the index of the token that raised the error is stored here
	message        string                 // The current error message, if any
	programPointer int
	tokenStack     []Token
	tokenPointer   int
	g              *Game
}

// Init initializes the Interpreter.
func (i *Interpreter) Init(g *Game) {
	i.store = make(map[string]interface{})
	i.program = make(map[int]string)
	i.currentTokens = []Token{}
	i.errorCode = Success
	i.lineNumber = -1
	i.badTokenIndex = -1
	i.message = ""
	i.programPointer = 0
	i.tokenStack = []Token{}
	i.g = g
}

// Tokenize receives a line of code, generates tokens and stores them in currentTokens.
func (i *Interpreter) Tokenize(code string) {
	s := &Scanner{}
	i.currentTokens = s.Scan(code)
}

// IsOperator receives a token and returns true if the token represents an operator
// otherwise false
func IsOperator(t Token) bool {
	operators := []int{Minus, Plus, ForwardSlash, Star, Exponential, BackSlash, Equal, InterestinglyEqual, LessThan,
		GreaterThan, LessThanEqualTo1, LessThanEqualTo2, GreaterThanEqualTo1, GreaterThanEqualTo2, Inequality1, Inequality2,
		AND, OR, XOR, NOT}
	for _, op := range operators {
		if op == t.TokenType {
			return true
		}
	}
	return false
}

// IsOperand receives a token and returns true if the token represents an operand
// otherwise false
func IsOperand(t Token) bool {
	operands := []int{NumericalLiteral, IdentifierLiteral, StringLiteral}
	for _, op := range operands {
		if op == t.TokenType {
			return true
		}
	}
	return false
}

// IsKeyword receives a token and returns true if the token's literal is a keyword
// otherwise false
func IsKeyword(t Token) bool {
	km := keywordMap()
	if _, ok := km[t.Literal]; ok {
		// is a keyword
		return true
	}
	return false
}

// Precedence returns the precedence of a token representing an operator
func Precedence(t Token) int {
	// As defined in the BASIC book p57
	precedences := map[int]int{}
	precedences[XOR] = 0
	precedences[OR] = 1
	precedences[AND] = 2
	precedences[NOT] = 3
	precedences[LessThan] = 4
	precedences[GreaterThan] = 4
	precedences[Inequality1] = 4
	precedences[Inequality2] = 4
	precedences[LessThanEqualTo1] = 4
	precedences[LessThanEqualTo2] = 4
	precedences[GreaterThanEqualTo1] = 4
	precedences[GreaterThanEqualTo2] = 4
	precedences[InterestinglyEqual] = 4
	precedences[Equal] = 4
	precedences[Plus] = 5
	precedences[Minus] = 5
	precedences[Star] = 6
	precedences[ForwardSlash] = 6
	precedences[BackSlash] = 6
	precedences[MOD] = 6
	precedences[Exponential] = 7
	return precedences[t.TokenType]
}

// GetType receives an arbitrary interface and returns the data type as a string if it is one of float64, int64 or string.
// If the type is not recognized it returns ""
func GetType(interfaceToTest interface{}) (dataType string) {
	ok := false
	_, ok = interfaceToTest.(float64)
	if ok {
		return "float64"
	}
	_, ok = interfaceToTest.(int64)
	if ok {
		return "int64"
	}
	_, ok = interfaceToTest.(string)
	if ok {
		return "string"
	}
	return ""
}

// IsTrue receives a float value and returns true if that value can represent the boolean true
// otherwise it returns false
func IsTrue(val float64) (result bool) {
	if val == math.Round(val) && math.Round(val) == -1 {
		return true
	} else {
		return false
	}
}

// WeighString receives a string and returns the sum of the ascii codes for each char
func WeighString(s string) (weight int) {
	// convert string to runs
	r := []rune(s)
	// add the codes and return total
	for _, val := range r {
		weight += int(val)
	}
	return weight
}

// RunSegment attempts to execute a segment of tokens and replies with an error code, the index
// of the token where parsing failed, and a message, or something.
func (i *Interpreter) RunSegment(tokens []Token) (ok bool) {
	// Load tokens onto the stack
	i.tokenStack = tokens
	i.tokenPointer = 0
	// 1. Pass if empty line
	if len(tokens) == 0 {
		return true
	}
	if tokens[0].TokenType == EndOfLine {
		return true
	}
	// 2. Try numeric variable assignment.  Must be at least 3 tokens.
	if len(tokens) >= 3 {
		// First 2 tokens must be identifier literal followed by = (equal) or := (assign)
		if tokens[0].TokenType == IdentifierLiteral &&
			(tokens[1].TokenType == Equal || tokens[1].TokenType == Assign) {
			// Hand over to rmAssign
			return i.rmAssign()
		}
	}
	// 3. Try built-in / keywords functions.
	if IsKeyword(tokens[0]) {
		switch tokens[0].TokenType {
		case REM:
			return true
		case PRINT:
			return i.rmPrint()
		case GOTO:
			return i.rmGoto()
		case RUN:
			return i.rmRun()
		case BYE:
			return i.rmBye()
		case LIST:
			return i.rmList()
		case SAVE:
			return i.rmSave()
		}
	}
	i.errorCode = UnknownCommandProcedure
	i.badTokenIndex = 0
	i.message = errorMessage(UnknownCommandProcedure)
	return false
}

// RunLine attempts to run a line of BASIC code and replies with an error code, the index
// of the token where parsing failed, and a message, or something.
func (i *Interpreter) RunLine(code string) (ok bool) {
	// tokenize the code
	i.Tokenize(code)
	// ensure no illegal chars
	for index, t := range i.currentTokens {
		if t.TokenType == Illegal {
			i.errorCode = EndOfInstructionExpected
			i.message = errorMessage(EndOfInstructionExpected)
			i.badTokenIndex = index
			return false
		}
	}
	// split the tokens into executable segments for each : token found
	segments := make([][]Token, 0)
	this_segment := make([]Token, 0)
	for _, token := range i.currentTokens {
		if token.TokenType != Colon {
			this_segment = append(this_segment, token)
		} else {
			// add EndOfLine to segment before adding to segments slice
			this_segment = append(this_segment, Token{EndOfLine, ""})
			segments = append(segments, this_segment)
			this_segment = make([]Token, 0)
		}
	}
	if len(this_segment) > 0 {
		segments = append(segments, this_segment)
	}
	// run each segment
	badTokenOffset := 0
	for index, segment := range segments {
		if index > 0 {
			badTokenOffset += len(segments[index-1])
		}
		ok := i.RunSegment(segment)
		if !ok {
			return false
		}
	}
	i.programPointer += 1
	return true
}

// ImmediateInput receives a string inputted by the REPL user, processes it and responds
// with a message, if any
func (i *Interpreter) ImmediateInput(code string) (response string) {
	// reset error status and tokenize code
	i.errorCode = Success
	i.message = ""
	i.badTokenIndex = 0
	i.lineNumber = -1
	i.Tokenize(code)
	// If the code begins with a line number then add it to the program otherwise try to execute it.
	if i.currentTokens[0].TokenType == NumericalLiteral {
		// TODO: This needs its own command-----
		// It starts with some kind of number so check if it's an integer, i.e. line number
		lineNumber, err := strconv.ParseFloat(i.currentTokens[0].Literal, 64)
		if err == nil {
			if lineNumber == math.Round(lineNumber) {
				// is a line number so format line and add to program
				i.program[int(lineNumber)] = i.FormatCode(code, -1, true)
				return response
			}
		}
	}
	// Does not begin with line number so try to execute
	ok := i.RunLine(code)
	if !ok {
		// There was an error so the response should include the error message
		if i.lineNumber == -1 {
			// immediate-mode syntax error without line number
			i.g.Print(fmt.Sprintf("Syntax error: %s", i.message))
			i.g.Print(fmt.Sprintf("  %s", i.FormatCode(code, i.badTokenIndex, false)))
			response = fmt.Sprintf("Syntax error: %s\n  %s", i.message, i.FormatCode(code, i.badTokenIndex, false))
		} else {
			// syntax error with line number
			i.g.Print(fmt.Sprintf("Syntax error in line %d: %s", i.lineNumber, i.message))
			i.g.Print(fmt.Sprintf("  %d %s", i.lineNumber, i.FormatCode(code, i.badTokenIndex, false)))
			response = fmt.Sprintf("Syntax error in line %d: %s\n  %d %s", i.lineNumber, i.message, 10, i.FormatCode(i.program[i.lineNumber], i.badTokenIndex, false))
		}
	}
	return response
}

// FormatCode receives a line of BASIC code and returns it formatted.  If a number
// > 0 is passed for highlightTokenIndex, the corresponding token is highlighted
// with arrows; this is used for printing error messages.
func (i *Interpreter) FormatCode(code string, highlightTokenIndex int, skipFirstToken bool) string {
	i.Tokenize(code)
	formattedCode := ""
	// handle skipFirstToken
	if skipFirstToken {
		i.currentTokens = i.currentTokens[1:]
	}
	// bump highlighter if it's pointing at :
	if highlightTokenIndex >= 0 {
		if i.currentTokens[highlightTokenIndex].TokenType == Colon && len(i.currentTokens) > highlightTokenIndex+1 {
			highlightTokenIndex += 1
		}
	}
	// format code and insert highlighter
	for index, t := range i.currentTokens {
		if index == highlightTokenIndex {
			formattedCode += "--> "
		}
		if t.TokenType == StringLiteral {
			formattedCode += "\""
			formattedCode += t.Literal
			formattedCode += "\""
		} else {
			formattedCode += t.Literal
		}
		formattedCode += " "
	}
	formattedCode = strings.TrimSpace(formattedCode)
	return formattedCode
}

// GetLineOrder returns a list of program line numbers ordered from smallest to greatest
func (i *Interpreter) GetLineOrder() (ordered []int) {
	for lineNumber, _ := range i.program {
		ordered = append(ordered, lineNumber)
	}
	sort.Ints(ordered)
	return ordered
}

// IsStringVar returns true if a token represents a string variable
func IsStringVar(t Token) bool {
	if t.TokenType == IdentifierLiteral && t.Literal[len(t.Literal)-1:] == "$" {
		return true
	} else {
		return false
	}
}

// IsIntVar returns true if a token represents a integer variable
func IsIntVar(t Token) bool {
	if t.TokenType == IdentifierLiteral && t.Literal[len(t.Literal)-1:] == "%" {
		return true
	} else {
		return false
	}
}

// IsFloatVar returns true if a token represents a float variable
func IsFloatVar(t Token) bool {
	if t.TokenType == IdentifierLiteral && (t.Literal[len(t.Literal)-1:] != "%" && t.Literal[len(t.Literal)-1:] != "$") {
		return true
	} else {
		return false
	}
}

// SetVar stores a variable
func (i *Interpreter) SetVar(variableName string, value interface{}) bool {
	switch variableName[len(variableName)-1:] {
	case "$":
		// set string variable
		if GetType(value) == "float64" {
			// cast float value to string and store
			i.store[variableName] = fmt.Sprintf("%e", value.(float64))
		} else {
			// store float value directly
			i.store[variableName] = value.(string)
			return true
		}
	case "%":
		// set integer variable
		if GetType(value) == "float64" {
			// round float and store
			i.store[variableName] = math.Round(value.(float64))
			return true
		} else {
			// try to parse float from string then round to int
			if valfloat64, err := strconv.ParseFloat(value.(string), 64); err == nil {
				i.store[variableName] = math.Round(valfloat64)
				return true
			} else {
				i.errorCode = CouldNotInterpretAsANumber
				i.message = fmt.Sprintf("%s%s", value.(string), errorMessage(CouldNotInterpretAsANumber))
				return false
			}
		}
	default:
		// set float variable
		if GetType(value) == "float64" {
			// store
			i.store[variableName] = value.(float64)
			return true
		} else {
			// try to parse float from string
			if valfloat64, err := strconv.ParseFloat(value.(string), 64); err == nil {
				i.store[variableName] = valfloat64
				return true
			} else {
				i.errorCode = CouldNotInterpretAsANumber
				i.message = fmt.Sprintf("%s%s", value.(string), errorMessage(CouldNotInterpretAsANumber))
				return false
			}
		}
	}
	// This should not happen therefore fatal
	log.Fatalf("Fatal error!")
	return false
}

// GetVar retrieves the value of a variable from the store
func (i *Interpreter) GetVar(variableName string) (value interface{}, ok bool) {
	val, ok := i.store[variableName]
	if !ok {
		i.errorCode = HasNotBeenDefined
		i.message = fmt.Sprintf("%s%s", variableName, errorMessage(HasNotBeenDefined))
		return 0, false
	} else {
		return val, true
	}
}

// GetValueFromToken receives a token representing either a variable or a literal, gets
// the value, and ~casts~ it to the required type
func (i *Interpreter) GetValueFromToken(t Token, castTo string) (value interface{}, ok bool) {
	if t.TokenType == IdentifierLiteral {
		value, ok = i.GetVar(t.Literal)
		if !ok {
			return 0, false
		}
	} else {
		value = t.Literal
	}
	switch GetType(value) {
	case "string":
		switch castTo {
		case "string":
			return value, true
		case "float64":
			// try to parse float from string
			if valfloat64, err := strconv.ParseFloat(value.(string), 64); err == nil {
				return valfloat64, true
			} else {
				i.errorCode = CouldNotInterpretAsANumber
				i.message = fmt.Sprintf("%s%s", value.(string), errorMessage(CouldNotInterpretAsANumber))
				return 0, false
			}
		case "int64":
			// try to parse float from string
			if valfloat64, err := strconv.ParseFloat(value.(string), 64); err == nil {
				return math.Round(valfloat64), true
			} else {
				i.errorCode = CouldNotInterpretAsANumber
				i.message = fmt.Sprintf("%s%s", value.(string), errorMessage(CouldNotInterpretAsANumber))
				return 0, false
			}
		case "":
			return value, true
		}
	case "float64":
		switch castTo {
		case "string":
			return fmt.Sprintf("%e", value.(float64)), true
		case "float64":
			return value, true
		case "int64":
			return math.Round(value.(float64)), true
		case "":
			return value, true
		}
	}
	log.Fatalf("Fatal error!")
	return 0, false
}

// AcceptAnyFloat checks if the current token represents a number and returns the value.  If
// it does not represent a number then the tokenPointer is reset to its original position.
func (i *Interpreter) AcceptAnyNumber() (acceptedValue float64, acceptOk bool) {
	originalPosition := i.tokenPointer
	val, ok := i.EvaluateExpression()
	if !ok {
		// broken expression so reset pointer and return false
		i.tokenPointer = originalPosition
		return 0, false
	} else {
		if GetType(val) != "float64" {
			// is not float
			// error?
			i.tokenPointer = originalPosition
			return 0, false
		} else {
			return val.(float64), true
		}
	}
}

// AcceptAnyString checks if the current token represents a string, returns the value
// and advances the pointer if so.
func (i *Interpreter) AcceptAnyString() (acceptedValue string, acceptOk bool) {
	originalPosition := i.tokenPointer
	val, ok := i.EvaluateExpression()
	if !ok {
		// broken expression so reset pointer and return false
		i.tokenPointer = originalPosition
		return "", false
	} else {
		if GetType(val) != "string" {
			// is not string
			// error?
			i.tokenPointer = originalPosition
			return "", false
		} else {
			return val.(string), true
		}
	}
}

// AcceptAnyOfTheseTokens checks if the current token matches any that are passed in a slice and, if so,
// returns the token and advances the pointer.
func (i *Interpreter) AcceptAnyOfTheseTokens(acceptableTokens []int) (acceptedToken Token, acceptOk bool) {
	for _, tokenType := range acceptableTokens {
		if tokenType == i.tokenStack[i.tokenPointer].TokenType {
			// found a match, advanced pointer and return the token type
			i.tokenPointer++
			return i.tokenStack[i.tokenPointer], true
		}
	}
	// no matches
	return Token{}, false
}

// IsAnyOfTheseTokens checks if the current token matches any that are passed in a slice and, if so,
// returns true but *does not advance the pointer*.
func (i *Interpreter) IsAnyOfTheseTokens(acceptableTokens []int) bool {
	for _, tokenType := range acceptableTokens {
		if tokenType == i.tokenStack[i.tokenPointer].TokenType {
			// found a match
			return true
		}
	}
	// no matches
	return false
}

// ExtractExpression receives a slice of tokens that represent an expression and returns
// all those tokens up to where the expression ends.
func (i *Interpreter) ExtractExpression() (expressionTokens []Token) {
	tokenStackSlice := i.tokenStack[i.tokenPointer:]
	for index, t := range tokenStackSlice {
		// The following tokens can delimit an expression
		if t.TokenType == Comma || t.TokenType == Semicolon || t.TokenType == EndOfLine || t.TokenType == Exclamation || t.TokenType == TO {
			break
		}
		// An expression is also delimited if one operand follows another directly
		if index < (len(tokenStackSlice) - 1) {
			if IsOperand(t) && IsOperand(tokenStackSlice[index+1]) {
				expressionTokens = append(expressionTokens, t)
				i.tokenPointer++
				break
			}
		}
		expressionTokens = append(expressionTokens, t)
		i.tokenPointer++
	}
	return expressionTokens
}

// EndOfTokens returns true if no more tokens are to be evaluated in the token stack
func (i *Interpreter) EndOfTokens() bool {
	if i.tokenPointer > len(i.tokenStack) {
		return true
	}
	if i.tokenStack[i.tokenPointer].TokenType == EndOfLine {
		return true
	}
	return false
}
