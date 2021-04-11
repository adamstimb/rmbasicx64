package main

import (
	"fmt"
	"math"
)

// rmGoto represents the GOTO command
func (i *Interpreter) rmGoto() (ok bool) {
	// Check that a parameter is passed
	if len(i.tokenStack) < 3 {
		i.errorCode = NotEnoughParametersFor
		i.badTokenIndex = 0
		i.message = fmt.Sprintf("%s%s", errorMessage(NotEnoughParametersFor), "GOTO")
		return false
	}
	// Don't accept string vars
	if IsStringVar(i.tokenStack[1]) {
		i.errorCode = LineNumberExpected
		i.message = errorMessage(LineNumberExpected)
		i.badTokenIndex = 1
		return false
	}
	// Get gotoLine
	if i.tokenStack[1].TokenType == NumericalLiteral || i.tokenStack[1].TokenType == IdentifierLiteral {
		i.tokenPointer++
		gotoLineExpression := i.ExtractExpression()
		// Validate no more tokens to evaluate
		if !i.EndOfTokens() {
			i.errorCode = TooManyParametersFor
			i.badTokenIndex = 2
			i.message = fmt.Sprintf("%s%s", errorMessage(TooManyParametersFor), "GOTO")
			return false
		}
		val, ok := i.EvaluateExpression(gotoLineExpression)
		if !ok {
			// broken expression
			return false
		}
		if GetType(val) == "string" {
			// string expressions not allowed
			i.errorCode = LineNumberExpected
			i.message = errorMessage(LineNumberExpected)
			i.badTokenIndex = 1
			return false
		}
		valfloat64 := val.(float64)
		if valfloat64 != math.Round(valfloat64) {
			// only whole numbers allowed
			i.errorCode = LineNumberExpected
			i.message = errorMessage(LineNumberExpected)
			i.badTokenIndex = 1
			return false
		}
		// scan for line number to goto
		gotoLine := int(valfloat64)
		lineOrder := i.GetLineOrder()
		for l := 0; l < len(lineOrder); l++ {
			if lineOrder[l] == gotoLine {
				// found the line so set pointer to one behind because RunLine advances it and return
				i.programPointer = l - 1
				return true
			}
		}
		// line does not exist
		i.errorCode = LineNumberDoesNotExist
		i.badTokenIndex = 2
		i.message = errorMessage(LineNumberDoesNotExist)
		return false
	}
	i.errorCode = LineNumberExpected
	i.badTokenIndex = 1
	i.message = errorMessage(LineNumberExpected)
	return false
}
