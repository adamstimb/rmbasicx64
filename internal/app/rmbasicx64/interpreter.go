package rmbasicx64

import (
	"fmt"
	"log"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
)

// Interpreter is the BASIC interpreter itself and behaves as a state machine that
// can receive, store and interpret BASIC code and execute the code to update its
// own state.
type Interpreter struct {
	Store          map[string]interface{} // A map for storing variables and array (the key is the variable name)
	Program        map[int]string         // A map for storing a program (the key is the line number)
	CurrentTokens  []token.Token          // A line of tokens for immediate execution
	ErrorCode      int                    // The current errorCode
	LineNumber     int                    // The current line number being executed (-1 indicates immediate-mode, therefore no line number)
	BadTokenIndex  int                    // If there was an error, the index of the token that raised the error is stored here
	Message        string                 // The current error message, if any
	ProgramPointer int
	TokenStack     []token.Token
	TokenPointer   int
	g              *Game
}

// Init initializes the Interpreter.
func (i *Interpreter) Init(g *Game) {
	i.Store = make(map[string]interface{})
	i.Program = make(map[int]string)
	i.CurrentTokens = []token.Token{}
	i.ErrorCode = syntaxerror.Success
	i.LineNumber = -1
	i.BadTokenIndex = -1
	i.Message = ""
	i.ProgramPointer = 0
	i.TokenStack = []token.Token{}
	i.g = g
}

// Tokenize receives a line of code, generates tokens and stores them in currentTokens.
func (i *Interpreter) Tokenize(code string) {
	s := &Scanner{}
	i.CurrentTokens = s.Scan(code)
}

// IsOperator receives a token and returns true if the token represents an operator
// otherwise false
func IsOperator(t token.Token) bool {
	operators := []int{token.Minus, token.Plus, token.ForwardSlash, token.Star, token.Exponential,
		token.BackSlash, token.Equal, token.InterestinglyEqual, token.LessThan,
		token.GreaterThan, token.LessThanEqualTo1, token.LessThanEqualTo2, token.GreaterThanEqualTo1,
		token.GreaterThanEqualTo2, token.Inequality1, token.Inequality2,
		token.AND, token.OR, token.XOR, token.NOT}
	for _, op := range operators {
		if op == t.TokenType {
			return true
		}
	}
	return false
}

// IsOperand receives a token and returns true if the token represents an operand
// otherwise false
func IsOperand(t token.Token) bool {
	operands := []int{token.NumericalLiteral, token.IdentifierLiteral, token.StringLiteral}
	for _, op := range operands {
		if op == t.TokenType {
			return true
		}
	}
	return false
}

// IsKeyword receives a token and returns true if the token's literal is a keyword
// otherwise false
func IsKeyword(t token.Token) bool {
	km := token.KeywordMap()
	if _, ok := km[t.Literal]; ok {
		// is a keyword
		return true
	}
	return false
}

