package rmbasicx64

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
)

// EvaluateExpression receives tokens that appear to represent an expression, tries to evaluate it
// and returns the result.
func (i *Interpreter) EvaluateExpression() (result interface{}, ok bool) {
	tokens := i.ExtractExpression()
	// If exactly one token representing a literal or variable we don't need to evaluate it
	if len(tokens) == 1 {
		switch tokens[0].TokenType {
		case token.StringLiteral:
			return tokens[0].Literal, true
		case token.NumericalLiteral:
			if valfloat64, err := strconv.ParseFloat(tokens[0].Literal, 64); err == nil {
				return valfloat64, true
			} // this should never fail unless the scanner parses numeric literals incorrectly
		case token.IdentifierLiteral:
			var val interface{}
			val, ok = i.GetVar(tokens[0].Literal)
			if !ok {
				i.BadTokenIndex = 0 + i.TokenPointer
				return 0, false
			} else {
				return val, true
			}
		}
	}
	// Make the postfix then evaluate it following Carrano's pseudocode:
	// http://www.solomonlrussell.com/spring16/cs2/ClassSource/Week6/stackcode.html
	// (this has been extended quite a lot to deal with expressions that mix numeric and string values)
	postfix := make([]token.Token, 0)
	operatorStack := make([]token.Token, 0)
	for index, t := range tokens {
		if IsOperand(t) {
			postfix = append(postfix, t)
			continue
		}
		if t.TokenType == token.LeftParen {
			// push
			operatorStack = append([]token.Token{t}, operatorStack...)
			continue
		}
		if t.TokenType == token.RightParen {
			// pop operator stack until matching LeftParen
			for operatorStack[0].TokenType != token.LeftParen {
				postfix = append(postfix, operatorStack[0])
				operatorStack = operatorStack[1:]
			}
			// pop and continue
			operatorStack = operatorStack[1:]
			continue
		}
		if IsOperator(t) {
			for len(operatorStack) > 0 &&
				operatorStack[0].TokenType != token.LeftParen &&
				Precedence(t) <= Precedence(operatorStack[0]) {
				postfix = append(postfix, operatorStack[0])
				// pop
				operatorStack = operatorStack[1:]
			}
			// push
			operatorStack = append([]token.Token{t}, operatorStack...)
			continue
		}
		i.ErrorCode = syntaxerror.InvalidExpressionFound
		i.Message = syntaxerror.ErrorMessage(syntaxerror.InvalidExpressionFound)
		i.BadTokenIndex = index
		return 0, false
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
			if t.TokenType == token.NumericalLiteral {
				// Is numeric but test it can be parsed before pushing token to operand stack
				if valfloat64, err := strconv.ParseFloat(t.Literal, 64); err == nil {
					// push
					operandStack = append([]interface{}{valfloat64}, operandStack...)
				}
			}
			if t.TokenType == token.StringLiteral {
				// push it as-is
				operandStack = append([]interface{}{t.Literal}, operandStack...)
			}
			if t.TokenType == token.IdentifierLiteral {
				// Is identifier, so first test it has been defined by looking in the store
				if _, ok := i.Store[t.Literal]; ok {
					if t.Literal[len(tokens[0].Literal)-1:] != "$" {
						// Represents a numeric value
						valfloat64, ok := i.Store[t.Literal].(float64)
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
					i.ErrorCode = syntaxerror.HasNotBeenDefined
					i.BadTokenIndex = 0
					i.Message = fmt.Sprintf("%s%s", t.Literal, syntaxerror.ErrorMessage(syntaxerror.HasNotBeenDefined))
					return 0, false
				}
			}
		} else {
			// Apply operator
			// First try unary operators, currently only NOT is implemented so:
			// Get operand 2 but *** DO NOT POP THE STACK ***
			operand2 := operandStack[0]
			if t.TokenType == token.NOT {
				// Is unary NOT but we can only apply this to rounded floats or ints
				if GetType(operand2) != "string" {
					op2 := operand2.(float64)
					op2 = math.Round(op2)
					result = float64(^int(op2))
					// pop the stack, push new result and skip to next item in the postfix
					operandStack = operandStack[1:]
					operandStack = append([]interface{}{result}, operandStack...)
					continue
					//}
				} else {
					i.ErrorCode = syntaxerror.NumericExpressionNeeded
					i.BadTokenIndex = 0
					i.Message = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
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
					op1 = RenderNumberAsString(operand1.(float64))
				}
				if GetType(operand2) == "string" {
					op2 = operand2.(string)
				}
				if GetType(operand2) == "float64" {
					op2 = RenderNumberAsString(operand2.(float64))
				}
				// Apply operation
				switch t.TokenType {
				case token.Plus:
					result = op1 + op2
				case token.Equal:
					if op1 == op2 {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case token.InterestinglyEqual:
					if strings.EqualFold(op1, op2) {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case token.LessThan:
					if WeighString(op1) < WeighString(op2) {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case token.GreaterThan:
					if WeighString(op1) > WeighString(op2) {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case token.LessThanEqualTo1:
					if WeighString(op1) <= WeighString(op2) {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case token.LessThanEqualTo2:
					if WeighString(op1) <= WeighString(op2) {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case token.GreaterThanEqualTo1:
					if WeighString(op1) >= WeighString(op2) {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case token.GreaterThanEqualTo2:
					if WeighString(op1) >= WeighString(op2) {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case token.Inequality1:
					if WeighString(op1) != WeighString(op2) {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case token.Inequality2:
					if WeighString(op1) != WeighString(op2) {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				default:
					i.ErrorCode = syntaxerror.InvalidExpressionFound
					i.BadTokenIndex = 0
					i.Message = syntaxerror.ErrorMessage(syntaxerror.InvalidExpressionFound)
					return 0, false
				}
			} else {
				// Numeric binary expression:
				// Can assume both operands are numeric, i.e. float64 so convert them directly
				op1 := operand1.(float64)
				op2 := operand2.(float64)
				// Apply operation
				switch t.TokenType {
				case token.Minus:
					result = op1 - op2
				case token.Plus:
					result = op1 + op2
				case token.ForwardSlash:
					if op2 == float64(0) {
						i.ErrorCode = syntaxerror.TryingToDivideByZero
						i.BadTokenIndex = 0
						i.Message = syntaxerror.ErrorMessage(syntaxerror.TryingToDivideByZero)
						return 0, false
					}
					result = op1 / op2
				case token.BackSlash:
					if op2 == float64(0) {
						i.ErrorCode = syntaxerror.TryingToDivideByZero
						i.BadTokenIndex = 0
						i.Message = syntaxerror.ErrorMessage(syntaxerror.TryingToDivideByZero)
						return 0, false
					}
					result = float64(int(op1) / int(op2))
				case token.Star:
					result = op1 * op2
				case token.Exponential:
					result = math.Pow(op1, op2)
				case token.MOD:
					result = float64(int(op1) % int(op2))
				case token.Equal:
					if op1 == op2 {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case token.InterestinglyEqual:
					if op1 == op2 {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case token.LessThan:
					if op1 < op2 {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case token.GreaterThan:
					if op1 > op2 {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case token.LessThanEqualTo1:
					if op1 <= op2 {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case token.LessThanEqualTo2:
					if op1 <= op2 {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case token.GreaterThanEqualTo1:
					if op1 >= op2 {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case token.GreaterThanEqualTo2:
					if op1 >= op2 {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case token.Inequality1:
					if op1 != op2 {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case token.Inequality2:
					if op1 != op2 {
						result = float64(-1)
					} else {
						result = float64(0)
					}
				case token.AND:
					op1 = math.Round(op1)
					op2 = math.Round(op2)
					result = float64(int(op1) & int(op2))
				case token.OR:
					op1 = math.Round(op1)
					op2 = math.Round(op2)
					result = float64(int(op1) | int(op2))
				case token.XOR:
					op1 = math.Round(op1)
					op2 = math.Round(op2)
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
