package rmbasicx64

import "github.com/adamstimb/rmbasicx64/internal/app/rmbasicx64/syntaxerror"

// rmCls represents the CLS command
// TODO: parameters
func (i *Interpreter) RmCls() (ok bool) {
	// Ensure no parameters
	i.TokenPointer++
	if !i.EndOfTokens() {
		i.ErrorCode = syntaxerror.EndOfInstructionExpected
		i.BadTokenIndex = i.TokenPointer
		return false
	}
	// execute
	i.g.Cls()
	i.g.SetCurpos(1, 1)
	return true
}