// Precedence returns the precedence of a token representing an operator
func Precedence(t token.Token) int {
	// As defined in the BASIC book p57
	precedences := map[int]int{}
	precedences[token.XOR] = 0
	precedences[token.OR] = 1
	precedences[token.AND] = 2
	precedences[token.NOT] = 3
	precedences[token.LessThan] = 4
	precedences[token.GreaterThan] = 4
	precedences[token.Inequality1] = 4
	precedences[token.Inequality2] = 4
	precedences[token.LessThanEqualTo1] = 4
	precedences[token.LessThanEqualTo2] = 4
	precedences[token.GreaterThanEqualTo1] = 4
	precedences[token.GreaterThanEqualTo2] = 4
	precedences[token.InterestinglyEqual] = 4
	precedences[token.Equal] = 4
	precedences[token.Plus] = 5
	precedences[token.Minus] = 5
	precedences[token.Star] = 6
	precedences[token.ForwardSlash] = 6
	precedences[token.BackSlash] = 6
	precedences[token.MOD] = 6
	precedences[token.Exponential] = 7
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
func (i *Interpreter) RunSegment(tokens []token.Token) (ok bool) {
	// Load tokens onto the stack
	i.TokenStack = tokens
	i.TokenPointer = 0
	// 1. Pass if empty line
	if len(tokens) == 0 {
		return true
	}
	if tokens[0].TokenType == token.EndOfLine {
		return true
	}
	// 2. Try numeric variable assignment.  Must be at least 3 tokens.
	if len(tokens) >= 3 {
		// First 2 tokens must be identifier literal followed by = (equal) or := (assign)
		if tokens[0].TokenType == token.IdentifierLiteral &&
			(tokens[1].TokenType == token.Equal || tokens[1].TokenType == token.Assign) {
			// Hand over to rmAssign
			return i.RmAssign()
		}
	}
	// 3. Try built-in / keywords functions.
	if IsKeyword(tokens[0]) {
		switch tokens[0].TokenType {
		case token.REM:
			return true
		case token.PRINT:
			return i.RmPrint()
		case token.GOTO:
			return i.RmGoto()
		case token.RUN:
			return i.RmRun()
		case token.BYE:
			return i.RmBye()
		case token.LIST:
			return i.RmList()
		case token.SAVE:
			return i.RmSave()
		}
	}
	i.ErrorCode = syntaxerror.UnknownCommandProcedure
	i.BadTokenIndex = 0
	i.Message = syntaxerror.ErrorMessage(syntaxerror.UnknownCommandProcedure)
	return false
}

// RunLine attempts to run a line of BASIC code and replies with an error code, the index
// of the token where parsing failed, and a message, or something.
func (i *Interpreter) RunLine(code string) (ok bool) {
	// tokenize the code
	i.Tokenize(code)
	// ensure no illegal chars
	for index, t := range i.CurrentTokens {
		if t.TokenType == token.Illegal {
			i.ErrorCode = syntaxerror.EndOfInstructionExpected
			i.Message = syntaxerror.ErrorMessage(syntaxerror.EndOfInstructionExpected)
			i.BadTokenIndex = index
			return false
		}
	}
	// split the tokens into executable segments for each : token found
	segments := make([][]token.Token, 0)
	this_segment := make([]token.Token, 0)
	for _, t := range i.CurrentTokens {
		if t.TokenType != token.Colon {
			this_segment = append(this_segment, t)
		} else {
			// add EndOfLine to segment before adding to segments slice
			this_segment = append(this_segment, token.Token{token.EndOfLine, ""})
			segments = append(segments, this_segment)
			this_segment = make([]token.Token, 0)
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
	i.ProgramPointer += 1
	return true
}

// ImmediateInput receives a string inputted by the REPL user, processes it and responds
// with a message, if any
func (i *Interpreter) ImmediateInput(code string) (response string) {
	// reset error status and tokenize code
	i.ErrorCode = syntaxerror.Success
	i.Message = ""
	i.BadTokenIndex = 0
	i.LineNumber = -1
	i.Tokenize(code)
	// If the code begins with a line number then add it to the program otherwise try to execute it.
	if i.CurrentTokens[0].TokenType == token.NumericalLiteral {
		// TODO: This needs its own command-----
		// It starts with some kind of number so check if it's an integer, i.e. line number
		lineNumber, err := strconv.ParseFloat(i.CurrentTokens[0].Literal, 64)
		if err == nil {
			if lineNumber == math.Round(lineNumber) {
				// is a line number so format line and add to program
				i.Program[int(lineNumber)] = i.FormatCode(code, -1, true)
				return response
			}
		}
	}
	// Does not begin with line number so try to execute
	ok := i.RunLine(code)
	if !ok {
		// There was an error so the response should include the error message
		if i.LineNumber == -1 {
			// immediate-mode syntax error without line number
			i.g.Print(fmt.Sprintf("Syntax error: %s", i.Message))
			i.g.Print(fmt.Sprintf("  %s", i.FormatCode(code, i.BadTokenIndex, false)))
			response = fmt.Sprintf("Syntax error: %s\n  %s", i.Message, i.FormatCode(code, i.BadTokenIndex, false))
		} else {
			// syntax error with line number
			i.g.Print(fmt.Sprintf("Syntax error in line %d: %s", i.LineNumber, i.Message))
			i.g.Print(fmt.Sprintf("  %d %s", i.LineNumber, i.FormatCode(code, i.BadTokenIndex, false)))
			response = fmt.Sprintf("Syntax error in line %d: %s\n  %d %s", i.LineNumber, i.Message, 10, i.FormatCode(i.Program[i.LineNumber], i.BadTokenIndex, false))
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
		i.CurrentTokens = i.CurrentTokens[1:]
	}
	// bump highlighter if it's pointing at :
	if highlightTokenIndex >= 0 {
		if i.CurrentTokens[highlightTokenIndex].TokenType == token.Colon && len(i.CurrentTokens) > highlightTokenIndex+1 {
			highlightTokenIndex += 1
		}
	}
	// format code and insert highlighter
	for index, t := range i.CurrentTokens {
		if index == highlightTokenIndex {
			formattedCode += "--> "
		}
		if t.TokenType == token.StringLiteral {
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
	for lineNumber, _ := range i.Program {
		ordered = append(ordered, lineNumber)
	}
	sort.Ints(ordered)
	return ordered
}

// IsStringVar returns true if a token represents a string variable
func IsStringVar(t token.Token) bool {
	if t.TokenType == token.IdentifierLiteral && t.Literal[len(t.Literal)-1:] == "$" {
		return true
	} else {
		return false
	}
}

// IsIntVar returns true if a token represents a integer variable
func IsIntVar(t token.Token) bool {
	if t.TokenType == token.IdentifierLiteral && t.Literal[len(t.Literal)-1:] == "%" {
		return true
	} else {
		return false
	}
}

// IsFloatVar returns true if a token represents a float variable
func IsFloatVar(t token.Token) bool {
	if t.TokenType == token.IdentifierLiteral && (t.Literal[len(t.Literal)-1:] != "%" && t.Literal[len(t.Literal)-1:] != "$") {
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
			i.Store[variableName] = fmt.Sprintf("%e", value.(float64))
		} else {
			// store float value directly
			i.Store[variableName] = value.(string)
			return true
		}
	case "%":
		// set integer variable
		if GetType(value) == "float64" {
			// round float and store
			i.Store[variableName] = math.Round(value.(float64))
			return true
		} else {
			// try to parse float from string then round to int
			if valfloat64, err := strconv.ParseFloat(value.(string), 64); err == nil {
				i.Store[variableName] = math.Round(valfloat64)
				return true
			} else {
				i.ErrorCode = syntaxerror.CouldNotInterpretAsANumber
				i.Message = fmt.Sprintf("%s%s", value.(string), syntaxerror.ErrorMessage(syntaxerror.CouldNotInterpretAsANumber))
				return false
			}
		}
	default:
		// set float variable
		if GetType(value) == "float64" {
			// store
			i.Store[variableName] = value.(float64)
			return true
		} else {
			// try to parse float from string
			if valfloat64, err := strconv.ParseFloat(value.(string), 64); err == nil {
				i.Store[variableName] = valfloat64
				return true
			} else {
				i.ErrorCode = syntaxerror.CouldNotInterpretAsANumber
				i.Message = fmt.Sprintf("%s%s", value.(string), syntaxerror.ErrorMessage(syntaxerror.CouldNotInterpretAsANumber))
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
	val, ok := i.Store[variableName]
	if !ok {
		i.ErrorCode = syntaxerror.HasNotBeenDefined
		i.Message = fmt.Sprintf("%s%s", variableName, syntaxerror.ErrorMessage(syntaxerror.HasNotBeenDefined))
		return 0, false
	} else {
		return val, true
	}
}

// GetValueFromToken receives a token representing either a variable or a literal, gets
// the value, and ~casts~ it to the required type
func (i *Interpreter) GetValueFromToken(t token.Token, castTo string) (value interface{}, ok bool) {
	if t.TokenType == token.IdentifierLiteral {
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
				i.ErrorCode = syntaxerror.CouldNotInterpretAsANumber
				i.Message = fmt.Sprintf("%s%s", value.(string), syntaxerror.ErrorMessage(syntaxerror.CouldNotInterpretAsANumber))
				return 0, false
			}
		case "int64":
			// try to parse float from string
			if valfloat64, err := strconv.ParseFloat(value.(string), 64); err == nil {
				return math.Round(valfloat64), true
			} else {
				i.ErrorCode = syntaxerror.CouldNotInterpretAsANumber
				i.Message = fmt.Sprintf("%s%s", value.(string), syntaxerror.ErrorMessage(syntaxerror.CouldNotInterpretAsANumber))
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
	originalPosition := i.TokenPointer
	val, ok := i.EvaluateExpression()
	if !ok {
		// broken expression so reset pointer and return false
		i.TokenPointer = originalPosition
		return 0, false
	} else {
		if GetType(val) != "float64" {
			// is not float
			// error?
			i.TokenPointer = originalPosition
			return 0, false
		} else {
			return val.(float64), true
		}
	}
}

// AcceptAnyString checks if the current token represents a string, returns the value
// and advances the pointer if so.
func (i *Interpreter) AcceptAnyString() (acceptedValue string, acceptOk bool) {
	originalPosition := i.TokenPointer
	val, ok := i.EvaluateExpression()
	if !ok {
		// broken expression so reset pointer and return false
		i.TokenPointer = originalPosition
		return "", false
	} else {
		if GetType(val) != "string" {
			// is not string
			// error?
			i.TokenPointer = originalPosition
			return "", false
		} else {
			return val.(string), true
		}
	}
}

// AcceptAnyOfTheseTokens checks if the current token matches any that are passed in a slice and, if so,
// returns the token and advances the pointer.
func (i *Interpreter) AcceptAnyOfTheseTokens(acceptableTokens []int) (acceptedToken token.Token, acceptOk bool) {
	for _, tokenType := range acceptableTokens {
		if tokenType == i.TokenStack[i.TokenPointer].TokenType {
			// found a match, advanced pointer and return the token type
			i.TokenPointer++
			return i.TokenStack[i.TokenPointer], true
		}
	}
	// no matches
	return token.Token{}, false
}

// IsAnyOfTheseTokens checks if the current token matches any that are passed in a slice and, if so,
// returns true but *does not advance the pointer*.
func (i *Interpreter) IsAnyOfTheseTokens(acceptableTokens []int) bool {
	for _, tokenType := range acceptableTokens {
		if tokenType == i.TokenStack[i.TokenPointer].TokenType {
			// found a match
			return true
		}
	}
	// no matches
	return false
}

// ExtractExpression receives a slice of tokens that represent an expression and returns
// all those tokens up to where the expression ends.
func (i *Interpreter) ExtractExpression() (expressionTokens []token.Token) {
	tokenStackSlice := i.TokenStack[i.TokenPointer:]
	for index, t := range tokenStackSlice {
		// The following tokens can delimit an expression
		if t.TokenType == token.Comma || t.TokenType == token.Semicolon || t.TokenType == token.EndOfLine || t.TokenType == token.Exclamation || t.TokenType == token.TO {
			break
		}
		// An expression is also delimited if one operand follows another directly
		if index < (len(tokenStackSlice) - 1) {
			if IsOperand(t) && IsOperand(tokenStackSlice[index+1]) {
				expressionTokens = append(expressionTokens, t)
				i.TokenPointer++
				break
			}
		}
		expressionTokens = append(expressionTokens, t)
		i.TokenPointer++
	}
	return expressionTokens
}

// EndOfTokens returns true if no more tokens are to be evaluated in the token stack
func (i *Interpreter) EndOfTokens() bool {
	if i.TokenPointer > len(i.TokenStack) {
		return true
	}
	if i.TokenStack[i.TokenPointer].TokenType == token.EndOfLine {
		return true
	}
	return false
}
