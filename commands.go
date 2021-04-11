package main

import (
	"fmt"
	"math"
)

// rmAssign represents a variable assignment (var = expr or var := expr)
func (i *Interpreter) rmAssign() (ok bool) {
	// Catch case where a keyword has been used as a variable name to assign to
	if IsKeyword(i.tokenStack[0]) &&
		(i.tokenStack[1].TokenType == Equal || i.tokenStack[1].TokenType == Assign) {
		i.errorCode = IsAKeywordAndCannotBeUsedAsAVariableName
		i.badTokenIndex = 0
		i.message = fmt.Sprintf("%s%s", i.tokenStack[0].Literal, errorMessage(IsAKeywordAndCannotBeUsedAsAVariableName))
		return false
	}
	// advance token point, extract expression, evaluate result then store
	i.tokenPointer += 2
	result, ok := i.EvaluateExpression(i.ExtractExpression())
	if ok {
		// Evaluation was successful so check data type and store
		if i.SetVar(i.tokenStack[0].Literal, result) {
			return true
		} else {
			return false
		}
	} else {
		// Something went wrong in the evaluation
		return false
	}
}

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
	// Validate is integer or variable representing an integer
	// Don't accept string vars
	if IsStringVar(tokens[1]) {
		i.errorCode = LineNumberExpected
		i.message = errorMessage(LineNumberExpected)
		i.badTokenIndex = 1
		return false
	}
	// Get gotoLine
	if tokens[1].TokenType == NumericalLiteral || tokens[1].TokenType == IdentifierLiteral {
		val, ok := i.GetValueFromToken(tokens[1], "float64")
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

// rmPrint represents the Print command
func (i *Interpreter) rmPrint(tokens []Token) (ok bool) {
	// PRINT with no args
	if len(tokens) == 1 {
		fmt.Println("")
		return true
	}
	if len(tokens) > 1 {
		if tokens[1].TokenType == EndOfLine {
			// Also PRINT with no args
			fmt.Println("")
			return true
		}
		if tokens[1].TokenType == StringLiteral || tokens[1].TokenType == IdentifierLiteral || tokens[1].TokenType == NumericalLiteral {
			i.tokenPointer++
			toPrint, ok := i.EvaluateExpression(i.ExtractExpression())
			if !ok {
				i.badTokenIndex = 1
				return false
			} else {
				switch GetType(toPrint) {
				case "string":
					fmt.Println(toPrint.(string))
					return true
				case "float64":
					fmt.Println(toPrint.(float64))
					return true
				}
			}
		}
	}
	// set error status here
	return false
}
