package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

// Interpreter is the BASIC interpreter itself and behaves as a state machine that
// can receieve, store and interpret BASIC code and execute the code to update its
// own state.
type Interpreter struct {
	store         map[string]interface{} // A map for storing variables and array (the key is the variable name)
	program       map[int]string         // A map for storing a program (the key is the line number)
	currentTokens []Token                // A line of tokens for immediate execution
}

// Init initializes the Interpreter.
func (i *Interpreter) Init() {
	i.store = make(map[string]interface{})
	i.program = make(map[int]string)
	i.currentTokens = []Token{}
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
		AND, OR}
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
func (i *Interpreter) EvaluateExpression(tokens []Token) (errorCode, badTokenIndex int, message string, result interface{}) {
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
	for index, t := range postfix {
		if IsOperand(t) {
			// Get the value and data type represented by the token.
			if t.TokenType == NumericalLiteral {
				// Is numeric but test it can be parsed before pushing token to operand stack
				if valfloat64, err := strconv.ParseFloat(t.Literal, 64); err == nil {
					// push
					operandStack = append([]interface{}{valfloat64}, operandStack...)
				} else {
					// Is meant to represent a numeric value but it can't be parsed (this should never actually happen...maybe remove it?)
					return CouldNotInterpretAsANumber, index, fmt.Sprintf("%s%s", t.Literal, errorMessage(CouldNotInterpretAsANumber)), 0
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
					return HasNotBeenDefined, index, fmt.Sprintf("%s%s", t.Literal, errorMessage(HasNotBeenDefined)), 0
				}
			}
		} else {
			// Get operands 1 and 2, and their operator
			operand2 := operandStack[0]
			// pop
			operandStack = operandStack[1:]
			operand1 := operandStack[0]
			// pop
			operandStack = operandStack[1:]
			// Now decide if it's a string expression, numeric expression or invalid expression
			// If one or both operands are string type then it's a string expression
			// If both operands are numeric then it's a numeric expression (however this will
			// not always be true, for example if the operator is a function that returns a string...)
			if GetType(operand1) == "string" || GetType(operand2) == "string" {
				// Is valid string expression
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
					return InvalidExpression, index, fmt.Sprintf("%s%s", t.Literal, errorMessage(InvalidExpression)), 0
				}
			} else {
				// Numeric expression:
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
						return CannotPerformBitwiseOperationsOnFloatValues, index, errorMessage(CannotPerformBitwiseOperationsOnFloatValues), 0
					}
					result = float64(int(op1) & int(op2))
				case OR:
					if op1 != math.Round(op1) || op2 != math.Round(op2) {
						return CannotPerformBitwiseOperationsOnFloatValues, index, errorMessage(CannotPerformBitwiseOperationsOnFloatValues), 0
					}
					result = float64(int(op1) | int(op2))
				}
			}
			// push
			operandStack = append([]interface{}{result}, operandStack...)
		}
	}
	// Evaluation successful, errorCode = 0
	return 0, 0, "", result
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
func (i *Interpreter) RunSegment(tokens []Token) (errorCode, badTokenIndex int, message string) {
	// 1. Pass if empty line
	if len(tokens) == 0 {
		return Success, 0, ""
	}
	if tokens[0].TokenType == EndOfLine {
		return Success, 0, ""
	}
	// 2. Try string variable assignment.  Must be at least 3 tokens.

	// 3. Try numeric variable assignment.  Must be at least 3 tokens.
	if len(tokens) >= 3 {
		// First 2 tokens must be identifier literal followed by = (equal) or := (assign)
		if tokens[0].TokenType == IdentifierLiteral &&
			(tokens[1].TokenType == Equal || tokens[1].TokenType == Assign) {
			// If exactly four tokens and the 3rd token is a numerical literal then we don't
			// have anything to evaluate
			if len(tokens) == 4 && tokens[2].TokenType == NumericalLiteral {
				if val, err := strconv.ParseFloat(tokens[2].Literal, 64); err == nil {
					// round val if variable is integer type, i.e. ends with %
					if tokens[0].Literal[len(tokens[0].Literal)-1:] == "%" {
						val = math.Round(val)
					}
					i.store[tokens[0].Literal] = val
					return Success, -1, ""
				} else {
					return CouldNotInterpretAsANumber, 2, fmt.Sprintf("%s%s", tokens[2].Literal, errorMessage(CouldNotInterpretAsANumber))
				}
			} else {
				// evaluate result then store
				errorCode, _, message, result := i.EvaluateExpression(tokens[2:])
				if errorCode == Success {
					// Evaluation was successful so check data type and store
					if GetType(result) == "string" {
						// Store the result
						i.store[tokens[0].Literal] = result.(string)
					} else {
						// round val if variable is integer type, i.e. ends with %
						if tokens[0].Literal[len(tokens[0].Literal)-1:] == "%" {
							result = math.Round(result.(float64))
						}
						// Store the result
						i.store[tokens[0].Literal] = result
					}
					return Success, -1, ""
				} else {
					// Something went wrong so return error info
					return errorCode, 2, message
				}
			}
		}
		// Catch case where a keyword has been used as a variable name to assign to
		if IsKeyword(tokens[0]) &&
			(tokens[1].TokenType == Equal || tokens[1].TokenType == Assign) {
			return IsAKeywordAndCannotBeUsedAsAVariableName, 0, fmt.Sprintf("%s%s", tokens[0].Literal, errorMessage(IsAKeywordAndCannotBeUsedAsAVariableName))
		}
	}
	// 3. Try built-in / keywords functions.
	if IsKeyword(tokens[0]) {
		// Try PRINT
		if tokens[0].TokenType == PRINT {
			// PRINT with no args
			if len(tokens) == 1 {
				fmt.Println("")
				return Success, -1, ""
			}
			if len(tokens) > 1 {
				if tokens[1].TokenType == EndOfLine {
					// Still PRINT with no args
					fmt.Println("")
					return Success, -1, ""
				}
				if tokens[1].TokenType == StringLiteral {
					// PRINT "hello"
					fmt.Println(tokens[1].Literal)
					return Success, -1, ""
				}
			}
		}
	}
	return ExpectedAKeywordLineNumberExpressionVariableAssignmentOrProcedureCall, 0, errorMessage(ExpectedAKeywordLineNumberExpressionVariableAssignmentOrProcedureCall)
}

// RunLine attempts to run a line of BASIC code and replies with an error code, the index
// of the token where parsing failed, and a message, or something.
func (i *Interpreter) RunLine(code string) (errorCode, badTokenIndex int, message string) {
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
		errorCode, badTokenIndex, message = i.RunSegment(segment)
		if errorCode != 0 {
			break
		}
	}
	return errorCode, badTokenIndex + badTokenOffset, message
}

// FormatCode receives a line of BASIC code and returns it formatted.  If a number
// > 0 is passed for highlightTokenIndex, the corresponding token is highlighted
// with arrows; this is used for printing error messages.
func (i *Interpreter) FormatCode(code string, highlightTokenIndex int) string {
	i.Tokenize(code)
	formattedCode := ""
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
