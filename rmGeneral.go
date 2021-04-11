package main

import (
	"fmt"
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
