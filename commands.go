package main

import (
	"fmt"
	"math"
	"strconv"
)

// rmRun represents the RUN command
// TODO: Run can accept one optional parameter for start-from line number
func (i *Interpreter) rmRun() (ok bool) {
	i.programPointer = 0
	lineOrder := i.GetLineOrder()
	for i.programPointer < len(lineOrder) {
		ok := i.RunLine(i.program[lineOrder[i.programPointer]])
		if !ok {
			i.programPointer = 0
			return false
		}
	}
	i.programPointer = 0
	return true
}

// rmGoto represents the GOTO command
func (i *Interpreter) rmGoto(tokens []Token) (ok bool) {
	// GOTO must be followed by one integer literal that represents a stored line number.
	// Validate only 1 parameter
	if len(tokens) > 3 {
		i.errorCode = TooManyParametersFor
		i.badTokenIndex = 2
		i.message = fmt.Sprintf("%s%s", errorMessage(TooManyParametersFor), "GOTO")
		return false
	}
	if len(tokens) < 3 {
		i.errorCode = NotEnoughParametersFor
		i.badTokenIndex = 0
		i.message = fmt.Sprintf("%s%s", errorMessage(NotEnoughParametersFor), "GOTO")
		return false
	}
	// Validate is integer
	if tokens[1].TokenType == NumericalLiteral {
		if valfloat64, err := strconv.ParseFloat(tokens[1].Literal, 64); err == nil {
			if valfloat64 == math.Round(valfloat64) {
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
		}
	}
	i.errorCode = LineNumberExpected
	i.badTokenIndex = 1
	i.message = errorMessage(LineNumberExpected)
	return false
}

// rmPrint represents the Print command
func (i *Interpreter) rmPrint(tokens []Token) (ok bool) {
	// PRINT with no args
	if len(tokens) == 1 {
		fmt.Println("")
		return true
	}
	if len(tokens) > 1 {
		if tokens[1].TokenType == EndOfLine {
			// Still PRINT with no args
			fmt.Println("")
			return true
		}
		if tokens[1].TokenType == StringLiteral {
			// PRINT "hello"
			fmt.Println(tokens[1].Literal)
			return true
		}
	}
	// set error status here
	return false
}
