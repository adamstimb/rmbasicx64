package rmbasicx64

import (
	"math"

	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"
	"github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/token"
)

// RmSetCursor represents the SET CURSOR command
// SET CURSOR e1[,e2[,e3][,e4]]]
// e1 === cursorMode, e2 === cursorChar, e3 === cursorCharSet
func (i *Interpreter) RmSetCursor() (ok bool) {
	paramCount := 0
	// Ensure a parameter is passed
	i.TokenPointer += 2
	if i.EndOfTokens() {
		i.ErrorCode = syntaxerror.NumericExpressionNeeded
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	// Get required cursorMode
	val, ok := i.AcceptAnyNumber()
	if !ok {
		i.ErrorCode = syntaxerror.NumericExpressionNeeded
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	// TODO: ensure within range --> have a function to do that
	cursorMode := int(math.Round(val))
	var cursorChar int
	var cursorCharSet int
	paramCount++
	// Get optional cursorChar
	if !i.EndOfTokens() {
		// Get comma
		_, ok = i.AcceptAnyOfTheseTokens([]int{token.Comma})
		if !ok {
			i.ErrorCode = syntaxerror.CommaSeparatorIsNeeded
			i.BadTokenIndex = i.TokenPointer
			return false
		}
		// Get cursorChar
		val, ok = i.AcceptAnyNumber()
		if !ok {
			i.ErrorCode = syntaxerror.NumericExpressionNeeded
			i.BadTokenIndex = i.TokenPointer
			return false
		}
		cursorChar = int(math.Round(val))
		paramCount++
	}
	// Get optional cursorCharSet
	if !i.EndOfTokens() {
		// Get comma
		_, ok = i.AcceptAnyOfTheseTokens([]int{token.Comma})
		if !ok {
			i.ErrorCode = syntaxerror.CommaSeparatorIsNeeded
			i.BadTokenIndex = i.TokenPointer
			return false
		}
		// Get cursorCharSet
		val, ok = i.AcceptAnyNumber()
		if !ok {
			i.ErrorCode = syntaxerror.NumericExpressionNeeded
			i.BadTokenIndex = i.TokenPointer
			return false
		}
		cursorCharSet = int(math.Round(val))
		paramCount++
	}
	// Ensure no more parameters
	if !i.EndOfTokens() {
		i.ErrorCode = syntaxerror.EndOfInstructionExpected
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	// Execute
	switch paramCount {
	case 1:
		i.g.SetCursor(cursorMode)
	case 2:
		i.g.SetCursor(cursorMode, cursorChar)
	case 3:
		i.g.SetCursor(cursorMode, cursorChar, cursorCharSet)
	}
	return true
}
