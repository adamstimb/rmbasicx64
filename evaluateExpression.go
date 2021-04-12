package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

// EvaluateExpression receives tokens that appear to represent an expression, tries to evaluate it
// and returns the result.
func (i *Interpreter) EvaluateExpression() (result interface{}, ok bool) {
	tokens := i.ExtractExpression()
	// If exactly one token representing a literal or variable we don't need to evaluate it
	if len(tokens) == 1 {
		switch tokens[0].TokenType {
		case StringLiteral:
			return tokens[0].Literal, true
		case NumericalLiteral:
			if valfloat64, err := strconv.ParseFloat(tokens[0].Literal, 64); err == nil {
				return valfloat64, true
			} // this should never fail unless the scanner parses numeric literals incorrectly
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
	for index, t := range tokens {
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
		i.errorCode = InvalidExpressionFound
		i.message = errorMessage(InvalidExpressionFound)
		i.badTokenIndex = index
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
			if t.TokenType == NumericalLiteral {
				// Is numeric but test it can be parsed before pushing token to operand stack
				if valfloat64, err := strconv.ParseFloat(t.Literal, 64); err == nil {
					// push
					operandStack = append([]interface{}{valfloat64}, operandStack...)
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
					op2 = math.Round(op2)
					result = float64(^int(op2))
					// pop the stack, push new result and skip to next item in the postfix
					operandStack = operandStack[1:]
					operandStack = append([]interface{}{result}, operandStack...)
					continue
					//}
				} else {
					i.errorCode = NumericExpressionNeeded
					i.badTokenIndex = 0
					i.message = errorMessage(NumericExpressionNeeded)
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
					i.errorCode = InvalidExpressionFound
					i.badTokenIndex = 0
					i.message = errorMessage(InvalidExpressionFound)
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
					if op2 == float64(0) {
						i.errorCode = TryingToDivideByZero
						i.badTokenIndex = 0
						i.message = errorMessage(TryingToDivideByZero)
						return 0, false
					}
					result = op1 / op2
				case BackSlash:
					if op2 == float64(0) {
						i.errorCode = TryingToDivideByZero
						i.badTokenIndex = 0
						i.message = errorMessage(TryingToDivideByZero)
						return 0, false
					}
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
					op1 = math.Round(op1)
					op2 = math.Round(op2)
					result = float64(int(op1) & int(op2))
				case OR:
					op1 = math.Round(op1)
					op2 = math.Round(op2)
					result = float64(int(op1) | int(op2))
				case XOR:
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
