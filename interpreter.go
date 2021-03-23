package main

import (
	"math"
	"strconv"
)

// Interpreter is the BASIC interpreter itself and behaves as a state machine that
// can receieve, store and interpret BASIC code and execute the code to update its
// own state.
type Interpreter struct {
	store         map[string]interface{} // A map for storing variables and array (the key is the variable name)
	program       map[int]string         // A map for storing a program (the key is the line number)
	currentTokens []Token                // A line of tokens for immediate execution
	operandStack  []float64              // The operand stack for expression evaluation
	operatorStack []Token                // The operator stack for expressin evaluation
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
	operators := []int{Minus, Plus, ForwardSlash, Star, Exponential, BackSlash}
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
	operands := []int{NumericalLiteral, IdentifierLiteral}
	for _, op := range operands {
		if op == t.TokenType {
			return true
		}
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
	precedences[Plus] = 5
	precedences[Minus] = 5
	precedences[Star] = 6
	precedences[ForwardSlash] = 6
	precedences[BackSlash] = 6
	precedences[MOD] = 6
	precedences[Exponential] = 7

	return precedences[t.TokenType]
}

// Evaluate receives tokens that appear to represent an expression, tries to evaluate it
// and returns the result.
func (i *Interpreter) Evaluate(tokens []Token) (errorCode, badTokenIndex int, message string, result float64) {
	// Make the postfix then evaluate it following Carrano's pseudocode:
	// http://www.solomonlrussell.com/spring16/cs2/ClassSource/Week6/stackcode.html
	postfix := make([]Token, 0)
	i.operatorStack = []Token{}
	for _, t := range tokens {
		if IsOperand(t) {
			postfix = append(postfix, t)
			continue
		}
		if t.TokenType == LeftParen {
			// push
			i.operatorStack = append([]Token{t}, i.operatorStack...)
			continue
		}
		if t.TokenType == RightParen {
			// pop operator stack until matching LeftParen
			for i.operatorStack[0].TokenType != LeftParen {
				postfix = append(postfix, i.operatorStack[0])
				i.operatorStack = i.operatorStack[1:]
			}
			// pop and continue
			i.operatorStack = i.operatorStack[1:]
			continue
		}
		if IsOperator(t) {
			for len(i.operatorStack) > 0 &&
				i.operatorStack[0].TokenType != LeftParen &&
				Precedence(t) <= Precedence(i.operatorStack[0]) {
				postfix = append(postfix, i.operatorStack[0])
				// pop
				i.operatorStack = i.operatorStack[1:]
			}
			// push
			i.operatorStack = append([]Token{t}, i.operatorStack...)
			continue
		}
	}
	for len(i.operatorStack) > 0 {
		postfix = append(postfix, i.operatorStack[0])
		// pop
		i.operatorStack = i.operatorStack[1:]
	}

	// Now evaluate the postfix:
	i.operandStack = []float64{}
	for _, t := range postfix {
		if IsOperand(t) {
			// Get the value represented by the token.  If it's a numerical
			// literal then we have to convert it to float64.  If it's an
			// identifier literal then we have to retrieve the value from the
			// store and convert to float64. Anything else returns an error.
			operand := float64(0)
			if t.TokenType == NumericalLiteral {
				if val, err := strconv.ParseFloat(t.Literal, 64); err == nil {
					operand = val
				} else {
					return 1, 0, "Could not interpret this as a number", 0
				}
			}
			if t.TokenType == IdentifierLiteral {
				if _, ok := i.store[t.Literal]; ok {
					valfloat64, ok := i.store[t.Literal].(float64)
					if !ok {
						return 1, 0, "Could not interpret stored value as a number", 0
					} else {
						operand = valfloat64
					}
				}
			}
			// push
			i.operandStack = append([]float64{operand}, i.operandStack...)
		} else {
			operand2 := i.operandStack[0]
			// pop
			i.operandStack = i.operandStack[1:]
			operand1 := i.operandStack[0]
			// pop
			i.operandStack = i.operandStack[1:]
			// Apply operation to operand1 and operand2
			switch t.TokenType {
			case Minus:
				result = operand1 - operand2
			case Plus:
				result = operand1 + operand2
			case ForwardSlash:
				result = operand1 / operand2
			case BackSlash:
				result = float64(int(operand1) / int(operand2))
			case Star:
				result = operand1 * operand2
			case Exponential:
				result = math.Pow(operand1, operand2)
			case MOD:
				result = float64(int(operand1) % int(operand2))
				// TODO: Comparitors will also just cast result to float64
			}
			// push
			i.operandStack = append([]float64{result}, i.operandStack...)
		}
	}
	return 0, 0, "", result
}

// Execute attempts to execute the currentTokens and replies with an error code, the index
// of the token where parsing failed, and a message, or something.
func (i *Interpreter) Execute() (errorCode, badTokenIndex int, message string) {
	// 1. Try variable assignment.  Must be at least 4 tokens.
	if len(i.currentTokens) >= 4 {
		// First 2 tokens must be identifier literal followed by = (equal) or := (assign)
		if i.currentTokens[0].TokenType == IdentifierLiteral &&
			(i.currentTokens[1].TokenType == Equal || i.currentTokens[1].TokenType == Assign) {
			// If exactly four tokens and the 3rd token is a numerical literal then we don't
			// have anything to evaluate
			if len(i.currentTokens) == 4 && i.currentTokens[2].TokenType == NumericalLiteral {
				if val, err := strconv.ParseFloat(i.currentTokens[2].Literal, 64); err == nil {
					i.store[i.currentTokens[0].Literal] = val
				} else {
					return 1, 0, "Could not interpret this as a number"
				}
			} else {
				// evaluate result then store
				_, _, _, result := i.Evaluate(i.currentTokens[2:])
				i.store[i.currentTokens[0].Literal] = result
			}
		}
	}
	return 0, 0, "Expected a keyword, line number, expression, variable assignment or procedure call"
}
