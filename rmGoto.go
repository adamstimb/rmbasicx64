package main

import (
	"fmt"
	"math"
)

// rmGoto represents the GOTO command
func (i *Interpreter) rmGoto() (ok bool) {
	// GOTO must be followed by one integer literal that represents a stored line number.
	// Validate only 1 parameter
	if len(i.tokenStack) > 3 {
		i.errorCode = TooManyParametersFor
		i.badTokenIndex = 2
		i.message = fmt.Sprintf("%s%s", errorMessage(TooManyParametersFor), "GOTO")
		return false
	}
	if len(i.tokenStack) < 3 {
		i.errorCode = NotEnoughParametersFor
		i.badTokenIndex = 0
		i.message = fmt.Sprintf("%s%s", errorMessage(NotEnoughParametersFor), "GOTO")
		return false
	}
	// Validate is integer or variable representing an integer
	// Don't accept string vars
	if IsStringVar(i.tokenStack[1]) {
		i.errorCode = LineNumberExpected
		i.message = errorMessage(LineNumberExpected)
		i.badTokenIndex = 1
		return false
	}
	// Get gotoLine
	if i.tokenStack[1].TokenType == NumericalLiteral || i.tokenStack[1].TokenType == IdentifierLiteral {
		val, ok := i.GetValueFromToken(i.tokenStack[1], "float64")
		if !ok {
			i.errorCode = LineNumberExpected
			i.message = errorMessage(LineNumberExpected)
			i.badTokenIndex = 1
			return false
		}
		valfloat64 := val.(float64)
		if valfloat64 != math.Round(valfloat64) {
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
