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
	result, ok := i.EvaluateExpression()
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
func (i *Interpreter) rmRun() (ok bool) {
	// pass thru is no program
	if len(i.program) == 0 {
		return true
	}
	i.programPointer = 0
	lineOrder := i.GetLineOrder()
	// Check for optional startFrom parameter
	startFrom := lineOrder[i.programPointer]
	if len(i.tokenStack) > 2 {
		i.tokenPointer++
		_, ok := i.AcceptAnyOfTheseTokens([]int{NumericalLiteral, IdentifierLiteral})
		if ok {
			i.tokenPointer--
			val, ok := i.AcceptAnyNumber()
			if ok {
				startFrom = int(math.Round(val))
			} else {
				return false
			}
		} else {
			return false
		}
	}
	// Run the program and if startFrom was passed scan ahead to it before executing
	scanAhead := true
	for i.programPointer < len(lineOrder) {
		i.lineNumber = lineOrder[i.programPointer]
		if scanAhead && (i.lineNumber != startFrom) {
			i.programPointer++
			continue
		}
		if scanAhead && (i.lineNumber == startFrom) {
			scanAhead = false
		}
		ok := i.RunLine(i.program[lineOrder[i.programPointer]])
		if !ok {
			i.programPointer = 0
			return false
		}
	}
	// Catch line number not found
	if scanAhead {
		i.errorCode = SpecifiedLineNotFound
		i.message = errorMessage(SpecifiedLineNotFound)
		i.badTokenIndex = 1
		i.lineNumber = -1
		return false
	}
	// Otherwise ok
	i.programPointer = 0
	return true
}
