package rmbasicx64

import (
	"math"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
)

// RmSetMode represents the SET MODE command
func (i *Interpreter) RmSetMode() (ok bool) {
	// Ensure a parameter is passed
	i.TokenPointer += 2
	if i.EndOfTokens() {
		i.ErrorCode = syntaxerror.NumericExpressionNeeded
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	// Get the paper colour
	val, ok := i.AcceptAnyNumber()
	if !ok {
		i.ErrorCode = syntaxerror.NumericExpressionNeeded
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	// TODO: ensure within range --> have a function to do that
	mode := int(math.Round(val))
	// Ensure no more parameters
	if !i.EndOfTokens() {
		i.ErrorCode = syntaxerror.EndOfInstructionExpected
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	// Execute
	i.g.SetMode(mode)
	return true
}
