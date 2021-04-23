package rmbasicx64

import "github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"

// RmNew represents the NEW command
func (i *Interpreter) RmNew() (ok bool) {
	i.TokenPointer++
	// No more parameters accepted
	if !i.EndOfTokens() {
		i.ErrorCode = syntaxerror.EndOfInstructionExpected
		i.BadTokenIndex = 1
		return false
	}
	// just initialize interpreter
	i.Init(i.g)
	return true
}
