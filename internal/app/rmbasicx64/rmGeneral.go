package rmbasicx64

import (
	"fmt"
	"math"
	"os"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
)

// rmBye represents the BYE command
func (i *Interpreter) RmBye() (ok bool) {
	// Ensure no parameters
	i.TokenPointer++
	if !i.EndOfTokens() {
		i.ErrorCode = syntaxerror.EndOfInstructionExpected
		i.Message = syntaxerror.ErrorMessage(syntaxerror.EndOfInstructionExpected)
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	// Exit application with success code
	os.Exit(0)
	// Return never actually happens
	return true
}

// rmAssign represents a variable assignment (var = expr or var := expr)
func (i *Interpreter) RmAssign() (ok bool) {
	// Catch case where a keyword has been used as a variable name to assign to
	if IsKeyword(i.TokenStack[0]) &&
		(i.TokenStack[1].TokenType == token.Equal || i.TokenStack[1].TokenType == token.Assign) {
		i.ErrorCode = syntaxerror.IsAKeywordAndCannotBeUsedAsAVariableName
		i.BadTokenIndex = 0
		i.Message = fmt.Sprintf("%s%s", i.TokenStack[0].Literal, syntaxerror.ErrorMessage(syntaxerror.IsAKeywordAndCannotBeUsedAsAVariableName))
		return false
	}
	// advance token point, extract expression, evaluate result then store
	i.TokenPointer += 2
	result, ok := i.EvaluateExpression()
	if ok {
		// Evaluation was successful so check data type and store
		if i.SetVar(i.TokenStack[0].Literal, result) {
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
func (i *Interpreter) RmRun() (ok bool) {
	// pass thru is no program
	if len(i.Program) == 0 {
		return true
	}
	i.ProgramPointer = 0
	lineOrder := i.GetLineOrder()
	// Check for optional startFrom parameter
	startFrom := lineOrder[i.ProgramPointer]
	if len(i.TokenStack) > 2 {
		i.TokenPointer++
		_, ok := i.AcceptAnyOfTheseTokens([]int{token.NumericalLiteral, token.IdentifierLiteral})
		if ok {
			i.TokenPointer--
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
	for i.ProgramPointer < len(lineOrder) {
		i.LineNumber = lineOrder[i.ProgramPointer]
		if scanAhead && (i.LineNumber != startFrom) {
			i.ProgramPointer++
			continue
		}
		if scanAhead && (i.LineNumber == startFrom) {
			scanAhead = false
		}
		ok := i.RunLine(i.Program[lineOrder[i.ProgramPointer]])
		if !ok {
			i.ProgramPointer = 0
			return false
		}
	}
	// Catch line number not found
	if scanAhead {
		i.ErrorCode = syntaxerror.SpecifiedLineNotFound
		i.Message = syntaxerror.ErrorMessage(syntaxerror.SpecifiedLineNotFound)
		i.BadTokenIndex = 1
		i.LineNumber = -1
		return false
	}
	// Otherwise ok
	i.ProgramPointer = 0
	return true
}
