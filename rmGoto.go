package main

import (
	"math"
)

// rmGoto represents the GOTO command
func (i *Interpreter) rmGoto() (ok bool) {
	// Check that a parameter is passed
	if len(i.tokenStack) < 3 {
		i.errorCode = NumericExpressionNeeded
		i.badTokenIndex = 1
		i.message = errorMessage(NumericExpressionNeeded)
		return false
	}
	// Get gotoLine
	if i.tokenStack[1].TokenType == NumericalLiteral || i.tokenStack[1].TokenType == IdentifierLiteral {
		i.tokenPointer++
		val, ok := i.AcceptAnyNumber()
		if !ok {
			return false
		}
		gotoLine := int(math.Round(val))
		// scan for line number to goto
		lineOrder := i.GetLineOrder()
		for l := 0; l < len(lineOrder); l++ {
			if lineOrder[l] == gotoLine {
				// found the line so set pointer to one behind because RunLine advances it and return
				i.programPointer = l - 1
				return true
			}
		}
		// line does not exist
		i.errorCode = SpecifiedLineNotFound
		i.badTokenIndex = 2
		i.message = errorMessage(SpecifiedLineNotFound)
		return false
	}
	i.errorCode = NumericExpressionNeeded
	i.badTokenIndex = 1
	i.message = errorMessage(NumericExpressionNeeded)
	return false
}
