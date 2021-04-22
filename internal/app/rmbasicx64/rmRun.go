package rmbasicx64

import (
	"math"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
)

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
		i.BadTokenIndex = 1
		i.LineNumber = -1
		return false
	}
	// Otherwise ok
	i.ProgramPointer = 0
	return true
}
