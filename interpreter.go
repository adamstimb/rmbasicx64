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
}

// Init initializes the Interpreter.
func (i *Interpreter) Init() {
	i.store = make(map[string]interface{})
	i.program = make(map[int]string)
	i.currentTokens = []Token{}
	i.errorCode = Success
	i.lineNumber = -1
	i.badTokenIndex = -1
	i.message = ""
	i.programPointer = 0
	i.tokenStack = []Token{}
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

// EvaluateExpression receives tokens that appear to represent an expression, tries to evaluate it
// and returns the result.
func (i *Interpreter) EvaluateExpression(tokens []Token) (result interface{}, ok bool) {
	// If exactly one token representing a literal or variable we don't need to evaluate it
	if len(tokens) == 1 {
		switch tokens[0].TokenType {
		case StringLiteral:
			return tokens[0].Literal, true
		case NumericalLiteral:
			if valfloat64, err := strconv.ParseFloat(tokens[0].Literal, 64); err == nil {
				return valfloat64, true
			} else {
				i.errorCode = CouldNotInterpretAsANumber
				i.message = fmt.Sprintf("%s%s", tokens[0].Literal, errorMessage(CouldNotInterpretAsANumber))
				i.badTokenIndex = 0 + i.tokenPointer
				return 0, false
			}
		case IdentifierLiteral:
			var val interface{}
			val, ok = i.GetVar(tokens[0].Literal)
			if !ok {
				i.badTokenIndex = 0 + i.tokenPointer
				return 0, false
			} else {
				return val, true
			}
		}
	}
	// Make the postfix then evaluate it following Carrano's pseudocode:
	// http://www.solomonlrussell.com/spring16/cs2/ClassSource/Week6/stackcode.html
	// (this has been extended quite a lot to deal with expressions that mix numeric and string values)
	postfix := make([]Token, 0)
	operatorStack := make([]Token, 0)
	for _, t := range tokens {
		if IsOperand(t) {
			postfix = append(postfix, t)
			continue
		}
		if t.TokenType == LeftParen {
			// push
			operatorStack = append([]Token{t}, operatorStack...)
			continue
		}
		if t.TokenType == RightParen {
			// pop operator stack until matching LeftParen
			for operatorStack[0].TokenType != LeftParen {
				postfix = append(postfix, operatorStack[0])
				operatorStack = operatorStack[1:]
			}
			// pop and continue
			operatorStack = operatorStack[1:]
			continue
		}
		if IsOperator(t) {
			for len(operatorStack) > 0 &&
				operatorStack[0].TokenType != LeftParen &&
				Precedence(t) <= Precedence(operatorStack[0]) {
				postfix = append(postfix, operatorStack[0])
				// pop
				operatorStack = operatorStack[1:]
			}
			// push
			operatorStack = append([]Token{t}, operatorStack...)
			continue
		}
	}
	for len(operatorStack) > 0 {
		postfix = append(postfix, operatorStack[0])
		// pop
		operatorStack = operatorStack[1:]
	}

	// Now evaluate the postfix:
	operandStack := make([]interface{}, 0)
	for _, t := range postfix {
		if IsOperand(t) {
			// Handle operand
			// Get the value and data type represented by the token.
			if t.TokenType == NumericalLiteral {
				// Is numeric but test it can be parsed before pushing token to operand stack
				if valfloat64, err := strconv.ParseFloat(t.Literal, 64); err == nil {
					// push
					operandStack = append([]interface{}{valfloat64}, operandStack...)
				} else {
					// Is meant to represent a numeric value but it can't be parsed (this should never actually happen...maybe remove it?)
					i.errorCode = CouldNotInterpretAsANumber
					i.badTokenIndex = 0
					i.message = fmt.Sprintf("%s%s", t.Literal, errorMessage(CouldNotInterpretAsANumber))
					return 0, false
				}
			}
			if t.TokenType == StringLiteral {
				// push it as-is
				operandStack = append([]interface{}{t.Literal}, operandStack...)
			}
			if t.TokenType == IdentifierLiteral {
				// Is identifier, so first test it has been defined by looking in the store
				if _, ok := i.store[t.Literal]; ok {
					if t.Literal[len(tokens[0].Literal)-1:] != "$" {
						// Represents a numeric value
						valfloat64, ok := i.store[t.Literal].(float64)
						if !ok {
							// This should not happen therefore fatal
							log.Fatalf("Fatal error!")
						} else {
							// push
							operandStack = append([]interface{}{valfloat64}, operandStack...)
						}
					}
				} else {
					// Variable not defined
					i.errorCode = HasNotBeenDefined
					i.badTokenIndex = 0
					i.message = fmt.Sprintf("%s%s", t.Literal, errorMessage(HasNotBeenDefined))
					return 0, false
				}
			}
		} else {
			// Apply operator
			// First try unary operators, currently only NOT is implemented so:
			// Get operand 2 but *** DO NOT POP THE STACK ***
			operand2 := operandStack[0]
			if t.TokenType == NOT {
				// Is unary NOT but we can only apply this to rounded floats or ints
				if GetType(operand2) != "string" {
					op2 := operand2.(float64)
					if op2 != math.Round(op2) {
						i.errorCode = CannotPerformBitwiseOperationsOnFloatValues
						i.badTokenIndex = -1
						i.message = fmt.Sprintf("%s%s", t.Literal, errorMessage(HasNotBeenDefined))
						return 0, false
					} else {
						result = float64(^int(op2))
						// pop the stack, push new result and skip to next item in the postfix
						operandStack = operandStack[1:]
						operandStack = append([]interface{}{result}, operandStack...)
						continue
					}
				} else {
					i.errorCode = CannotPerformBitwiseOperationsOnStringValues
					i.badTokenIndex = -1
					i.message = errorMessage(CannotPerformBitwiseOperationsOnStringValues)
					return 0, false
				}
			}
			// Binary operator
			// pop the stack and get operand 1
			operandStack = operandStack[1:]
			operand1 := operandStack[0]
			// pop the stack again
			operandStack = operandStack[1:]
			// Now decide if it's a string expression, numeric expression or invalid expression
			// If one or both operands are string type then it's a string expression
			// If both operands are numeric then it's a numeric expression (however this will
			// not always be true, for example if the operator is a function that returns a string...)
			if GetType(operand1) == "string" || GetType(operand2) == "string" {
				// String binary expression
				// If either operand is numeric, convert it to string before applying the operator
				op1 := ""
				op2 := ""
				if GetType(operand1) == "string" {
					op1 = operand1.(string)
				}
				if GetType(operand1) == "float64" {
					op1 = fmt.Sprintf("%e", operand1.(float64))
				}
				if GetType(operand2) == "string" {
					op2 = operand2.(string)
				}
				if GetType(operand2) == "float64" {
					// Use scientific notation if...what? Check manual.
					if operand2.(float64) == math.Round(operand2.(float64)) {
						op2 = fmt.Sprintf("%.0f", operand2.(float64))
					} else {
						op2 = fmt.Sprintf("%e", operand2.(float64))
					}
				}
				// Apply operation
				switch t.TokenType {
				case Plus:
					result = op1 + op2
				case Equal:
					if op1 == op2 {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case InterestinglyEqual:
					if strings.EqualFold(op1, op2) {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case LessThan:
					if WeighString(op1) < WeighString(op2) {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case GreaterThan:
					if WeighString(op1) > WeighString(op2) {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case LessThanEqualTo1:
					if WeighString(op1) <= WeighString(op2) {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case LessThanEqualTo2:
					if WeighString(op1) <= WeighString(op2) {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case GreaterThanEqualTo1:
					if WeighString(op1) >= WeighString(op2) {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case GreaterThanEqualTo2:
					if WeighString(op1) >= WeighString(op2) {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case Inequality1:
					if WeighString(op1) != WeighString(op2) {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case Inequality2:
					if WeighString(op1) != WeighString(op2) {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				default:
					i.errorCode = InvalidExpression
					i.badTokenIndex = 0
					i.message = fmt.Sprintf("%s%s", t.Literal, errorMessage(InvalidExpression))
					return 0, false
				}
			} else {
				// Numeric binary expression:
				// Can assume both operands are numeric, i.e. float64 so convert them directly
				op1 := operand1.(float64)
				op2 := operand2.(float64)
				// Apply operation
				switch t.TokenType {
				case Minus:
					result = op1 - op2
				case Plus:
					result = op1 + op2
				case ForwardSlash:
					result = op1 / op2
				case BackSlash:
					result = float64(int(op1) / int(op2))
				case Star:
					result = op1 * op2
				case Exponential:
					result = math.Pow(op1, op2)
				case MOD:
					result = float64(int(op1) % int(op2))
				case Equal:
					if op1 == op2 {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case InterestinglyEqual:
					if op1 == op2 {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case LessThan:
					if op1 < op2 {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case GreaterThan:
					if op1 > op2 {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case LessThanEqualTo1:
					if op1 <= op2 {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case LessThanEqualTo2:
					if op1 <= op2 {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case GreaterThanEqualTo1:
					if op1 >= op2 {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case GreaterThanEqualTo2:
					if op1 >= op2 {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case Inequality1:
					if op1 != op2 {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case Inequality2:
					if op1 != op2 {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case AND:
					if op1 != math.Round(op1) || op2 != math.Round(op2) {
						i.errorCode = CannotPerformBitwiseOperationsOnFloatValues
						i.badTokenIndex = 0
						i.message = errorMessage(CannotPerformBitwiseOperationsOnFloatValues)
						return 0, false
					}
					result = float64(int(op1) & int(op2))
				case OR:
					if op1 != math.Round(op1) || op2 != math.Round(op2) {
						i.errorCode = CannotPerformBitwiseOperationsOnFloatValues
						i.badTokenIndex = 0
						i.message = errorMessage(CannotPerformBitwiseOperationsOnFloatValues)
						return 0, false
					}
					result = float64(int(op1) | int(op2))
				case XOR:
					if op1 != math.Round(op1) || op2 != math.Round(op2) {
						i.errorCode = CannotPerformBitwiseOperationsOnFloatValues
						i.badTokenIndex = 0
						i.message = errorMessage(CannotPerformBitwiseOperationsOnFloatValues)
						return 0, false
					}
					result = float64(int(op1) ^ int(op2))
				}
			}
			// push
			operandStack = append([]interface{}{result}, operandStack...)
		}
	}
	// Evaluation complete
	return result, true
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
		case PRINT:
			return i.rmPrint(tokens)
		case GOTO:
			return i.rmGoto(tokens)
		case RUN:
			return i.rmRun()
		}
	}
	i.errorCode = ExpectedAKeywordLineNumberExpressionVariableAssignmentOrProcedureCall
	i.badTokenIndex = 0
	i.message = errorMessage(ExpectedAKeywordLineNumberExpressionVariableAssignmentOrProcedureCall)
	return false
}

// RunLine attempts to run a line of BASIC code and replies with an error code, the index
// of the token where parsing failed, and a message, or something.
func (i *Interpreter) RunLine(code string) (ok bool) {
	// tokenize the code
	i.Tokenize(code)
	// split the tokens into executable segments for each : token found
	segments := make([][]Token, 0)
	this_segment := make([]Token, 0)
	for _, token := range i.currentTokens {
		if token.TokenType != Colon {
			this_segment = append(this_segment, token)
		} else {
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
			response = fmt.Sprintf("Syntax error: %s\n%s", i.message, i.FormatCode(code, i.badTokenIndex, false))
		} else {
			// syntax error with line number
			response = fmt.Sprintf("Syntax error in line %d: %s\n%s", i.lineNumber, i.message, i.FormatCode(code, i.badTokenIndex, false))
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
		}
	case "float64":
		switch castTo {
		case "string":
			return fmt.Sprintf("%e", value.(float64)), true
		case "float64":
			return value, true
		case "int64":
			return math.Round(value.(float64)), true
		}
	}
	log.Fatalf("Fatal error!")
	return 0, false
}

// ExtractExpression receives a slice of tokens that represent an expression and returns
// all those tokens up to where the expression ends.
func (i *Interpreter) ExtractExpression() (expressionTokens []Token) {
	for _, t := range i.tokenStack[i.tokenPointer:] {
		if t.TokenType == Comma || t.TokenType == Semicolon || t.TokenType == EndOfLine || IsKeyword(t) {
			break
		}
		expressionTokens = append(expressionTokens, t)
		i.tokenPointer++
	}
	return expressionTokens
}
