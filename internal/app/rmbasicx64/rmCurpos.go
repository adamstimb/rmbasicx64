package rmbasicx64

import (
	"math"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
)

// RmSetCurpos represents the SET CURPOS command
func (i *Interpreter) RmSetCurpos() (ok bool) {
	// Ensure a parameter is passed
	i.TokenPointer += 2
	if i.EndOfTokens() {
		i.ErrorCode = syntaxerror.NumericExpressionNeeded
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	// Get the col
	val, ok := i.AcceptAnyNumber()
	if !ok {
		i.ErrorCode = syntaxerror.NumericExpressionNeeded
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	// TODO: ensure within range --> have a function to do that
	col := int(math.Round(val))
	// Get the comma
	_, ok = i.AcceptAnyOfTheseTokens([]int{int(token.Comma)})
	if !ok {
		i.ErrorCode = syntaxerror.CommaSeparatorIsNeeded
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	// Get the row
	val, ok = i.AcceptAnyNumber()
	if !ok {
		i.ErrorCode = syntaxerror.NumericExpressionNeeded
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	// TODO: ensure within range --> have a function to do that
	row := int(math.Round(val))
	// Ensure no more parameters
	if !i.EndOfTokens() {
		i.ErrorCode = syntaxerror.EndOfInstructionExpected
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	// Execute
	i.g.SetCurpos(col, row)
	return true
}
