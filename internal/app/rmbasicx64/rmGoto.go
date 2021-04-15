package rmbasicx64

import (
	"math"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
)

// rmGoto represents the GOTO command
func (i *Interpreter) RmGoto() (ok bool) {
	// Check that a parameter is passed
	if len(i.TokenStack) < 3 {
		i.ErrorCode = syntaxerror.NumericExpressionNeeded
		i.BadTokenIndex = 1
		i.Message = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
		return false
	}
	// Get gotoLine
	if i.TokenStack[1].TokenType == token.NumericalLiteral || i.TokenStack[1].TokenType == token.IdentifierLiteral {
		i.TokenPointer++
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
				i.ProgramPointer = l - 1
				return true
			}
		}
		// line does not exist
		i.ErrorCode = syntaxerror.SpecifiedLineNotFound
		i.BadTokenIndex = 2
		i.Message = syntaxerror.ErrorMessage(syntaxerror.SpecifiedLineNotFound)
		return false
	}
	i.ErrorCode = syntaxerror.NumericExpressionNeeded
	i.BadTokenIndex = 1
	i.Message = syntaxerror.ErrorMessage(syntaxerror.NumericExpressionNeeded)
	return false
}
